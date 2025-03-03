package models

type PaginationParam struct {
	Page     int
	PageSize int
}

type Pagination struct {
	Page             int   `json:"page"`           // The current page
	PageSize         int   `json:"limit"`          // The size of the page
	TotalPage        int64 `json:"total_pages"`    // The total number of pages
	TotalDataPerPage int64 `json:"total_per_page"` // The total number of data per page
	TotalData        int64 `json:"total"`          // The total number of data
}

type PaginationData[T any] struct {
	Page             int   `json:"page"`               // The current page
	PageSize         int   `json:"limit"`              // The size of the page
	TotalPage        int64 `json:"total_pages"`        // The total number of pages
	TotalDataPerPage int64 `json:"total_row_per_page"` // The total number of data per page
	TotalData        int64 `json:"total_rows"`         // The total number of data
	Data             []*T  `json:"data"`               // The actual data
}