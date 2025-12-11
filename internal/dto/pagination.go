package dto

type Pagination struct {
	HasNextPage bool  `json:"hasNextPage"`
	Total       int64 `json:"total"`
	Status      int   `json:"status"`
}
