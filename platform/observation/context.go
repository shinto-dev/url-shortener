package observation

import (
	"context"
)

type FieldFn func(ctx context.Context)

func Add(ctx context.Context, fieldFns ...FieldFn) {
	for _, fieldFn := range fieldFns {
		fieldFn(ctx)
	}
}
