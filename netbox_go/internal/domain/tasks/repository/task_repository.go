// Package repository содержит интерфейсы репозиториев домена Tasks
package repository

import (
	"context"

	"netbox_go/internal/domain/tasks/entity"
)

// WorkTypeRepository определяет интерфейс для работы с видами работ
type WorkTypeRepository interface {
	Create(ctx context.Context, wt *entity.WorkType) error
	GetByID(ctx context.Context, id string) (*entity.WorkType, error)
	Update(ctx context.Context, wt *entity.WorkType) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter WorkTypeFilter) ([]*entity.WorkType, int64, error)
}

// WorkTypeFilter параметры фильтрации видов работ
type WorkTypeFilter struct {
	Limit  int
	Offset int
	Search string // Поиск по названию или описанию
}

// GroupRepository определяет интерфейс для работы с группами
type GroupRepository interface {
	Create(ctx context.Context, g *entity.Group) error
	GetByID(ctx context.Context, id string) (*entity.Group, error)
	Update(ctx context.Context, g *entity.Group) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter GroupFilter) ([]*entity.Group, int64, error)
	
	// Методы для управления участниками группы
	AddMember(ctx context.Context, member *entity.GroupMember) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	GetMembers(ctx context.Context, groupID string) ([]*entity.GroupMember, error)
	
	// Методы для управления компетенциями группы
	AddWorkType(ctx context.Context, gwt *entity.GroupWorkType) error
	RemoveWorkType(ctx context.Context, groupID, workTypeID string) error
	GetWorkTypes(ctx context.Context, groupID string) ([]*entity.GroupWorkType, error)
	CanPerformWork(ctx context.Context, groupID, workTypeID string) (bool, error)
}

// GroupFilter параметры фильтрации групп
type GroupFilter struct {
	Limit      int
	Offset     int
	Type       entity.GroupType // Фильтр по типу группы
	Search     string           // Поиск по названию
	WorkTypeID string           // Фильтр по компетенции (вид работ)
}

// TaskRepository определяет интерфейс для работы с задачами
type TaskRepository interface {
	Create(ctx context.Context, t *entity.Task) error
	GetByID(ctx context.Context, id string) (*entity.Task, error)
	Update(ctx context.Context, t *entity.Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter TaskFilter) ([]*entity.Task, int64, error)
}

// TaskFilter параметры фильтрации задач
type TaskFilter struct {
	Limit        int
	Offset       int
	Status       entity.TaskStatus // Фильтр по статусу
	WorkTypeID   string            // Фильтр по виду работ
	CreatedByID  string            // Фильтр по создателю
	AssigneeID   string            // Фильтр по ответственному (user или group)
	ReviewerID   string            // Фильтр по проверяющему
	Priority     int               // Фильтр по приоритету
	Search       string            // Поиск по заголовку или описанию
	DueDateFrom  string            // Фильтр по дедлайну (от)
	DueDateTo    string            // Фильтр по дедлайну (до)
}

// TaskAssignmentRepository определяет интерфейс для назначений задач
type TaskAssignmentRepository interface {
	Create(ctx context.Context, ta *entity.TaskAssignment) error
	GetByID(ctx context.Context, id string) (*entity.TaskAssignment, error)
	Update(ctx context.Context, ta *entity.TaskAssignment) error
	Delete(ctx context.Context, id string) error
	
	// Получить все назначения для задачи
	GetByTaskID(ctx context.Context, taskID string) ([]*entity.TaskAssignment, error)
	
	// Получить назначение по задаче и роли
	GetByTaskAndRole(ctx context.Context, taskID string, role entity.TaskRole) (*entity.TaskAssignment, error)
	
	// Получить задачи по пользователю и роли
	GetByUserAndRole(ctx context.Context, userID string, role entity.TaskRole) ([]*entity.TaskAssignment, error)
	
	// Получить задачи по группе и роли
	GetByGroupAndRole(ctx context.Context, groupID string, role entity.TaskRole) ([]*entity.TaskAssignment, error)
	
	// Установить фактического исполнителя (когда задача назначена группе)
	SetActualExecutor(ctx context.Context, assignmentID, executorUserID string) error
}

// TaskCommentRepository определяет интерфейс для комментариев к задачам
type TaskCommentRepository interface {
	Create(ctx context.Context, tc *entity.TaskComment) error
	GetByID(ctx context.Context, id string) (*entity.TaskComment, error)
	Update(ctx context.Context, tc *entity.TaskComment) error
	Delete(ctx context.Context, id string) error
	ListByTaskID(ctx context.Context, taskID string, limit, offset int) ([]*entity.TaskComment, int64, error)
}

// TaskAttachmentRepository определяет интерфейс для вложений задач
type TaskAttachmentRepository interface {
	Create(ctx context.Context, ta *entity.TaskAttachment) error
	GetByID(ctx context.Context, id string) (*entity.TaskAttachment, error)
	Delete(ctx context.Context, id string) error
	ListByTaskID(ctx context.Context, taskID string, limit, offset int) ([]*entity.TaskAttachment, int64, error)
}
