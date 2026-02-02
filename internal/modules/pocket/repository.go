package pocket

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
	pockets *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		pockets: db.Collection("pockets"),
	}
}

func (r *Repository) CreatePocket(ctx context.Context, pocket *Pocket) error {
	pocket.ID = primitive.NewObjectID()
	pocket.CreatedAt = time.Now()
	pocket.UpdatedAt = time.Now()
	_, err := r.pockets.InsertOne(ctx, pocket)
	return err
}

func (r *Repository) CreatePocketBulk(ctx context.Context, pockets []*Pocket) error {
	for _, pocket := range pockets {
		pocket.ID = primitive.NewObjectID()
		pocket.CreatedAt = time.Now()
		pocket.UpdatedAt = time.Now()
	}

	// Convert []*Pocket to []interface{}
	docs := make([]interface{}, len(pockets))
	for i, pocket := range pockets {
		docs[i] = pocket
	}

	_, err := r.pockets.InsertMany(ctx, docs)
	return err
}

func (r *Repository) GetPocketByID(ctx context.Context, id primitive.ObjectID) (*Pocket, error) {
	var pocket Pocket
	err := r.pockets.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&pocket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("pocket not found")
		}
		return nil, err
	}
	return &pocket, nil
}

func (r *Repository) GetPocketByUserIDAndType(ctx context.Context, userID primitive.ObjectID, pocketType string) (*Pocket, error) {
	var pocket Pocket
	err := r.pockets.FindOne(ctx, bson.M{"user_id": userID, "type": pocketType, "deleted_at": nil}).Decode(&pocket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pocket, nil
}

func (r *Repository) GetMainPocketByUserID(ctx context.Context, userID primitive.ObjectID) (*Pocket, error) {
	var pocket Pocket
	err := r.pockets.FindOne(ctx, bson.M{"user_id": userID, "type": string(TypeMain), "deleted_at": nil}).Decode(&pocket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &pocket, nil
}

func (r *Repository) GetPocketsByUserID(ctx context.Context, userID primitive.ObjectID, limit int64, skip int64) ([]*Pocket, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.pockets.Find(ctx, bson.M{"user_id": userID, "deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pockets []*Pocket
	if err = cursor.All(ctx, &pockets); err != nil {
		return nil, err
	}
	return pockets, nil
}

func (r *Repository) GetActivePocketsByUserID(ctx context.Context, userID primitive.ObjectID, limit int64, skip int64) ([]*Pocket, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.pockets.Find(ctx, bson.M{"user_id": userID, "is_active": true, "deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pockets []*Pocket
	if err = cursor.All(ctx, &pockets); err != nil {
		return nil, err
	}
	return pockets, nil
}

func (r *Repository) UpdatePocket(ctx context.Context, id primitive.ObjectID, pocket *Pocket) error {
	pocket.UpdatedAt = time.Now()
	result, err := r.pockets.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, bson.M{"$set": pocket})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("pocket not found")
	}
	return nil
}

func (r *Repository) DeletePocket(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.pockets.UpdateOne(
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
		return errors.New("pocket not found")
	}
	return nil
}

func (r *Repository) GetAllPockets(ctx context.Context, limit int64, skip int64) ([]*Pocket, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.pockets.Find(ctx, bson.M{"deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pockets []*Pocket
	if err = cursor.All(ctx, &pockets); err != nil {
		return nil, err
	}
	return pockets, nil
}

func (r *Repository) CountUserPockets(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := r.pockets.CountDocuments(ctx, bson.M{"user_id": userID, "deleted_at": nil})
	return count, err
}
