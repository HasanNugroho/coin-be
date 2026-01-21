package dto

import "time"

type CreateTargetRequest struct {
	AllocationID string    `json:"allocation_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	TargetAmount float64   `json:"target_amount" binding:"required,gt=0"`
	Deadline     time.Time `json:"deadline" binding:"required"`
}

type UpdateTargetRequest struct {
	Name         string    `json:"name"`
	TargetAmount float64   `json:"target_amount" binding:"gt=0"`
	Deadline     time.Time `json:"deadline"`
}
