package transaction

import (
	"net/http"
	"strconv"

	"github.com/HasanNugroho/coin-be/internal/core/utils"
	"github.com/HasanNugroho/coin-be/internal/modules/transaction/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

// CreateIncome godoc
// @Summary Create income transaction
// @Description Create an income transaction and automatically distribute to allocations
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Income transaction details"
// @Success 201 {object} map[string]interface{} "Income created and distributed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions/income [post]
func (c *Controller) CreateIncome(ctx *gin.Context) {
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

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	result, err := c.service.CreateIncome(ctx, userOID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Income created and distributed successfully", result)
	ctx.JSON(http.StatusCreated, resp)
}

// CreateExpense godoc
// @Summary Create expense transaction
// @Description Create an expense transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Expense transaction details"
// @Success 201 {object} map[string]interface{} "Expense created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions/expense [post]
func (c *Controller) CreateExpense(ctx *gin.Context) {
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

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	transaction, err := c.service.CreateExpense(ctx, userOID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Expense created successfully", transaction)
	ctx.JSON(http.StatusCreated, resp)
}

// GetTransactions godoc
// @Summary Get all transactions
// @Description Get all transactions for the authenticated user with pagination
// @Tags Transactions
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param skip query int false "Skip" default(0)
// @Success 200 {object} map[string]interface{} "Transactions retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions [get]
func (c *Controller) GetTransactions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "50"), 10, 64)
	skip, _ := strconv.ParseInt(ctx.DefaultQuery("skip", "0"), 10, 64)

	transactions, err := c.service.GetTransactions(ctx, userOID, limit, skip)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Transactions retrieved successfully", transactions)
	ctx.JSON(http.StatusOK, resp)
}

// FilterTransactions godoc
// @Summary Filter transactions
// @Description Filter transactions by type, category, allocation, and date range
// @Tags Transactions
// @Accept json
// @Produce json
// @Param type query string false "Transaction type (income or expense)"
// @Param category_id query string false "Category ID"
// @Param allocation_id query string false "Allocation ID"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Limit" default(50)
// @Param skip query int false "Skip" default(0)
// @Success 200 {object} map[string]interface{} "Transactions filtered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions/filter [get]
func (c *Controller) FilterTransactions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		resp := utils.NewErrorResponse(http.StatusUnauthorized, "unauthorized")
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userOID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid user ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req dto.FilterTransactionRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	transactions, err := c.service.FilterTransactions(ctx, userOID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Transactions filtered successfully", transactions)
	ctx.JSON(http.StatusOK, resp)
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Get a specific transaction by its ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid transaction ID"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Security BearerAuth
// @Router /transactions/{id} [get]
func (c *Controller) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	transactionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid transaction ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	transaction, err := c.service.GetTransactionByID(ctx, transactionID)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusNotFound, err.Error())
		ctx.JSON(http.StatusNotFound, resp)
		return
	}

	resp := utils.NewSuccessResponse("Transaction retrieved successfully", transaction)
	ctx.JSON(http.StatusOK, resp)
}

// UpdateTransaction godoc
// @Summary Update transaction
// @Description Update an existing transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param request body dto.UpdateTransactionRequest true "Updated transaction details"
// @Success 200 {object} map[string]interface{} "Transaction updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions/{id} [put]
func (c *Controller) UpdateTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	transactionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid transaction ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	var req dto.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, err.Error())
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	transaction, err := c.service.UpdateTransaction(ctx, transactionID, &req)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Transaction updated successfully", transaction)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteTransaction godoc
// @Summary Delete transaction
// @Description Delete a transaction by ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid transaction ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /transactions/{id} [delete]
func (c *Controller) DeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	transactionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp := utils.NewErrorResponse(http.StatusBadRequest, "invalid transaction ID")
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.service.DeleteTransaction(ctx, transactionID); err != nil {
		resp := utils.NewErrorResponse(http.StatusInternalServerError, err.Error())
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := utils.NewSuccessResponse("Transaction deleted successfully", nil)
	ctx.JSON(http.StatusOK, resp)
}
