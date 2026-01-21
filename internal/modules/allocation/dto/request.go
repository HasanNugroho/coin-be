package dto

type CreateAllocationRequest struct {
	Name         string   `json:"name" binding:"required"`
	Priority     int      `json:"priority" binding:"required,gt=0"`
	Percentage   float64  `json:"percentage" binding:"required,gt=0,lte=100"`
	TargetAmount *float64 `json:"target_amount,omitempty"`
}

type UpdateAllocationRequest struct {
	Name         string   `json:"name"`
	Priority     int      `json:"priority" binding:"gt=0"`
	Percentage   float64  `json:"percentage" binding:"gt=0,lte=100"`
	TargetAmount *float64 `json:"target_amount,omitempty"`
	IsActive     *bool    `json:"is_active"`
}
