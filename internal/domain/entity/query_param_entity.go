package entity

type QueryParamEntity struct {
	Search    string
	Page      int64
	Limit     int64
	OrderBy   string
	OrderType string
}
