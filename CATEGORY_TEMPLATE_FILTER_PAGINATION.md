# Category Template Filter and Pagination

## Overview

The Category Template API now supports filtering by transaction type and pagination for the FindAll endpoint. This allows clients to retrieve category templates with flexible filtering and efficient data retrieval.

## Features

### 1. Filter by Transaction Type
- Filter templates by `transaction_type` (income or expense)
- Optional parameter - if not provided, returns all templates
- Case-sensitive matching

### 2. Pagination
- Page-based pagination with configurable page size
- Default page size: 10 items
- Maximum page size: 100 items (enforced)
- Includes pagination metadata in response

## API Endpoint

### Get All Category Templates with Filter and Pagination

**Endpoint:** `GET /v1/category-templates`

**Authentication:** Required (admin only)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | No | - | Filter by transaction type: "income" or "expense" |
| `page` | integer | No | 1 | Page number (starts from 1) |
| `page_size` | integer | No | 10 | Number of items per page (max 100) |

**Request Examples:**

```bash
# Get all templates
curl -X GET http://localhost:8080/v1/category-templates \
  -H "Authorization: Bearer <token>"

# Get only income templates
curl -X GET "http://localhost:8080/v1/category-templates?type=income" \
  -H "Authorization: Bearer <token>"

# Get expense templates with pagination
curl -X GET "http://localhost:8080/v1/category-templates?type=expense&page=1&page_size=20" \
  -H "Authorization: Bearer <token>"

# Get second page with 15 items per page
curl -X GET "http://localhost:8080/v1/category-templates?page=2&page_size=15" \
  -H "Authorization: Bearer <token>"
```

## Response Format

### Success Response

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category templates retrieved successfully",
  "data": {
    "data": [
      {
        "id": "507f1f77bcf86cd799439011",
        "name": "Salary",
        "transaction_type": "income",
        "is_default": true,
        "color": "#4CAF50",
        "icon": "💰",
        "description": "Monthly salary income",
        "parent_id": null,
        "user_id": null,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z",
        "deleted_at": null
      },
      {
        "id": "507f1f77bcf86cd799439012",
        "name": "Freelance",
        "transaction_type": "income",
        "is_default": true,
        "color": "#2196F3",
        "icon": "💼",
        "description": "Freelance income",
        "parent_id": null,
        "user_id": null,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z",
        "deleted_at": null
      }
    ],
    "meta": {
      "total": 25,
      "page": 1,
      "page_size": 10,
      "total_pages": 3
    }
  }
}
```

### Pagination Metadata

The `meta` object contains:

| Field | Type | Description |
|-------|------|-------------|
| `total` | integer | Total number of templates matching the filter |
| `page` | integer | Current page number |
| `page_size` | integer | Number of items per page |
| `total_pages` | integer | Total number of pages |

### Error Response

```json
{
  "success": false,
  "statusCode": 400,
  "message": "Error description",
  "data": null
}
```

## Implementation Details

### Repository Layer

**File:** `internal/modules/category_template/repository.go`

#### FindAllWithFilter Method

```go
func (r *Repository) FindAllWithFilter(
    ctx context.Context,
    transactionType *string,
    page int64,
    pageSize int64,
) ([]*CategoryTemplate, int64, error)
```

**Features:**
- Builds dynamic filter based on transaction type
- Counts total matching documents
- Applies skip and limit for pagination
- Returns templates and total count

**Filter Logic:**
- Base filter: `{"user_id": nil, "is_deleted": false}`
- If transaction type provided: adds `"transaction_type": <value>`

### Service Layer

**File:** `internal/modules/category_template/service.go`

#### FindAllWithFilter Method

```go
func (s *Service) FindAllWithFilter(
    ctx context.Context,
    transactionType *string,
    page int64,
    pageSize int64,
) ([]*CategoryTemplate, int64, error)
```

**Validation:**
- Page < 1 defaults to 1
- PageSize < 1 defaults to 10
- PageSize > 100 capped at 100

### Controller Layer

**File:** `internal/modules/category_template/controller.go`

#### FindAll Handler

**Query Parameter Parsing:**
```go
transactionType := ctx.Query("type")
page := ctx.DefaultQuery("page", "1")
pageSize := ctx.DefaultQuery("page_size", "10")
```

**Integer Parsing:**
```go
if p, err := strconv.ParseInt(page, 10, 64); err == nil && p > 0 {
    pageNum = p
}
```

**Response Building:**
- Converts templates to DTOs
- Calculates total pages: `(total + pageSize - 1) / pageSize`
- Includes metadata in response

## Usage Examples

### Example 1: Get First Page of Income Templates

```bash
curl -X GET "http://localhost:8080/v1/category-templates?type=income&page=1&page_size=10" \
  -H "Authorization: Bearer eyJhbGc..."
