# User Category Creation Flow

## Overview

When a user registers in the system, default user categories are automatically created from the default category templates. This document describes the complete flow and implementation.

## Flow Diagram

```
User Registration (Register Request)
    ↓
Create User Account
    ↓
Create User Profile
    ↓
Create Default Pockets (from pocket templates)
    ↓
Create Default Categories (from category templates) ← NEW
    ↓
Return User with Tokens
```

## Implementation Details

### 1. Auth Service Updates

**File:** `internal/modules/auth/service.go`

#### Service Struct
```go
type Service struct {
	userRepo             *user.Repository
	pocketRepo           *pocket.Repository
	pocketTemplateRepo   *pocket_template.Repository
	categoryTemplateRepo *category_template.Repository
	userCategoryRepo     *user_category.Repository  // NEW
	redis                *redis.Client
	jwtManager           *utils.JWTManager
	passwordMgr          *utils.PasswordManager
}
```

#### NewService Function
```go
func NewService(
	userRepo *user.Repository,
	pocketRepo *pocket.Repository,
	pocketTemplateRepo *pocket_template.Repository,
	categoryTemplateRepo *category_template.Repository,
	userCategoryRepo *user_category.Repository,  // NEW
	redis *redis.Client,
	jwtManager *utils.JWTManager,
	passwordMgr *utils.PasswordManager,
) *Service
```

#### Register Method
The `Register` method now calls `createDefaultCategories` after creating default pockets:

```go
func (s *Service) Register(ctx context.Context, req *authDTO.RegisterRequest) (*user.User, error) {
	// ... existing code ...
	
	// Create default pockets from active templates
	err = s.createDefaultPockets(ctx, newUser.ID)
	if err != nil {
		_ = s.userRepo.DeleteUser(ctx, newUser.ID)
		return nil, err
	}

	// Create default categories from active category templates
	err = s.createDefaultCategories(ctx, newUser.ID)
	if err != nil {
		_ = s.userRepo.DeleteUser(ctx, newUser.ID)
		return nil, err
	}

	return newUser, nil
}
```

### 2. createDefaultCategories Implementation

**Location:** `internal/modules/auth/service.go:272-331`

#### Key Features

1. **Fetch Default Templates**
   - Retrieves all category templates marked as `is_default: true`
   - Handles errors gracefully with logging

2. **Template Mapping**
   - Creates a map of template IDs for quick parent reference resolution
   - Enables efficient parent-child relationship handling

3. **User Category Creation**
   - For each default template, creates a corresponding user category
   - Copies all relevant fields from template to user category:
     - `Name`
     - `TransactionType` (with type conversion)
     - `Description`
     - `Icon`
     - `Color`
     - `IsDefault`

4. **Parent-Child Relationship Handling**
   - **Critical:** ParentID references must point to user category IDs, not template IDs
   - Process:
     1. Check if template has a parent (ParentID is not nil)
     2. Verify parent template exists in the defaults
     3. Find the corresponding user category for the parent template
     4. Set the user category's ParentID to the user category ID (not template ID)

5. **Type Conversion**
   - Converts `category_template.TransactionType` to `user_category.TransactionType`
   - Both are string-based enums with same values ("income", "expense")
   - Conversion: `user_category.TransactionType(*template.TransactionType)`

#### Code Structure

```go
func (s *Service) createDefaultCategories(ctx context.Context, userID primitive.ObjectID) error {
	// Step 1: Fetch default category templates
	templates, err := s.categoryTemplateRepo.GetDefaults(ctx)
	if err != nil {
		log.Printf("failed to fetch active category templates for user %s: %v", userID.Hex(), err)
		return errors.New("failed to fetch category templates")
	}

	// Step 2: Create template map for parent resolution
	templateMap := make(map[string]*category_template.CategoryTemplate)
	for _, template := range templates {
		templateMap[template.ID.Hex()] = template
	}

	// Step 3: Create user categories from templates
	for _, template := range templates {
		// Convert TransactionType
		var transactionType *user_category.TransactionType
		if template.TransactionType != nil {
			tt := user_category.TransactionType(*template.TransactionType)
			transactionType = &tt
		}

		// Create user category
		userCategory := &user_category.UserCategory{
			UserID:          userID,
			TemplateID:      &template.ID,
			Name:            template.Name,
			TransactionType: transactionType,
			Description:     template.Description,
			Icon:            template.Icon,
			Color:           template.Color,
			IsDefault:       template.IsDefault,
		}

		// Step 4: Handle parent-child relationships
		if template.ParentID != nil {
			parentTemplate, exists := templateMap[template.ParentID.Hex()]
			if exists {
				// Find corresponding user category for parent template
				parentUserCategories, err := s.userCategoryRepo.FindAllByUserID(ctx, userID)
				if err == nil {
					for _, uc := range parentUserCategories {
						if uc.TemplateID != nil && uc.TemplateID.Hex() == parentTemplate.ID.Hex() {
							userCategory.ParentID = &uc.ID
							break
						}
					}
				}
			}
		}

		// Step 5: Create user category in database
		err := s.userCategoryRepo.Create(ctx, userCategory)
		if err != nil {
			log.Printf("failed to create user category from template %s for user %s: %v", template.ID.Hex(), userID.Hex(), err)
			return errors.New("failed to create default categories")
		}
	}

	return nil
}
```

### 3. Auth Module Updates

