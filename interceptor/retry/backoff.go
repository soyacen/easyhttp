package easyhttpretry

import (
	"math"
	"math/rand"
	"time"
)

// BackoffFunc denotes a family of functions that control the backoff duration between call retries.
//
// They are called with an identifier of the attempt, and should return a time the system client should
// hold off for. If the time returned is longer than the `context.Context.Deadline` of the request
// the deadline of the request takes precedence and the wait will be interrupted before proceeding
// with the next iteration.
type BackoffFunc func(attempt uint) time.Duration

// BackoffLinear is very simple: it waits for a fixed period of time between calls.
func BackoffLinear(waitBetween time.Duration) BackoffFunc {
	return func(attempt uint) time.Duration {
		return waitBetween
	}
}

// BackoffLinearWithJitter waits a set period of time, allowing for jitter (fractional adjustment).
//
// For example waitBetween=1s and jitter=0.10 can generate waits between 900ms and 1100ms.
func BackoffLinearWithJitter(waitBetween time.Duration, jitterFraction float64) BackoffFunc {
	return func(attempt uint) time.Duration {
		return JitterUp(waitBetween, jitterFraction)
	}
}

// BackoffExponential produces increasing intervals for each attempt.
//
// The scalar is multiplied times 2 raised to the current attempt. So the first
// retry with a scalar of 100ms is 100ms, while the 5th attempt would be 1.6s.
func BackoffExponential(scalar time.Duration) BackoffFunc {
	return func(attempt uint) time.Duration {
		return exponential(attempt, scalar)
	}
}

func exponential(attempt uint, scalar time.Duration) time.Duration {
	return scalar * time.Duration(math.Exp2(float64(attempt)))
}

// BackoffExponentialWithJitter creates an exponential backoff like
// BackoffExponential does, but adds jitter.
func BackoffExponentialWithJitter(scalar time.Duration, jitterFraction float64) BackoffFunc {
	return func(attempt uint) time.Duration {
		return JitterUp(exponential(attempt, scalar), jitterFraction)
	}
}

// JitterUp adds random jitter to the duration.
//
// This adds or subtracts time from the duration within a given jitter fraction.
// For example for 10s and jitter 0.1, it will return a time within [9s, 11s])
func JitterUp(duration time.Duration, jitter float64) time.Duration {
	multiplier := jitter * (rand.Float64()*2 - 1)
	return time.Duration(float64(duration) * (1 + multiplier))
}
