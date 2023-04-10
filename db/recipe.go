package db

import (
	structure "bcraft/api/structures"
	structure2 "bcraft/api/structures/requests"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

const (
	getRecipes            = "SELECT r.id, r.name, r.description, r.img_path, (SELECT AVG(rating) FROM user_recipe_ratings WHERE recipe_id = r.id) AS avgRating from recipes r"
	getRecipeIngredients  = "SELECT i.id, i.name FROM ingredients i INNER JOIN recipe_ingredients ri ON i.id = ri.ingredient_id WHERE ri.recipe_id = $1"
	getRecipeCookingSteps = "SELECT id, recipe_id, step, description, cook_min_time, img_path FROM cooking_steps WHERE recipe_id = $1 ORDER BY step;"
	getRecipeById         = "SELECT id, name, description, img_path, (SELECT AVG(rating) FROM user_recipe_ratings WHERE recipe_id = recipes.id) AS avgRating FROM recipes WHERE id = $1"
	getIngredientByName   = "SELECT id, name FROM ingredients WHERE name=$1"

	updateRecipe            = "UPDATE recipes SET name=$1, description=$2, img_path=$3 WHERE id=$4"
	clearRecipeIngredients  = "DELETE FROM recipe_ingredients WHERE recipe_id=$1"
	clearRecipeCookingSteps = "DELETE FROM cooking_steps WHERE recipe_id=$1"

	createIngredient         = "INSERT INTO ingredients(name) VALUES($1) RETURNING id"
	linkIngredientWithRecipe = "INSERT INTO recipe_ingredients(recipe_id, ingredient_id) VALUES ($1, $2)"
	createCookingStep        = "INSERT INTO cooking_steps(recipe_id, step, description, cook_min_time, img_path) VALUES($1, $2, $3, $4, $5)"
	createRecipe             = "INSERT INTO recipes(name, description, img_path) VALUES ($1, $2, $3) RETURNING id"

	deleteRecipe = "DELETE FROM recipes WHERE id=$1"

	cookingStepsAndRecipeJoin = "left join cooking_steps cs on cs.recipe_id = r.id"
	ingredientsAndRecipeJoin  = "left join recipe_ingredients ri on r.id = ri.recipe_id"

	setRecipeRating             = "INSERT INTO user_recipe_ratings(recipe_id, user_id, rating) VALUES ($1, $2, $3)"
	userAlreadyRatedRecipeCheck = "SELECT EXISTS(SELECT FROM user_recipe_ratings WHERE recipe_id=$1 AND user_id=$2)"
	getAvgRecipeRating          = "SELECT AVG(rating) FROM user_recipe_ratings WHERE recipe_id = $1"
)

func GetFilteredRecipeRecords(filters structure2.RecipeFilters) ([]structure.Recipe, error) {
	sqlFilters := structure.SQLFilters{
		Params:           nil,
		Joins:            make(map[string]string),
		Groupings:        make(map[string]string),
		HavingConditions: nil,
		WhereConditions:  nil,
		Sorting:          nil,
	}

	sqlFilters = setSqlFilters(filters, sqlFilters)

	query := getFilterQuery(sqlFilters)

	rows, err := Conn.Query(query, sqlFilters.Params...)
	if err != nil {
		return nil, err
	}

	var recipes []structure.Recipe
	for rows.Next() {
		var recipe structure.Recipe
		err = rows.Scan(&recipe.Id, &recipe.Name, &recipe.Description, &recipe.ImgPathRecord, &recipe.AvgRatingRecord)
		if err != nil {
			return nil, err
		}
		recipe.ImgPath = recipe.ImgPathRecord.String
		recipe.AvgRating = recipe.AvgRatingRecord.Float64
		recipe.Ingredients, err = getRecipeIngredientRecords(recipe.Id)
		if err != nil {
			return nil, err
		}
		recipe.CookingSteps, err = getRecipeCookingStepRecords(recipe.Id)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func setSqlFilters(filters structure2.RecipeFilters, sqlFilters structure.SQLFilters) structure.SQLFilters {
	if filters.TotalCookingTimeSortDirection != "" {
		sqlFilters.Joins["cooking_steps"] = cookingStepsAndRecipeJoin
		sqlFilters.Sorting = append(sqlFilters.Sorting, "sum(cs.cook_min_time) "+filters.TotalCookingTimeSortDirection)
		sqlFilters.Groupings["r.id"] = "r.id"
	}

	if filters.TotalCookingTime != 0 {
		sqlFilters.Joins["cooking_steps"] = cookingStepsAndRecipeJoin
		sqlFilters.Groupings["r.id"] = "r.id"
		sqlFilters.Params = append(sqlFilters.Params, filters.TotalCookingTime)
		sqlFilters.HavingConditions = append(sqlFilters.HavingConditions, fmt.Sprintf("sum(cs.cook_min_time) = $%d", len(sqlFilters.Params)))
	}

	if len(filters.IngredientsIds) > 0 {
		sqlFilters.Joins["ingredients"] = ingredientsAndRecipeJoin
		sqlFilters.Groupings["r.id"] = "r.id"
		bufferedQuery := bytes.NewBufferString("ri.ingredient_id in (")
		for i, v := range filters.IngredientsIds {
			if i > 0 {
				bufferedQuery.WriteString(",")
			}
			bufferedQuery.WriteString(strconv.Itoa(v))
		}
		bufferedQuery.WriteString(")")
		sqlFilters.WhereConditions = append(sqlFilters.WhereConditions, bufferedQuery.String())
	}
	return sqlFilters
}

func getFilterQuery(sqlFilters structure.SQLFilters) string {
	query := getRecipes
	for _, join := range sqlFilters.Joins {
		query += " " + join + " "
	}
	if len(sqlFilters.WhereConditions) > 0 {
		query += " where " + strings.Join(sqlFilters.WhereConditions, " AND ")
	}

	if len(sqlFilters.Groupings) > 0 {
		query += " group by"
		for _, group := range sqlFilters.Groupings {
			query += " " + group + ","
		}
		query = strings.TrimRight(query, ",")
	}

	if len(sqlFilters.HavingConditions) > 0 {
		query += " having " + strings.Join(sqlFilters.HavingConditions, ", ")
	}
	if len(sqlFilters.Sorting) > 0 {
		query += " order by " + strings.Join(sqlFilters.Sorting, ",")
	}

	logrus.Info(query)
	return query
}

func RateRecipeRecord(userId, recipeId, rate int) (float32, error) {
	var alreadyRated bool
	logrus.Infof(fmt.Sprintf("rate: %d", rate))
	err := Conn.QueryRow(userAlreadyRatedRecipeCheck, recipeId, userId).Scan(&alreadyRated)
	if err != nil {
		return 0, err
	}

	if alreadyRated {
		return 0, errors.New("user has already rated this recipe")
	}

	_, err = Conn.Exec(setRecipeRating, recipeId, userId, rate)
	if err != nil {
		return 0, err
	}
	avgRate, err := getAvgRecipeRate(recipeId)
	if err != nil {
		return 0, err
	}
	return avgRate, nil
}

func getAvgRecipeRate(recipeId int) (float32, error) {
	var avgRate float32
	err := Conn.QueryRow(getAvgRecipeRating, recipeId).Scan(&avgRate)
	if err != nil {
		return 0, err
	}
	return avgRate, nil
}

func CreateRecipeRecord(ctx context.Context, recipe structure.Recipe) (int, error) {
	tx, err := Conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	var recipeId int
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx, createRecipe, recipe.Name, recipe.Description, recipe.ImgPath).Scan(&recipeId)
	if err != nil {
		return 0, err
	}
	err = createIngredientRecordsForRecipeRecord(ctx, recipe.Ingredients, tx, recipeId)
	if err != nil {
		return 0, err
	}

	err = createCookingStepRecordsForRecipe(ctx, recipe.CookingSteps, tx, recipeId)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return recipeId, nil
}

func UpdateRecipeRecord(ctx context.Context, recipe structure.Recipe) error {
	tx, err := Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, updateRecipe, recipe.Name, recipe.Description, recipe.ImgPath, recipe.Id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, clearRecipeIngredients, recipe.Id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, clearRecipeCookingSteps, recipe.Id)
	if err != nil {
		return err
	}

	err = createIngredientRecordsForRecipeRecord(ctx, recipe.Ingredients, tx, recipe.Id)
	if err != nil {
		return err
	}

	err = createCookingStepRecordsForRecipe(ctx, recipe.CookingSteps, tx, recipe.Id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func DeleteRecipeRecord(id int) error {
	_, err := Conn.Exec(deleteRecipe, id)
	return err
}

func createCookingStepRecordsForRecipe(ctx context.Context, cookingSteps []structure.CookingSteps, tx *sql.Tx, recipeId int) error {
	for _, step := range cookingSteps {
		_, err := tx.ExecContext(ctx, createCookingStep, recipeId, step.Step, step.Description, step.CookingTimeInMinutes, step.ImgPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func createIngredientRecordsForRecipeRecord(ctx context.Context, ingredients []structure.Ingredients, tx *sql.Tx, recipeId int) error {
	for _, ingredient := range ingredients {
		if ingredient.Id == 0 {
			tx.QueryRow(getIngredientByName, ingredient.Name).Scan(&ingredient.Id, &ingredient.Name)
		}

		if ingredient.Id == 0 {
			err := tx.QueryRowContext(ctx, createIngredient, ingredient.Name).Scan(&ingredient.Id)
			if err != nil {
				return err
			}
		}

		_, err := tx.ExecContext(ctx, linkIngredientWithRecipe, recipeId, ingredient.Id)

		if err != nil {
			return err
		}
	}
	return nil
}

func getRecipeCookingStepRecords(recipeId int) ([]structure.CookingSteps, error) {
	cookingSteps := []structure.CookingSteps{}

	rows, err := Conn.Query(getRecipeCookingSteps, recipeId)
	if err != nil {
		return []structure.CookingSteps{}, err
	}
	for rows.Next() {
		var cookingStep structure.CookingSteps
		err = rows.Scan(&cookingStep.Id, &cookingStep.RecipeId, &cookingStep.Step, &cookingStep.Description, &cookingStep.CookingTimeInMinutes, &cookingStep.ImgPathRecord)
		if err != nil {
			return []structure.CookingSteps{}, err
		}
		cookingStep.ImgPath = cookingStep.ImgPathRecord.String
		cookingSteps = append(cookingSteps, cookingStep)
	}
	return cookingSteps, nil
}

func getRecipeIngredientRecords(recipeId int) ([]structure.Ingredients, error) {
	ingredients := []structure.Ingredients{}

	rows, err := Conn.Query(getRecipeIngredients, recipeId)
	if err != nil {
		return []structure.Ingredients{}, err
	}
	for rows.Next() {
		var ingredient structure.Ingredients
		err = rows.Scan(&ingredient.Id, &ingredient.Name)
		if err != nil {
			return []structure.Ingredients{}, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func GetRecipeRecords(pagination structure.Pagination) ([]structure.Recipe, error) {
	recipes := []structure.Recipe{}

	query := fmt.Sprintf("%s limit %d offset %d", getRecipes, pagination.Limit, pagination.Offset)
	rows, err := Conn.Query(query)
	if err != nil {
		return []structure.Recipe{}, err
	}
	for rows.Next() {
		var recipe structure.Recipe
		err = rows.Scan(&recipe.Id, &recipe.Name, &recipe.Description, &recipe.ImgPathRecord, &recipe.AvgRatingRecord)
		if err != nil {
			return []structure.Recipe{}, err
		}
		recipe.ImgPath = recipe.ImgPathRecord.String
		recipe.AvgRating = recipe.AvgRatingRecord.Float64
		recipe.Ingredients, err = getRecipeIngredientRecords(recipe.Id)
		if err != nil {
			return []structure.Recipe{}, err
		}
		recipe.CookingSteps, err = getRecipeCookingStepRecords(recipe.Id)
		if err != nil {
			return []structure.Recipe{}, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func GetRecipeRecord(recipeId int) (*structure.Recipe, error) {
	var recipe structure.Recipe

	err := Conn.QueryRow(getRecipeById, recipeId).Scan(&recipe.Id, &recipe.Name, &recipe.Description, &recipe.ImgPathRecord, &recipe.AvgRatingRecord)
	if err != nil {
		return nil, err
	}
	recipe.ImgPath = recipe.ImgPathRecord.String
	recipe.AvgRating = recipe.AvgRatingRecord.Float64
	recipe.Ingredients, err = getRecipeIngredientRecords(recipeId)
	if err != nil {
		return nil, err
	}
	recipe.CookingSteps, err = getRecipeCookingStepRecords(recipeId)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}
