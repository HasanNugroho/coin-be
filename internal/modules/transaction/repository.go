package transaction

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	transactions *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		transactions: db.Collection("transactions"),
	}
}

func (r *Repository) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	transaction.ID = primitive.NewObjectID()
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
	_, err := r.transactions.InsertOne(ctx, transaction)
	return err
}

func (r *Repository) GetTransactionByID(ctx context.Context, id primitive.ObjectID) (*Transaction, error) {
	var transaction Transaction
	err := r.transactions.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *Repository) GetTransactionsByUserID(ctx context.Context, userID primitive.ObjectID, limit int64, skip int64) ([]*Transaction, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"date": -1})
	cursor, err := r.transactions.Find(ctx, bson.M{"user_id": userID, "deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *Repository) GetTransactionsByUserIDWithSort(
	ctx context.Context,
	userID primitive.ObjectID,
	txType *string,
	search *string,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*Transaction, int64, error) {

	filter := bson.M{
		"user_id":    userID,
		"deleted_at": nil,
	}

	if txType != nil && *txType != "" {
		filter["type"] = *txType
	}

	if search != nil && *search != "" {
		keyword := regexp.QuoteMeta(strings.TrimSpace(*search))

		filter["$or"] = []bson.M{
			{"note": bson.M{"$regex": keyword, "$options": "i"}},
			{"ref": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	// total count
	total, err := r.transactions.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// sort order
	sortValue := int32(-1)
	if sortOrder == "asc" {
		sortValue = 1
	}

	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	sort := bson.D{{Key: sortBy, Value: sortValue}}

	opts := options.Find().
		SetLimit(pageSize).
		SetSkip(skip).
		SetSort(sort)

	cursor, err := r.transactions.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var transactions []*Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *Repository) GetTransactionsByPocketID(ctx context.Context, pocketID primitive.ObjectID, limit int64, skip int64) ([]*Transaction, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"date": -1})
	filter := bson.M{
		"deleted_at": nil,
		"$or": []bson.M{
			{"pocket_from": pocketID},
			{"pocket_to": pocketID},
		},
	}
	cursor, err := r.transactions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *Repository) GetTransactionsByPocketIDWithSort(ctx context.Context, pocketID primitive.ObjectID, page int64, pageSize int64, sortBy string, sortOrder string) ([]*Transaction, int64, error) {
	filter := bson.M{
		"deleted_at": nil,
		"$or": []bson.M{
			{"pocket_from": pocketID},
			{"pocket_to": pocketID},
		},
	}

	// Get total count
	total, err := r.transactions.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort order
	sortValue := int32(-1) // default descending
	if sortOrder == "asc" {
		sortValue = 1
	}

	// Calculate skip
	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	opts := options.Find().
		SetLimit(pageSize).
		SetSkip(skip).
		SetSort(bson.M{sortBy: sortValue})

	cursor, err := r.transactions.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var transactions []*Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, 0, err
	}
	return transactions, total, nil
}

func (r *Repository) CountUserTransactions(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	count, err := r.transactions.CountDocuments(ctx, bson.M{"user_id": userID, "deleted_at": nil})
	return count, err
}

func (r *Repository) DeleteTransaction(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	result, err := r.transactions.UpdateOne(
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
		return errors.New("transaction not found")
	}
	return nil
}

func (r *Repository) GetTransactionsByUserIDAndType(ctx context.Context, userID primitive.ObjectID, txType string, limit int64, skip int64) ([]*Transaction, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"date": -1})
	cursor, err := r.transactions.Find(ctx, bson.M{"user_id": userID, "type": txType, "deleted_at": nil}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "note", Value: "text"},
				{Key: "ref", Value: "text"},
			},
			Options: options.Index().
				SetName("idx_transactions_note_ref_text"),
		},
	}

	_, err := r.transactions.Indexes().CreateMany(ctx, indexes)
	return err
}
