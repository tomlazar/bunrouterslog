package bunrouterslog

import (
	"context"

	"go.jetpack.io/typeid"
)

type requestIdKey struct{}

type requestIdPrefix struct{}

func (requestIdPrefix) Prefix() string { return "request" }

type RequestID struct {
	typeid.TypeID[requestIdPrefix]
}

func NewRequestID() (RequestID, error) {
	return typeid.New[RequestID]()
}

func MustNewRequestID() RequestID {
	return typeid.Must(NewRequestID())
}

func ContextWithRequestID(ctx context.Context, id RequestID) context.Context {
	return context.WithValue(ctx, requestIdKey{}, id)
}

func RequestIDFromContext(ctx context.Context) (RequestID, bool) {
	id, ok := ctx.Value(requestIdKey{}).(RequestID)
	return id, ok
}
