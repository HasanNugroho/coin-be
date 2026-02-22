package user_category

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	categories *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		categories: db.Collection("user_categories"),
	}
}

func (r *Repository) Create(ctx context.Context, category *UserCategory) error {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.IsDeleted = false
	_, err := r.categories.InsertOne(ctx, category)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*UserCategory, error) {
	var category UserCategory
	err := r.categories.FindOne(ctx, bson.M{"_id": id, "user_id": userID, "is_deleted": false}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]*UserCategory, error) {
	cursor, err := r.categories.Find(ctx, bson.M{"user_id": userID, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*UserCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) FindAllWithFilters(
	ctx context.Context,
	userID primitive.ObjectID,
	txType *string,
	search *string,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*UserCategory, int64, error) {
	match := bson.M{
		"user_id":    userID,
		"is_deleted": false,
	}

	if txType != nil && *txType != "" {
		match["transaction_type"] = *txType
	}

	if search != nil && *search != "" {
		match["name"] = bson.M{"$regex": *search, "$options": "i"}
	}

	total, err := r.categories.CountDocuments(ctx, match)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	sortValue := -1
	if sortOrder == "asc" {
		sortValue = 1
	}

	allowedSort := map[string]bool{
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}

	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}

	opts := options.Find().
		SetLimit(pageSize).
		SetSkip(skip).
		SetSort(bson.M{sortBy: sortValue})

	cursor, err := r.categories.Find(ctx, match, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var categories []*UserCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *Repository) FindAllParent(ctx context.Context, userID primitive.ObjectID, transactionType *string) ([]*UserCategory, error) {
	filter := bson.M{
		"user_id":    userID,
		"is_deleted": false,
		"parent_id":  nil,
	}

	if transactionType != nil && *transactionType != "" {
		filter["transaction_type"] = *transactionType
	}

	opts := options.Find().SetProjection(bson.M{
		"_id":              1,
		"user_id":          1,
		"parent_id":        1,
		"name":             1,
		"transaction_type": 1,
		"color":            1,
		"icon":             1,
	})

	cursor, err := r.categories.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*UserCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) FindAllDropdown(ctx context.Context, userID primitive.ObjectID, transactionType *string) ([]*UserCategory, error) {
	filter := bson.M{
		"user_id":    userID,
		"is_deleted": false,
	}

	if transactionType != nil && *transactionType != "" {
		filter["transaction_type"] = *transactionType
	}

	opts := options.Find().SetProjection(bson.M{
		"_id":              1,
		"user_id":          1,
		"parent_id":        1,
		"name":             1,
		"transaction_type": 1,
		"color":            1,
		"icon":             1,
	})

	cursor, err := r.categories.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*UserCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID, category *UserCategory) error {
	category.UpdatedAt = time.Now()
	result, err := r.categories.UpdateOne(ctx, bson.M{"_id": id, "user_id": userID, "is_deleted": false}, bson.M{"$set": category})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user category not found")
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) error {
	now := time.Now()
	result, err := r.categories.UpdateOne(
		ctx,
		bson.M{"_id": id, "user_id": userID, "is_deleted": false},
		bson.M{
			"$set": bson.M{
				"is_deleted": true,
				"deleted_at": now,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user category not found")
	}
	return nil
}

func (r *Repository) FindByNamesSimilarity(ctx context.Context, userID primitive.ObjectID, names []string, txType *string) ([]*UserCategory, error) {
	if len(names) == 0 {
		return nil, nil
	}

	nameFilters := make([]bson.M, len(names))
	for i, name := range names {
		// Case-insensitive regex match for each name
		nameFilters[i] = bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
	}

	filter := bson.M{
		"user_id":    userID,
		"is_deleted": false,
		"$or":        nameFilters,
	}

	if txType != nil && *txType != "" {
		filter["transaction_type"] = *txType
	}

	cursor, err := r.categories.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*UserCategory
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}
