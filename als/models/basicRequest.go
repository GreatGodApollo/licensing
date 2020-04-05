package models

type BasicRequest struct {
	Key string `json:"key" form:"key" binding:"required"`
}
