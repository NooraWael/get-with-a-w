package downloader

//controls the speed - rate limiter flag

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"golang.org/x/time/rate"
)

// adjustRateLimit parses and converts user-provided rate limit to bytes/sec
func adjustRateLimit(rateLimit string) (int, error) {
	if rate, err := strconv.Atoi(rateLimit); err == nil {
		return rate, nil
	}

	unit := rateLimit[len(rateLimit)-1]
	multiplier := 1
	switch unit {
	case 'k':
		multiplier = 1_000
	case 'M':
		multiplier = 1_000_000
	case 'G':
		multiplier = 1_000_000_000
	default:
		return 0, fmt.Errorf("invalid rate limit format")
	}

	rate, err := strconv.Atoi(rateLimit[:len(rateLimit)-1])
	if err != nil {
		return 0, err
	}

	return rate * multiplier * 9 / 10, nil // apply 90% overhead factor
}

// Setup rate limiter if provided
var limiter *rate.Limiter

// rateLimitedReader wraps an io.ReadCloser and applies rate limiting
type rateLimitedReader struct {
	io.ReadCloser
	limiter *rate.Limiter
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if err != nil {
		return n, err
	}
	if err := r.limiter.WaitN(context.Background(), n); err != nil {
		return n, err
	}
	return n, nil
}
