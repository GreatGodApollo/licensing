package models

type LicenseRequest struct {
	Email   string `json:"email" form:"email" binding:"required"`
	Product string `json:"product" form:"product" binding:"required"`
}
