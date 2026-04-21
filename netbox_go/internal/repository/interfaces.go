// Package repository содержит интерфейсы репозиториев для всех доменов
package repository

import (
	"context"
	"time"

	account_entity "github.com/AlekseyPromet/netbox_go/internal/domain/account/entity"
	core_entity "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
	dcim_entity "github.com/AlekseyPromet/netbox_go/internal/domain/dcim/entity"
	extras_entity "github.com/AlekseyPromet/netbox_go/internal/domain/extras/entity"
	users_entity "github.com/AlekseyPromet/netbox_go/internal/domain/users/entity"
	"github.com/AlekseyPromet/netbox_go/pkg/types"

	circuits_entity "github.com/AlekseyPromet/netbox_go/internal/domain/circuits/entity"
)

// SiteRepository определяет интерфейс для работы с сайтами
type SiteRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Site, error)
	GetBySlug(ctx context.Context, slug string) (*dcim_entity.Site, error)
	List(ctx context.Context, filter SiteFilter) ([]*dcim_entity.Site, int64, error)
	Create(ctx context.Context, site *dcim_entity.Site) error
	Update(ctx context.Context, site *dcim_entity.Site) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// SiteFilter представляет фильтры для поиска сайтов
type SiteFilter struct {
	Status   *string
	RegionID *string
	GroupID  *string
	TenantID *string
	Limit    int
	Offset   int
}

// RackRepository определяет интерфейс для работы со стойками
type RackRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Rack, error)
	List(ctx context.Context, filter RackFilter) ([]*dcim_entity.Rack, int64, error)
	Create(ctx context.Context, rack *dcim_entity.Rack) error
	Update(ctx context.Context, rack *dcim_entity.Rack) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// RackFilter представляет фильтры для поиска стоек
type RackFilter struct {
	SiteID     *string
	LocationID *string
	Status     *string
	TypeID     *string
	Limit      int
	Offset     int
}

// DeviceRepository определяет интерфейс для работы с устройствами
type DeviceRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Device, error)
	GetByName(ctx context.Context, name string) (*dcim_entity.Device, error)
	List(ctx context.Context, filter DeviceFilter) ([]*dcim_entity.Device, int64, error)
	Create(ctx context.Context, device *dcim_entity.Device) error
	Update(ctx context.Context, device *dcim_entity.Device) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// DeviceFilter представляет фильтры для поиска устройств
type DeviceFilter struct {
	SiteID     *string
	RackID     *string
	TypeID     *string
	RoleID     *string
	TenantID   *string
	Status     *string
	PlatformID *string
	Limit      int
	Offset     int
}

// UserTokenRepository определяет интерфейс управления API-токенами пользователя
// Операции ограничены пользователем (scoped)
type UserTokenRepository interface {
	ListByUser(ctx context.Context, userID types.ID) ([]*account_entity.UserToken, error)
	Get(ctx context.Context, id types.ID, userID types.ID) (*account_entity.UserToken, error)
	Create(ctx context.Context, token *account_entity.UserToken) error
	Update(ctx context.Context, token *account_entity.UserToken) error
	Delete(ctx context.Context, id types.ID, userID types.ID) error
}

// BookmarkRepository определяет интерфейс для закладок пользователя
// В упрощённой модели используем только закладки текущего пользователя
type BookmarkRepository interface {
	ListBookmarksByUser(ctx context.Context, userID types.ID) ([]*account_entity.Bookmark, error)
}

// NotificationRepository определяет интерфейс для уведомлений пользователя
type NotificationRepository interface {
	ListNotificationsByUser(ctx context.Context, userID types.ID) ([]*account_entity.Notification, error)
}

// SubscriptionRepository определяет интерфейс для подписок пользователя
type SubscriptionRepository interface {
	ListSubscriptionsByUser(ctx context.Context, userID types.ID) ([]*account_entity.Subscription, error)
}

// UserConfigRepository определяет интерфейс для пользовательских настроек
type UserConfigRepository interface {
	GetByUser(ctx context.Context, userID types.ID) (*account_entity.UserConfig, error)
	Upsert(ctx context.Context, config *account_entity.UserConfig) error
}

// CableRepository определяет интерфейс для работы с кабелями
type CableRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.Cable, error)
	List(ctx context.Context, filter CableFilter) ([]*dcim_entity.Cable, int64, error)
	Create(ctx context.Context, cable *dcim_entity.Cable) error
	Update(ctx context.Context, cable *dcim_entity.Cable) error
	Delete(ctx context.Context, id string) error
	GetTerminations(ctx context.Context, terminationType string, terminationID string) ([]*dcim_entity.Cable, error)
}

