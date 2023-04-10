package recipes

import (
	structure "bcraft/api/structures"
	request "bcraft/api/structures/requests"
	"bcraft/db"
	"context"
	"errors"
	"strings"
)

func GetRecipes(pagination structure.Pagination) ([]structure.Recipe, error) {
	if pagination.Limit == 0 {
		pagination.Limit = 100
		pagination.Offset = 0
	}
	return db.GetRecipeRecords(pagination)
}

func GetRecipeById(recipeId int) (*structure.Recipe, error) {
	return db.GetRecipeRecord(recipeId)
}

func CreateRecipe(recipeReq request.CreateRecipeRequest) (int, error) {
	ctx := context.TODO()

	recipe := structure.Recipe{
		Name: recipeReq.Name, Description: recipeReq.Description,
		CookingSteps: recipeReq.CookingSteps, Ingredients: recipeReq.Ingredients,
		ImgPath: recipeReq.ImgPath,
	}

	return db.CreateRecipeRecord(ctx, recipe)
}

func GetFilteredRecipes(filters request.RecipeFilters) ([]structure.Recipe, error) {
	if filters.TotalCookingTimeSortDirection != "" && strings.ToLower(filters.TotalCookingTimeSortDirection) != "asc" &&
		strings.ToLower(filters.TotalCookingTimeSortDirection) != "desc" {
		return nil, errors.New("total cooking time sort direction must be either asc or desc")
	}
	return db.GetFilteredRecipeRecords(filters)
}

func RateRecipe(userId, recipeId int, rateRequest request.RateRecipe) (float32, error) {
	if rateRequest.Rate < 0 || rateRequest.Rate > 5 {
		return 0, errors.New("rate must be between 0 and 5")
	}

	return db.RateRecipeRecord(userId, recipeId, rateRequest.Rate)
}

func UpdateRecipe(id int, recipeReq request.CreateRecipeRequest) error {
	ctx := context.TODO()

	recipe := structure.Recipe{
		Id: id, Name: recipeReq.Name, Description: recipeReq.Description,
		CookingSteps: recipeReq.CookingSteps, Ingredients: recipeReq.Ingredients,
		ImgPath: recipeReq.ImgPath,
	}

	return db.UpdateRecipeRecord(ctx, recipe)
}

func DeleteRecipe(id int) error {
	return db.DeleteRecipeRecord(id)
}
