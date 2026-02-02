package payroll

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	payrollRecords *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		payrollRecords: db.Collection("payroll_records"),
	}
}

func (r *Repository) CreatePayrollRecord(ctx context.Context, record *PayrollRecord) error {
	record.ID = primitive.NewObjectID()
	record.CreatedAt = time.Now()
	_, err := r.payrollRecords.InsertOne(ctx, record)
	return err
}

func (r *Repository) CreatePayrollRecordBulk(ctx context.Context, records []*PayrollRecord) error {
	if len(records) == 0 {
		return nil
	}

	// Set IDs and timestamps for all records
	recordInterfaces := make([]interface{}, len(records))
	for i, record := range records {
		record.ID = primitive.NewObjectID()
		record.CreatedAt = time.Now()
		recordInterfaces[i] = record
	}

	// Insert all records in bulk
	_, err := r.payrollRecords.InsertMany(ctx, recordInterfaces)
	return err
}

func (r *Repository) GetPayrollRecord(ctx context.Context, userID primitive.ObjectID, year, month, day int) (*PayrollRecord, error) {
	var record PayrollRecord
	err := r.payrollRecords.FindOne(ctx, bson.M{
		"user_id": userID,
		"year":    year,
		"month":   month,
		"day":     day,
	}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *Repository) GetPayrollRecordByUserIDs(ctx context.Context, userIDs []primitive.ObjectID, year, month, day int) ([]*PayrollRecord, error) {
	var records []*PayrollRecord
	cursor, err := r.payrollRecords.Find(ctx, bson.M{
		"user_id": bson.M{"$in": userIDs},
		"year":    year,
		"month":   month,
		"day":     day,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *Repository) GetUserPayrollRecords(ctx context.Context, userID primitive.ObjectID, limit int64) ([]*PayrollRecord, error) {
	cursor, err := r.payrollRecords.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []*PayrollRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}
