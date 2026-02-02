# Payroll & Allocation System Guide

## Overview

This document describes the payroll auto-input and allocation system implemented in the Coin backend. The system automates salary distribution across user platforms and pockets based on predefined allocation rules.

---

## Core Concepts

### 1. Platform (Global/Admin)

**Purpose**: Reference-only master platform data  
**Key Field**: `is_default` (boolean)

- Platform has **NO balance**
- Platform is **global** (shared across all users)
- `is_default = true` means the platform will be auto-generated for new users

**Example Platforms**:
- BRI (Bank) - is_default: true
- BCA (Bank) - is_default: true
- GoPay (E-Wallet) - is_default: true
- Cash - is_default: true

---

### 2. UserPlatform (User-Owned)

**Purpose**: User-specific platform holding real-time balance

**Key Fields**:
- `user_id`: Owner of this platform instance
- `platform_id`: Reference to global Platform
- `alias_name`: User-friendly name (e.g., "BRI - Salary", "BRI - Saving")
- `balance`: Real money holder (updated via transactions only)
- `is_active`: Status flag

**Important Rules**:
- User can have **MULTIPLE UserPlatforms** referencing the same Platform
- Each UserPlatform has its own balance
- If `alias_name` is null, fallback to `Platform.name` in responses
- All transactions reference `user_platform_id`, NOT `platform_id`

**Example**:
```
User: John
- UserPlatform 1: platform_id=BRI, alias_name="BRI - Salary", balance=5,000,000
- UserPlatform 2: platform_id=BRI, alias_name="BRI - Emergency", balance=2,000,000
- UserPlatform 3: platform_id=GoPay, alias_name="GoPay - Daily", balance=500,000
```

---

### 3. Auto-Generated UserPlatforms

**Trigger**: User registration or admin user creation

**Flow**:
1. System fetches all Platforms where:
   - `is_active = true`
   - `is_default = true`
2. For EACH platform:
   - Create UserPlatform with:
     - `user_id` = new user
     - `platform_id` = platform.id
     - `alias_name` = Platform.name (initial default)
     - `balance` = 0
     - `is_active` = true

**Purpose**: Ensure every new user has basic platform accounts ready

---

### 4. UserProfile.auto_input_payroll

**Purpose**: Opt-in feature for automated payroll processing

**Field**: `auto_input_payroll` (boolean, default: false)

**Behavior**:
- When `true`: System automatically processes payroll on salary day
- When `false`: User must manually input salary transactions

**Update via**: User profile update endpoint

---

### 5. Allocation Module

**Purpose**: Define salary allocation rules

**Fields**:
- `user_id`: Owner
- `pocket_id`: Target pocket (optional)
- `user_platform_id`: Target user platform (optional)
- `priority`: 1 (HIGH), 2 (MEDIUM), 3 (LOW)
- `allocation_type`: PERCENTAGE or NOMINAL
- `nominal`: Percentage (0-100) or fixed amount
- `is_active`: Status flag

**Rules**:
- At least one target (pocket_id OR user_platform_id) must be provided
- Both can be provided for platform+pocket transfers
- Percentage cannot exceed 100
- Allocations are executed in priority order: HIGH → MEDIUM → LOW

**Example Allocations**:
```
Priority 1 (HIGH):
- Pocket: Emergency Fund, Type: PERCENTAGE, Nominal: 10%

Priority 2 (MEDIUM):
- Pocket: Savings, Type: PERCENTAGE, Nominal: 20%
- UserPlatform: BRI - Investment, Type: NOMINAL, Nominal: 1,000,000

Priority 3 (LOW):
- Pocket: Entertainment, Type: PERCENTAGE, Nominal: 5%

Remaining: Goes to Main Wallet (Free Cash)
```

---

## Payroll Auto Input Flow

### Prerequisites

1. User has `auto_input_payroll = true`
2. User has `base_salary > 0`
3. User has `salary_day` configured
4. User has at least one default UserPlatform
5. User has active allocations (optional)

### Flow Diagram

```
Salary Day Trigger
    ↓
Check: user_profile.auto_input_payroll == true?
    ↓ YES
Create INCOME Transaction
    - Type: INCOME
    - Amount: base_salary
    - user_platform_to: DEFAULT UserPlatform
    - Date: Current date
    ↓
Money enters DEFAULT UserPlatform
    ↓
Execute Allocation Rules (Priority Order)
    ↓
Priority 1 (HIGH) Allocations
    - Generate transactions for each allocation
    - Deduct from DEFAULT UserPlatform
    - Transfer to target pocket/platform
    ↓
Priority 2 (MEDIUM) Allocations
    - Same process
    ↓
Priority 3 (LOW) Allocations
    - Same process
    ↓
Remaining Balance
    - Stays in DEFAULT UserPlatform (Free Cash)
    ↓
Complete (All atomic)
```

