# Coin Backend API Documentation

## Overview

This document describes the complete API for the Coin financial management backend system. The system manages multiple financial platforms and pockets (money allocations) per user, with real-time balance tracking through atomic transactions.

### Key Concepts

- **Platform (Admin)**: Global master platform definitions (BANK, E_WALLET, CASH, ATM). These are NOT used directly in transactions but serve as reference data for categorizing money sources.
- **Pocket (Kantong)**: User-owned money allocations with real-time balance. Each pocket tracks money for a specific purpose (main, allocation, saving, debt, or system).
- **Transaction**: Single source of truth for all balance changes. Every balance update is recorded as a transaction (income, expense, or transfer).

---

## 1. Platform APIs (Admin Only)

Platforms are global master data that define types of money sources. They are **NOT used directly in transactions** but serve as reference information.

### 1.1 Create Platform

**HTTP Method**: `POST`  
**URL**: `/v1/platforms/admin`  
**Authentication**: Required (Admin only)

**Request Body**:
```json
{
  "name": "BCA Bank",
  "type": "BANK",
  "is_active": true
}
```

**Field Descriptions**:
- `name` (string, required): Platform name (1-255 characters)
- `type` (string, required): One of `BANK`, `E_WALLET`, `CASH`, `ATM`
- `is_active` (boolean, required): Whether platform is active

**Example Response** (201 Created):
```json
{
  "success": true,
  "message": "Platform created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "BCA Bank",
    "type": "BANK",
    "is_active": true,
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:30:00Z",
    "deleted_at": null
  }
}
```

---

### 1.2 Get Platform by ID

**HTTP Method**: `GET`  
**URL**: `/v1/platforms/{id}`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Platform ID (24-character hex string)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Platform retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "BCA Bank",
    "type": "BANK",
    "is_active": true,
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:30:00Z",
    "deleted_at": null
  }
}
```

---

### 1.3 List All Platforms

**HTTP Method**: `GET`  
**URL**: `/v1/platforms`  
**Authentication**: Required

**Query Parameters**:
- `limit` (integer, optional): Number of results per page (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Platforms retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "BCA Bank",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-02-02T10:30:00Z",
      "updated_at": "2024-02-02T10:30:00Z",
      "deleted_at": null
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "name": "GCash",
      "type": "E_WALLET",
      "is_active": true,
      "created_at": "2024-02-02T10:31:00Z",
      "updated_at": "2024-02-02T10:31:00Z",
      "deleted_at": null
    }
  ]
}
```

---

### 1.4 List Active Platforms

**HTTP Method**: `GET`  
**URL**: `/v1/platforms/active`  
**Authentication**: Required

**Query Parameters**:
- `limit` (integer, optional): Number of results per page (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Active platforms retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "BCA Bank",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-02-02T10:30:00Z",
      "updated_at": "2024-02-02T10:30:00Z",
      "deleted_at": null
    }
  ]
}
```

---

### 1.5 List Platforms by Type

**HTTP Method**: `GET`  
**URL**: `/v1/platforms/type/{type}`  
**Authentication**: Required

**Path Parameters**:
- `type` (string): Platform type filter (`BANK`, `E_WALLET`, `CASH`, or `ATM`)

**Query Parameters**:
- `limit` (integer, optional): Number of results per page (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Platforms retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "BCA Bank",
      "type": "BANK",
      "is_active": true,
      "created_at": "2024-02-02T10:30:00Z",
      "updated_at": "2024-02-02T10:30:00Z",
      "deleted_at": null
    }
  ]
}
```

---

### 1.6 Update Platform

**HTTP Method**: `PUT`  
**URL**: `/v1/platforms/{id}`  
**Authentication**: Required (Admin only)

**Path Parameters**:
- `id` (string): Platform ID

**Request Body** (all fields optional):
```json
{
  "name": "BCA Bank Updated",
  "type": "BANK",
  "is_active": false
}
```

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Platform updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "BCA Bank Updated",
    "type": "BANK",
    "is_active": false,
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:35:00Z",
    "deleted_at": null
  }
}
```

---

### 1.7 Delete Platform

**HTTP Method**: `DELETE`  
**URL**: `/v1/platforms/{id}`  
**Authentication**: Required (Admin only)

**Path Parameters**:
- `id` (string): Platform ID

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Platform deleted successfully",
  "data": null
}
```

