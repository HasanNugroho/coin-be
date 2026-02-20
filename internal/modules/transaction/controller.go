package transaction

import (
	"log"
	"net/http"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction/dto"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction for the authenticated user
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Transaction details"
// @Success 201 {object} map[string]interface{} "Transaction created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/transactions [post]
func (c *Controller) CreateTransaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req dto.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := utils.ValidateRequest(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	transaction, err := c.service.CreateTransaction(ctx, userID.(string), &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	txResp := c.mapToResponse(transaction)
	resp := utils.NewSuccessResponse("Transaction created successfully", txResp)
	ctx.JSON(http.StatusCreated, resp)
}

// GetTransaction godoc
// @Summary Get transaction by ID
// @Description Get a specific transaction by ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Security BearerAuth
// @Router /v1/transactions/{id} [get]
func (c *Controller) GetTransaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	transaction, err := c.service.GetTransactionByID(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	txResp := c.mapToResponse(transaction)
	resp := utils.NewSuccessResponse("Transaction retrieved successfully", txResp)
	ctx.JSON(http.StatusOK, resp)
}

// ListUserTransactions godoc
// @Summary List user transactions
// @Description Get a list of user transactions with pagination
// @Tags Transactions
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/transactions [get]
func (c *Controller) ListUserTransactions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	// Get query parameters
	transactionType := ctx.Query("type")

	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(ctx, 10)

	// Parse sorting parameters (allowed fields: date, amount)
	allowedFields := []string{"date", "amount"}
	sorting := utils.ParseSortParams(ctx, allowedFields, "date")

	// Prepare filter
	var typeFilter *string
	if transactionType != "" {
		typeFilter = &transactionType
	}

	search := ctx.Query("search")
	var searchFilter *string
	if search != "" {
		searchFilter = &search
	}

	// Fetch transactions with filter, pagination and sorting
	transactions, total, err := c.service.GetUserTransactionsWithSort(ctx, userID.(string), typeFilter, searchFilter, pagination.Page, pagination.PageSize, sorting.SortBy, sorting.SortOrder)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// Calculate pagination metadata
	meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)

	// Build paginated response
	respData := utils.BuildPaginatedResponse(transactions, meta)
	resp := utils.NewSuccessResponse("Transactions retrieved successfully", respData)
	ctx.JSON(http.StatusOK, resp)
}

// ListPocketTransactions godoc
// @Summary List pocket transactions
// @Description Get a list of transactions for a specific pocket
// @Tags Transactions
// @Accept json
// @Produce json
// @Param pocket_id path string true "Pocket ID"
// @Param limit query int false "Limit (default: 10, max: 1000)"
// @Param skip query int false "Skip (default: 0)"
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Security BearerAuth
// @Router /v1/transactions/pocket/{pocket_id} [get]
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
	transactions, total, err := c.service.GetPocketTransactionsWithSort(ctx, userID.(string), pocketID, pagination.Page, pagination.PageSize, sorting.SortBy, sorting.SortOrder)
	log.Printf("Transaction: %s", transactions)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	// Calculate pagination metadata
	meta := utils.CalculatePaginationMeta(total, pagination.Page, pagination.PageSize)

	// Build paginated response
	respData := utils.BuildPaginatedResponse(transactions, meta)
	resp := utils.NewSuccessResponse("Transactions retrieved successfully", respData)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteTransaction godoc
// @Summary Delete transaction
// @Description Soft delete a transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Security BearerAuth
// @Router /v1/transactions/{id} [delete]
func (c *Controller) DeleteTransaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	id := ctx.Param("id")

	err := c.service.DeleteTransaction(ctx, userID.(string), id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	resp := utils.NewSuccessResponse("Transaction deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) mapToResponse(transaction *Transaction) *dto.TransactionResponse {
	var pocketFromID *string
	if transaction.PocketFromID != nil {
		id := transaction.PocketFromID.Hex()
		pocketFromID = &id
	}

	var pocketToID *string
	if transaction.PocketToID != nil {
		id := transaction.PocketToID.Hex()
		pocketToID = &id
	}

	var userPlatformFrom *string
	if transaction.UserPlatformFromID != nil {
		id := transaction.UserPlatformFromID.Hex()
		userPlatformFrom = &id
	}

	var userPlatformTo *string
	if transaction.UserPlatformToID != nil {
		id := transaction.UserPlatformToID.Hex()
		userPlatformTo = &id
	}

	var categoryID *string
	if transaction.CategoryID != nil {
		id := transaction.CategoryID.Hex()
		categoryID = &id
	}

	return &dto.TransactionResponse{
		ID:                 transaction.ID.Hex(),
		UserID:             transaction.UserID.Hex(),
		Type:               transaction.Type,
		Amount:             transaction.Amount,
		PocketFromID:       pocketFromID,
		PocketToID:         pocketToID,
		UserPlatformFromID: userPlatformFrom,
		UserPlatformToID:   userPlatformTo,
		CategoryID:         categoryID,
		Note:               transaction.Note,
		Date:               transaction.Date,
		Ref:                transaction.Ref,
		CreatedAt:          transaction.CreatedAt,
		UpdatedAt:          transaction.UpdatedAt,
		DeletedAt:          transaction.DeletedAt,
	}
}

func (c *Controller) mapToResponseList(transactions []*Transaction) []*dto.TransactionResponse {
	responses := make([]*dto.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = c.mapToResponse(transaction)
	}
	return responses
}