### Allocation Execution Details

**For Each Allocation** (in priority order):

1. **Calculate Amount**:
   - If PERCENTAGE: `amount = (remaining_balance * nominal / 100)`
   - If NOMINAL: `amount = nominal`

2. **Validate**:
   - Check sufficient balance in DEFAULT UserPlatform
   - Check target pocket/platform is active
   - Check ownership

3. **Create Transaction**:
   - Type: TRANSFER
   - Amount: calculated amount
   - user_platform_from: DEFAULT UserPlatform
   - user_platform_to: target (if platform allocation)
   - pocket_to: target (if pocket allocation)
   - Note: "Auto allocation - [allocation_name]"

4. **Update Balances** (via BalanceProcessor):
   - Decrease DEFAULT UserPlatform balance
   - Increase target balance

5. **Move to Next Allocation**

---

## Transaction Lifecycle

### Single Source of Truth

**All balance changes MUST go through transactions**

- No direct balance mutations allowed
- Every balance change is auditable
- Transactions are atomic

### Transaction Types & Balance Rules

#### INCOME
- **Increases**: `pocket_to` and/or `user_platform_to`
- **Cannot have**: `pocket_from` or `user_platform_from`
- **Example**: Salary, bonus, gift

#### EXPENSE
- **Decreases**: `pocket_from` and/or `user_platform_from`
- **Cannot have**: `pocket_to` or `user_platform_to`
- **Example**: Food, transport, bills

#### TRANSFER
Three scenarios:

1. **Pocket-to-Pocket**: Reallocate between pockets only
   - `pocket_from` → `pocket_to`
   - Platform balance unchanged

2. **Platform-to-Platform**: Move between user platforms only
   - `user_platform_from` → `user_platform_to`
   - Pocket balance unchanged

3. **Platform+Pocket**: Move between platforms and reassign pockets
   - Both pocket and platform balances change
   - `pocket_from` + `user_platform_from` → `pocket_to` + `user_platform_to`

---

## Atomicity & Safety

### Database Transactions

All payroll and allocation operations are wrapped in database transactions:

```
BEGIN TRANSACTION
    1. Create INCOME transaction (payroll)
    2. Update DEFAULT UserPlatform balance
    3. For each allocation:
        a. Create TRANSFER transaction
        b. Update source balance
        c. Update target balance
    4. Validate all balances
COMMIT or ROLLBACK
```

### Rollback Scenarios

System rolls back if:
- Invalid allocation (inactive pocket/platform)
- Insufficient balance
- Ownership mismatch
- Any database error

### Validation Checks

Before processing:
- User owns all pockets and platforms
- All resources are active
- Sufficient balance for all allocations
- No circular references

---

## API Endpoints

### Platform Management

```
POST   /v1/platforms              # Create platform (admin)
GET    /v1/platforms              # List platforms
GET    /v1/platforms/:id          # Get platform
PUT    /v1/platforms/:id          # Update platform (admin)
DELETE /v1/platforms/:id          # Delete platform (admin)
```

### UserPlatform Management

```
POST   /v1/user-platforms         # Create user platform
GET    /v1/user-platforms         # List user platforms
GET    /v1/user-platforms/dropdown/list  # Dropdown with platform data
GET    /v1/user-platforms/:id     # Get user platform
PUT    /v1/user-platforms/:id     # Update user platform (alias_name, is_active)
DELETE /v1/user-platforms/:id     # Delete user platform
```

### Allocation Management

```
POST   /v1/allocations            # Create allocation
GET    /v1/allocations            # List allocations
GET    /v1/allocations/:id        # Get allocation
PUT    /v1/allocations/:id        # Update allocation
DELETE /v1/allocations/:id        # Delete allocation
```

### User Profile

```
PUT    /v1/users/:id              # Update user (includes auto_input_payroll)
```

---

## Example Scenarios

### Scenario 1: New User Registration

```
1. User registers with email/password
2. System creates User record
3. System creates UserProfile (auto_input_payroll = false by default)
4. System creates default Pockets from templates
5. System creates default UserCategories from templates
6. System auto-generates UserPlatforms:
   - Fetches Platforms where is_active=true, is_default=true
   - Creates UserPlatform for each (balance=0, alias_name=Platform.name)
7. User is ready to use the system
```

### Scenario 2: User Enables Payroll Auto Input