---

## 2. Pocket APIs (User Pockets)

Pockets are user-owned money allocations with real-time balance. Balance is updated atomically through transactions.

### 2.1 Create Pocket

**HTTP Method**: `POST`  
**URL**: `/v1/pockets`  
**Authentication**: Required

**Request Body**:
```json
{
  "name": "Monthly Savings",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439013",
  "target_balance": 1000000,
  "icon": "piggy-bank",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5"
}
```

**Field Descriptions**:
- `name` (string, required): Pocket name (2-255 characters)
- `type` (string, required): One of `main`, `allocation`, `saving`, `debt`
- `category_id` (string, optional): Category ID (24-character hex)
- `target_balance` (number, optional): Target balance goal (must be > 0)
- `icon` (string, optional): Icon identifier (max 100 characters)
- `icon_color` (string, optional): Icon color hex (max 50 characters)
- `background_color` (string, optional): Background color hex (max 50 characters)

**Pocket Types**:
- `main`: Default pocket for the user (one per user, cannot be deleted)
- `allocation`: Money allocated for specific purposes
- `saving`: Money set aside for savings goals
- `debt`: Money owed or debt tracking

**Example Response** (201 Created):
```json
{
  "success": true,
  "message": "Pocket created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "user_id": "507f1f77bcf86cd799439001",
    "name": "Monthly Savings",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439013",
    "balance": 0,
    "target_balance": 1000000,
    "is_default": false,
    "is_active": true,
    "is_locked": false,
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:30:00Z",
    "deleted_at": null
  }
}
```

---

### 2.2 Get Pocket by ID

**HTTP Method**: `GET`  
**URL**: `/v1/pockets/{id}`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Pocket ID

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pocket retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "user_id": "507f1f77bcf86cd799439001",
    "name": "Monthly Savings",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439013",
    "balance": 500000,
    "target_balance": 1000000,
    "is_default": false,
    "is_active": true,
    "is_locked": false,
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:35:00Z",
    "deleted_at": null
  }
}
```

---

### 2.3 List User Pockets

**HTTP Method**: `GET`  
**URL**: `/v1/pockets`  
**Authentication**: Required

**Query Parameters**:
- `limit` (integer, optional): Number of results per page (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pockets retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "user_id": "507f1f77bcf86cd799439001",
      "name": "Main Pocket",
      "type": "main",
      "category_id": null,
      "balance": 5000000,
      "target_balance": null,
      "is_default": true,
      "is_active": true,
      "is_locked": false,
      "icon": "wallet",
      "icon_color": "#4CAF50",
      "background_color": "#E8F5E9",
      "created_at": "2024-02-01T10:00:00Z",
      "updated_at": "2024-02-02T10:35:00Z",
      "deleted_at": null
    },
    {
      "id": "507f1f77bcf86cd799439021",
      "user_id": "507f1f77bcf86cd799439001",
      "name": "Monthly Savings",
      "type": "saving",
      "category_id": "507f1f77bcf86cd799439013",
      "balance": 500000,
      "target_balance": 1000000,
      "is_default": false,
      "is_active": true,
      "is_locked": false,
      "icon": "piggy-bank",
      "icon_color": "#FF6B6B",
      "background_color": "#FFE5E5",
      "created_at": "2024-02-02T10:30:00Z",
      "updated_at": "2024-02-02T10:35:00Z",
      "deleted_at": null
    }
  ]
}
```

---

### 2.4 List Active Pockets

**HTTP Method**: `GET`  
**URL**: `/v1/pockets/active`  
**Authentication**: Required

