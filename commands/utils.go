package commands

import (
	"context"
	"errors"
	"fmt"
)

var errCancelled = errors.New("cancelled")

func wasCancelled(ctx context.Context, done chan struct{}) bool {
	select {
	case <-ctx.Done():
	case <-done:
		fmt.Println("* cancelled *")
		return true
	default:
	}
	return false
}
