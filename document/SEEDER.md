# Database Seeder Documentation

## Overview

The seeder is a utility tool that populates the database with default data needed to set up and test the application. It creates:

- **12 Default Categories** (4 income, 8 expense)
- **4 Default Allocations** (Bills, Emergency Fund, Investment, Savings)

## Default Categories

### Income Categories
1. **Salary** 💼 - Regular employment income
2. **Bonus** 🎁 - Bonus payments
3. **Investment** 📈 - Investment returns
4. **Freelance** 💻 - Freelance work income

### Expense Categories
1. **Food & Dining** 🍔 - Restaurants, groceries
2. **Transportation** 🚗 - Gas, public transport, parking
3. **Shopping** 🛍️ - Retail purchases
4. **Bills & Utilities** 💡 - Electricity, water, internet
5. **Entertainment** 🎬 - Movies, games, hobbies
6. **Healthcare** 🏥 - Medical expenses
7. **Education** 📚 - Courses, books, tuition
8. **Personal Care** 💅 - Haircuts, skincare

## Default Allocations

| Name | Priority | Percentage | Target Amount | Purpose |
|------|----------|-----------|----------------|---------|
| Bills & Utilities | 1 | 40% | - | Essential expenses |
| Emergency Fund | 2 | 10% | Rp 10,000,000 | Emergency savings |
| Investment | 3 | 30% | - | Investment growth |
| Savings | 4 | 20% | Rp 5,000,000 | General savings |

## Usage

### Running the Seeder

#### Option 1: Using Make Command
```bash
make seed
```

#### Option 2: Direct Command
```bash
go run cmd/seeder/main.go
```

#### Option 3: Build and Run
```bash
go build -o bin/seeder cmd/seeder/main.go
./bin/seeder
```

### Environment Setup

Ensure your `.env` file is configured with:
```env
MONGO_URI=mongodb://localhost:27017
MONGODB_NAME=coin_db
```

### Execution Flow

1. **Connect to MongoDB** - Establishes connection using MONGO_URI
2. **Verify Connection** - Pings MongoDB to ensure connectivity
3. **Seed Categories** - Inserts 12 default categories (if not already present)
4. **Seed Allocations** - Inserts 4 default allocations (if not already present)
5. **Complete** - Displays success message

### Idempotent Seeding

The seeder is **idempotent**, meaning:
- Running it multiple times is safe
- It checks if data already exists before inserting
- Existing data is never overwritten
- You can run it multiple times without issues

### Output Example

```
Connected to MongoDB successfully
Starting database seeding...
Seeding categories...
Inserted 12 categories
Seeding allocations...
Inserted 4 allocations
Database seeding completed successfully!

✅ Database seeding completed successfully!

Default data has been set up:
- 12 default categories (4 income, 8 expense)
- 4 default allocations (Bills, Emergency Fund, Investment, Savings)

You can now start using the application!
```

## Seeder Structure

### Package: `internal/seeder`

#### `seeder.go`
Main seeder logic:
- `NewSeeder(db *mongo.Database)` - Create seeder instance
- `Seed(ctx context.Context)` - Run all seeders
- `seedCategories(ctx context.Context)` - Seed categories
- `seedAllocations(ctx context.Context)` - Seed allocations

#### `data.go`
Default data definitions:
- `getDefaultCategories()` - Returns 12 default categories
- `getDefaultAllocations()` - Returns 4 default allocations

### Command: `cmd/seeder/main.go`

Entry point for the seeder:
- Loads configuration
- Connects to MongoDB
- Runs seeder
- Displays results

## Customization

### Adding More Seed Data

To add more categories or allocations:

1. **Edit `internal/seeder/data.go`**
   ```go
   func getDefaultCategories() []Category {
       // Add new categories here
   }
   ```

2. **Add seeding method in `internal/seeder/seeder.go`**
   ```go
   func (s *Seeder) seedYourData(ctx context.Context) error {
       // Your seeding logic
   }
   ```

3. **Call in `Seed()` method**
   ```go
   if err := s.seedYourData(ctx); err != nil {
       return fmt.Errorf("error seeding data: %w", err)
   }
   ```

### Modifying Default Values

Edit the data in `internal/seeder/data.go`:

```go
{
    Name:          "Bills & Utilities",
    Priority:      1,
    Percentage:    40,      // Change percentage
    CurrentAmount: 0,
    TargetAmount:  5000000, // Change target
    IsActive:      true,
    CreatedAt:     now,
}
```

## Troubleshooting

### Connection Failed
```
Failed to connect to MongoDB: connection refused
```
**Solution**: Ensure MongoDB is running and MONGO_URI is correct

### Already Seeded
```
Categories already exist, skipping...
Allocations already exist, skipping...
```
**Solution**: This is normal. Data already exists in database.

### Permission Denied
```
Failed to connect to MongoDB: permission denied
```
**Solution**: Check MongoDB credentials in `.env` file

### Database Not Found
```
Failed to ping MongoDB: server selection timeout
```
**Solution**: Verify MongoDB is running and accessible

## Integration with Application Setup

### First Time Setup
```bash
# 1. Start MongoDB
mongod

# 2. Run migrations (if any)
# make migrate

# 3. Run seeder
make seed

# 4. Start application
make dev
```

### Docker Setup
```bash
# Start MongoDB in Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Run seeder
make seed

# Start application
make dev
```

## Database Schema

### Categories Collection
```json
{
  "_id": ObjectId,
  "user_id": null,
  "name": "Salary",
  "type": "income",
  "icon": "💼",
  "color": "#3498db",
  "is_default": true,
  "created_at": ISODate
}
```

### Allocations Collection
```json
{
  "_id": ObjectId,
  "user_id": null,
  "name": "Bills & Utilities",
  "priority": 1,
  "percentage": 40,
  "current_amount": 0,
  "target_amount": null,
  "is_active": true,
  "created_at": ISODate
}
```

## Best Practices

1. **Run seeder after database creation** - Always seed after setting up MongoDB
2. **Run before first deployment** - Ensure default data exists
3. **Don't modify seeded data manually** - Keep defaults consistent
4. **Back up before seeding production** - Always backup production database
5. **Test in development first** - Test seeder in dev environment before production

## Advanced Usage

### Seed Specific Data Only
Create a custom seeder script:

```go
package main

import (
    "context"
    "github.com/HasanNugroho/coin-be/internal/seeder"
)

func main() {
    // ... setup code ...
    s := seeder.NewSeeder(db)
    s.seedCategories(ctx)  // Only seed categories
}
```

### Clear and Reseed
```bash
# Clear collections
mongo coin_db --eval "db.categories.deleteMany({}); db.allocations.deleteMany({})"

# Reseed
make seed
```

## Support

For issues or questions:
1. Check troubleshooting section above
2. Verify MongoDB connection
3. Check `.env` configuration
4. Review logs in `internal/seeder/seeder.go`
