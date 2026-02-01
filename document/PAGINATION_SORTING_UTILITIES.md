# Pagination and Sorting Utilities

## Overview

A comprehensive set of utilities for handling pagination and sorting across the application. These utilities provide consistent pagination parameters parsing, sorting validation, and response formatting.

## Location

**File:** `internal/core/utils/pagination.go`

## Core Components

### 1. PaginationParams

Holds pagination parameters extracted from query strings.

```go
type PaginationParams struct {
	Page     int64
	PageSize int64
}
```

### 2. SortParams

Holds sorting parameters with field name and order.

```go
type SortParams struct {
	SortBy    string // field name to sort by
	SortOrder string // "asc" or "desc"
}
```

### 3. PaginationMeta

Holds pagination metadata for API responses.

```go
type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
}
```

## Utility Functions

### ParsePaginationParams

Extracts and validates pagination parameters from query strings.

```go
func ParsePaginationParams(ctx *gin.Context, defaultPageSize int64) PaginationParams
```

**Parameters:**
- `ctx` - Gin context
- `defaultPageSize` - Default page size if not provided

**Query Parameters:**
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: defaultPageSize, max: 100)

**Validation:**
- Page < 1 defaults to 1
- PageSize < 1 defaults to defaultPageSize
- PageSize > 100 capped at 100

**Returns:** `PaginationParams` with validated values

### ParseSortParams

Extracts and validates sorting parameters from query strings.

```go
func ParseSortParams(ctx *gin.Context, allowedFields []string, defaultField string) SortParams
```

**Parameters:**
- `ctx` - Gin context
- `allowedFields` - List of allowed field names to sort by
- `defaultField` - Default field if not provided or invalid

**Query Parameters:**
- `sort_by` - Field name to sort by (default: defaultField)
- `sort_order` - Sort order: "asc" or "desc" (default: "desc")

**Validation:**
- Only allows fields in `allowedFields` list
- Falls back to `defaultField` if invalid
- Validates sort order is "asc" or "desc"

**Returns:** `SortParams` with validated values

### CalculatePaginationMeta

Calculates pagination metadata including total pages.

```go
func CalculatePaginationMeta(total int64, page int64, pageSize int64) PaginationMeta
```

**Parameters:**
- `total` - Total number of items
- `page` - Current page number
- `pageSize` - Items per page

**Calculation:** `totalPages = (total + pageSize - 1) / pageSize`

**Returns:** `PaginationMeta` with all pagination information

### BuildPaginatedResponse

Builds a standardized paginated response structure.

```go
func BuildPaginatedResponse(data interface{}, meta PaginationMeta) map[string]interface{}
```

**Parameters:**
- `data` - Response data (array of items)
- `meta` - Pagination metadata

**Returns:** Map with `data` and `meta` keys

## Usage Examples

### Example 1: Transaction Controller

```go
func (c *Controller) ListPocketTransactions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	pocketID := ctx.Param("pocket_id")

	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(ctx, 10)

	// Parse sorting parameters (allowed fields: date, amount)
	allowedFields := []string{"date", "amount"}
	sorting := utils.ParseSortParams(ctx, allowedFields, "date")

	// Fetch transactions with pagination and sorting
	transactions, total, err := c.service.GetPocketTransactionsWithSort(
		ctx, 
		userID.(string), 
		pocketID, 
		pagination.Page, 
		pagination.PageSize, 
		sorting.SortBy, 
		sorting.SortOrder,
	)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// Convert to response DTOs
	txsResp := c.mapToResponseList(transactions)

	// Calculate pagination metadata
	meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)

	// Build paginated response
	respData := utils.BuildPaginatedResponse(txsResp, meta)
	resp := utils.NewSuccessResponse("Transactions retrieved successfully", respData)
	ctx.JSON(http.StatusOK, resp)
}
```

### Example 2: Category Template Controller

```go
func (c *Controller) FindAll(ctx *gin.Context) {
	// ... auth check ...

	// Get query parameters
	transactionType := ctx.Query("type")

	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(ctx, 10)

	// Prepare filter
	var typeFilter *string
	if transactionType != "" {
		typeFilter = &transactionType
	}

	// Fetch templates with filter and pagination
	templates, total, err := c.service.FindAllWithFilter(ctx, typeFilter, pagination.Page, pagination.PageSize)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	templateResps := make([]*dto.CategoryTemplateResponse, len(templates))
	for i, template := range templates {
		templateResps[i] = c.mapToResponse(template)
	}

	// Calculate pagination metadata
	meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)

	resp := utils.NewSuccessResponse("Category templates retrieved successfully", gin.H{
		"data": templateResps,
		"meta": meta,
	})
	ctx.JSON(http.StatusOK, resp)
}
```

