package types

type ListRequest struct {
	Limit  int                    `json:"limit,omitempty" form:"limit,default=20" binding:"min=1,max=10000"`
	Offset int                    `json:"offset,omitempty" form:"offset,default=0" binding:"min=0"`
	Sort   []string               `json:"sort,omitempty" form:"sort" binding:""`
	Filter map[string]interface{} `json:"filter,omitempty" form:"filter"`
}
