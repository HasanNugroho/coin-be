# Role System Documentation

## Overview

The application now uses a simplified role system with only **2 roles**: `admin` and `user`.

## Role Constants

```go
const (
    RoleAdmin = "admin"
    RoleUser  = "user"
)
```

## Role Assignment

### Registration Flow
- **First user registered** → Automatically assigned `admin` role
- **Subsequent users** → Automatically assigned `user` role
- **Cannot register admin if admin exists** → System prevents multiple admins

### User Model
The `User` model now includes a `Role` field:

```go
type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Email        string             `bson:"email" json:"email"`
    Phone        string             `bson:"phone" json:"phone"`
    PasswordHash string             `bson:"password_hash" json:"-"`
    Name         string             `bson:"name" json:"name"`
    Role         string             `bson:"role" json:"role"` // "admin" or "user"
    IsActive     bool               `bson:"is_active" json:"is_active"`
    CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
```

## Admin Capabilities

### 1. User Management
- **List all users**: `GET /api/users`
- **Get specific user**: `GET /api/users/:id`
- **Delete user**: `DELETE /api/users/:id`
- **Disable user**: `POST /api/users/:id/disable`
- **Enable user**: `POST /api/users/:id/enable`

### 2. Dashboard Access
- **Admin Dashboard**: `GET /api/reports/dashboard`
  - View summary of all system activities
  - See total balance across all allocations
  - Monitor income and expenses
  - Track allocation distribution

### 3. User Control
- Disable/enable user accounts
- Delete user accounts
- View all user information

## Middleware

### AuthMiddleware
- Validates JWT token
- Fetches user role from database
- Sets `user_id`, `email`, and `role` in context

```go
middleware.AuthMiddleware(jwtManager, db)
```

### AdminMiddleware
- Checks if user has `admin` role
- Returns 403 Forbidden if user is not admin
- Must be used after AuthMiddleware

```go
middleware.AdminMiddleware()
```

## API Endpoints

### Protected Routes (All authenticated users)
- `/api/categories` - Category management
- `/api/transactions` - Transaction management
- `/api/allocations` - Allocation management
- `/api/targets` - Saving target management
- `/api/reports` - Financial reports

### Admin-Only Routes (Requires admin role)
- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user details
- `DELETE /api/users/:id` - Delete user
- `POST /api/users/:id/disable` - Disable user
- `POST /api/users/:id/enable` - Enable user

### Public Routes (No authentication required)
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh token
- `POST /api/auth/logout` - User logout
- `GET /api/health` - Health check

## Database Changes

### Users Collection
The `users` collection now includes a `role` field:

```json
{
  "_id": ObjectId,
  "email": "user@example.com",
  "phone": "081234567890",
  "password_hash": "hashed_password",
  "name": "John Doe",
  "role": "admin",
  "is_active": true,
  "created_at": ISODate,
  "updated_at": ISODate
}
```

### Removed Collections
The following collections are no longer used:
- `roles` - Removed (roles are now stored in User model)
- `user_roles` - Removed (roles are now stored in User model)

## Service Methods

### User Service
```go
// Enable/Disable user (admin only)
DisableUser(ctx context.Context, userID string) error
EnableUser(ctx context.Context, userID string) error
```

### User Repository
```go
// Count users by role
CountUsersByRole(ctx context.Context, role string) (int64, error)
```

## Authentication Flow

### Registration
1. User submits registration request
2. System checks if admin exists
3. If no admin exists → assign `admin` role
4. If admin exists → assign `user` role
5. Create default categories and allocations
6. Return success response

### Login
1. User submits email and password
2. System validates credentials
3. System fetches user role from database
4. Generate JWT token with user info
5. Return token and user details

### Protected Request
1. Client sends request with Bearer token
2. AuthMiddleware validates token
3. AuthMiddleware fetches user role from database
4. Sets `user_id`, `email`, and `role` in context
5. Request proceeds to handler

### Admin Request
1. Client sends request with Bearer token
2. AuthMiddleware validates and sets role
3. AdminMiddleware checks if role is `admin`
4. If not admin → return 403 Forbidden
5. If admin → request proceeds to handler

## Migration Notes

### For Existing Databases
If migrating from the old role system:

1. Add `role` field to all users in the `users` collection
2. Set first user's role to `admin`
3. Set all other users' roles to `user`
4. Drop `roles` collection
5. Drop `user_roles` collection

### Migration Script (MongoDB)
```javascript
// Add role field to all users
db.users.updateMany(
  { role: { $exists: false } },
  { $set: { role: "user" } }
);

// Set first user as admin
db.users.updateOne(
  {},
  { $set: { role: "admin" } },
  { sort: { created_at: 1 } }
);

// Drop old collections
db.roles.drop();
db.user_roles.drop();
```

## Security Considerations

1. **Role Immutability**: User roles cannot be changed via API (only during registration)
2. **Admin Protection**: Only first user can be admin
3. **User Isolation**: Users can only access their own financial data
4. **Admin Scope**: Admins can manage users but cannot access user financial data
5. **Token Validation**: Role is fetched from database on each request (not cached in token)

## Testing

### Test Admin Registration
```bash
# First registration - should be admin
POST /api/auth/register
{
  "email": "admin@example.com",
  "phone": "081234567890",
  "password": "password123",
  "name": "Admin User"
}

# Second registration - should be user
POST /api/auth/register
{
  "email": "user@example.com",
  "phone": "081234567891",
  "password": "password123",
  "name": "Regular User"
}
```

### Test Admin Endpoints
```bash
# Login as admin
POST /api/auth/login
{
  "email": "admin@example.com",
  "password": "password123"
}

# List all users (admin only)
GET /api/users
Authorization: Bearer <admin_token>

# Disable user (admin only)
POST /api/users/<user_id>/disable
Authorization: Bearer <admin_token>
```

### Test User Isolation
```bash
# Login as regular user
POST /api/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Try to access admin endpoints (should fail)
GET /api/users
Authorization: Bearer <user_token>
# Response: 403 Forbidden - admin access required
```