```
1. User updates profile: auto_input_payroll = true
2. User sets base_salary = 10,000,000
3. User sets salary_day = 25
4. User creates allocations:
   - Priority 1: Emergency Fund (10%, PERCENTAGE)
   - Priority 2: Savings (20%, PERCENTAGE)
   - Priority 3: Investment (1,000,000, NOMINAL)
5. On day 25, system automatically:
   - Creates INCOME transaction (10,000,000 → DEFAULT UserPlatform)
   - Executes allocations:
     * Emergency: 1,000,000 (10%)
     * Savings: 2,000,000 (20%)
     * Investment: 1,000,000
   - Remaining: 6,000,000 stays in DEFAULT UserPlatform
```

### Scenario 3: User Creates Multiple Platforms

```
1. User has default "BRI" UserPlatform (auto-generated)
2. User creates additional UserPlatforms:
   - POST /v1/user-platforms
   - Body: { platform_id: "BRI_ID", alias_name: "BRI - Emergency" }
3. User now has:
   - UserPlatform 1: BRI (default, alias_name="BRI")
   - UserPlatform 2: BRI (alias_name="BRI - Emergency")
4. User can allocate salary to different BRI accounts
5. Each has independent balance tracking
```

---

## Migration Guide

### For Existing Systems

If migrating from old system:

1. **Add is_default to existing Platforms**:
   ```sql
   UPDATE platforms SET is_default = true WHERE name IN ('Cash', 'BRI', 'BCA');
   ```

2. **Create user_platforms collection**:
   - For each user + platform combination
   - Migrate platform balances to UserPlatform.Balance
   - Set alias_name = Platform.name initially

3. **Update existing transactions**:
   - Add user_platform_from and user_platform_to fields
   - Map from old platform_id references

4. **Add auto_input_payroll to user_profiles**:
   ```sql
   UPDATE user_profiles SET auto_input_payroll = false;
   ```

5. **Validate balances**:
   - Ensure UserPlatform balances match sum of transactions

---

## Testing Checklist

### Unit Tests

- [ ] Platform CRUD with is_default field
- [ ] UserPlatform CRUD with alias_name
- [ ] Multiple UserPlatforms per user per platform
- [ ] Allocation CRUD with validation
- [ ] Auto-generate UserPlatforms on user creation
- [ ] Allocation execution logic (priority order)
- [ ] Balance calculations (percentage vs nominal)

### Integration Tests

- [ ] End-to-end user registration with auto-generated UserPlatforms
- [ ] Payroll auto input with allocations
- [ ] Transaction creation with UserPlatform references
- [ ] Balance consistency across multiple transactions
- [ ] Rollback on allocation failure
- [ ] Concurrent transaction handling

### Manual Tests

- [ ] Create user and verify auto-generated UserPlatforms
- [ ] Enable auto_input_payroll and verify behavior
- [ ] Create allocations with different priorities
- [ ] Simulate salary day and verify distribution
- [ ] Create multiple UserPlatforms for same Platform
- [ ] Update alias_name and verify response

---

## Troubleshooting

### Issue: UserPlatforms not auto-generated

**Check**:
1. Are there Platforms with is_default=true?
2. Are those Platforms active (is_active=true)?
3. Check auth service logs for errors

### Issue: Payroll not auto-processing

**Check**:
1. Is auto_input_payroll = true?
2. Is base_salary > 0?
3. Is salary_day configured?
4. Does user have a default UserPlatform?
5. Check cron job is running

### Issue: Allocation fails

**Check**:
1. Is target pocket/platform active?
2. Is there sufficient balance in source?
3. Does user own all resources?
4. Check allocation priority order
5. Review transaction logs

---

## Best Practices

1. **Always use UserPlatform references in transactions**, never direct Platform references
2. **Set meaningful alias_name** for multiple platforms of same type
3. **Order allocations by priority** (HIGH for essentials, LOW for discretionary)
4. **Test allocations** before enabling auto_input_payroll
5. **Monitor balance consistency** via transaction audit logs
6. **Use percentage allocations** for flexibility with varying salaries
7. **Keep at least one default UserPlatform** for free cash

---

## Future Enhancements

Potential improvements:

- [ ] Scheduled allocation execution (not just payroll day)
- [ ] Allocation templates for quick setup
- [ ] Allocation history and analytics
- [ ] Multi-currency support for UserPlatforms
- [ ] Allocation rules based on conditions (e.g., if balance > X)
- [ ] Notification system for payroll processing
- [ ] Rollback mechanism for manual correction

---

## Summary

The payroll and allocation system provides:

✓ **Automated salary distribution** based on user-defined rules  
✓ **Multiple platform instances** per user for better organization  
✓ **User-friendly naming** with alias_name support  
✓ **Priority-based allocation** for flexible money management  
✓ **Atomic transactions** ensuring data consistency  
✓ **Auto-generated platforms** for new users  
✓ **Opt-in automation** via auto_input_payroll flag  

All balance changes flow through the centralized transaction system, ensuring auditability and consistency.
