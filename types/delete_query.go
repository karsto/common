package types

type DeleteQuery struct {
	// ID  int   `form:"id"` TODO: https://github.com/gin-gonic/gin/pull/1061
	IDs []int `form:"ids"`
}