**Query Parameters**:
- `limit` (integer, optional): Number of results per page (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Active pockets retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "user_id": "507f1f77bcf86cd799439001",
      "name": "Main Pocket",
      "type": "main",
      "category_id": null,
      "balance": 5000000,
      "target_balance": null,
      "is_default": true,
      "is_active": true,
      "is_locked": false,
      "icon": "wallet",
      "icon_color": "#4CAF50",
      "background_color": "#E8F5E9",
      "created_at": "2024-02-01T10:00:00Z",
      "updated_at": "2024-02-02T10:35:00Z",
      "deleted_at": null
    }
  ]
}
```

---

### 2.5 Get Main Pocket

**HTTP Method**: `GET`  
**URL**: `/v1/pockets/main`  
**Authentication**: Required

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Main pocket retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439020",
    "user_id": "507f1f77bcf86cd799439001",
    "name": "Main Pocket",
    "type": "main",
    "category_id": null,
    "balance": 5000000,
    "target_balance": null,
    "is_default": true,
    "is_active": true,
    "is_locked": false,
    "icon": "wallet",
    "icon_color": "#4CAF50",
    "background_color": "#E8F5E9",
    "created_at": "2024-02-01T10:00:00Z",
    "updated_at": "2024-02-02T10:35:00Z",
    "deleted_at": null
  }
}
```

---

### 2.6 Update Pocket

**HTTP Method**: `PUT`  
**URL**: `/v1/pockets/{id}`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Pocket ID

**Request Body** (all fields optional):
```json
{
  "name": "Updated Savings",
  "type": "saving",
  "category_id": "507f1f77bcf86cd799439013",
  "target_balance": 2000000,
  "icon": "piggy-bank",
  "icon_color": "#FF6B6B",
  "background_color": "#FFE5E5",
  "is_active": true
}
```

**Restrictions**:
- Cannot update main pocket (type = `main`)
- Cannot update locked pocket
- Only non-main pockets can be updated

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pocket updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439021",
    "user_id": "507f1f77bcf86cd799439001",
    "name": "Updated Savings",
    "type": "saving",
    "category_id": "507f1f77bcf86cd799439013",
    "balance": 500000,
    "target_balance": 2000000,
    "is_default": false,
    "is_active": true,
    "is_locked": false,
    "icon": "piggy-bank",
    "icon_color": "#FF6B6B",
    "background_color": "#FFE5E5",
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:40:00Z",
    "deleted_at": null
  }
}
```

---

### 2.7 Lock Pocket

**HTTP Method**: `PUT`  
**URL**: `/v1/pockets/{id}/lock`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Pocket ID

**Restrictions**:
- Pocket balance must be zero
- Cannot lock already locked pocket

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pocket locked successfully",
  "data": null
}
```

---

### 2.8 Unlock Pocket

**HTTP Method**: `PUT`  
**URL**: `/v1/pockets/{id}/unlock`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Pocket ID

**Restrictions**:
- Pocket balance must be zero
- Cannot unlock already unlocked pocket

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pocket unlocked successfully",
  "data": null
}
```

---

### 2.9 Delete Pocket

**HTTP Method**: `DELETE`  
**URL**: `/v1/pockets/{id}`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Pocket ID

**Restrictions**:
- Cannot delete main pocket
- Cannot delete locked pocket
- Pocket balance must be zero

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pocket deleted successfully",
  "data": null
}
```

---

### 2.10 Create System Pocket (Admin Only)

**HTTP Method**: `POST`  
**URL**: `/v1/pockets/admin/{user_id}`  
**Authentication**: Required (Admin only)

**Path Parameters**:
- `user_id` (string): Target user ID

**Request Body**:
```json
{
  "name": "System Pocket",
  "category_id": "507f1f77bcf86cd799439013",
  "icon": "lock",
  "icon_color": "#999999",
  "background_color": "#F5F5F5"
}
```

**Field Descriptions**:
- `name` (string, required): Pocket name (2-255 characters)
- `category_id` (string, optional): Category ID
- `icon` (string, optional): Icon identifier
- `icon_color` (string, optional): Icon color hex
- `background_color` (string, optional): Background color hex

**System Pocket Notes**:
- Type is automatically set to `system`
- Automatically locked (cannot be modified by user)
- Used for system-level money tracking

**Example Response** (201 Created):
```json
{
  "success": true,
  "message": "System pocket created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439022",
    "user_id": "507f1f77bcf86cd799439001",
    "name": "System Pocket",
    "type": "system",
    "category_id": "507f1f77bcf86cd799439013",
    "balance": 0,
    "target_balance": null,
    "is_default": false,
    "is_active": true,
    "is_locked": true,
    "icon": "lock",
    "icon_color": "#999999",
    "background_color": "#F5F5F5",
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:30:00Z",
    "deleted_at": null
  }
}
```

