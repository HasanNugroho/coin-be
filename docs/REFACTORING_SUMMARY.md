# Allocation Module Refactoring Summary

## Overview
This document summarizes the comprehensive refactoring of the allocation module to separate allocation execution from payroll processing and implement scheduled allocation execution with timezone support.

## Changes Made

### 1. Allocation Schema & Models
**File**: `internal/modules/allocation/models.go`

**Changes**:
- Added `ExecuteDate *time.Time` field to the `Allocation` struct
- This field stores the scheduled execution date for allocations
- Supports both immediate allocations (nil) and scheduled allocations (with date)

### 2. Data Transfer Objects (DTOs)
**Files**: 
- `internal/modules/allocation/dto/request.go`
- `internal/modules/allocation/dto/response.go`

**Changes**:
- Added `ExecuteDate *time.Time` field to `CreateAllocationRequest`
- Added `ExecuteDate *time.Time` field to `UpdateAllocationRequest`
- Added `ExecuteDate *time.Time` field to `AllocationResponse`
- Allows clients to specify and retrieve execution dates via API

### 3. Repository Layer
**File**: `internal/modules/allocation/repository.go`

**New Methods**:
- `GetAllocationsByExecuteDate(ctx context.Context, executeDate time.Time)`: Fetches allocations scheduled for a specific date
- `GetAllocationsByExecuteDateWithLookup(ctx context.Context, executeDate time.Time)`: Aggregation pipeline that fetches allocations with joined user and user_profile data
  - Uses MongoDB aggregation with `$lookup` stages
  - Filters by execute_date, is_active, and non-deleted records
  - Joins with users and user_profiles collections
  - Returns enriched allocation data for processing

### 4. Service Layer
**File**: `internal/modules/allocation/service.go`

**Major Changes**:
- Updated `NewService()` constructor to accept additional dependencies:
  - `userRepo *user.Repository`
  - `transactionRepo *transaction.Repository`
  - `db *mongo.Database`

**New Methods**:
- `ProcessDailyAllocations(ctx context.Context, executeDate time.Time)`: Main entry point for daily allocation processing
  - Fetches all allocations scheduled for the given date
  - Iterates through each allocation and processes execution
  - Tracks success/failure counts
  - Logs results

- `processAllocationExecution(ctx context.Context, userID, allocationID primitive.ObjectID, allocData map[string]interface{})`: Executes a single allocation
  - Validates allocation is active
  - Validates target pocket/platform
  - Creates TRANSFER transaction
  - Updates balances within database transaction
  - Ensures atomicity with MongoDB sessions

- `applyBalanceUpdates(ctx context.Context, updates []balanceUpdate)`: Batch balance updates
  - Consolidates updates by entity
  - Applies updates to pockets and user platforms
  - Avoids N+1 query patterns

- `getMainPocket(ctx context.Context, userID primitive.ObjectID)`: Helper to retrieve user's main pocket

**Updated Methods**:
- `CreateAllocation()`: Now sets `ExecuteDate` from request
- `UpdateAllocation()`: Now handles `ExecuteDate` updates

### 5. Controller Layer
**File**: `internal/modules/allocation/controller.go`

**Changes**:
- Updated `mapToResponse()` to include `ExecuteDate` field in response mapping
- All existing endpoints now support execute_date parameter

### 6. Cron Job Implementation
**File**: `internal/modules/allocation/cron.go` (NEW)

**Features**:
- `CronJob` struct manages scheduled allocation execution
- `Start()`: Initializes cron job to run daily at 01:00 AM (Asia/Jakarta timezone)
- `Stop()`: Gracefully stops the cron job
- `getJakartaLocation()`: Helper function to load Asia/Jakarta timezone
  - Falls back to UTC if timezone loading fails
  - Ensures consistent timezone handling

### 7. Payroll Module Cleanup
**File**: `internal/modules/payroll/service.go`

**Removed**:
- Allocation execution logic (lines 252-359 in original)
- Allocation repository dependency
- `balanceUpdate` type (moved to allocation module)
- `applyBalanceUpdates()` method (moved to allocation module)

**Simplified**:
- `processUserPayroll()` now only handles:
  1. Create INCOME transaction
  2. Update balances for income
  3. Commit transaction
- Removed allocation-related imports

**File**: `internal/modules/payroll/module.go`

**Changes**:
- Removed `allocationRepository` injection
- Updated `NewService()` call to exclude allocation repository

**File**: `internal/modules/payroll/cron.go`

**Changes**:
- Added timezone support (Asia/Jakarta)
- Updated cron initialization with `cron.WithLocation(getJakartaLocation())`
- Added `getJakartaLocation()` helper function
- Cron job now runs at 00:01 AM (Asia/Jakarta timezone)

### 8. Allocation Module Registration
**File**: `internal/modules/allocation/module.go`

**Changes**:
- Updated service registration to include new dependencies:
  - `userRepository`
  - `transactionRepository`
  - `mongo.Database`
- Service now has full context for allocation execution

### 9. Main Application Integration
**File**: `cmd/api/main.go`

**Changes**:
- Added allocation cron job initialization:
  ```go
  allocationService := appContainer.Get("allocationService").(*allocation.Service)
  allocationCronJob := allocation.NewCronJob(allocationService)
  allocationCronJob.Start()
  defer allocationCronJob.Stop()
  ```
- Cron job runs alongside payroll and dashboard jobs
- Ensures graceful shutdown with defer

### 10. API Documentation
**File**: `docs/ALLOCATION_API.md` (NEW)

**Contents**:
- Complete API specification for allocation endpoints
- Data model documentation
- All 5 endpoints documented:
  1. POST /v1/allocations - Create allocation
  2. GET /v1/allocations/{id} - Get allocation
  3. GET /v1/allocations - List allocations
  4. PUT /v1/allocations/{id} - Update allocation
  5. DELETE /v1/allocations/{id} - Delete allocation
