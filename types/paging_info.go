package types

// PagedResult - this is the model for all the list endpoints, each endpoint however has its own type that has typed records
type PagedResult struct {
	Records []interface{} `json:"records"`
	Page    PagingInfo    `json:"pagingInfo"`
}

func GetPageInfo(offset, limit, total int, sort []string, filter map[string]interface{}) PagingInfo {
	return PagingInfo{
		Offset:       offset,
		Limit:        limit,
		TotalRecords: total,
		Sort:         sort,
		Filter:       filter,
	}
}

// PagingInfo - this is used to describe the page in a list
type PagingInfo struct {
	Offset       int                    `json:"offset,omitempty" example:"1"`
	Limit        int                    `json:"limit,omitempty" example:"1000"`        // the page size
	TotalRecords int                    `json:"totalRecords,omitempty" example:"6530"` // the total amount of records
	Sort         []string               `json:"sort,omitempty" example:"id"`
	Filter       map[string]interface{} `json:"filter,omitempty"`
	// TODO: needed? should be elsewhere?
	// ? Path string `json:""`
	// ? Params       map[string]interface{} `json:"params,omitempty"`
	// ? RequestBody  map[string]interface{} `json:"requestBody,omitempty"`
}
