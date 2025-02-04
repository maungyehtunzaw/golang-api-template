package utils

import (
	"fmt"
	"strconv"

	"golang-api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds the request query parameters for paging.
type PaginationParams struct {
	Page  int
	Limit int
	Query string
}

// ParsePagination reads `page` and `limit` from the query string
// (e.g., ?page=2&limit=10). Defaults are page=1, limit=10 if not provided.
func ParsePagination(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// Offset calculates the OFFSET for SQL queries based on the current page & limit.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// LaravelPaginatedResponse is a helper to return a JSON response
// resembling Laravel's pagination structure.
//
// - httpCode: HTTP status code (e.g., 200, 400).
// - message: short message (e.g., "List of users").
// - items: the slice/array of data items for the current page.
// - currentPage: the current page number (e.g. p.Page).
// - perPage: items shown per page (e.g. p.Limit).
// - total: total number of records across all pages.
// - basePath: the base URL path for constructing prev/next links (e.g., "/api/v1/users").
//
// Additional fields included:
// - `from`: starting record number in this page
// - `to`: last record number in this page
// - `last_page`: total number of pages
// - `prev_page_url` / `next_page_url`
// - `first_page_url` / `last_page_url`
func PaginatedResponse(
	c *gin.Context,
	httpCode int,
	message string,
	items interface{},
	currentPage int,
	perPage int,
	total int,
	basePath string,
) {
	// Avoid division by zero
	if perPage < 1 {
		perPage = 1
	}

	lastPage := (total + perPage - 1) / perPage // simple ceiling division

	// from and to represent the item indices (in a 1-based index sense)
	from := (currentPage-1)*perPage + 1
	to := from + perPage - 1
	if to > total {
		to = total
	}
	if from > total {
		from = 0
		to = 0
	}

	// Build URLs
	// We'll assume query params remain the same except for `page`.
	// You can get the query string from c.Request.URL.Query().
	// For simplicity, let's just create minimal links.
	var prevPageURL, nextPageURL, firstPageURL, lastPageURL string

	// We'll keep the existing query parameters, but replace "page" only.
	queryParams := c.Request.URL.Query()

	// first page
	queryParams.Set("page", "1")
	firstPageURL = fmt.Sprintf("%s?%s", basePath, queryParams.Encode())

	// last page
	queryParams.Set("page", strconv.Itoa(lastPage))
	lastPageURL = fmt.Sprintf("%s?%s", basePath, queryParams.Encode())

	// prev page
	if currentPage > 1 {
		queryParams.Set("page", strconv.Itoa(currentPage-1))
		prevPageURL = fmt.Sprintf("%s?%s", basePath, queryParams.Encode())
	}

	// next page
	if currentPage < lastPage {
		queryParams.Set("page", strconv.Itoa(currentPage+1))
		nextPageURL = fmt.Sprintf("%s?%s", basePath, queryParams.Encode())
	}

	// Construct response data in Laravel style
	data := gin.H{
		"current_page":   currentPage,
		"data":           items,
		"first_page_url": firstPageURL,
		"from":           from,
		"last_page":      lastPage,
		"last_page_url":  lastPageURL,
		"next_page_url":  nextPageURL,
		"path":           basePath,
		"per_page":       perPage,
		"prev_page_url":  prevPageURL,
		"to":             to,
		"total":          total,
	}

	response.Success(c, httpCode, message, data)
}