- Validation rules and error responses
- Execution process documentation
- Timezone information
- Example workflows
- Status codes reference

## Architecture Changes

### Separation of Concerns
**Before**:
- Payroll module handled both payroll AND allocation execution
- Allocation logic tightly coupled with payroll processing
- Difficult to manage allocation execution independently

**After**:
- Payroll module: Only handles salary income transactions
- Allocation module: Handles allocation execution and scheduling
- Clear separation of responsibilities
- Each module can be tested and deployed independently

### Data Flow

#### Allocation Execution Flow
```
CronJob (01:00 AM Jakarta)
  ↓
AllocationService.ProcessDailyAllocations()
  ↓
Repository.GetAllocationsByExecuteDateWithLookup()
  ↓
For each allocation:
  - Validate allocation & targets
  - Create TRANSFER transaction
  - Update balances (atomic)
  ↓
Log results
```

#### Payroll Processing Flow
```
CronJob (00:01 AM Jakarta)
  ↓
PayrollService.ProcessDailyPayroll()
  ↓
Repository.getEligibleUsersForPayroll()
  ↓
For each eligible user:
  - Create INCOME transaction
  - Update balances
  ↓
Log results
```

## Timezone Implementation

### Consistent Timezone Usage
- **Timezone**: Asia/Jakarta (UTC+7)
- **Applied to**:
  - Payroll cron job
  - Allocation cron job
  - All scheduled execution logic
- **Fallback**: UTC if timezone loading fails
- **API**: All timestamps in ISO 8601 UTC format

### Cron Schedule
- **Payroll**: 00:01 AM (Asia/Jakarta)
- **Allocation**: 01:00 AM (Asia/Jakarta)
- Staggered to avoid resource contention

## Database Transactions

### Allocation Execution
- Uses MongoDB sessions for ACID compliance
- Wraps entire execution in transaction:
  1. Start transaction
  2. Create transaction record
  3. Update balances
  4. Commit or abort
- Ensures data consistency

## Error Handling

### Allocation Processing
- Validates allocation is active
- Validates target pocket/platform exists and is active
- Validates user has default platform configured
- Validates user has main pocket
- Logs errors without failing entire batch
- Continues processing remaining allocations

### Balance Updates
- Consolidates updates by entity
- Applies in batched operations
- Avoids N+1 query patterns
- Handles concurrent updates safely

## Performance Optimizations

### Aggregation Pipeline
- Uses MongoDB aggregation for efficient data retrieval
- Single query with joins instead of multiple queries
- Filters at database level
- Reduces network overhead

### Batch Operations
- Consolidates balance updates
- Single update per entity instead of multiple
- Reduces database round trips

### Transaction Handling
- Allocations processed individually within transactions
- Prevents partial updates
- Maintains data integrity

## Testing Considerations

### Unit Tests Should Cover
- `ProcessDailyAllocations()`: Correct allocation fetching and processing
- `processAllocationExecution()`: Transaction creation and balance updates
- `applyBalanceUpdates()`: Correct balance consolidation
- Repository methods: Query correctness
- Cron job: Correct scheduling and timezone handling

### Integration Tests Should Cover
- End-to-end allocation execution
- Database transaction rollback on error
- Balance consistency across entities
- Timezone-aware scheduling

## Migration Notes

### Database Schema
- Add `execute_date` field to allocations collection
- Optional field (can be null for non-scheduled allocations)
- Index recommended on `(execute_date, is_active, deleted_at)` for query performance

### Backward Compatibility
- Existing allocations without `execute_date` continue to work
- Only allocations with `execute_date` are processed by cron job
- API accepts null `execute_date` for backward compatibility

## Future Enhancements

### Potential Improvements
1. Allocation execution history tracking
2. Failed allocation retry mechanism
3. Allocation execution status field
4. Bulk allocation scheduling
5. Allocation templates for common patterns
6. Execution notifications/webhooks
7. Allocation analytics and reporting

## Files Modified/Created

### Modified Files
- `internal/modules/allocation/models.go`
- `internal/modules/allocation/dto/request.go`
- `internal/modules/allocation/dto/response.go`
- `internal/modules/allocation/repository.go`
- `internal/modules/allocation/service.go`
- `internal/modules/allocation/controller.go`
- `internal/modules/allocation/module.go`
- `internal/modules/payroll/service.go`
- `internal/modules/payroll/cron.go`
- `internal/modules/payroll/module.go`
- `cmd/api/main.go`

### New Files
- `internal/modules/allocation/cron.go`
- `docs/ALLOCATION_API.md`
- `docs/REFACTORING_SUMMARY.md`

## Verification Checklist

- [x] Allocation schema updated with execute_date
- [x] DTOs updated for execute_date
- [x] Repository methods for execute_date queries
- [x] Service layer allocation processing logic
- [x] Allocation cron job implemented
- [x] Payroll service cleaned up (allocation logic removed)
- [x] Timezone support added to both cron jobs
- [x] Main application integration complete
- [x] API documentation created
- [x] Module dependencies properly configured

## Summary

The refactoring successfully separates allocation execution from payroll processing while introducing scheduled allocation execution capabilities. The allocation module now:

1. **Manages allocations** with optional scheduled execution dates
2. **Processes allocations** independently via dedicated cron job
3. **Maintains data integrity** through database transactions
4. **Handles timezones** consistently across the application
5. **Provides clear API** for allocation management

The payroll module is now focused solely on payroll processing, improving code maintainability and testability. Both modules operate on Asia/Jakarta timezone with staggered execution times to optimize resource usage.
