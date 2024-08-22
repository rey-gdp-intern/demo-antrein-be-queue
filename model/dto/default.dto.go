package dto

type DefaultDTOResponseWrapper struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type NoBodyDTOResponseWrapper struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type PaginationDTOWrapper struct {
	PageSize  int         `json:"page_size"`
	Page      int         `json:"page"`
	TotalPage int         `json:"total_page"`
	Data      interface{} `json:"page_data"`
}