// CableFilter представляет фильтры для поиска кабелей
type CableFilter struct {
	Status *string
	Type   *string
	Limit  int
	Offset int
}

// PowerPanelRepository определяет интерфейс для работы с панелями питания
type PowerPanelRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.PowerPanel, error)
	List(ctx context.Context, filter PowerPanelFilter) ([]*dcim_entity.PowerPanel, int64, error)
	Create(ctx context.Context, panel *dcim_entity.PowerPanel) error
	Update(ctx context.Context, panel *dcim_entity.PowerPanel) error
	Delete(ctx context.Context, id string) error
}

// PowerPanelFilter представляет фильтры для поиска панелей питания
type PowerPanelFilter struct {
	SiteID     *string
	LocationID *string
	Limit      int
	Offset     int
}

// PowerFeedRepository определяет интерфейс для работы с фидерами питания
type PowerFeedRepository interface {
	GetByID(ctx context.Context, id string) (*dcim_entity.PowerFeed, error)
	List(ctx context.Context, filter PowerFeedFilter) ([]*dcim_entity.PowerFeed, int64, error)
	Create(ctx context.Context, feed *dcim_entity.PowerFeed) error
	Update(ctx context.Context, feed *dcim_entity.PowerFeed) error
	Delete(ctx context.Context, id string) error
}

// PowerFeedFilter представляет фильтры для поиска фидеров питания
type PowerFeedFilter struct {
	PowerPanelID *string
	RackID       *string
	Status       *string
	Limit        int
	Offset       int
}

// DashboardRepository определяет интерфейс для работы с панелями управления
type DashboardRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*extras_entity.Dashboard, error)
	Create(ctx context.Context, dashboard *extras_entity.Dashboard) error
	Update(ctx context.Context, dashboard *extras_entity.Dashboard) error
	Delete(ctx context.Context, id int64) error
}

// ProviderFilter задает параметры фильтрации провайдеров
type ProviderFilter struct {
	Name   *string
	Slug   *string
	Tenant *string
	Limit  int
	Offset int
}

// ProviderRepository операции над провайдерами
type ProviderRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.Provider, error)
	List(ctx context.Context, filter ProviderFilter) ([]*circuits_entity.Provider, int64, error)
	Create(ctx context.Context, provider *circuits_entity.Provider) error
	Update(ctx context.Context, provider *circuits_entity.Provider) error
	Delete(ctx context.Context, id string) error
}

// ProviderAccountFilter фильтр аккаунтов провайдера
type ProviderAccountFilter struct {
	ProviderID *string
	Account    *string
	Limit      int
	Offset     int
}

// ProviderAccountRepository операции над аккаунтами провайдеров
type ProviderAccountRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.ProviderAccount, error)
	List(ctx context.Context, filter ProviderAccountFilter) ([]*circuits_entity.ProviderAccount, int64, error)
	Create(ctx context.Context, account *circuits_entity.ProviderAccount) error
	Update(ctx context.Context, account *circuits_entity.ProviderAccount) error
	Delete(ctx context.Context, id string) error
}

// ProviderNetworkFilter фильтр сетей провайдера
type ProviderNetworkFilter struct {
	ProviderID *string
	Name       *string
	Limit      int
	Offset     int
}

// ProviderNetworkRepository операции над сетями провайдера
type ProviderNetworkRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.ProviderNetwork, error)
	List(ctx context.Context, filter ProviderNetworkFilter) ([]*circuits_entity.ProviderNetwork, int64, error)
	Create(ctx context.Context, network *circuits_entity.ProviderNetwork) error
	Update(ctx context.Context, network *circuits_entity.ProviderNetwork) error
	Delete(ctx context.Context, id string) error
}

// CircuitTypeFilter фильтр типов цепей
type CircuitTypeFilter struct {
	Slug   *string
	Name   *string
	Limit  int
	Offset int
}

// CircuitTypeRepository операции над типами цепей
type CircuitTypeRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.CircuitType, error)
	List(ctx context.Context, filter CircuitTypeFilter) ([]*circuits_entity.CircuitType, int64, error)
	Create(ctx context.Context, ct *circuits_entity.CircuitType) error
	Update(ctx context.Context, ct *circuits_entity.CircuitType) error
	Delete(ctx context.Context, id string) error
}

// CircuitFilter задает параметры фильтрации цепей
type CircuitFilter struct {
	ProviderID *string
	TypeID     *string
	Status     *string
	TenantID   *string
	Limit      int
	Offset     int
}

