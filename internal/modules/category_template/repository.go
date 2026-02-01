package category_template

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	templates *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		templates: db.Collection("category_templates"),
	}
}

func (r *Repository) Create(ctx context.Context, template *CategoryTemplate) error {
	template.ID = primitive.NewObjectID()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.IsDeleted = false
	_, err := r.templates.InsertOne(ctx, template)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*CategoryTemplate, error) {
	var template CategoryTemplate
	err := r.templates.FindOne(ctx, bson.M{"_id": id, "is_deleted": false}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category template not found")
		}
		return nil, err
	}
	return &template, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]*CategoryTemplate, error) {
	cursor, err := r.templates.Find(ctx, bson.M{"user_id": nil, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*CategoryTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *Repository) FindAllWithFilter(ctx context.Context, transactionType *string, page int64, pageSize int64) ([]*CategoryTemplate, int64, error) {
	filter := bson.M{"user_id": nil, "is_deleted": false}

	if transactionType != nil && *transactionType != "" {
		filter["transaction_type"] = *transactionType
	}

	// Get total count
	total, err := r.templates.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip
	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	opts := options.Find().SetSkip(skip).SetLimit(pageSize)
	cursor, err := r.templates.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var templates []*CategoryTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, 0, err
	}
	return templates, total, nil
}

func (r *Repository) FindAllParent(ctx context.Context, transactionType *string) ([]*CategoryTemplate, error) {
	filter := bson.M{
		// "is_default": true,
		"is_deleted": false,
		"parent_id":  nil,
	}

	if transactionType != nil && *transactionType != "" {
		filter["transaction_type"] = *transactionType
	}

	fmt.Println(filter)
	// Projection: hanya ambil id, name, transaction_type
	opts := options.Find().SetProjection(bson.M{
		"_id":              1,
		"name":             1,
		"transaction_type": 1,
	})

	cursor, err := r.templates.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*CategoryTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *Repository) GetDefaults(ctx context.Context) ([]*CategoryTemplate, error) {
	cursor, err := r.templates.Find(ctx, bson.M{"is_default": true, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*CategoryTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]*CategoryTemplate, error) {
	cursor, err := r.templates.Find(ctx, bson.M{"user_id": userID, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*CategoryTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, template *CategoryTemplate) error {
	template.UpdatedAt = time.Now()
	result, err := r.templates.UpdateOne(ctx, bson.M{"_id": id, "is_deleted": false}, bson.M{"$set": template})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("category template not found")
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.templates.UpdateOne(
		ctx,
		bson.M{"_id": id, "is_deleted": false},
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
		return errors.New("category template not found")
	}
	return nil
}