---

### 2.11 List All Pockets (Admin Only)

**HTTP Method**: `GET`  
**URL**: `/v1/pockets/admin`  
**Authentication**: Required (Admin only)

**Query Parameters**:
- `limit` (integer, optional): Number of results per page (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Pockets retrieved successfully",
  "data": [
    {
      "id": "507f1f77bcf86cd799439020",
      "user_id": "507f1f77bcf86cd799439001",
      "name": "Main Pocket",
      "type": "main",
      "category_id": null,
      "balance": 5000000,
      "target_balance": null,
      "is_default": true,
      "is_active": true,
      "is_locked": false,
      "icon": "wallet",
      "icon_color": "#4CAF50",
      "background_color": "#E8F5E9",
      "created_at": "2024-02-01T10:00:00Z",
      "updated_at": "2024-02-02T10:35:00Z",
      "deleted_at": null
    }
  ]
}
```

---

## 3. Transaction APIs

Transactions are the single source of truth for all balance changes. Every balance update is recorded as a transaction.

### 3.1 Create Transaction

**HTTP Method**: `POST`  
**URL**: `/v1/transactions`  
**Authentication**: Required

**Request Body** (varies by transaction type):

#### Income Transaction
```json
{
  "type": "income",
  "amount": 500000,
  "pocket_to": "507f1f77bcf86cd799439020",
  "category_id": "507f1f77bcf86cd799439013",
  "platform_id": "507f1f77bcf86cd799439011",
  "note": "Monthly salary",
  "date": "2024-02-02T10:30:00Z",
  "ref": "SALARY-2024-02"
}
```

#### Expense Transaction
```json
{
  "type": "expense",
  "amount": 50000,
  "pocket_from": "507f1f77bcf86cd799439020",
  "category_id": "507f1f77bcf86cd799439013",
  "platform_id": "507f1f77bcf86cd799439011",
  "note": "Lunch",
  "date": "2024-02-02T12:00:00Z",
  "ref": "EXP-2024-02-001"
}
```

#### Transfer Transaction (Between Pockets)
```json
{
  "type": "transfer",
  "amount": 100000,
  "pocket_from": "507f1f77bcf86cd799439020",
  "pocket_to": "507f1f77bcf86cd799439021",
  "note": "Transfer to savings",
  "date": "2024-02-02T14:00:00Z",
  "ref": "TRF-2024-02-001"
}
```

**Field Descriptions**:
- `type` (string, required): One of `income`, `expense`, `transfer`
- `amount` (number, required): Transaction amount (must be > 0)
- `pocket_from` (string, optional): Source pocket ID (24-character hex)
- `pocket_to` (string, optional): Destination pocket ID (24-character hex)
- `category_id` (string, optional): Category ID (24-character hex)
- `platform_id` (string, optional): Platform ID (24-character hex)
- `note` (string, optional): Transaction note (max 500 characters)
- `date` (string, required): Transaction date (RFC3339 format: `2024-02-02T10:30:00Z`)
- `ref` (string, optional): Reference number (max 100 characters)

**Example Response** (201 Created):
```json
{
  "success": true,
  "message": "Transaction created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439030",
    "user_id": "507f1f77bcf86cd799439001",
    "type": "income",
    "amount": 500000,
    "pocket_from": null,
    "pocket_to": "507f1f77bcf86cd799439020",
    "category_id": "507f1f77bcf86cd799439013",
    "platform_id": "507f1f77bcf86cd799439011",
    "note": "Monthly salary",
    "date": "2024-02-02T10:30:00Z",
    "ref": "SALARY-2024-02",
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:30:00Z",
    "deleted_at": null
  }
}
```

---

### 3.2 Get Transaction by ID

**HTTP Method**: `GET`  
**URL**: `/v1/transactions/{id}`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Transaction ID

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Transaction retrieved successfully",
  "data": {
    "id": "507f1f77bcf86cd799439030",
    "user_id": "507f1f77bcf86cd799439001",
    "type": "income",
    "amount": 500000,
    "pocket_from": null,
    "pocket_from_name": null,
    "pocket_to": "507f1f77bcf86cd799439020",
    "pocket_to_name": "Main Pocket",
    "category_id": "507f1f77bcf86cd799439013",
    "category_name": "Salary",
    "platform_id": "507f1f77bcf86cd799439011",
    "platform_name": "BCA Bank",
    "note": "Monthly salary",
    "date": "2024-02-02T10:30:00Z",
    "ref": "SALARY-2024-02",
    "created_at": "2024-02-02T10:30:00Z",
    "updated_at": "2024-02-02T10:30:00Z",
    "deleted_at": null
  }
}
```

