package utils

import (
	"context"
	"errors"
	"fmt"
)

var ErrCancelled = errors.New("cancelled")

func WasCancelled(ctx context.Context, done chan struct{}) bool {
	select {
	case <-ctx.Done():
	case <-done:
		fmt.Println("* cancelled *")
		return true
	default:
	}
	return false
}
