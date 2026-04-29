package data

import "cafe_store.hiyabnako/internal/validator"

type Filters struct {
	Page int
	PageSize int
	Sort string
	SortSafeList []string
}

type Metadata struct {
	CurrentPage int `json:"current_page,omitzero"`
	PageSize int `json:"page_size,omitzero"`
	FirstPage int `json:"first_page,omitzero"`
	LastPage int `json:"last_page,omitzero"`
	TotalRecords int `json:"total_records,omitzero"`
}

func ValidateFilters(v *validator.Validator, f Filters) {

	v.Check(0 < f.Page,"page", "page must be greater than zero")
	v.Check(f.Page < 10_000_000_000, "page", "page must be less than 10 million")
	v.Check(f.PageSize <= 100, "page_size", "page_size must be a maxumum of 100")
	v.Check(f.PageSize > 0, "page_size", "page_size must be greater than zero")
	v.Check(validator.PermittedValue(f.Sort,f.SortSafeList...),"sort", "invalid sort value" )
}

func (f Filters) Pagel() int {
	return f.Page
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	
	if totalRecords == 0 {
		return Metadata{}
	}
	
	return Metadata{
	CurrentPage: page,
	PageSize: pageSize,
	FirstPage: 1,
	LastPage: (totalRecords + pageSize - 1) / pageSize,
	TotalRecords: totalRecords,
}

}