package domain

import "context"

type Delivery interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
