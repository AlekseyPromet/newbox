// Package entity содержит сущности домена Virtualization
package entity

import (
	"time"

	dcim_entity "github.com/AlekseyPromet/netbox_go/internal/domain/dcim/entity"
	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/enum"
	ipam_entity "github.com/AlekseyPromet/netbox_go/internal/domain/ipam/entity"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// ClusterType представляет тип кластера виртуализации
type ClusterType struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность типа кластера
func (ct *ClusterType) Validate() error {
	if ct.Name == "" {
		return types.ErrNameRequired
	}
	if err := ct.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// ClusterGroup представляет группу кластеров
type ClusterGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Slug        types.Slug      `json:"slug"`
	Description types.Description `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы кластеров
func (cg *ClusterGroup) Validate() error {
	if cg.Name == "" {
		return types.ErrNameRequired
	}
	if err := cg.Slug.Validate(); err != nil {
		return err
	}
	return nil
}

// Cluster представляет кластер виртуализации
type Cluster struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	ClusterType *ClusterType    `json:"cluster_type"`
	GroupID     *types.ID       `json:"group_id,omitempty"`
	SiteID      *types.ID       `json:"site_id,omitempty"`
	TenantID    *types.ID       `json:"tenant_id,omitempty"`
	Status      enum.ClusterStatus `json:"status"`
	Description types.Description `json:"description,omitempty"`
	Comments    types.Comments  `json:"comments,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность кластера
func (c *Cluster) Validate() error {
	if c.Name == "" {
		return types.ErrNameRequired
	}
	if c.ClusterType == nil {
		return types.ErrValidationFailed
	}
	if err := c.Status.Validate(); err != nil {
		return err
	}
	return nil
}

// GetStatusColor возвращает цвет статуса кластера
func (c *Cluster) GetStatusColor() string {
	return c.Status.Color()
}

// ClusterDevice представляет устройство в кластере
type ClusterDevice struct {
	ID        types.ID `json:"id"`
	ClusterID types.ID `json:"cluster_id"`
	DeviceID  types.ID `json:"device_id"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

// Validate проверяет корректность устройства в кластере
func (cd *ClusterDevice) Validate() error {
	if cd.ClusterID.String() == "" {
		return types.ErrValidationFailed
	}
	if cd.DeviceID.String() == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// VirtualMachine представляет виртуальную машину
type VirtualMachine struct {
	ID             types.ID              `json:"id"`
	Name           string                `json:"name"`
	ClusterID      *types.ID             `json:"cluster_id,omitempty"`
	DeviceID       *types.ID             `json:"device_id,omitempty"` // For VMs hosted on specific device
	SiteID         *types.ID             `json:"site_id,omitempty"`
	TenantID       *types.ID             `json:"tenant_id,omitempty"`
	Platform       *dcim_entity.Platform      `json:"platform,omitempty"`
	RoleID         *types.ID             `json:"role_id,omitempty"`
	PrimaryIPv4    *ipam_entity.IPAddress     `json:"primary_ipv4,omitempty"`
	PrimaryIPv6    *ipam_entity.IPAddress     `json:"primary_ipv6,omitempty"`
	VCPUs          *int32                `json:"vcpus,omitempty"`
	Memory         *int32                `json:"memory,omitempty"` // MB
	Disk           *int32                `json:"disk,omitempty"`   // MB
	Status         enum.VirtualMachineStatus `json:"status"`
	Airflow        *enum.AirflowDirection `json:"airflow,omitempty"`
	ConfigTemplate *dcim_entity.ConfigTemplate `json:"config_template,omitempty"`
	Comments       types.Comments        `json:"comments,omitempty"`
	LocalContextData interface{}          `json:"local_context_data,omitempty"`
	Created        time.Time             `json:"created"`
	Updated        time.Time             `json:"updated"`
}

// Validate проверяет корректность виртуальной машины
func (vm *VirtualMachine) Validate() error {
	if vm.Name == "" {
		return types.ErrNameRequired
	}
	if vm.ClusterID == nil && vm.DeviceID == nil {
		return types.ErrValidationFailed
	}
	if err := vm.Status.Validate(); err != nil {
		return err
	}
	if vm.VCPUs != nil && *vm.VCPUs < 1 {
		return types.ErrValidationFailed
	}
	if vm.Memory != nil && *vm.Memory < 1 {
		return types.ErrValidationFailed
	}
	if vm.Disk != nil && *vm.Disk < 1 {
		return types.ErrValidationFailed
	}
	return nil
}

// GetStatusColor возвращает цвет статуса VM
func (vm *VirtualMachine) GetStatusColor() string {
	return vm.Status.Color()
}

// VMInterface представляет интерфейс виртуальной машины
type VMInterface struct {
	ID           types.ID        `json:"id"`
	VirtualMachineID types.ID    `json:"virtual_machine_id"`
	Name         string          `json:"name"`
	Label        string          `json:"label,omitempty"`
	Type         enum.InterfaceType `json:"type"`
	Enabled      bool            `json:"enabled"`
	MACAddress   string          `json:"mac_address,omitempty"`
	VLANMode     *enum.InterfaceMode `json:"vlan_mode,omitempty"`
	UntaggedVLAN *ipam_entity.VLAN    `json:"untagged_vlan,omitempty"`
	TaggedVLANs  []*ipam_entity.VLAN  `json:"tagged_vlans,omitempty"`
	Description  types.Description `json:"description,omitempty"`
	Comments     types.Comments  `json:"comments,omitempty"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
}

// Validate проверяет корректность интерфейса VM
func (vif *VMInterface) Validate() error {
	if vif.Name == "" {
		return types.ErrNameRequired
	}
	if vif.VirtualMachineID.String() == "" {
		return types.ErrValidationFailed
	}
	if err := vif.Type.Validate(); err != nil {
		return err
	}
	return nil
}

// VMDisk представляет диск виртуальной машины
type VMDisk struct {
	ID               types.ID `json:"id"`
	VirtualMachineID types.ID `json:"virtual_machine_id"`
	Name             string   `json:"name"`
	Size             int64    `json:"size"` // bytes
	DiskType         string   `json:"disk_type,omitempty"` // ssd, hdd, nvme
	BootOrder        *int     `json:"boot_order,omitempty"`
	Description      types.Description `json:"description,omitempty"`
	Created          time.Time `json:"created"`
	Updated          time.Time `json:"updated"`
}

// Validate проверяет корректность диска VM
func (vd *VMDisk) Validate() error {
	if vd.Name == "" {
		return types.ErrNameRequired
	}
	if vd.VirtualMachineID.String() == "" {
		return types.ErrValidationFailed
	}
	if vd.Size <= 0 {
		return types.ErrValidationFailed
	}
	return nil
}
