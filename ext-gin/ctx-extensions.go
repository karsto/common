package extgin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
)

// TODO: investigate new binding uri https://github.com/gin-gonic/gin/blob/f98b339b773105aad77f321d0baaa30475bf875d/binding/binding.go#L41
func ParamInt(ctx *gin.Context, paramName string) (int, error) {
	paramStr := ctx.Param(paramName)
	paramStr = strings.Trim(paramStr, "/") // TODO: workaround https://github.com/gin-gonic/gin/pull/1061
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	if param <= 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0, found %v", paramName, paramStr)
	}
	return param, nil
}

func ParamOptionalInt(ctx *gin.Context, paramName string, def int) (int, error) {
	paramStr := ctx.Param(paramName)
	paramStr = strings.Trim(paramStr, "/") // TODO: workaround https://github.com/gin-gonic/gin/pull/1061
	if len(paramStr) < 1 {
		return def, nil
	}

	return ParamInt(ctx, paramName)
}

func QueryInt(ctx *gin.Context, paramName string) (int, error) {
	paramStr := ctx.Query(paramName)
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	if param <= 0 {
		return 0, fmt.Errorf("invalid %s: must be greater than 0, found %v", paramName, paramStr)
	}
	return param, nil
}

//TODO: monitor default proposal use new binduri https://github.com/gin-gonic/gin/issues/1898
func QueryOptionalInt(ctx *gin.Context, paramName string, def int) (int, error) {
	paramStr := ctx.Query(paramName)
	if len(paramStr) < 1 {
		return def, nil
	}

	return QueryInt(ctx, paramName)
}

func QueryTime(ctx *gin.Context, paramName string) (time.Time, error) {
	paramStr := ctx.Query(paramName)
	t, err := dateparse.ParseAny(paramStr)
	if err != nil {
		return time.Time{}, err
	}

	const sanityPeriod = 10 * 356 * 24 * time.Hour
	if t.Before(time.Now().Add(-sanityPeriod)) || t.After(time.Now().Add(sanityPeriod)) {
		return time.Time{}, fmt.Errorf("invalid %s: must be a valid time format (yyyy-mm-dd) and within 10 years of this date to be considered valid", paramName)
	}
	return t, nil
}

func QueryOptionalTime(ctx *gin.Context, paramName string, def time.Time) (time.Time, error) {
	paramStr := ctx.Query(paramName)
	if len(paramStr) < 1 {
		return def, nil
	}
	return QueryTime(ctx, paramName)
}

func QueryString(ctx *gin.Context, paramName string) (string, error) {
	paramStr := ctx.Query(paramName)
	if len(paramStr) < 1 {
		return "", fmt.Errorf("invalid %s: must be a valid string that contains letters", paramStr)
	}
	return strings.TrimSpace(paramStr), nil
}

func ParamString(ctx *gin.Context, paramName string) (string, error) {
	paramStr := ctx.Param(paramName)
	if len(paramStr) < 1 {
		return "", fmt.Errorf("invalid %s: must be a valid string that contains letters", paramStr)
	}
	return strings.TrimSpace(paramStr), nil
}
