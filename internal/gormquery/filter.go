package gormquery

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// Filter represents a single WHERE clause condition
type Filter struct {
	Field    string      // field name
	Operator string      // operator: =, !=, <, >, ILIKE, IN, etc.
	Value    interface{} // value to compare
}

// FilterGroup is a group of conditions combined with AND logic
type FilterGroup []Filter

// OrderOption specifies ordering for a field
type OrderOption struct {
	Field     string // field name
	Direction string // ASC or DESC
}

// QueryOptions represents optional modifiers for a query (ordering, limit, offset)
type QueryOptions struct {
	OrderBy []OrderOption
	Limit   *int
	Offset  *int
}

// NewFilter creates a new Filter condition
func NewFilter(field string, operator string, value interface{}) Filter {
	return Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// NewFilterGroup creates a FilterGroup from multiple Filters
func NewFilterGroup(filters ...Filter) FilterGroup {
	return filters
}

// ApplyFilters applies multiple filter groups to a GORM query
// Each group is combined with AND inside, and all groups are combined with OR
func ApplyFilters(query *gorm.DB, groups []FilterGroup) *gorm.DB {
	for i, group := range groups {
		var conditions []string
		var values []interface{}

		for _, filter := range group {
			conditions = append(conditions, fmt.Sprintf("%s %s ?", filter.Field, filter.Operator))
			values = append(values, filter.Value)
		}

		clause := strings.Join(conditions, " AND ")

		if i == 0 {
			query = query.Where(clause, values...)
		} else {
			query = query.Or(clause, values...)
		}
	}
	return query
}

// ApplyQueryOptions applies ordering, limit, and offset to the query
func ApplyQueryOptions(query *gorm.DB, options QueryOptions) *gorm.DB {
	for _, order := range options.OrderBy {
		direction := "ASC"
		if strings.ToUpper(order.Direction) == "DESC" {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", order.Field, direction))
	}
	if options.Limit != nil {
		query = query.Limit(*options.Limit)
	}
	if options.Offset != nil {
		query = query.Offset(*options.Offset)
	}
	return query
}

// IntPtr is a helper to get pointer from int
func IntPtr(i int) *int {
	return &i
}
