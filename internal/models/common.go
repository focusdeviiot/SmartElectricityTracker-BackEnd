package models

type Pageable struct {
	PageSize      int   `json:"pageSize" validate:"required"`
	PageNumber    int   `json:"pageNumber" validate:"required"`
	TotalPages    int   `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}
