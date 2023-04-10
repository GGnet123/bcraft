package recipes

import (
	"bcraft/api/errs"
	structure "bcraft/api/structures"
	request "bcraft/api/structures/requests"
	response "bcraft/api/structures/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RecipesController struct{}

// @Summary Get Recipes
// @Tags recipes
// @Description Get recipes list
// @Accept json
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} structure.Recipe
// @Failure 500
// @Router /recipes [get]
func (r *RecipesController) Get(c *gin.Context) {
	var pagination structure.Pagination
	var err error

	pagination.Limit, err = strconv.Atoi(c.Query("limit"))
	pagination.Offset, err = strconv.Atoi(c.Query("offset"))

	recipes, err := GetRecipes(pagination)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// @Summary Get Filtered Recipes
// @Tags recipes
// @Description Get filtered recipes list. Leave empty if filter is not needed
// @Accept json
// @Produce json
// @Param Body body request.RecipeFilters true "Body"
// @Success 200 {array} structure.Recipe
// @Failure 500
// @Router /recipes/filter [post]
func (r *RecipesController) GetFiltered(c *gin.Context) {
	var filters request.RecipeFilters
	if err := c.BindJSON(&filters); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recipes, err := GetFilteredRecipes(filters)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recipes)
}

// @Summary Create recipe
// @Tags recipes
// @Description Create recipe. Do not set ingredient ID to create a new ingredient (Leave as 0 or remove it).
// @Accept json
// @Produce json
// @Param Body body request.CreateRecipeRequest true "Body"
// @Success 200 {object} response.SuccessCreateResponse
// @Failure 500
// @Failure 401
// @Failure 400
// @Security Bearer
// @Router /recipes [post]
func (r *RecipesController) Create(c *gin.Context) {
	var body request.CreateRecipeRequest

	if err := c.BindJSON(&body); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recipeId, err := CreateRecipe(body)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.SuccessCreateResponse{Success: true, Id: recipeId})
}

// @Summary Update recipe
// @Tags recipes
// @Description Update recipe
// @Accept json
// @Produce json
// @Param Body body request.CreateRecipeRequest true "Body"
// @Param id path int true "recipe id"
// @Success 200 {object} response.SuccessResponse
// @Failure 500
// @Failure 401
// @Failure 400
// @Security Bearer
// @Router /recipes/{id} [put]
func (r *RecipesController) Update(c *gin.Context) {
	var body request.CreateRecipeRequest

	if err := c.BindJSON(&body); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = UpdateRecipe(id, body)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{Success: true})
}

// @Summary Delete recipe
// @Tags recipes
// @Description Delete recipe
// @Accept json
// @Produce json
// @Param id path int true "recipe id"
// @Success 200 {object} response.SuccessResponse
// @Failure 500
// @Failure 401
// @Security Bearer
// @Router /recipes/{id} [delete]
func (r *RecipesController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = DeleteRecipe(id)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{Success: true})
}

// @Summary Get recipe
// @Tags recipes
// @Description Get recipe by id
// @Accept json
// @Produce json
// @Param id path int true "recipe id"
// @Success 200 {object} structure.Recipe
// @Failure 500
// @Router /recipes/{id} [get]
func (r *RecipesController) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recipe, err := GetRecipeById(id)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// @Summary Rate recipe
// @Tags recipes
// @Description Put a rating to recipe
// @Accept json
// @Produce json
// @Param id path int true "recipe id"
// @Param body body request.RateRecipe true "recipe rate"
// @Success 200 {object} response.SuccessRecipeRateResponse
// @Failure 500
// @Failure 401
// @Failure 400
// @Security Bearer
// @Router /recipes/rate/{id} [post]
func (r *RecipesController) Rate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var body request.RateRecipe
	if err = c.BindJSON(&body); err != nil {
		errs.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uid, _ := c.Get("uid")
	avgRate, err := RateRecipe(uid.(int), id, body)
	if err != nil {
		errs.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.SuccessRecipeRateResponse{Success: true, AvgRating: avgRate})
}