## API Response Format

### Request

```bash
GET /v1/transactions/pocket/:pocket_id?page=1&page_size=10&sort_by=date&sort_order=desc
```

### Response

```json
{
  "success": true,
  "statusCode": 200,
  "message": "Transactions retrieved successfully",
  "data": {
    "data": [
      {
        "id": "507f1f77bcf86cd799439011",
        "user_id": "507f1f77bcf86cd799439010",
        "type": "income",
        "amount": 5000.00,
        "pocket_from": null,
        "pocket_to": "507f1f77bcf86cd799439012",
        "date": "2024-01-15T10:30:00Z",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
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

## Query Parameters

### Pagination

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `page` | integer | 1 | - | Page number (starts from 1) |
| `page_size` | integer | 10 | 100 | Items per page |

### Sorting

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `sort_by` | string | field-specific | Field to sort by (must be in allowed list) |
| `sort_order` | string | desc | Sort order: "asc" or "desc" |

## Implementation Details

### Transaction Module

**Repository Method:** `GetTransactionsByPocketIDWithSort`
- Filters transactions by pocket ID (pocket_from or pocket_to)
- Supports sorting by "date" or "amount"
- Returns transactions and total count

**Service Method:** `GetPocketTransactionsWithSort`
- Validates user ID and pocket ID
- Checks authorization
- Calls repository method

**Controller Method:** `ListPocketTransactions`
- Parses pagination and sorting parameters
- Calls service method
- Returns paginated response with metadata

### Category Template Module

**Repository Method:** `FindAllWithFilter`
- Filters by transaction type
- Supports pagination
- Returns templates and total count

**Service Method:** `FindAllWithFilter`
- Validates pagination parameters
- Calls repository method

**Controller Method:** `FindAll`
- Parses pagination parameters
- Calls service method
- Returns paginated response with metadata

## Best Practices

### 1. Always Use ParsePaginationParams

```go
pagination := utils.ParsePaginationParams(ctx, 10) // 10 is default page size
```

### 2. Always Use ParseSortParams with Allowed Fields

```go
allowedFields := []string{"date", "amount", "created_at"}
sorting := utils.ParseSortParams(ctx, allowedFields, "date")
```

### 3. Always Calculate and Include Metadata

```go
meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)
```

### 4. Use BuildPaginatedResponse for Consistency

```go
respData := utils.BuildPaginatedResponse(items, meta)
resp := utils.NewSuccessResponse("Message", respData)
```

## Validation Rules

### Page Parameter
- Must be positive integer
- Defaults to 1 if invalid
- Minimum: 1

### Page Size Parameter
- Must be positive integer
- Defaults to provided defaultPageSize if invalid
- Minimum: 1
- Maximum: 100 (enforced)

### Sort By Parameter
- Must be in allowedFields list
- Defaults to defaultField if not in list
- Case-sensitive

### Sort Order Parameter
- Must be "asc" or "desc"
- Defaults to "desc" if invalid

## Performance Considerations

### Database Indexes

For optimal performance, create indexes on commonly sorted fields:

```javascript
// For transaction sorting
db.transactions.createIndex({ "date": -1, "deleted_at": 1 })
db.transactions.createIndex({ "amount": -1, "deleted_at": 1 })

// For category template filtering
db.category_templates.createIndex({ "transaction_type": 1, "is_deleted": 1 })
```

### Query Optimization

- Pagination prevents loading entire collections
- Sorting is applied at database level
- Filters are applied before pagination
- Count is efficient with indexed fields

## Error Handling

### Invalid Parameters

All invalid parameters are handled gracefully:
- Non-integer values default to defaults
- Out-of-range values are capped or defaulted
- Invalid sort fields fall back to default field
- Invalid sort orders default to "desc"

### Service Layer Errors

Service methods return errors for:
- Invalid user ID format
- Invalid pocket ID format
- Unauthorized access
- Resource not found

## Future Enhancements

1. **Cursor-Based Pagination** - For large datasets
2. **Advanced Filtering** - Multiple field filters
3. **Search** - Full-text search support
4. **Aggregation** - Summary statistics with pagination
5. **Export** - Paginated export to CSV/JSON

## Summary

The pagination and sorting utilities provide:
- ✅ Consistent parameter parsing and validation
- ✅ Standardized response format
- ✅ Automatic metadata calculation
- ✅ Security through field validation
- ✅ Performance through database-level operations
- ✅ Easy integration across modules
