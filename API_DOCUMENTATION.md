# Coin Backend API Documentation

## Overview
A comprehensive financial management system with smart allocation engine, transaction tracking, and detailed reports.

**Base URL:** `http://localhost:8080/api`  
**API Version:** v1  
**Authentication:** Bearer Token (JWT)

---

## Table of Contents
1. [Authentication Endpoints](#authentication-endpoints)
2. [User Endpoints](#user-endpoints)
3. [Error Responses](#error-responses)
4. [Data Models](#data-models)

---

## Authentication Endpoints

### 1. Register User
Create a new user account with email and password.

**Endpoint:** `POST /v1/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "phone": "+628123456789"
}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | Valid email address |
| password | string | Yes | Minimum 6 characters |
| name | string | Yes | User's full name |
| phone | string | No | Phone number |

**Response (201 Created):**
```json
{
  "success": true,
  "statusCode": 201,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "697ce57ea135e8c451bb2b46",
      "name": "John Doe",
      "email": "user@example.com",
      "phone": "+628123456789",
      "is_active": true,
      "created_at": "2026-01-31T00:51:00Z",
      "updated_at": "2026-01-31T00:51:00Z"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid email format or password too short
- `400 Bad Request` - Email already registered

---

### 2. User Login
Authenticate user with email and password, returns access and refresh tokens.

**Endpoint:** `POST /v1/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | User's email address |
| password | string | Yes | User's password |

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "697ce57ea135e8c451bb2b46",
      "name": "John Doe",
      "email": "user@example.com",
      "phone": "+628123456789",
      "is_active": true,
      "created_at": "2026-01-31T00:51:00Z",
      "updated_at": "2026-01-31T00:51:00Z"
    },
    "token_pair": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid email or password format
- `401 Unauthorized` - Invalid credentials

---

### 3. Refresh Access Token
Generate a new access token using a valid refresh token.

**Endpoint:** `POST /v1/auth/refresh-token`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| refresh_token | string | Yes | Valid refresh token from login |

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid refresh token format
- `401 Unauthorized` - Invalid or expired token

---

### 4. User Logout
Invalidate the user's refresh token and end the session.

**Endpoint:** `POST /v1/auth/logout`

**Authentication:** Required (Bearer Token)

**Request Body:** Empty

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Logged out successfully",
  "data": null
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `400 Bad Request` - Logout failed

---

### 5. Validate Token
Check if the provided token is valid.

**Endpoint:** `GET /v1/auth/validate`

**Authentication:** Required (Bearer Token)

**Request Parameters:** None

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "token is valid",
  "data": {
    "valid": true,
    "user_id": "697ce57ea135e8c451bb2b46"
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or expired token

---

### 6. Get Current User Profile
Retrieve the authenticated user's complete profile information including user and profile data.

**Endpoint:** `GET /v1/auth/me`

**Authentication:** Required (Bearer Token)

**Request Parameters:** None

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "fetch user profile successfully",
  "data": {
    "id": "697ce57ea135e8c451bb2b46",
    "name": "John Doe",
    "email": "user@example.com",
    "phone": "+628123456789",
    "telegramId": "@johndoe",
    "currency": "IDR",
    "baseSalary": 10000000,
    "salaryCycle": "monthly",
    "salaryDay": 25,
    "language": "id",
    "is_active": true,
    "created_at": "2026-01-31T00:51:00Z",
    "updated_at": "2026-01-31T00:51:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `400 Bad Request` - User profile not found

---

## User Endpoints

### 1. Get User Profile
Get the authenticated user's profile information.

**Endpoint:** `GET /v1/users/profile`

**Authentication:** Required (Bearer Token)

**Request Parameters:** None

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "697ce57ea135e8c451bb2b46",
    "name": "John Doe",
    "email": "user@example.com",
    "phone": "+628123456789",
    "telegramId": "@johndoe",
    "currency": "IDR",
    "baseSalary": 10000000,
    "salaryCycle": "monthly",
    "salaryDay": 25,
    "language": "id",
    "is_active": true,
    "created_at": "2026-01-31T00:51:00Z",
    "updated_at": "2026-01-31T00:51:00Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - User profile not found

---

### 2. Update User Profile
Update the authenticated user's profile information.

**Endpoint:** `PUT /v1/users/profile`

**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "name": "John Doe Updated",
  "email": "newemail@example.com",
  "phone": "+628987654321",
  "telegramId": "@johndoe_updated",
  "currency": "USD",
  "baseSalary": 15000000,
  "salaryCycle": "monthly",
  "salaryDay": 1,
  "language": "en"
}
```

**Request Parameters:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | No | User's full name |
| email | string | No | Email address |
| phone | string | No | Phone number |
| telegramId | string | No | Telegram username |
| currency | string | No | Currency code (IDR, USD) |
| baseSalary | number | No | Base salary amount |
| salaryCycle | string | No | Salary cycle (daily, weekly, monthly) |
| salaryDay | integer | No | Salary day (1-28) |
| language | string | No | Language preference (id, en) |

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": "697ce57ea135e8c451bb2b46",
    "name": "John Doe Updated",
    "email": "newemail@example.com",
    "phone": "+628987654321",
    "telegramId": "@johndoe_updated",
    "currency": "USD",
    "baseSalary": 15000000,
    "salaryCycle": "monthly",
    "salaryDay": 1,
    "language": "en",
    "is_active": true,
    "created_at": "2026-01-31T00:51:00Z",
    "updated_at": "2026-01-31T00:52:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid token

---

### 3. List All Users (Admin Only)
Get a paginated list of all users (admin access required).

**Endpoint:** `GET /v1/users`

**Authentication:** Required (Bearer Token - Admin)

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number for pagination |
| limit | integer | 10 | Number of items per page |

**Example Request:**
```
GET /v1/users?page=1&limit=10
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "697ce57ea135e8c451bb2b46",
      "name": "John Doe",
      "email": "user@example.com",
      "phone": "+628123456789",
      "telegramId": "@johndoe",
      "currency": "IDR",
      "baseSalary": 10000000,
      "salaryCycle": "monthly",
      "salaryDay": 25,
      "language": "id",
      "is_active": true,
      "created_at": "2026-01-31T00:51:00Z",
      "updated_at": "2026-01-31T00:51:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Admin access required

---

### 4. Get User by ID (Admin Only)
Get a specific user's information (admin access required).

**Endpoint:** `GET /v1/users/{id}`

**Authentication:** Required (Bearer Token - Admin)

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | User ID (MongoDB ObjectID) |

**Example Request:**
```
GET /v1/users/697ce57ea135e8c451bb2b46
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "697ce57ea135e8c451bb2b46",
    "name": "John Doe",
    "email": "user@example.com",
    "phone": "+628123456789",
    "telegramId": "@johndoe",
    "currency": "IDR",
    "baseSalary": 10000000,
    "salaryCycle": "monthly",
    "salaryDay": 25,
    "language": "id",
    "is_active": true,
    "created_at": "2026-01-31T00:51:00Z",
    "updated_at": "2026-01-31T00:51:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid user ID format
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Admin access required
- `404 Not Found` - User not found

---

### 5. Delete User (Admin Only)
Delete a user account (admin access required).

**Endpoint:** `DELETE /v1/users/{id}`

**Authentication:** Required (Bearer Token - Admin)

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | User ID (MongoDB ObjectID) |

**Example Request:**
```
DELETE /v1/users/697ce57ea135e8c451bb2b46
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User deleted successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request` - Invalid user ID
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Admin access required

---

### 6. Disable User (Admin Only)
Disable a user account (admin access required).

**Endpoint:** `POST /v1/users/{id}/disable`

**Authentication:** Required (Bearer Token - Admin)

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | User ID (MongoDB ObjectID) |

**Request Body:** Empty

**Example Request:**
```
POST /v1/users/697ce57ea135e8c451bb2b46/disable
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User disabled successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request` - Invalid user ID
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Admin access required

---

### 7. Enable User (Admin Only)
Enable a user account (admin access required).

**Endpoint:** `POST /v1/users/{id}/enable`

**Authentication:** Required (Bearer Token - Admin)

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | User ID (MongoDB ObjectID) |

**Request Body:** Empty

**Example Request:**
```
POST /v1/users/697ce57ea135e8c451bb2b46/enable
```

**Response (200 OK):**
```json
{
  "success": true,
  "statusCode": 200,
  "message": "User enabled successfully",
  "data": null
}
```

**Error Responses:**
- `400 Bad Request` - Invalid user ID
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Admin access required

---

## Error Responses

### Standard Error Response Format
All error responses follow this format:

```json
{
  "success": false,
  "statusCode": 400,
  "message": "Error message describing what went wrong",
  "data": null
}
```

### Common HTTP Status Codes
| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request parameters |
| 401 | Unauthorized - Missing or invalid authentication token |
| 403 | Forbidden - Insufficient permissions (admin required) |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Data Models

### User Response Model
```json
{
  "id": "string (MongoDB ObjectID)",
  "name": "string",
  "email": "string (email format)",
  "phone": "string",
  "telegramId": "string",
  "currency": "string (IDR, USD)",
  "baseSalary": "number",
  "salaryCycle": "string (daily, weekly, monthly)",
  "salaryDay": "integer (1-28)",
  "language": "string (id, en)",
  "is_active": "boolean",
  "created_at": "string (ISO 8601 datetime)",
  "updated_at": "string (ISO 8601 datetime)"
}
```

### Token Pair Model
```json
{
  "access_token": "string (JWT token)",
  "refresh_token": "string (JWT token)"
}
```

### Pagination Model
```json
{
  "page": "integer",
  "limit": "integer",
  "total": "integer",
  "totalPages": "integer"
}
```

---

## Authentication

### Bearer Token Format
All authenticated endpoints require the `Authorization` header with a Bearer token:

```
Authorization: Bearer <access_token>
```

### Example cURL Request
```bash
curl -X GET 'http://localhost:8080/api/v1/auth/me' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  -H 'Content-Type: application/json'
```

### Token Expiration
- **Access Token:** Expires in ~15 minutes
- **Refresh Token:** Expires in ~7 days

Use the refresh token endpoint to obtain a new access token before expiration.

---

## Rate Limiting
Currently no rate limiting is implemented. This may be added in future versions.

---

## Swagger UI
Interactive API documentation is available at:
```
http://localhost:8080/api/swagger/index.html
```

---

## Support
For API support, contact: support@swagger.io
