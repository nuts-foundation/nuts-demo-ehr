package sql

import "context"

// SQLContextGetter is an interface provided both by transaction and standard db connection
type SQLContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
