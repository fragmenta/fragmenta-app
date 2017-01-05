package users

import (
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-app/src/lib/status"
)

// Query returns a new query for users with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for users with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all users with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
