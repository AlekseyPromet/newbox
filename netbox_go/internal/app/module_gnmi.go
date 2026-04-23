package app

import (
	"go.uber.org/fx"
)

// ModuleGNMI provides telemetry collection dependencies
// Configuration is now centralized in Config struct (see config.go)
// This module is kept for fx dependency injection wiring
var ModuleGNMI = fx.Options(
	// Provide constructors
	fx.Provide(
		NewTelemetryConfigFromViper,
	),
)

// NewTelemetryConfigFromViper creates TelemetryConfig using Viper
func NewTelemetryConfigFromViper() *TelemetryConfig {
	cfg := GetConfig()
	return &TelemetryConfig{
		InfluxDB: InfluxDBConfig{
			URL:           cfg.Telemetry.InfluxDB.URL,
			Token:         cfg.Telemetry.InfluxDB.Token,
			Org:           cfg.Telemetry.InfluxDB.Org,
			Bucket:        cfg.Telemetry.InfluxDB.Bucket,
			BatchSize:     cfg.Telemetry.InfluxDB.BatchSize,
			FlushInterval: cfg.Telemetry.InfluxDB.FlushInterval,
			Timeout:       cfg.Telemetry.InfluxDB.Timeout,
		},
		Vault: VaultConfig{
			Address:  cfg.Telemetry.Vault.Address,
			RoleID:   cfg.Telemetry.Vault.RoleID,
			SecretID: cfg.Telemetry.Vault.SecretID,
			Timeout:  cfg.Telemetry.Vault.Timeout,
			CacheTTL: cfg.Telemetry.Vault.CacheTTL,
		},
		GNMI: GNMIConfig{
			Timeout:     cfg.Telemetry.GNMI.Timeout,
			TLSEnabled:  cfg.Telemetry.GNMI.TLSEnabled,
			TLSCertPath: cfg.Telemetry.GNMI.TLSCertPath,
			TLSKeyPath:  cfg.Telemetry.GNMI.TLSKeyPath,
			TLSCAPath:   cfg.Telemetry.GNMI.TLSCAPath,
		},
		Ping: PingerConfig{
			Timeout:    cfg.Telemetry.Ping.Timeout,
			Interval:   cfg.Telemetry.Ping.Interval,
			Count:      cfg.Telemetry.Ping.Count,
			PacketSize: cfg.Telemetry.Ping.PacketSize,
			Workers:    cfg.Telemetry.Ping.Workers,
		},
		DNS: DNSConfig{
			Timeout:   cfg.Telemetry.DNS.Timeout,
			Interval:  cfg.Telemetry.DNS.Interval,
			Workers:   cfg.Telemetry.DNS.Workers,
			Servers:   cfg.Telemetry.DNS.Servers,
			UseCustom: cfg.Telemetry.DNS.UseCustom,
		},
		NetFlow: NetFlowConfig{
			NetFlowPort: cfg.Telemetry.NetFlow.NetFlowPort,
			SFlowPort:   cfg.Telemetry.NetFlow.SFlowPort,
			Workers:     cfg.Telemetry.NetFlow.Workers,
			BufferSize:  cfg.Telemetry.NetFlow.BufferSize,
		},
	}
}