// CircuitRepository определяет операции над цепями
type CircuitRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.Circuit, error)
	List(ctx context.Context, filter CircuitFilter) ([]*circuits_entity.Circuit, int64, error)
	Create(ctx context.Context, circuit *circuits_entity.Circuit) error
	Update(ctx context.Context, circuit *circuits_entity.Circuit) error
	Delete(ctx context.Context, id string) error
}

// CircuitTerminationRepository определяет операции над точками завершения цепей
type CircuitTerminationRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.CircuitTermination, error)
	ListByCircuit(ctx context.Context, circuitID string) ([]*circuits_entity.CircuitTermination, error)
	Create(ctx context.Context, termination *circuits_entity.CircuitTermination) error
	Update(ctx context.Context, termination *circuits_entity.CircuitTermination) error
	Delete(ctx context.Context, id string) error
}

// CircuitGroupFilter фильтр групп цепей
type CircuitGroupFilter struct {
	TenantID *string
	Name     *string
	Limit    int
	Offset   int
}

// CircuitGroupRepository операции над группами цепей
type CircuitGroupRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.CircuitGroup, error)
	List(ctx context.Context, filter CircuitGroupFilter) ([]*circuits_entity.CircuitGroup, int64, error)
	Create(ctx context.Context, group *circuits_entity.CircuitGroup) error
	Update(ctx context.Context, group *circuits_entity.CircuitGroup) error
	Delete(ctx context.Context, id string) error
}

// CircuitGroupAssignmentFilter фильтр назначений в группы
type CircuitGroupAssignmentFilter struct {
	GroupID    *string
	MemberType *string
	MemberID   *string
	Priority   *string
	Limit      int
	Offset     int
}

// CircuitGroupAssignmentRepository операции над назначениями групп
type CircuitGroupAssignmentRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.CircuitGroupAssignment, error)
	List(ctx context.Context, filter CircuitGroupAssignmentFilter) ([]*circuits_entity.CircuitGroupAssignment, int64, error)
	Create(ctx context.Context, assignment *circuits_entity.CircuitGroupAssignment) error
	Update(ctx context.Context, assignment *circuits_entity.CircuitGroupAssignment) error
	Delete(ctx context.Context, id string) error
}

// VirtualCircuitTypeFilter фильтр типов виртуальных цепей
type VirtualCircuitTypeFilter struct {
	Slug   *string
	Name   *string
	Limit  int
	Offset int
}

// VirtualCircuitTypeRepository операции над типами виртуальных цепей
type VirtualCircuitTypeRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.VirtualCircuitType, error)
	List(ctx context.Context, filter VirtualCircuitTypeFilter) ([]*circuits_entity.VirtualCircuitType, int64, error)
	Create(ctx context.Context, vct *circuits_entity.VirtualCircuitType) error
	Update(ctx context.Context, vct *circuits_entity.VirtualCircuitType) error
	Delete(ctx context.Context, id string) error
}

// VirtualCircuitFilter фильтр виртуальных цепей
type VirtualCircuitFilter struct {
	ProviderNetworkID *string
	TypeID            *string
	Status            *string
	TenantID          *string
	Limit             int
	Offset            int
}

// VirtualCircuitRepository операции над виртуальными цепями
type VirtualCircuitRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.VirtualCircuit, error)
	List(ctx context.Context, filter VirtualCircuitFilter) ([]*circuits_entity.VirtualCircuit, int64, error)
	Create(ctx context.Context, vc *circuits_entity.VirtualCircuit) error
	Update(ctx context.Context, vc *circuits_entity.VirtualCircuit) error
	Delete(ctx context.Context, id string) error
}

// VirtualCircuitTerminationRepository операции над точками завершения виртуальных цепей
type VirtualCircuitTerminationRepository interface {
	GetByID(ctx context.Context, id string) (*circuits_entity.VirtualCircuitTermination, error)
	ListByVirtualCircuit(ctx context.Context, virtualCircuitID string) ([]*circuits_entity.VirtualCircuitTermination, error)
	Create(ctx context.Context, termination *circuits_entity.VirtualCircuitTermination) error
	Update(ctx context.Context, termination *circuits_entity.VirtualCircuitTermination) error
	Delete(ctx context.Context, id string) error
}

// ObjectTypeFilter фильтр типов объектов
type ObjectTypeFilter struct {
	AppLabel *string
	Model    *string
	Public   *bool
	Feature  *string
	Limit    int
	Offset   int
}

