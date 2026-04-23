package telemetry

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	CircuitStateClosed CircuitState = iota
	CircuitStateOpen
	CircuitStateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitStateClosed:
		return "closed"
	case CircuitStateOpen:
		return "open"
	case CircuitStateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

var (
	// ErrCircuitOpen is returned when the circuit breaker is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrMaxRetriesExceeded is returned when max retries have been exceeded
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name             string
	failureThreshold int32
	successThreshold int32
	timeout          time.Duration
	openState        atomic.Value // stores time.Time when circuit was opened
	mu               sync.RWMutex
	failureCount     int32
	successCount     int32
	state            CircuitState
	lastFailure      time.Time
	logger           *zap.Logger
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Name             string
	FailureThreshold int32         // failures before opening circuit (default: 5)
	SuccessThreshold int32         // successes in half-open before closing (default: 3)
	Timeout          time.Duration // time before trying half-open (default: 30s)
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig(name string) *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		Name:             name,
		FailureThreshold: 5,
		SuccessThreshold: 3,
		Timeout:          30 * time.Second,
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg *CircuitBreakerConfig, logger *zap.Logger) *CircuitBreaker {
	if cfg == nil {
		cfg = DefaultCircuitBreakerConfig("default")
	}

	cb := &CircuitBreaker{
		name:             cfg.Name,
		failureThreshold: cfg.FailureThreshold,
		successThreshold: cfg.SuccessThreshold,
		timeout:          cfg.Timeout,
		state:            CircuitStateClosed,
		logger:           logger,
	}

	cb.openState.Store(time.Time{})
	return cb
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if !cb.canExecute() {
		return ErrCircuitOpen
	}

	err := fn()
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// canExecute checks if the circuit allows execution
func (cb *CircuitBreaker) canExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitStateClosed:
		return true
	case CircuitStateOpen:
		// Check if timeout has passed
		openTime := cb.openState.Load().(time.Time)
		if time.Since(openTime) > cb.timeout {
			cb.toHalfOpen()
			return true
		}
		return false
	case CircuitStateHalfOpen:
		return true
	default:
		return false
	}
}

// recordFailure records a failure and potentially opens the circuit
func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailure = time.Now()

	if cb.state == CircuitStateHalfOpen {
		// Any failure in half-open opens the circuit
		cb.toOpen()
		return
	}

	if cb.failureCount >= cb.failureThreshold {
		cb.toOpen()
	}
}

// recordSuccess records a success and potentially closes the circuit
func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount = 0

	if cb.state == CircuitStateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.toClosed()
		}
	}
}

// toOpen transitions to open state
func (cb *CircuitBreaker) toOpen() {
	cb.state = CircuitStateOpen
	cb.openState.Store(time.Now())
	cb.failureCount = 0
	cb.successCount = 0

	if cb.logger != nil {
		cb.logger.Warn("circuit breaker opened",
			zap.String("name", cb.name),
			zap.Duration("timeout", cb.timeout),
		)
	}
}

// toHalfOpen transitions to half-open state
func (cb *CircuitBreaker) toHalfOpen() {
	cb.state = CircuitStateHalfOpen
	cb.successCount = 0

	if cb.logger != nil {
		cb.logger.Info("circuit breaker half-open",
			zap.String("name", cb.name),
		)
	}
}

// toClosed transitions to closed state
func (cb *CircuitBreaker) toClosed() {
	cb.state = CircuitStateClosed
	cb.failureCount = 0
	cb.successCount = 0

	if cb.logger != nil {
		cb.logger.Info("circuit breaker closed",
			zap.String("name", cb.name),
		)
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Stats returns circuit breaker statistics
func (cb *CircuitBreaker) Stats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"name":          cb.name,
		"state":         cb.state.String(),
		"failure_count": atomic.LoadInt32(&cb.failureCount),
		"success_count": cb.successCount,
		"last_failure":  cb.lastFailure,
	}
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	Jitter          bool
	JitterFactor    float64
	RetryableErrors func(error) bool
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		JitterFactor: 0.1,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// Retry executes a function with retry logic
func Retry(ctx context.Context, cfg *RetryConfig, fn RetryableFunc) error {
	if cfg == nil {
		cfg = DefaultRetryConfig()
	}

	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}

			// Calculate next delay
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}

			// Add jitter
			if cfg.Jitter {
				jitter := time.Duration(float64(delay) * cfg.JitterFactor * (2*float64(attempt%2) - 1))
				delay = delay + jitter
			}
		}

		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}

		// Check if error is retryable
		if cfg.RetryableErrors != nil && !cfg.RetryableErrors(lastErr) {
			return lastErr
		}
	}

	return errors.Join(ErrMaxRetriesExceeded, lastErr)
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for temporary/network errors
	var temporary interface {
		Temporary() bool
	}
	if errors.As(err, &temporary) {
		return temporary.Temporary()
	}

	// Check for context errors (not retryable)
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		return false
	}

	// Default: retryable
	return true
}

// Backoff calculates exponential backoff delay
func Backoff(attempt int, initialDelay, maxDelay time.Duration, multiplier float64) time.Duration {
	delay := time.Duration(float64(initialDelay) * pow(multiplier, float64(attempt)))
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}
