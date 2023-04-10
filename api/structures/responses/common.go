package structures

type SuccessResponse struct {
	Success bool `json:"success"`
}

type SuccessRecipeRateResponse struct {
	Success   bool    `json:"success"`
	AvgRating float32 `json:"avgRecipeRate"`
}

type SuccessCreateResponse struct {
	Success bool `json:"success"`
	Id      int  `json:"id"`
}