---

### 3.3 List User Transactions

**HTTP Method**: `GET`  
**URL**: `/v1/transactions`  
**Authentication**: Required

**Query Parameters**:
- `type` (string, optional): Filter by transaction type (`income`, `expense`, `transfer`)
- `search` (string, optional): Search in note and reference fields
- `sort_by` (string, optional): Sort field (`date` or `amount`, default: `date`)
- `sort_order` (string, optional): Sort order (`asc` or `desc`, default: `desc`)
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Results per page (default: 10, max: 1000)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Transactions retrieved successfully",
  "data": {
    "items": [
      {
        "id": "507f1f77bcf86cd799439030",
        "user_id": "507f1f77bcf86cd799439001",
        "type": "income",
        "amount": 500000,
        "pocket_from": null,
        "pocket_from_name": null,
        "pocket_to": "507f1f77bcf86cd799439020",
        "pocket_to_name": "Main Pocket",
        "category_id": "507f1f77bcf86cd799439013",
        "category_name": "Salary",
        "platform_id": "507f1f77bcf86cd799439011",
        "platform_name": "BCA Bank",
        "note": "Monthly salary",
        "date": "2024-02-02T10:30:00Z",
        "ref": "SALARY-2024-02",
        "created_at": "2024-02-02T10:30:00Z",
        "updated_at": "2024-02-02T10:30:00Z",
        "deleted_at": null
      }
    ],
    "meta": {
      "total": 1,
      "page": 1,
      "page_size": 10,
      "total_pages": 1
    }
  }
}
```

---

### 3.4 List Pocket Transactions

**HTTP Method**: `GET`  
**URL**: `/v1/transactions/pocket/{pocket_id}`  
**Authentication**: Required

**Path Parameters**:
- `pocket_id` (string): Pocket ID

**Query Parameters**:
- `sort_by` (string, optional): Sort field (`date` or `amount`, default: `date`)
- `sort_order` (string, optional): Sort order (`asc` or `desc`, default: `desc`)
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Results per page (default: 10, max: 1000)

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Transactions retrieved successfully",
  "data": {
    "items": [
      {
        "id": "507f1f77bcf86cd799439030",
        "user_id": "507f1f77bcf86cd799439001",
        "type": "income",
        "amount": 500000,
        "pocket_from": null,
        "pocket_from_name": null,
        "pocket_to": "507f1f77bcf86cd799439020",
        "pocket_to_name": "Main Pocket",
        "category_id": "507f1f77bcf86cd799439013",
        "category_name": "Salary",
        "platform_id": "507f1f77bcf86cd799439011",
        "platform_name": "BCA Bank",
        "note": "Monthly salary",
        "date": "2024-02-02T10:30:00Z",
        "ref": "SALARY-2024-02",
        "created_at": "2024-02-02T10:30:00Z",
        "updated_at": "2024-02-02T10:30:00Z",
        "deleted_at": null
      }
    ],
    "meta": {
      "total": 1,
      "page": 1,
      "page_size": 10,
      "total_pages": 1
    }
  }
}
```

---

### 3.5 Delete Transaction

**HTTP Method**: `DELETE`  
**URL**: `/v1/transactions/{id}`  
**Authentication**: Required

**Path Parameters**:
- `id` (string): Transaction ID

**Note**: Deletion is soft delete. Balance is NOT automatically reversed. This is a record deletion only.

**Example Response** (200 OK):
```json
{
  "success": true,
  "message": "Transaction deleted successfully",
  "data": null
}
```

