package models

type License struct {
	Id         int    `json:"id"`
	LicenseKey string `json:"key"`
	Product    string `json:"product"`
	Email      string `json:"email"`
	Valid      bool   `json:"valid"`
	Code       int    `json:"code"`
}

type Licenses struct {
	Code     int       `json:"code"`
	Licenses []License `json:"licenses"`
}
