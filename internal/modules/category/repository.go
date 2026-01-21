package category

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	categories *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		categories: db.Collection("categories"),
	}
}

func (r *Repository) Create(ctx context.Context, category *Category) error {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	_, err := r.categories.InsertOne(ctx, category)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*Category, error) {
	var category Category
	err := r.categories.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*Category, error) {
	cursor, err := r.categories.Find(ctx, bson.M{"user_id": userID})
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

func (r *Repository) GetByUserIDAndType(ctx context.Context, userID primitive.ObjectID, categoryType CategoryType) ([]*Category, error) {
	cursor, err := r.categories.Find(ctx, bson.M{"user_id": userID, "type": categoryType})
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

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, category *Category) error {
	result, err := r.categories.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": category})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.categories.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *Repository) CreateDefaultCategories(ctx context.Context, userID primitive.ObjectID) error {
	defaultCategories := []Category{
		{UserID: userID, Name: "Salary", Type: CategoryTypeIncome, Icon: "💰", Color: "#4CAF50", IsDefault: true},
		{UserID: userID, Name: "Bonus", Type: CategoryTypeIncome, Icon: "🎁", Color: "#8BC34A", IsDefault: true},
		{UserID: userID, Name: "Investment", Type: CategoryTypeIncome, Icon: "📈", Color: "#00BCD4", IsDefault: true},
		{UserID: userID, Name: "Food & Dining", Type: CategoryTypeExpense, Icon: "🍔", Color: "#FF9800", IsDefault: true},
		{UserID: userID, Name: "Transportation", Type: CategoryTypeExpense, Icon: "🚗", Color: "#2196F3", IsDefault: true},
		{UserID: userID, Name: "Shopping", Type: CategoryTypeExpense, Icon: "🛍️", Color: "#E91E63", IsDefault: true},
		{UserID: userID, Name: "Bills & Utilities", Type: CategoryTypeExpense, Icon: "💡", Color: "#9C27B0", IsDefault: true},
		{UserID: userID, Name: "Entertainment", Type: CategoryTypeExpense, Icon: "🎬", Color: "#FF5722", IsDefault: true},
	}

	for _, cat := range defaultCategories {
		if err := r.Create(ctx, &cat); err != nil {
			return err
		}
	}
	return nil
}
