package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetLimitOffsetParam(ctx *gin.Context) (limit, offset int) {
	var err error
	queryParamLimit := ctx.Query("_limit")
	queryParamOffset := ctx.Query("_offset")
	if queryParamLimit != "" {
		limit, err = strconv.Atoi(queryParamLimit)
		if err != nil {
			return
		}
	} else {
		limit = 4
	}
	if queryParamOffset != "" {
		offset, err = strconv.Atoi(queryParamOffset)
		if err != nil {
			return
		}
	}
	return
}
