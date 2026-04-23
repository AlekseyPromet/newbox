package app

import (
	"go.uber.org/fx"

	"netbox_go/internal/domain/telemetry/service"
	"netbox_go/internal/infrastructure/telemetry"
)

// ModuleTelemetry provides telemetry services and infrastructure
var ModuleTelemetry = fx.Options(
	// Provide telemetry service
	fx.Provide(
		service.NewTelemetryService,
	),

	// Provide coordinator
	fx.Provide(
		telemetry.NewCoordinator,
	),

	// Provide collector
	fx.Provide(
		telemetry.NewCollector,
	),

	// Provide circuit breaker
	fx.Provide(
		telemetry.NewCircuitBreaker,
	),

	// Provide metrics
	fx.Provide(
		telemetry.NewTelemetryMetrics,
	),
)

// ModuleTelemetryAPI provides telemetry HTTP handlers
var ModuleTelemetryAPI = fx.Options(
	// Provide handlers
	fx.Provide(
		NewTelemetryHandler,
	),
)

// TelemetryHandlerDeps holds dependencies for telemetry handlers
type TelemetryHandlerDeps struct {
	fx.In
	TelemetryService *service.TelemetryService
}

// NewTelemetryHandler creates a new telemetry handler
func NewTelemetryHandler(deps TelemetryHandlerDeps) *TelemetryHandlerWrapper {
	return &TelemetryHandlerWrapper{
		service: deps.TelemetryService,
	}
}

// TelemetryHandlerWrapper wraps the handler with logging
type TelemetryHandlerWrapper struct {
	service *service.TelemetryService
}
