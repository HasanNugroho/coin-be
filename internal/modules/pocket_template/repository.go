package pocket_template

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
	pocketTemplates *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		pocketTemplates: db.Collection("pocket_templates"),
	}
}

func (r *Repository) CreatePocketTemplate(ctx context.Context, template *PocketTemplate) error {
	template.ID = primitive.NewObjectID()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	_, err := r.pocketTemplates.InsertOne(ctx, template)
	return err
}

func (r *Repository) GetPocketTemplateByID(ctx context.Context, id primitive.ObjectID) (*PocketTemplate, error) {
	var template PocketTemplate
	err := r.pocketTemplates.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("pocket template not found")
		}
		return nil, err
	}
	return &template, nil
}

func (r *Repository) GetPocketTemplateByName(ctx context.Context, name string) (*PocketTemplate, error) {
	var template PocketTemplate
	err := r.pocketTemplates.FindOne(ctx, bson.M{"name": name, "deleted_at": nil}).Decode(&template)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("pocket template not found")
		}
		return nil, err
	}
	return &template, nil
}

func (r *Repository) UpdatePocketTemplate(ctx context.Context, id primitive.ObjectID, template *PocketTemplate) error {
	template.UpdatedAt = time.Now()
	result, err := r.pocketTemplates.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, bson.M{"$set": template})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("pocket template not found")
	}
	return nil
}

func (r *Repository) DeletePocketTemplate(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.pocketTemplates.UpdateOne(
		ctx,
		bson.M{"_id": id, "deleted_at": nil},
		bson.M{
			"$set": bson.M{
				"deleted_at": now,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("pocket template not found")
	}
	return nil
}

func (r *Repository) ListPocketTemplates(ctx context.Context, limit int64, skip int64) ([]*PocketTemplate, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"order": 1})
	cursor, err := r.pocketTemplates.Find(ctx, bson.M{"deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*PocketTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *Repository) ListActivePocketTemplates(ctx context.Context, limit int64, skip int64) ([]*PocketTemplate, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"order": 1})
	cursor, err := r.pocketTemplates.Find(ctx, bson.M{"deleted_at": nil, "is_active": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*PocketTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *Repository) ListPocketTemplatesByType(ctx context.Context, templateType string, limit int64, skip int64) ([]*PocketTemplate, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"order": 1})
	cursor, err := r.pocketTemplates.Find(ctx, bson.M{"deleted_at": nil, "type": templateType}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []*PocketTemplate
	if err = cursor.All(ctx, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}
