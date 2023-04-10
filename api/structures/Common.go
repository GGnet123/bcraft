package structure

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type SQLFilters struct {
	Params           []interface{}
	Joins            map[string]string
	Groupings        map[string]string
	HavingConditions []string
	WhereConditions  []string
	Sorting          []string
}
