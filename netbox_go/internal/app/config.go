package app

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration using Viper
type Config struct {
	Database  DatabaseConfig
	Server    ServerConfig
	Telemetry TelemetryConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL      string `mapstructure:"DATABASE_URL"`
	Host     string `mapstructure:"DATABASE_HOST"`
	Port     int    `mapstructure:"DATABASE_PORT"`
	Name     string `mapstructure:"DATABASE_NAME"`
	User     string `mapstructure:"DATABASE_USER"`
	Password string `mapstructure:"DATABASE_PASSWORD"`
	SSLMode  string `mapstructure:"DATABASE_SSLMODE"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    int    `mapstructure:"PORT"`
	Host    string `mapstructure:"HOST"`
	Timeout int    `mapstructure:"TIMEOUT"`
}

// TelemetryConfig holds telemetry configuration
type TelemetryConfig struct {
	InfluxDB InfluxDBConfig
	Vault    VaultConfig
	GNMI     GNMIConfig
	Ping     PingerConfig
	DNS      DNSConfig
	NetFlow  NetFlowConfig
}

// InfluxDBConfig holds InfluxDB configuration
type InfluxDBConfig struct {
	URL           string `mapstructure:"INFLUXDB_URL"`
	Token         string `mapstructure:"INFLUXDB_TOKEN"`
	Org           string `mapstructure:"INFLUXDB_ORG"`
	Bucket        string `mapstructure:"INFLUXDB_BUCKET"`
	BatchSize     int    `mapstructure:"INFLUXDB_BATCH_SIZE"`
	FlushInterval int    `mapstructure:"INFLUXDB_FLUSH_INTERVAL"`
	Timeout       int    `mapstructure:"INFLUXDB_TIMEOUT"`
}

// VaultConfig holds HashiCorp Vault configuration
type VaultConfig struct {
	Address  string `mapstructure:"VAULT_ADDR"`
	RoleID   string `mapstructure:"VAULT_APPROLE_ROLE_ID"`
	SecretID string `mapstructure:"VAULT_APPROLE_SECRET_ID"`
	Timeout  int    `mapstructure:"VAULT_TIMEOUT"`
	CacheTTL int    `mapstructure:"VAULT_CACHE_TTL"`
}

// GNMIConfig holds gNMI client configuration
type GNMIConfig struct {
	Timeout     int    `mapstructure:"GNMI_TIMEOUT"`
	TLSEnabled  bool   `mapstructure:"GNMI_TLS_ENABLED"`
	TLSCertPath string `mapstructure:"GNMI_TLS_CERT_PATH"`
	TLSKeyPath  string `mapstructure:"GNMI_TLS_KEY_PATH"`
	TLSCAPath   string `mapstructure:"GNMI_TLS_CA_PATH"`
}

// PingerConfig holds ICMP pinger configuration
type PingerConfig struct {
	Timeout    int `mapstructure:"PING_TIMEOUT"`
	Interval   int `mapstructure:"PING_INTERVAL"`
	Count      int `mapstructure:"PING_COUNT"`
	PacketSize int `mapstructure:"PING_PACKET_SIZE"`
	Workers    int `mapstructure:"PING_WORKERS"`
}

// DNSConfig holds DNS resolver configuration
type DNSConfig struct {
	Timeout   int    `mapstructure:"DNS_TIMEOUT"`
	Interval  int    `mapstructure:"DNS_INTERVAL"`
	Workers   int    `mapstructure:"DNS_WORKERS"`
	Servers   string `mapstructure:"DNS_SERVERS"`
	UseCustom bool   `mapstructure:"DNS_USE_CUSTOM"`
}

// NetFlowConfig holds NetFlow collector configuration
type NetFlowConfig struct {
	NetFlowPort int `mapstructure:"NETFLOW_PORT"`
	SFlowPort   int `mapstructure:"SFLOW_PORT"`
	Workers     int `mapstructure:"NETFLOW_WORKERS"`
	BufferSize  int `mapstructure:"NETFLOW_BUFFER_SIZE"`
}

// InitConfig initializes Viper and reads configuration
func InitConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Enable environment variable support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set config file paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/netbox/")
	v.AddConfigPath("$HOME/.netbox")

	// Try to read config file (optional - won't fail if missing)
	_ = v.ReadInConfig()

	// Unmarshal configuration
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for all configuration
func setDefaults(v *viper.Viper) {
	// Database defaults
	v.SetDefault("DATABASE_HOST", "localhost")
	v.SetDefault("DATABASE_PORT", 5432)
	v.SetDefault("DATABASE_NAME", "netbox")
	v.SetDefault("DATABASE_USER", "netbox")
	v.SetDefault("DATABASE_PASSWORD", "netbox")
	v.SetDefault("DATABASE_SSLMODE", "disable")
	v.SetDefault("DATABASE_URL", "postgres://netbox:netbox@localhost:5432/netbox?sslmode=disable")

	// Server defaults
	v.SetDefault("PORT", 8080)
	v.SetDefault("HOST", "")
	v.SetDefault("TIMEOUT", 30)

	// InfluxDB defaults
	v.SetDefault("INFLUXDB_URL", "http://localhost:8086")
	v.SetDefault("INFLUXDB_ORG", "netbox")
	v.SetDefault("INFLUXDB_BUCKET", "telemetry")
	v.SetDefault("INFLUXDB_BATCH_SIZE", 5000)
	v.SetDefault("INFLUXDB_FLUSH_INTERVAL", 10)
	v.SetDefault("INFLUXDB_TIMEOUT", 30)

	// Vault defaults
	v.SetDefault("VAULT_ADDR", "http://localhost:8200")
	v.SetDefault("VAULT_TIMEOUT", 10)
	v.SetDefault("VAULT_CACHE_TTL", 300)

	// GNMI defaults
	v.SetDefault("GNMI_TIMEOUT", 30)
	v.SetDefault("GNMI_TLS_ENABLED", true)

	// Ping defaults
	v.SetDefault("PING_TIMEOUT", 5)
	v.SetDefault("PING_INTERVAL", 30)
	v.SetDefault("PING_COUNT", 5)
	v.SetDefault("PING_PACKET_SIZE", 64)
	v.SetDefault("PING_WORKERS", 50)

	// DNS defaults
	v.SetDefault("DNS_TIMEOUT", 10)
	v.SetDefault("DNS_INTERVAL", 300)
	v.SetDefault("DNS_WORKERS", 20)
	v.SetDefault("DNS_SERVERS", "")
	v.SetDefault("DNS_USE_CUSTOM", false)

	// NetFlow defaults
	v.SetDefault("NETFLOW_PORT", 2055)
	v.SetDefault("SFLOW_PORT", 6343)
	v.SetDefault("NETFLOW_WORKERS", 10)
	v.SetDefault("NETFLOW_BUFFER_SIZE", 10000)
}

// GetConfig returns the global config instance
func GetConfig() *Config {
	cfg, err := InitConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize config: %v", err))
	}
	return cfg
}
