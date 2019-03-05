package common

import (
	"fmt"
	"strconv"

	"github.com/bytom/errors"
	"github.com/gin-gonic/gin"
)

const (
	defaultSatrtStr = "0"
	defaultLimitStr = "100"
	maxPageLimit    = 1000
)

var (
	errParsePaginationStart = fmt.Errorf("parse pagination start")
	errParsePaginationLimit = fmt.Errorf("parse pagination limit")
)

type PaginationQuery struct {
	Start uint64 `json:"start"`
	Limit uint64 `json:"limit"`
}

type PaginationFun func(start uint64, limit uint64) (interface{}, int, error)

// ParsePagination request meets the standard on https://developer.atlassian.com/server/confluence/pagination-in-the-rest-api/
func ParsePagination(c *gin.Context) (*PaginationQuery, error) {
	startStr := c.DefaultQuery("start", defaultSatrtStr)
	limitStr := c.DefaultQuery("limit", defaultLimitStr)

	start, err := strconv.ParseUint(startStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, errParsePaginationStart)
	}

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, errParsePaginationLimit)
	}

	if limit > maxPageLimit {
		limit = maxPageLimit
	}

	return &PaginationQuery{
		Start: start,
		Limit: limit,
	}, nil
}

type PaginationInfo struct {
	Start   uint64
	Limit   uint64
	HasNext bool
}