```

**Response:**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Category templates retrieved successfully",
  "data": {
    "data": [
      {"id": "...", "name": "Salary", "transaction_type": "income", ...},
      {"id": "...", "name": "Freelance", "transaction_type": "income", ...}
    ],
    "meta": {
      "total": 15,
      "page": 1,
      "page_size": 10,
      "total_pages": 2
    }
  }
}
```

### Example 2: Get All Expense Templates (No Pagination)

```bash
curl -X GET "http://localhost:8080/v1/category-templates?type=expense" \
  -H "Authorization: Bearer eyJhbGc..."
```

### Example 3: Get Second Page with Custom Page Size

```bash
curl -X GET "http://localhost:8080/v1/category-templates?page=2&page_size=25" \
  -H "Authorization: Bearer eyJhbGc..."
```

## Validation Rules

### Page Parameter
- Must be a positive integer
- Defaults to 1 if invalid or not provided
- Minimum: 1

### Page Size Parameter
- Must be a positive integer
- Defaults to 10 if invalid or not provided
- Minimum: 1
- Maximum: 100 (enforced in service layer)

### Type Parameter
- Optional
- Valid values: "income", "expense"
- Case-sensitive
- If not provided or empty, no type filter applied

## Performance Considerations

### Database Indexes

Recommended indexes for optimal performance:

```javascript
// Index for filtering by transaction type
db.category_templates.createIndex({ "transaction_type": 1, "is_deleted": 1 })

// Index for system templates (user_id: null)
db.category_templates.createIndex({ "user_id": 1, "is_deleted": 1 })

// Compound index for filtered pagination
db.category_templates.createIndex({ 
  "user_id": 1, 
  "transaction_type": 1, 
  "is_deleted": 1 
})
```

### Query Optimization

- Filters are applied at database level
- Only requested fields are returned
- Pagination prevents loading entire collection
- Count is efficient with indexed filters

## Error Handling

### Invalid Page Parameter
- Non-integer values default to page 1
- Negative values default to page 1
- Zero defaults to page 1

### Invalid Page Size Parameter
- Non-integer values default to page size 10
- Negative values default to page size 10
- Zero defaults to page size 10
- Values > 100 are capped at 100

### Invalid Type Parameter
- Non-matching values are ignored
- Empty string is treated as no filter
- Case-sensitive matching

## Testing Recommendations

### Test Cases

1. **No Filters**
   - Request: `GET /v1/category-templates`
   - Expected: All templates with default pagination

2. **Filter by Type**
   - Request: `GET /v1/category-templates?type=income`
   - Expected: Only income templates

3. **Pagination**
   - Request: `GET /v1/category-templates?page=2&page_size=5`
   - Expected: Second page with 5 items

4. **Combined Filter and Pagination**
   - Request: `GET /v1/category-templates?type=expense&page=1&page_size=20`
   - Expected: Income templates, first page, 20 items

5. **Invalid Parameters**
   - Request: `GET /v1/category-templates?page=abc&page_size=xyz`
   - Expected: Defaults applied (page=1, page_size=10)

6. **Page Size Limit**
   - Request: `GET /v1/category-templates?page_size=200`
   - Expected: Page size capped at 100

7. **Out of Range Page**
   - Request: `GET /v1/category-templates?page=999`
   - Expected: Empty data array, correct metadata

## API Changes Summary

### Before
- GET `/v1/category-templates` returned all templates without pagination

### After
- GET `/v1/category-templates` supports:
  - Optional `type` query parameter for filtering
  - Optional `page` query parameter (default: 1)
  - Optional `page_size` query parameter (default: 10, max: 100)
  - Response includes pagination metadata

### Backward Compatibility
- Existing requests without parameters still work
- Default pagination (page 1, size 10) applied automatically
- Response structure enhanced with `meta` object

## Files Modified

1. **`internal/modules/category_template/repository.go`**
   - Added `FindAllWithFilter` method
   - Added `mongo/options` import

2. **`internal/modules/category_template/service.go`**
   - Added `FindAllWithFilter` method with validation

3. **`internal/modules/category_template/controller.go`**
   - Updated `FindAll` handler with filter and pagination logic
   - Added `strconv` import
   - Enhanced response with metadata

## Future Enhancements

1. **Additional Filters**
   - Filter by `is_default` status
   - Filter by `parent_id` (subcategories)
   - Search by name

2. **Sorting**
   - Sort by name, creation date, update date
   - Ascending/descending order

3. **Advanced Pagination**
   - Cursor-based pagination
   - Offset-based pagination

4. **Response Optimization**
   - Selective field projection
   - Nested resource expansion

## Summary

The Category Template API now provides powerful filtering and pagination capabilities:
- ✅ Filter by transaction type
- ✅ Page-based pagination
- ✅ Configurable page size (1-100)
- ✅ Pagination metadata in response
- ✅ Backward compatible
- ✅ Validated parameters with sensible defaults
- ✅ Efficient database queries with recommended indexes
