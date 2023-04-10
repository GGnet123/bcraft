package structure

import (
	structure "bcraft/api/structures"
)

type CreateRecipeRequest struct {
	Name         string                   `json:"name" binding:"required"`
	Description  string                   `json:"description" binding:"required"`
	ImgPath      string                   `json:"imgPath"`
	Ingredients  []structure.Ingredients  `json:"ingredients" binding:"required"`
	CookingSteps []structure.CookingSteps `json:"cookingSteps" binding:"required"`
}

type RecipeFilters struct {
	IngredientsIds                []int  `json:"ingredientsIds"`
	TotalCookingTime              int    `json:"totalCookingTime"`
	TotalCookingTimeSortDirection string `json:"totalCookingTimeSortDirection"`
}

type RateRecipe struct {
	Rate int `json:"rate" binding:"required"`
}
