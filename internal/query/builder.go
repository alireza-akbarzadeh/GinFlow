package query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ===========================================
// QUERY BUILDER
// ===========================================

// QueryBuilder helps build database queries with a fluent API
// Supports pagination, filtering, sorting, and search
type QueryBuilder struct {
	request        *QueryParams
	db             *gorm.DB
	allowedFilters map[string]bool
	allowedSorts   map[string]bool
	defaultSort    []SortField
	searchColumns  []string
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{
		db:             db,
		allowedFilters: make(map[string]bool),
		allowedSorts:   make(map[string]bool),
		defaultSort:    []SortField{{Field: "created_at", Direction: SortDesc}},
	}
}

// Alias for backward compatibility
func NewPaginationBuilder(db *gorm.DB) *QueryBuilder {
	return NewQueryBuilder(db)
}

// ===========================================
// BUILDER CONFIGURATION METHODS
// ===========================================

// WithRequest sets the query params
func (qb *QueryBuilder) WithRequest(req *QueryParams) *QueryBuilder {
	qb.request = req
	return qb
}

// AllowFilters sets allowed filter fields (for security)
func (qb *QueryBuilder) AllowFilters(fields ...string) *QueryBuilder {
	for _, f := range fields {
		qb.allowedFilters[f] = true
	}
	return qb
}

// AllowSorts sets allowed sort fields (for security)
func (qb *QueryBuilder) AllowSorts(fields ...string) *QueryBuilder {
	for _, f := range fields {
		qb.allowedSorts[f] = true
	}
	return qb
}

// DefaultSort sets the default sort order when no sort is specified
func (qb *QueryBuilder) DefaultSort(field string, direction SortDirection) *QueryBuilder {
	qb.defaultSort = []SortField{{Field: field, Direction: direction}}
	return qb
}

// SearchColumns sets the columns to search in
func (qb *QueryBuilder) SearchColumns(columns ...string) *QueryBuilder {
	qb.searchColumns = columns
	return qb
}

// ===========================================
// BUILD METHODS
// ===========================================

// Build applies pagination, filtering, and sorting to the query
func (qb *QueryBuilder) Build() *gorm.DB {
	qb.ensureRequest()

	query := qb.db
	query = qb.applyFilters(query)
	query = qb.applySearch(query)
	query = qb.applySorting(query)
	query = qb.applyPagination(query)

	return query
}

// BuildWithCount applies pagination and returns both query and total count
func (qb *QueryBuilder) BuildWithCount(model interface{}) (*gorm.DB, int64) {
	qb.ensureRequest()

	var total int64

	// Count query (without pagination)
	countQuery := qb.db.Model(model)
	countQuery = qb.applyFilters(countQuery)
	countQuery = qb.applySearch(countQuery)
	countQuery.Count(&total)

	// Main query with pagination
	query := qb.Build()

	return query, total
}

func (qb *QueryBuilder) ensureRequest() {
	if qb.request == nil {
		qb.request = NewQueryParams()
	}
}

// ===========================================
// FILTER APPLICATION
// ===========================================

func (qb *QueryBuilder) applyFilters(query *gorm.DB) *gorm.DB {
	for _, filter := range qb.request.Filters {
		if qb.isFilterAllowed(filter.Field) {
			query = ApplyFilter(query, filter)
		}
	}
	return query
}

func (qb *QueryBuilder) isFilterAllowed(field string) bool {
	// If no allowed filters specified, allow all
	if len(qb.allowedFilters) == 0 {
		return true
	}
	return qb.allowedFilters[field]
}

// ===========================================
// SEARCH APPLICATION
// ===========================================

func (qb *QueryBuilder) applySearch(query *gorm.DB) *gorm.DB {
	if qb.request.Search == "" {
		return query
	}

	searchFields := qb.getSearchFields()
	if len(searchFields) == 0 {
		return query
	}

	return ApplySearch(query, qb.request.Search, searchFields...)
}

func (qb *QueryBuilder) getSearchFields() []string {
	// Use request search fields if provided, otherwise use builder's
	if len(qb.request.SearchFields) > 0 {
		return qb.request.SearchFields
	}
	return qb.searchColumns
}

// ===========================================
// SORT APPLICATION
// ===========================================

func (qb *QueryBuilder) applySorting(query *gorm.DB) *gorm.DB {
	sorts := qb.request.Sort
	if len(sorts) == 0 {
		sorts = qb.defaultSort
	}

	for _, sort := range sorts {
		if qb.isSortAllowed(sort.Field) {
			query = ApplySortField(query, sort)
		}
	}

	return query
}

func (qb *QueryBuilder) isSortAllowed(field string) bool {
	// If no allowed sorts specified, allow all
	if len(qb.allowedSorts) == 0 {
		return true
	}
	return qb.allowedSorts[field]
}

// ===========================================
// PAGINATION APPLICATION
// ===========================================

func (qb *QueryBuilder) applyPagination(query *gorm.DB) *gorm.DB {
	limit := qb.getLimit()

	if qb.request.Type == CursorPagination {
		return qb.applyCursorPagination(query, limit)
	}

	return qb.applyOffsetPagination(query, limit)
}

func (qb *QueryBuilder) getLimit() int {
	limit := qb.request.PageSize
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return limit
}

func (qb *QueryBuilder) applyCursorPagination(query *gorm.DB, limit int) *gorm.DB {
	if qb.request.Cursor != "" {
		cursor, err := DecodeCursor(qb.request.Cursor)
		if err == nil {
			query = query.Where("id > ?", cursor.ID)
		}
	}
	return query.Limit(limit)
}

func (qb *QueryBuilder) applyOffsetPagination(query *gorm.DB, limit int) *gorm.DB {
	page := qb.request.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return query.Offset(offset).Limit(limit)
}

// ===========================================
// STANDALONE FILTER/SEARCH/SORT FUNCTIONS
// ===========================================

// ApplyFilter applies a single filter to a GORM query
func ApplyFilter(query *gorm.DB, filter Filter) *gorm.DB {
	switch filter.Operator {
	case OpEqual:
		return query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case OpNotEqual:
		return query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case OpGreaterThan:
		return query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case OpGreaterEqual:
		return query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case OpLessThan:
		return query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case OpLessEqual:
		return query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case OpLike:
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")
	case OpILike:
		return query.Where(fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", filter.Field), "%"+fmt.Sprint(filter.Value)+"%")
	case OpIn:
		return query.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Values)
	case OpNotIn:
		return query.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Values)
	case OpIsNull:
		return query.Where(fmt.Sprintf("%s IS NULL", filter.Field))
	case OpIsNotNull:
		return query.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	case OpBetween:
		if len(filter.Values) >= 2 {
			return query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field), filter.Values[0], filter.Values[1])
		}
	}
	return query
}

// ApplySearch applies search to a GORM query across multiple columns
func ApplySearch(query *gorm.DB, term string, columns ...string) *gorm.DB {
	if term == "" || len(columns) == 0 {
		return query
	}

	searchTerm := "%" + strings.ToLower(term) + "%"
	var conditions []string
	var args []interface{}

	for _, col := range columns {
		conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE ?", col))
		args = append(args, searchTerm)
	}

	return query.Where(strings.Join(conditions, " OR "), args...)
}

// ApplySortField applies a single sort field to a GORM query
func ApplySortField(query *gorm.DB, sort SortField) *gorm.DB {
	direction := "ASC"
	if sort.Direction == SortDesc {
		direction = "DESC"
	}
	return query.Order(fmt.Sprintf("%s %s", sort.Field, direction))
}
