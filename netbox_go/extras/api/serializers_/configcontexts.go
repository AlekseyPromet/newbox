// Package serializers содержит DTO и функции сериализации для API
package serializers

import (
	"time"
)

// ConfigContextProfile представляет профиль контекста конфигурации
type ConfigContextProfile struct {
	ID           string      `json:"id"`
	URL          string      `json:"url"`
	DisplayURL   string      `json:"display_url,omitempty"`
	Display      string      `json:"display"`
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Schema       interface{} `json:"schema,omitempty"`
	Tags         []Tag       `json:"tags,omitempty"`
	Owner        interface{} `json:"owner,omitempty"`
	Comments     string      `json:"comments,omitempty"`
	DataSource   *DataSource `json:"data_source,omitempty"`
	DataPath     string      `json:"data_path,omitempty"`
	DataFile     *DataFile   `json:"data_file,omitempty"`
	DataSynced   *time.Time  `json:"data_synced,omitempty"`
	Created      *time.Time  `json:"created,omitempty"`
	LastUpdated  *time.Time  `json:"last_updated,omitempty"`
}

// ConfigContext представляет контекст конфигурации
type ConfigContext struct {
	ID            string        `json:"id"`
	URL           string        `json:"url"`
	DisplayURL    string        `json:"display_url,omitempty"`
	Display       string        `json:"display"`
	Name          string        `json:"name"`
	Weight        int           `json:"weight"`
	Profile       *ConfigContextProfile `json:"profile,omitempty"`
	Description   string        `json:"description,omitempty"`
	IsActive      bool          `json:"is_active"`
	Regions       []Region      `json:"regions,omitempty"`
	SiteGroups    []SiteGroup   `json:"site_groups,omitempty"`
	Sites         []Site        `json:"sites,omitempty"`
	Locations     []Location    `json:"locations,omitempty"`
	DeviceTypes   []DeviceType  `json:"device_types,omitempty"`
	Roles         []DeviceRole  `json:"roles,omitempty"`
	Platforms     []Platform    `json:"platforms,omitempty"`
	ClusterTypes  []ClusterType `json:"cluster_types,omitempty"`
	ClusterGroups []ClusterGroup `json:"cluster_groups,omitempty"`
	Clusters      []Cluster     `json:"clusters,omitempty"`
	TenantGroups  []TenantGroup `json:"tenant_groups,omitempty"`
	Tenants       []Tenant      `json:"tenants,omitempty"`
	Owner         interface{}   `json:"owner,omitempty"`
	Tags          []Tag         `json:"tags,omitempty"`
	DataSource    *DataSource   `json:"data_source,omitempty"`
	DataPath      string        `json:"data_path,omitempty"`
	DataFile      *DataFile     `json:"data_file,omitempty"`
	DataSynced    *time.Time    `json:"data_synced,omitempty"`
	Data          interface{}   `json:"data"`
	Created       *time.Time    `json:"created,omitempty"`
	LastUpdated   *time.Time    `json:"last_updated,omitempty"`
}

// ConfigContextProfileBrief краткое представление профиля контекста
type ConfigContextProfileBrief struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Display     string `json:"display"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ConfigContextBrief краткое представление контекста конфигурации
type ConfigContextBrief struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Display     string `json:"display"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// DataSource представляет источник данных
type DataSource struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// DataFile представляет файл данных
type DataFile struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Tag представляет тег
type Tag struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

// Region представляет регион
type Region struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// SiteGroup представляет группу сайтов
type SiteGroup struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// Site представляет сайт
type Site struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// Location представляет местоположение
type Location struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// DeviceType представляет тип устройства
type DeviceType struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Model       string `json:"model,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// DeviceRole представляет роль устройства
type DeviceRole struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// Platform представляет платформу
type Platform struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// ClusterType представляет тип кластера
type ClusterType struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// ClusterGroup представляет группу кластеров
type ClusterGroup struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// Cluster представляет кластер
type Cluster struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// TenantGroup представляет группу арендаторов
type TenantGroup struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// Tenant представляет арендатора
type Tenant struct {
	ID          string `json:"id"`
	URL         string `json:"url,omitempty"`
	Display     string `json:"display,omitempty"`
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}