---

## 4. Transaction Rules

This section defines the exact rules for each transaction type.

### 4.1 Income Transaction

| Aspect | Rule |
|--------|------|
| **Required Fields** | `type`, `amount`, `pocket_to`, `date` |
| **Forbidden Fields** | `pocket_from` must be null |
| **Balance Changes** | `pocket_to` balance increases by `amount` |
| **Validation** | `pocket_to` must exist, be active, and not locked |
| **Use Case** | Money entering the system (salary, refund, gift) |

**Example**:
```json
{
  "type": "income",
  "amount": 500000,
  "pocket_to": "507f1f77bcf86cd799439020",
  "date": "2024-02-02T10:30:00Z"
}
```

---

### 4.2 Expense Transaction

| Aspect | Rule |
|--------|------|
| **Required Fields** | `type`, `amount`, `pocket_from`, `date` |
| **Forbidden Fields** | `pocket_to` must be null |
| **Balance Changes** | `pocket_from` balance decreases by `amount` |
| **Validation** | `pocket_from` must exist, be active, not locked, and have sufficient balance |
| **Use Case** | Money leaving the system (purchase, payment, withdrawal) |

**Example**:
```json
{
  "type": "expense",
  "amount": 50000,
  "pocket_from": "507f1f77bcf86cd799439020",
  "date": "2024-02-02T12:00:00Z"
}
```

---

### 4.3 Transfer Transaction

| Aspect | Rule |
|--------|------|
| **Required Fields** | `type`, `amount`, `pocket_from`, `pocket_to`, `date` |
| **Forbidden Fields** | None (optional: `category_id`, `platform_id`) |
| **Balance Changes** | `pocket_from` decreases by `amount`, `pocket_to` increases by `amount` |
| **Validation** | Both pockets must exist, be active, not locked; `pocket_from` must have sufficient balance; `pocket_from` ≠ `pocket_to` |
| **Use Case** | Money moving between user's pockets (reallocation, savings transfer) |

**Example**:
```json
{
  "type": "transfer",
  "amount": 100000,
  "pocket_from": "507f1f77bcf86cd799439020",
  "pocket_to": "507f1f77bcf86cd799439021",
  "date": "2024-02-02T14:00:00Z"
}
```

---

### 4.4 Transaction Validation Summary

| Validation | Income | Expense | Transfer |
|-----------|--------|---------|----------|
| `pocket_from` required | ❌ | ✅ | ✅ |
| `pocket_to` required | ✅ | ❌ | ✅ |
| `pocket_from` != `pocket_to` | N/A | N/A | ✅ |
| `pocket_from` active | N/A | ✅ | ✅ |
| `pocket_to` active | ✅ | N/A | ✅ |
| `pocket_from` not locked | N/A | ✅ | ✅ |
| `pocket_to` not locked | ✅ | N/A | ✅ |
| Sufficient balance in `pocket_from` | N/A | ✅ | ✅ |
| Amount > 0 | ✅ | ✅ | ✅ |

---

## 5. Balance Consistency

### 5.1 Source of Truth

**Transactions are the single source of truth for all balance changes.**

- Every balance update is recorded as a transaction
- Balances are derived from transactions, not stored independently
- No balance can be modified outside of a transaction

### 5.2 Real-Time Balance Updates

**Pocket balances are updated atomically with transaction creation:**

1. Transaction is created and persisted to database
2. Pocket balance(s) are updated based on transaction type
3. If balance update fails, transaction is rolled back
4. Client receives transaction confirmation with updated balances

### 5.3 Balance Consistency Rules

- **Income**: Only `pocket_to` balance changes (increases)
- **Expense**: Only `pocket_from` balance changes (decreases)
- **Transfer**: Both `pocket_from` (decreases) and `pocket_to` (increases) change
- **No orphaned transactions**: Every transaction must update at least one pocket balance
- **No balance without transaction**: Every balance change must have a corresponding transaction record

### 5.4 Why Balances Must Never Be Updated Outside Transactions

1. **Audit trail**: Every balance change is traceable to a transaction
2. **Consistency**: Balance = sum of all transactions for that pocket
3. **Reconciliation**: Easy to identify discrepancies
4. **Atomicity**: Balance and transaction are updated together or not at all
5. **Reversibility**: Transaction deletion can reverse balance changes if needed

