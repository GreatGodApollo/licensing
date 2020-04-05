package models

type LicenseResponse struct {
	LicenseKey string `json:"license_key"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Code       int    `json:"code"`
}