// ObjectTypeRepository операции над типами объектов
type ObjectTypeRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.ObjectType, error)
	List(ctx context.Context, filter ObjectTypeFilter) ([]*core_entity.ObjectType, int64, error)
	Create(ctx context.Context, ot *core_entity.ObjectType) error
	Update(ctx context.Context, ot *core_entity.ObjectType) error
	Delete(ctx context.Context, id string) error
	GetByAppAndModel(ctx context.Context, appLabel string, model string) (*core_entity.ObjectType, error)
	GetForModel(ctx context.Context, model string) ([]*core_entity.ObjectType, error)
	Public() []string
	WithFeature(feature string) []string
}

// ObjectChangeFilter фильтр записей журнала изменений
type ObjectChangeFilter struct {
	ChangedObjectType *string
	ChangedObjectID   *string
	UserID            *string
	Action            *string
	RequestID         *string
	Since             *time.Time
	Until             *time.Time
	Limit             int
	Offset            int
}

// ObjectChangeRepository операции над журналом изменений
type ObjectChangeRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.ObjectChange, error)
	List(ctx context.Context, filter ObjectChangeFilter) ([]*core_entity.ObjectChange, int64, error)
	Create(ctx context.Context, change *core_entity.ObjectChange) error
	LogChange(ctx context.Context, action types.Status, objectType string, objectID string, objectRepr string, objectData interface{}, userID *types.ID, requestID *string) error
	GetChangesForObject(ctx context.Context, objectType string, objectID string, limit int, offset int) ([]*core_entity.ObjectChange, int64, error)
}

// ConfigRevisionFilter фильтр ревизий конфигурации
type ConfigRevisionFilter struct {
	Active       *bool
	CreatedSince *time.Time
	CreatedUntil *time.Time
	Limit        int
	Offset       int
}

// ConfigRevisionRepository операции над ревизиями конфигурации
type ConfigRevisionRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.ConfigRevision, error)
	List(ctx context.Context, filter ConfigRevisionFilter) ([]*core_entity.ConfigRevision, int64, error)
	Create(ctx context.Context, revision *core_entity.ConfigRevision) error
	Update(ctx context.Context, revision *core_entity.ConfigRevision) error
	Activate(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	GetActive(ctx context.Context) (*core_entity.ConfigRevision, error)
	GetLatest(ctx context.Context) (*core_entity.ConfigRevision, error)
}

// DataSourceFilter фильтр источников данных
type DataSourceFilter struct {
	Name         *string
	Type         *string
	Status       *string
	Enabled      *bool
	SyncInterval *int
	Limit        int
	Offset       int
}

// DataSourceRepository операции над источниками данных
type DataSourceRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.DataSource, error)
	List(ctx context.Context, filter DataSourceFilter) ([]*core_entity.DataSource, int64, error)
	Create(ctx context.Context, ds *core_entity.DataSource) error
	Update(ctx context.Context, ds *core_entity.DataSource) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string, lastSynced *time.Time) error
	Sync(ctx context.Context, id string) error
	Exists(ctx context.Context, name string) (bool, error)
	GetByName(ctx context.Context, name string) (*core_entity.DataSource, error)
}

// DataFileFilter фильтр файлов данных
type DataFileFilter struct {
	SourceID *string
	Path     *string
	Limit    int
	Offset   int
}

// DataFileRepository операции над файлами данных
type DataFileRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.DataFile, error)
	List(ctx context.Context, filter DataFileFilter) ([]*core_entity.DataFile, int64, error)
	Create(ctx context.Context, df *core_entity.DataFile) error
	Update(ctx context.Context, df *core_entity.DataFile) error
	Delete(ctx context.Context, id string) error
	BulkCreate(ctx context.Context, files []*core_entity.DataFile) error
	BulkUpdate(ctx context.Context, files []*core_entity.DataFile) error
	BulkDelete(ctx context.Context, ids []string) error
}

// JobFilter фильтр задач (jobs)
type JobFilter struct {
	ObjectType  *string
	ObjectID    *string
	Status      *string
	QueueName   *string
	ScheduledAt *time.Time
	Limit       int
	Offset      int
}

// JobRepository операции над задачами
type JobRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.Job, error)
	List(ctx context.Context, filter JobFilter) ([]*core_entity.Job, int64, error)
	Create(ctx context.Context, job *core_entity.Job) error
	Update(ctx context.Context, job *core_entity.Job) error
	Delete(ctx context.Context, id string) error
	Start(ctx context.Context, id string, queueName string, jobID string) error
	Complete(ctx context.Context, id string, hasError bool, errorMsg *string) error
	Log(ctx context.Context, id string, message string) error
}

