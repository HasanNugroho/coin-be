package category

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
		categories: db.Collection("categories"),
	}
}

func (r *Repository) CreateCategory(ctx context.Context, category *Category) error {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	category.IsDeleted = false
	_, err := r.categories.InsertOne(ctx, category)
	return err
}

func (r *Repository) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (*Category, error) {
	var category Category
	err := r.categories.FindOne(ctx, bson.M{"_id": id, "is_deleted": false}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryByName(ctx context.Context, name string) (*Category, error) {
	var category Category
	err := r.categories.FindOne(ctx, bson.M{"name": name, "is_deleted": false}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, id primitive.ObjectID, category *Category) error {
	category.UpdatedAt = time.Now()
	result, err := r.categories.UpdateOne(ctx, bson.M{"_id": id, "is_deleted": false}, bson.M{"$set": category})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *Repository) DeleteCategory(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.categories.UpdateOne(
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
		return errors.New("category not found")
	}
	return nil
}

func (r *Repository) ListCategories(ctx context.Context, limit int64, skip int64) ([]*Category, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := r.categories.Find(ctx, bson.M{"is_deleted": false}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) ListCategoriesByType(ctx context.Context, categoryType string, limit int64, skip int64) ([]*Category, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := r.categories.Find(ctx, bson.M{"type": categoryType, "is_deleted": false}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) ListSubcategories(ctx context.Context, parentID primitive.ObjectID) ([]*Category, error) {
	cursor, err := r.categories.Find(ctx, bson.M{"parent_id": parentID, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) CountCategories(ctx context.Context) (int64, error) {
	return r.categories.CountDocuments(ctx, bson.M{"is_deleted": false})
}
