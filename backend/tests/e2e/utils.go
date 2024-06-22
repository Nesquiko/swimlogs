package e2e

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// waitForReady calls the specified endpoint until it gets a 200
// response or until the context is cancelled or the timeout is
// reached.
func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	interval time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			nestedErr := errors.Unwrap(errors.Unwrap(err))
			if nestedErr.Error() == "connect: connection refused" {
				time.Sleep(interval)
				continue
			}
			return fmt.Errorf("error while waiting for ready: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			slog.Info("server is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(interval)
		}
	}
}

// compareTimes compares two time.Time values, ignoring seconds and nanoseconds.
func compareTimes(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	h1, min1, _ := t1.Clock()
	h2, min2, _ := t2.Clock()

	return y1 == y2 && m1 == m2 && d1 == d2 && h1 == h2 && min1 == min2
}

func asPtr[T any](v T) *T {
	return &v
}
