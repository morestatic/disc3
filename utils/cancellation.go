package utils

import (
	"context"
	"errors"
	"fmt"
)

var ErrCancelled = errors.New("cancelled")

func WasCancelled(ctx context.Context, interrupt chan struct{}, hasErr chan error, done chan struct{}) bool {
	select {
	case <-ctx.Done():
	case <-interrupt:
		fmt.Println("* aborted *")
		return true
	case err := <-hasErr:
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		close(done)
		return true
	case <-done:
		return true
	default:
	}
	return false
}
