package models

type CheckRequest struct {
	Key     string `json:"key" form:"key" binding:"required"`
	Product string `json:"product" form:"product" binding:"required"`
}
