# Coin Finance System - ERD & Module Documentation

## Table of Contents
1. [Entity Relationship Diagram](#entity-relationship-diagram)
2. [Database Collections](#database-collections)
3. [Module Architecture](#module-architecture)
4. [Data Flow](#data-flow)
5. [API Endpoints by Module](#api-endpoints-by-module)

---

## Entity Relationship Diagram

### Collections Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         COIN FINANCE SYSTEM                         │
└─────────────────────────────────────────────────────────────────────┘

                              ┌──────────────┐
                              │    USERS     │
                              │              │
                              │ • id (PK)    │
                              │ • email      │
                              │ • name       │
                              │ • role       │
                              │ • is_active  │
                              └──────┬───────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │                │                │
                    ▼                ▼                ▼
          ┌──────────────────┐  ┌──────────────┐  ┌─────────────────┐
          │  USER_PROFILES   │  │   POCKETS    │  │ USER_CATEGORIES │
          │                  │  │              │  │                 │
          │ • id (PK)        │  │ • id (PK)    │  │ • id (PK)       │
          │ • user_id (FK)   │  │ • user_id(FK)│  │ • user_id (FK)  │
          │ • phone          │  │ • name       │  │ • name          │
          │ • telegram_id    │  │ • type       │  │ • type          │
          │ • base_salary    │  │ • balance    │  │ • parent_id     │
          │ • salary_cycle   │  │ • is_active  │  │ • is_deleted    │
          │ • pay_currency   │  │ • is_locked  │  └─────────────────┘
          │ • lang           │  └──────┬───────┘
          └──────────────────┘         │
                                       │
                    ┌──────────────────┴──────────────────┐
                    │                                     │
                    ▼                                     ▼
          ┌──────────────────┐                ┌──────────────────┐
          │  TRANSACTIONS    │                │   PLATFORMS      │
          │                  │                │                  │
          │ • id (PK)        │                │ • id (PK)        │
          │ • user_id (FK)   │                │ • name           │
          │ • type           │                │ • type           │
          │ • amount         │                │ • is_active      │
          │ • pocket_from    │                └──────────────────┘
          │ • pocket_to      │
          │ • category_id    │
          │ • platform_id    │
          │ • date           │
          │ • note           │
          │ • ref            │
          └──────────────────┘

          ┌──────────────────┐        ┌──────────────────┐
          │CATEGORY_TEMPLATES│        │ POCKET_TEMPLATES │
          │                  │        │                  │
          │ • id (PK)        │        │ • id (PK)        │
          │ • name           │        │ • name           │
          │ • type           │        │ • type           │
          │ • parent_id      │        │ • category_id    │
          │ • user_id        │        │ • is_default     │
          │ • is_default     │        │ • order          │
          └──────────────────┘        └──────────────────┘
```

---

## Database Collections

### 1. **users** (Core User Entity)
**Purpose**: Store user account information

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `email` | String | User email | Unique, Required |
| `password_hash` | String | Hashed password | Required |
| `salt` | String | Password salt | Required |
| `name` | String | User full name | Required |
| `role` | String | User role (admin/user) | Enum: admin, user |
| `is_active` | Boolean | Account status | Default: true |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |

**Indexes**:
- `{email: 1}` - Unique index for login
- `{is_active: 1, created_at: -1}` - For active users list

---

### 2. **user_profiles** (Extended User Information)
**Purpose**: Store user profile and preferences

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `user_id` | ObjectID | Reference to users | FK, Required |
| `phone` | String | Phone number | Optional |
| `telegram_id` | String | Telegram user ID | Optional |
| `base_salary` | Float | Base salary amount | Default: 0 |
| `salary_cycle` | String | Salary frequency | Enum: daily, weekly, monthly |
| `salary_day` | Integer | Day of salary payment | 1-31 |
| `pay_currency` | String | Salary currency | Enum: IDR, USD |
| `lang` | String | Preferred language | Enum: id, en |
| `is_active` | Boolean | Profile status | Default: true |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |

**Indexes**:
- `{user_id: 1}` - Unique index for profile lookup

---

### 3. **pockets** (Money Storage Accounts)
**Purpose**: Represent different money storage accounts/wallets

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `user_id` | ObjectID | Reference to users | FK, Required |
| `name` | String | Pocket name | Required |
| `type` | String | Pocket type | Enum: main, allocation, saving, debt, system |
| `category_id` | ObjectID | Reference to category | FK, Optional |
| `balance` | Decimal128 | Current balance | Precision: 2 decimals |
| `is_default` | Boolean | Default pocket flag | Default: false |
| `is_active` | Boolean | Active status | Default: true |
| `is_locked` | Boolean | Locked status | Default: false |
| `icon` | String | Icon identifier | Optional |
| `icon_color` | String | Icon color | Optional |
| `background_color` | String | Background color | Optional |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |
| `deleted_at` | DateTime | Soft delete timestamp | Optional |

**Indexes**:
- `{user_id: 1, is_active: 1}` - For user's active pockets
- `{user_id: 1, type: 1}` - For pockets by type
- `{user_id: 1, deleted_at: 1}` - For soft-deleted pockets

**Pocket Types**:
- `main` - Primary wallet for daily transactions
- `allocation` - Budget allocation pockets
- `saving` - Savings accounts
- `debt` - Debt tracking
- `system` - System-managed pockets

---

### 4. **transactions** (Financial Records)
**Purpose**: Record all financial movements

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `user_id` | ObjectID | Reference to users | FK, Required |
| `type` | String | Transaction type | Enum: income, expense, transfer, dp, withdraw |
| `amount` | Float | Transaction amount | Required, > 0 |
| `pocket_from` | ObjectID | Source pocket | FK, Optional |
| `pocket_to` | ObjectID | Destination pocket | FK, Optional |
| `category_id` | ObjectID | Reference to category | FK, Optional |
| `platform_id` | ObjectID | Reference to platform | FK, Optional |
| `note` | String | Transaction note | Optional |
| `date` | DateTime | Transaction date | Required |
| `ref` | String | Reference number | Optional |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |
| `deleted_at` | DateTime | Soft delete timestamp | Optional |

**Indexes**:
- `{user_id: 1, date: -1}` - For user's transactions timeline
- `{user_id: 1, type: 1, date: -1}` - For transactions by type
- `{user_id: 1, category_id: 1, date: -1}` - For category analysis
- `{user_id: 1, deleted_at: 1}` - For soft-deleted transactions

**Transaction Types**:
- `income` - Money received
- `expense` - Money spent
- `transfer` - Money moved between pockets
- `dp` - Debt payment
- `withdraw` - Cash withdrawal

---

### 5. **user_categories** (User-Defined Categories)
**Purpose**: Store user's custom transaction categories

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `user_id` | ObjectID | Reference to users | FK, Required |
| `template_id` | ObjectID | Reference to template | FK, Optional |
| `name` | String | Category name | Required |
| `transaction_type` | String | Category type | Enum: income, expense |
| `parent_id` | ObjectID | Parent category | FK, Optional (for hierarchy) |
| `description` | String | Category description | Optional |
| `icon` | String | Icon identifier | Optional |
| `color` | String | Category color | Optional |
| `is_default` | Boolean | Default category flag | Default: false |
| `is_deleted` | Boolean | Deleted flag | Default: false |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |
| `deleted_at` | DateTime | Soft delete timestamp | Optional |

**Indexes**:
- `{user_id: 1, transaction_type: 1}` - For categories by type
- `{user_id: 1, parent_id: 1}` - For category hierarchy
- `{user_id: 1, is_deleted: 1}` - For active categories

---

### 6. **platforms** (Payment Methods)
**Purpose**: Store payment platform/method information

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `name` | String | Platform name | Required |
| `type` | String | Platform type | Enum: BANK, E_WALLET, CASH, ATM |
| `is_active` | Boolean | Active status | Default: true |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |
| `deleted_at` | DateTime | Soft delete timestamp | Optional |

**Indexes**:
- `{type: 1, is_active: 1}` - For platforms by type

**Platform Types**:
- `BANK` - Bank accounts
- `E_WALLET` - Digital wallets (GCash, OVO, etc.)
- `CASH` - Physical cash
- `ATM` - ATM withdrawals

---

### 7. **category_templates** (Category Templates)
**Purpose**: Provide default category templates for users

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `name` | String | Template name | Required |
| `transaction_type` | String | Category type | Enum: income, expense |
| `parent_id` | ObjectID | Parent template | FK, Optional |
| `user_id` | ObjectID | Creator user | FK, Optional (null = system) |
| `description` | String | Template description | Optional |
| `icon` | String | Icon identifier | Optional |
| `color` | String | Category color | Optional |
| `is_default` | Boolean | Default template flag | Default: false |
| `is_deleted` | Boolean | Deleted flag | Default: false |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |
| `deleted_at` | DateTime | Soft delete timestamp | Optional |

**Indexes**:
- `{transaction_type: 1, is_default: 1}` - For default templates
- `{user_id: 1}` - For user-created templates

---

### 8. **pocket_templates** (Pocket Templates)
**Purpose**: Provide default pocket templates for users

| Field | Type | Description | Constraints |
|-------|------|-------------|-------------|
| `_id` | ObjectID | Primary key | Unique |
| `name` | String | Template name | Required |
| `type` | String | Pocket type | Enum: main, saving, allocation |
| `category_id` | ObjectID | Reference to category | FK, Optional |
| `is_default` | Boolean | Default template flag | Default: false |
| `is_active` | Boolean | Active status | Default: true |
| `order` | Integer | Display order | Default: 0 |
| `icon` | String | Icon identifier | Optional |
| `icon_color` | String | Icon color | Optional |
| `background_color` | String | Background color | Optional |
| `created_at` | DateTime | Creation timestamp | Auto-set |
| `updated_at` | DateTime | Last update timestamp | Auto-set |
| `deleted_at` | DateTime | Soft delete timestamp | Optional |

**Indexes**:
- `{type: 1, is_default: 1}` - For default templates by type
- `{is_active: 1, order: 1}` - For display ordering

---

## Module Architecture

### Module Structure

Each module follows a consistent architecture:

```
module/
├── models.go          # Data structures and constants
├── controller.go      # HTTP request handlers
├── service.go         # Business logic
├── repository.go      # Database operations
├── routes.go          # Route definitions
└── register.go        # DI container registration
```

### Module Descriptions

#### 1. **Auth Module**
**Purpose**: User authentication and authorization

**Key Features**:
- User login/logout
- JWT token generation
- Password hashing with salt
- Role-based access control

**Dependencies**: User module

---

#### 2. **User Module**
**Purpose**: User account management

**Key Features**:
- User registration
- Profile management
- Account settings
- User activation/deactivation

**Collections**: `users`, `user_profiles`

**Dependencies**: None (core module)

---

#### 3. **Pocket Module**
**Purpose**: Manage money storage accounts

**Key Features**:
- Create/update/delete pockets
- Balance tracking
- Pocket locking
- Pocket type management

**Collections**: `pockets`

**Dependencies**: User module, Category module (optional)

**Pocket Types**:
- Main wallet
- Savings accounts
- Budget allocations
- Debt tracking
- System accounts

---

#### 4. **Transaction Module**
**Purpose**: Record and manage financial transactions

**Key Features**:
- Create transactions (income, expense, transfer)
- Transaction history
- Soft delete transactions
- Transaction filtering and search

**Collections**: `transactions`

**Dependencies**: User module, Pocket module, Category module, Platform module

**Transaction Types**:
- Income: Money received
- Expense: Money spent
- Transfer: Between pockets
- Debt Payment (dp): Debt repayment
- Withdraw: Cash withdrawal

---

#### 5. **User Category Module**
**Purpose**: Manage user-defined transaction categories

**Key Features**:
- Create custom categories
- Category hierarchy (parent-child)
- Category templates
- Category filtering by type

**Collections**: `user_categories`

**Dependencies**: User module, Category Template module (optional)

**Category Types**:
- Income categories
- Expense categories
- Hierarchical organization

---

#### 6. **Platform Module**
**Purpose**: Manage payment platforms/methods

**Key Features**:
- Platform management
- Platform type classification
- Platform activation/deactivation

**Collections**: `platforms`

**Dependencies**: None

**Platform Types**:
- Bank accounts
- E-wallets
- Cash
- ATM

---

#### 7. **Category Template Module**
**Purpose**: Provide default category templates

**Key Features**:
- System default templates
- User-created templates
- Template hierarchy
- Template usage tracking

**Collections**: `category_templates`

**Dependencies**: User module (optional)

---

#### 8. **Pocket Template Module**
**Purpose**: Provide default pocket templates

**Key Features**:
- System default templates
- Template ordering
- Template activation
- Quick pocket creation

**Collections**: `pocket_templates`

**Dependencies**: Category module (optional)

---

## Data Flow

### User Registration Flow
```
1. User submits registration
   ↓
2. Auth module validates input
   ↓
3. User module creates user account
   ↓
4. User module creates user profile
   ↓
5. Pocket module creates default pockets (from templates)
   ↓
6. User category module creates default categories (from templates)
   ↓
7. User account ready for use
```

### Transaction Creation Flow
```
1. User submits transaction
   ↓
2. Transaction module validates input
   ↓
3. Pocket module updates source pocket balance
   ↓
4. Pocket module updates destination pocket balance (if transfer)
   ↓
5. Transaction module records transaction
   ↓
6. Transaction complete
```

### Category Management Flow
```
1. User views/creates category
   ↓
2. User category module checks templates
   ↓
3. User can create custom category
   ↓
4. Category available for transactions
   ↓
5. Category can be organized hierarchically
```

---

## API Endpoints by Module

### Auth Module
```
POST   /api/v1/auth/register          - User registration
POST   /api/v1/auth/login             - User login
POST   /api/v1/auth/logout            - User logout
POST   /api/v1/auth/refresh-token     - Refresh JWT token
```

### User Module
```
GET    /api/v1/users/profile          - Get user profile
PUT    /api/v1/users/profile          - Update user profile
GET    /api/v1/users/settings         - Get user settings
PUT    /api/v1/users/settings         - Update user settings
```

### Pocket Module
```
GET    /api/v1/pockets                - List all pockets
POST   /api/v1/pockets                - Create pocket
GET    /api/v1/pockets/:id            - Get pocket details
PUT    /api/v1/pockets/:id            - Update pocket
DELETE /api/v1/pockets/:id            - Delete pocket
PUT    /api/v1/pockets/:id/lock       - Lock pocket
PUT    /api/v1/pockets/:id/unlock     - Unlock pocket
```

### Transaction Module
```
GET    /api/v1/transactions           - List transactions
POST   /api/v1/transactions           - Create transaction
GET    /api/v1/transactions/:id       - Get transaction details
PUT    /api/v1/transactions/:id       - Update transaction
DELETE /api/v1/transactions/:id       - Delete transaction
GET    /api/v1/transactions/filter    - Filter transactions
```

### User Category Module
```
GET    /api/v1/user-categories        - List categories
POST   /api/v1/user-categories        - Create category
GET    /api/v1/user-categories/:id    - Get category details
PUT    /api/v1/user-categories/:id    - Update category
DELETE /api/v1/user-categories/:id    - Delete category
```

### Platform Module
```
GET    /api/v1/platforms              - List platforms
POST   /api/v1/platforms              - Create platform
GET    /api/v1/platforms/:id          - Get platform details
PUT    /api/v1/platforms/:id          - Update platform
DELETE /api/v1/platforms/:id          - Delete platform
```

### Category Template Module
```
GET    /api/v1/category-templates     - List templates
POST   /api/v1/category-templates     - Create template
GET    /api/v1/category-templates/:id - Get template details
PUT    /api/v1/category-templates/:id - Update template
DELETE /api/v1/category-templates/:id - Delete template
```

### Pocket Template Module
```
GET    /api/v1/pocket-templates       - List templates
POST   /api/v1/pocket-templates       - Create template
GET    /api/v1/pocket-templates/:id   - Get template details
PUT    /api/v1/pocket-templates/:id   - Update template
DELETE /api/v1/pocket-templates/:id   - Delete template
```

---

## Key Relationships

### User-Centric Design
All core entities are tied to a user:
- Each pocket belongs to one user
- Each transaction belongs to one user
- Each category belongs to one user
- User profile extends user account

### Transaction Flow
```
User
  ├── Pockets (multiple)
  │   └── Transactions (multiple)
  │       ├── Category (one)
  │       ├── Platform (one)
  │       ├── Source Pocket (one)
  │       └── Destination Pocket (one)
  └── Categories (multiple)
      └── Transactions (multiple)
```

### Template System
```
Templates (System-wide)
  ├── Category Templates
  │   └── Used by User Categories
  └── Pocket Templates
      └── Used by Pockets
```

---

## Data Integrity Rules

### Soft Deletes
- Transactions, Pockets, Users, and Categories support soft deletes
- Use `deleted_at` field for soft deletion
- Queries should filter out soft-deleted records

### Balance Integrity
- Pocket balance is updated on transaction creation
- Balance must always be non-negative (configurable)
- Transfer transactions update both source and destination

### Category Hierarchy
- Categories can have parent categories
- Supports unlimited nesting depth
- Parent category deletion cascades to children (configurable)

### User Isolation
- All queries must filter by `user_id`
- No cross-user data access
- User deletion cascades to all related data

---

## Performance Considerations

### Indexing Strategy
- All `user_id` fields are indexed
- Date ranges indexed for transaction queries
- Type fields indexed for filtering
- Composite indexes for common query patterns

### Query Optimization
- Use indexes for user_id + date range queries
- Aggregate transactions at read time for reports
- Cache category and platform lists
- Denormalize frequently accessed data

### Scalability
- Sharding by `user_id` for horizontal scaling
- Archive old transactions to separate collections
- Use read replicas for reporting queries
- Implement caching layer for dashboard data

---

## Summary

The Coin Finance System uses a modular architecture with 8 core modules managing:
- **User Management**: Authentication and profiles
- **Financial Tracking**: Transactions and pockets
- **Organization**: Categories and platforms
- **Templates**: Default configurations for quick setup

All modules follow consistent patterns for controllers, services, repositories, and routes, making the system maintainable and scalable.