// RoleRepository операции над ролями
type RoleRepository interface {
	GetByID(ctx context.Context, id string) (*users_entity.Role, error)
	GetByName(ctx context.Context, name string) (*users_entity.Role, error)
	List(ctx context.Context, filter RoleFilter) ([]*users_entity.Role, int64, error)
	Create(ctx context.Context, role *users_entity.Role) error
	Update(ctx context.Context, role *users_entity.Role) error
	Delete(ctx context.Context, id string) error
}

// PermissionRepository операции над разрешениями
type PermissionRepository interface {
	GetByID(ctx context.Context, id string) (*users_entity.Permission, error)
	GetByCode(ctx context.Context, code string) (*users_entity.Permission, error)
	List(ctx context.Context, filter PermissionFilter) ([]*users_entity.Permission, int64, error)
	Create(ctx context.Context, permission *users_entity.Permission) error
	Update(ctx context.Context, permission *users_entity.Permission) error
	Delete(ctx context.Context, id string) error
}

// UserRoleRepository операции над связями пользователь-роль
type UserRoleRepository interface {
	GetByID(ctx context.Context, id string) (*users_entity.UserRole, error)
	ListByUser(ctx context.Context, userID string) ([]*users_entity.UserRole, error)
	ListByRole(ctx context.Context, roleID string) ([]*users_entity.UserRole, error)
	Create(ctx context.Context, userRole *users_entity.UserRole) error
	Delete(ctx context.Context, id string) error
	DeleteByUserAndRole(ctx context.Context, userID, roleID string) error
}

// JobAssignmentRepository операции над назначениями задач
type JobAssignmentRepository interface {
	GetByID(ctx context.Context, id string) (*users_entity.JobAssignment, error)
	ListByJob(ctx context.Context, jobID string) ([]*users_entity.JobAssignment, error)
	ListByUser(ctx context.Context, userID string) ([]*users_entity.JobAssignment, error)
	ListByGroup(ctx context.Context, groupID string) ([]*users_entity.JobAssignment, error)
	Create(ctx context.Context, assignment *users_entity.JobAssignment) error
	Update(ctx context.Context, assignment *users_entity.JobAssignment) error
	Delete(ctx context.Context, id string) error
}

// RoleFilter фильтр для поиска ролей
type RoleFilter struct {
	Name   *string
	Limit  int
	Offset int
}

// PermissionFilter фильтр для поиска разрешений
type PermissionFilter struct {
	Code       *string
	ObjectType *string
	Action     *string
	Limit      int
	Offset     int
}

// AutoSyncRecordRepository операции над записями авто-синхронизации
type AutoSyncRecordRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.AutoSyncRecord, error)
	List(ctx context.Context, filter AutoSyncRecordFilter) ([]*core_entity.AutoSyncRecord, int64, error)
	Create(ctx context.Context, record *core_entity.AutoSyncRecord) error
	Update(ctx context.Context, record *core_entity.AutoSyncRecord) error
	Delete(ctx context.Context, id string) error
	GetByObjectTypeAndID(ctx context.Context, objectType string, objectID string) (*core_entity.AutoSyncRecord, error)
	GetByDataFileID(ctx context.Context, datafileID string) ([]*core_entity.AutoSyncRecord, error)
}

// AutoSyncRecordFilter фильтр для поиска записей авто-синхронизации
type AutoSyncRecordFilter struct {
	DataFileID *string
	ObjectType *string
	ObjectID   *string
	Limit      int
	Offset     int
}

// ManagedFileRepository операции над управляемыми файлами
type ManagedFileRepository interface {
	GetByID(ctx context.Context, id string) (*core_entity.ManagedFile, error)
	List(ctx context.Context, filter ManagedFileFilter) ([]*core_entity.ManagedFile, int64, error)
	Create(ctx context.Context, file *core_entity.ManagedFile) error
	Update(ctx context.Context, file *core_entity.ManagedFile) error
	Delete(ctx context.Context, id string) error
	GetByRootAndPath(ctx context.Context, fileRoot string, filePath string) (*core_entity.ManagedFile, error)
	GetByDataFileID(ctx context.Context, datafileID string) ([]*core_entity.ManagedFile, error)
}

// ManagedFileFilter фильтр для поиска управляемых файлов
type ManagedFileFilter struct {
	FileRoot   *string
	FilePath   *string
	DataFileID *string
	Limit      int
	Offset     int
}
