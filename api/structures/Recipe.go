package structure

import "database/sql"

type Recipe struct {
	Id              int             `json:"id" db:"id"`
	Name            string          `json:"name" db:"name"`
	Description     string          `json:"description" db:"description"`
	ImgPathRecord   sql.NullString  `json:"-" db:"img_path"`
	ImgPath         string          `json:"imgPath"`
	AvgRatingRecord sql.NullFloat64 `json:"-" db:"avgRating"`
	AvgRating       float64         `json:"avgRating"`
	Ingredients     []Ingredients   `json:"ingredients"`
	CookingSteps    []CookingSteps  `json:"cookingSteps"`
}

type Ingredients struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CookingSteps struct {
	Id                   int            `json:"-" db:"id"`
	RecipeId             int            `json:"-" db:"recipe_id"`
	Step                 int            `json:"step" db:"step"`
	Description          string         `json:"description" db:"description"`
	ImgPathRecord        sql.NullString `json:"-" db:"img_path"`
	ImgPath              string         `json:"imgPath"`
	CookingTimeInMinutes int            `json:"cookingTimeInMinutes" db:"cook_min_time"`
}