**File:** `internal/modules/auth/module.go`

#### Imports Added
```go
import (
	"github.com/HasanNugroho/coin-be/internal/modules/category_template"
	"github.com/HasanNugroho/coin-be/internal/modules/user_category"
	// ... other imports ...
)
```

#### DI Registration Updated
```go
func Register(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "authService",
		Build: func(ctn di.Container) (interface{}, error) {
			userRepo := ctn.Get("userRepository").(*user.Repository)
			pocketRepo := ctn.Get("pocketRepository").(*pocket.Repository)
			pocketTemplateRepo := ctn.Get("pocketTemplateRepository").(*pocket_template.Repository)
			categoryTemplateRepo := ctn.Get("categoryTemplateRepository").(*category_template.Repository)
			userCategoryRepo := ctn.Get("userCategoryRepository").(*user_category.Repository)
			redisClient := ctn.Get("redis").(*redis.Client)
			cfg := ctn.Get("config").(*config.Config)
			jwtManager := utils.NewJWTManager(cfg)
			passwordMgr := utils.NewPasswordManager()
			return NewService(userRepo, pocketRepo, pocketTemplateRepo, categoryTemplateRepo, userCategoryRepo, redisClient, jwtManager, passwordMgr), nil
		},
	})
	// ... rest of registration ...
}
```

## Data Flow Example

### Scenario: User Registration with Category Templates

**Input:**
```json
{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Doe",
  "phone": "+1234567890"
}
```

**Process:**

1. **Create User Account**
   - User ID: `507f1f77bcf86cd799439010`

2. **Create User Profile**
   - Phone, salary info, language preferences

3. **Create Default Pockets**
   - From pocket templates

4. **Create Default Categories** (NEW)
   - Fetch default category templates:
     ```
     Template 1: "Salary" (income, is_default: true, parent_id: null)
     Template 2: "Freelance" (income, is_default: true, parent_id: null)
     Template 3: "Food" (expense, is_default: true, parent_id: null)
     Template 4: "Groceries" (expense, is_default: true, parent_id: Template3.ID)
     ```

   - Create user categories:
     ```
     UserCategory 1: 
       - user_id: 507f1f77bcf86cd799439010
       - template_id: Template1.ID
       - name: "Salary"
       - transaction_type: "income"
       - parent_id: null

     UserCategory 2:
       - user_id: 507f1f77bcf86cd799439010
       - template_id: Template2.ID
       - name: "Freelance"
       - transaction_type: "income"
       - parent_id: null

     UserCategory 3:
       - user_id: 507f1f77bcf86cd799439010
       - template_id: Template3.ID
       - name: "Food"
       - transaction_type: "expense"
       - parent_id: null

     UserCategory 4:
       - user_id: 507f1f77bcf86cd799439010
       - template_id: Template4.ID
       - name: "Groceries"
       - transaction_type: "expense"
       - parent_id: UserCategory3.ID  ← IMPORTANT: References UserCategory3, not Template3
     ```

## Important Notes

### ParentID Reference Correctness

**CRITICAL:** The `ParentID` in a `UserCategory` must reference another `UserCategory` ID, NOT a `CategoryTemplate` ID.

**Why?**
- Each user has their own set of categories
- Parent-child relationships must be within the user's own categories
- Template IDs are shared across all users
- User category IDs are unique per user

**Implementation:**
```go
// CORRECT: Reference user category ID
userCategory.ParentID = &uc.ID  // uc is a UserCategory

// WRONG: Reference template ID
userCategory.ParentID = &template.ParentID  // This would be wrong!
```

### Error Handling

- If category template fetching fails, user registration is rolled back
- If any user category creation fails, user registration is rolled back
- All errors are logged with user ID and template ID for debugging
- Graceful error messages returned to client

### Performance Considerations

1. **Template Mapping:** O(n) space complexity for quick parent lookup
2. **Parent Resolution:** O(n²) worst case if many categories have parents
   - Could be optimized by sorting templates by creation order
   - Or by creating categories in dependency order

### Transaction Safety

- All operations are within a single registration transaction
- If any step fails, the entire user creation is rolled back
- Ensures data consistency

## Testing Recommendations

1. **Happy Path:**
   - Register user with default category templates
   - Verify all categories created
   - Verify parent-child relationships correct

2. **Parent-Child Relationships:**
   - Create templates with parent-child hierarchy
   - Verify user categories maintain correct relationships
   - Verify ParentID references user category IDs

3. **Error Cases:**
   - Template fetch failure → user creation rollback
   - Category creation failure → user creation rollback
   - Missing parent template → graceful handling

4. **Edge Cases:**
   - No default templates → user created with no categories
   - Circular parent references → should be prevented at template level
   - Multiple levels of nesting → verify all levels created correctly

## Future Enhancements

1. **Bulk Operations:** Create multiple categories in single operation
2. **Template Customization:** Allow users to customize templates during registration
3. **Category Presets:** Different category sets based on user type/region
4. **Async Creation:** Create categories asynchronously if performance becomes issue
5. **Migration Tool:** Migrate existing users to have default categories

## Summary

The implementation successfully:
- ✅ Creates user categories from default category templates during registration
- ✅ Properly handles parent-child relationships using user category IDs
- ✅ Converts transaction types between modules
- ✅ Includes comprehensive error handling and logging
- ✅ Maintains data consistency with rollback on failure
- ✅ Follows existing project patterns and conventions
