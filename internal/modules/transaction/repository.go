package transaction

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
	transactions *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		transactions: db.Collection("transactions"),
	}
}

func (r *Repository) Create(ctx context.Context, transaction *Transaction) error {
	transaction.ID = primitive.NewObjectID()
	transaction.CreatedAt = time.Now()
	_, err := r.transactions.InsertOne(ctx, transaction)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*Transaction, error) {
	var transaction Transaction
	err := r.transactions.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID primitive.ObjectID, limit, skip int64) ([]*Transaction, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "transaction_date", Value: -1}})

	cursor, err := r.transactions.Find(ctx, bson.M{"user_id": userID}, opts)
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

func (r *Repository) Filter(ctx context.Context, userID primitive.ObjectID, filter bson.M, limit, skip int64) ([]*Transaction, error) {
	filter["user_id"] = userID
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "transaction_date", Value: -1}})

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

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, transaction *Transaction) error {
	result, err := r.transactions.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": transaction})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.transactions.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *Repository) GetTotalByUserAndType(ctx context.Context, userID primitive.ObjectID, transactionType TransactionType, startDate, endDate time.Time) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"user_id": userID,
			"type":    transactionType,
			"transaction_date": bson.M{
				"$gte": startDate,
				"$lte": endDate,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
		}}},
	}

	cursor, err := r.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	total, ok := result[0]["total"].(float64)
	if !ok {
		return 0, nil
	}

	return total, nil
}

func (r *Repository) GetByCategory(ctx context.Context, userID primitive.ObjectID, categoryID primitive.ObjectID, startDate, endDate time.Time) ([]*Transaction, error) {
	filter := bson.M{
		"user_id":     userID,
		"category_id": categoryID,
		"transaction_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.transactions.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "transaction_date", Value: -1}}))
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

func (r *Repository) GetByAllocation(ctx context.Context, userID primitive.ObjectID, allocationID primitive.ObjectID, startDate, endDate time.Time) ([]*Transaction, error) {
	filter := bson.M{
		"user_id":       userID,
		"allocation_id": allocationID,
		"transaction_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := r.transactions.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "transaction_date", Value: -1}}))
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