---

## 6. Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "message": "Error description",
  "data": null
}
```

### Common HTTP Status Codes

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad request (validation error, insufficient balance, etc.) |
| 401 | Unauthorized (missing or invalid token) |
| 403 | Forbidden (admin-only endpoint, insufficient permissions) |
| 404 | Not found (resource doesn't exist) |
| 500 | Server error |

### Common Error Messages

| Error | Cause |
|-------|-------|
| `invalid user id` | User ID format is invalid |
| `invalid pocket id` | Pocket ID format is invalid |
| `unauthorized` | User is not authenticated or doesn't own the resource |
| `pocket not found` | Pocket doesn't exist or is deleted |
| `insufficient balance` | Pocket balance is less than transaction amount |
| `pocket is locked` | Cannot perform transaction on locked pocket |
| `pocket is not active` | Cannot perform transaction on inactive pocket |
| `invalid transaction type` | Transaction type is not `income`, `expense`, or `transfer` |
| `pocket_from and pocket_to cannot be the same` | Transfer source and destination are identical |
| `pocket_from is required for EXPENSE transactions` | Expense missing source pocket |
| `pocket_to is required for INCOME transactions` | Income missing destination pocket |
| `both pocket_from and pocket_to are required for TRANSFER transactions` | Transfer missing source or destination |

---

## 7. Authentication

All endpoints require Bearer token authentication in the `Authorization` header:

```
Authorization: Bearer <token>
```

Admin-only endpoints additionally require the user to have admin privileges.

---

## 8. Pagination

List endpoints support pagination with the following parameters:

- `page` (integer, optional): Page number (default: 1, starts at 1)
- `page_size` (integer, optional): Results per page (default: 10, max: 1000)

Or legacy pagination:

- `limit` (integer, optional): Number of results (default: 10, max: 1000)
- `skip` (integer, optional): Number of results to skip (default: 0)

Paginated responses include metadata:

```json
{
  "success": true,
  "message": "...",
  "data": {
    "items": [...],
    "meta": {
      "total": 100,
      "page": 1,
      "page_size": 10,
      "total_pages": 10
    }
  }
}
```

---

## 9. Date Format

All dates use RFC3339 format with timezone:

```
2024-02-02T10:30:00Z
2024-02-02T10:30:00+07:00
```

When sending dates in requests, use the same format. When receiving dates in responses, they are always in UTC (Z timezone).

---

## 10. Summary Table

| Resource | Create | Read | Update | Delete | List |
|----------|--------|------|--------|--------|------|
| **Platform** | ✅ Admin | ✅ | ✅ Admin | ✅ Admin | ✅ |
| **Pocket** | ✅ User | ✅ User | ✅ User | ✅ User | ✅ User |
| **Pocket (System)** | ✅ Admin | ✅ Admin | ❌ | ❌ | ✅ Admin |
| **Transaction** | ✅ User | ✅ User | ❌ | ✅ User | ✅ User |

---

## 11. Quick Reference

### Create Income
```bash
curl -X POST http://localhost:8080/v1/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "income",
    "amount": 500000,
    "pocket_to": "507f1f77bcf86cd799439020",
    "date": "2024-02-02T10:30:00Z"
  }'
```

### Create Expense
```bash
curl -X POST http://localhost:8080/v1/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "expense",
    "amount": 50000,
    "pocket_from": "507f1f77bcf86cd799439020",
    "date": "2024-02-02T12:00:00Z"
  }'
```

### Create Transfer
```bash
curl -X POST http://localhost:8080/v1/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "transfer",
    "amount": 100000,
    "pocket_from": "507f1f77bcf86cd799439020",
    "pocket_to": "507f1f77bcf86cd799439021",
    "date": "2024-02-02T14:00:00Z"
  }'
```

### List User Transactions
```bash
curl -X GET "http://localhost:8080/v1/transactions?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

### Get Pocket Balance
```bash
curl -X GET http://localhost:8080/v1/pockets/507f1f77bcf86cd799439020 \
  -H "Authorization: Bearer <token>"
```
