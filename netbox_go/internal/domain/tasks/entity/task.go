// Package entity содержит сущности домена Tasks для управления задачами, ролями и группами
package entity

import (
	"errors"
	"time"

	"netbox_go/pkg/types"
)

// TaskRole определяет роль пользователя в задаче
type TaskRole string

const (
	// TaskRoleCreator - создатель задачи (формирует задачу, указывает вид работ и назначает ответственных)
	TaskRoleCreator TaskRole = "creator"
	// TaskRoleAssignee - ответственный за выполнение задачи
	TaskRoleAssignee TaskRole = "assignee"
	// TaskRoleReviewer - проверяющий результат выполнения задачи
	TaskRoleReviewer TaskRole = "reviewer"
)

// IsValid проверяет, является ли роль допустимой
func (r TaskRole) IsValid() bool {
	switch r {
	case TaskRoleCreator, TaskRoleAssignee, TaskRoleReviewer:
		return true
	default:
		return false
	}
}

// WorkType представляет вид работ в системе
type WorkType struct {
	ID          types.ID  `json:"id"`
	Name        string    `json:"name"`         // Название вида работ (например: "Физическое подключение сервера")
	Description string    `json:"description"`  // Описание вида работ
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// Validate проверяет корректность вида работ
func (wt *WorkType) Validate() error {
	if len(wt.Name) == 0 || len(wt.Name) > 100 {
		return errors.New("name is required and must be <= 100 characters")
	}
	if len(wt.Description) > 500 {
		return errors.New("description too long (max 500 characters)")
	}
	return nil
}

// GroupType определяет тип группы
type GroupType string

const (
	// GroupTypeAssignee - группа ответственных (выполняет задачи)
	GroupTypeAssignee GroupType = "assignee"
	// GroupTypeReviewer - группа проверяющих (принимает результат выполнения)
	GroupTypeReviewer GroupType = "reviewer"
)

// IsValid проверяет, является ли тип группы допустимым
func (gt GroupType) IsValid() bool {
	switch gt {
	case GroupTypeAssignee, GroupTypeReviewer:
		return true
	default:
		return false
	}
}

// Group представляет группу пользователей (может быть группой ответственных или проверяющих)
type Group struct {
	ID          types.ID   `json:"id"`
	Name        string     `json:"name"`         // Название группы
	Type        GroupType  `json:"type"`         // Тип группы: assignee или reviewer
	Description string     `json:"description"`  // Описание группы
	ShiftStart  *time.Time `json:"shift_start"`  // Начало смены (для групп-смен)
	ShiftEnd    *time.Time `json:"shift_end"`    // Конец смены (для групп-смен)
	WorkDays    []string   `json:"work_days"`    // Дни недели работы (понедельник, вторник, ...)
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

// Validate проверяет корректность группы
func (g *Group) Validate() error {
	if len(g.Name) == 0 || len(g.Name) > 100 {
		return errors.New("name is required and must be <= 100 characters")
	}
	if !g.Type.IsValid() {
		return errors.New("invalid group type")
	}
	if len(g.Description) > 500 {
		return errors.New("description too long (max 500 characters)")
	}
	// Если заданы время смены, проверяем корректность
	if g.ShiftStart != nil && g.ShiftEnd != nil {
		if g.ShiftEnd.Before(*g.ShiftStart) {
			return errors.New("shift_end must be after shift_start")
		}
	}
	return nil
}

// IsShiftGroup проверяет, является ли группа группой-сменой
func (g *Group) IsShiftGroup() bool {
	return g.ShiftStart != nil && g.ShiftEnd != nil
}

// GroupMember представляет участника группы
type GroupMember struct {
	ID        types.ID  `json:"id"`
	GroupID   types.ID  `json:"group_id"`
	UserID    types.ID  `json:"user_id"`
	Created   time.Time `json:"created"`
	AddedBy   types.ID  `json:"added_by"` // Кто добавил участника
}

// Validate проверяет корректность участника группы
func (gm *GroupMember) Validate() error {
	if gm.GroupID.String() == "" {
		return errors.New("group_id is required")
	}
	if gm.UserID.String() == "" {
		return errors.New("user_id is required")
	}
	return nil
}

// GroupWorkType представляет компетенцию группы (связь группы с видами работ, которые она может выполнять)
type GroupWorkType struct {
	ID         types.ID  `json:"id"`
	GroupID    types.ID  `json:"group_id"`
	WorkTypeID types.ID  `json:"work_type_id"`
	Created    time.Time `json:"created"`
	AddedBy    types.ID  `json:"added_by"` // Кто добавил компетенцию
}

// Validate проверяет корректность компетенции группы
func (gwt *GroupWorkType) Validate() error {
	if gwt.GroupID.String() == "" {
		return errors.New("group_id is required")
	}
	if gwt.WorkTypeID.String() == "" {
		return errors.New("work_type_id is required")
	}
	return nil
}

// Task представляет задачу в системе
type Task struct {
	ID              types.ID   `json:"id"`
	Title           string     `json:"title"`            // Заголовок задачи
	Description     string     `json:"description"`      // Описание задачи
	WorkTypeID      types.ID   `json:"work_type_id"`     // Вид работ
	Status          TaskStatus `json:"status"`           // Статус задачи
	Priority        int        `json:"priority"`         // Приоритет (1-5, где 5 - наивысший)
	CreatedByID     types.ID   `json:"created_by_id"`    // Создатель задачи
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DueDate         *time.Time `json:"due_date,omitempty"` // Дедлайн
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	ReviewedAt      *time.Time `json:"reviewed_at,omitempty"`
	ReviewComment   string     `json:"review_comment,omitempty"`
}

// Validate проверяет корректность задачи
func (t *Task) Validate() error {
	if len(t.Title) == 0 || len(t.Title) > 200 {
		return errors.New("title is required and must be <= 200 characters")
	}
	if len(t.Description) > 2000 {
		return errors.New("description too long (max 2000 characters)")
	}
	if t.WorkTypeID.String() == "" {
		return errors.New("work_type_id is required")
	}
	if !t.Status.IsValid() {
		return errors.New("invalid task status")
	}
	if t.Priority < 1 || t.Priority > 5 {
		return errors.New("priority must be between 1 and 5")
	}
	if t.CreatedByID.String() == "" {
		return errors.New("created_by_id is required")
	}
	return nil
}

// TaskStatus определяет статус задачи
type TaskStatus string

const (
	// TaskStatusDraft - черновик (задача создана, но ещё не назначена)
	TaskStatusDraft TaskStatus = "draft"
	// TaskStatusAssigned - назначена (ответственный назначен и задача ожидает выполнения)
	TaskStatusAssigned TaskStatus = "assigned"
	// TaskStatusInProgress - в работе (ответственный начал выполнение)
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusCompleted - выполнена (ответственный завершил работу)
	TaskStatusCompleted TaskStatus = "completed"
	// TaskStatusUnderReview - на проверке (ожидает проверки проверяющим)
	TaskStatusUnderReview TaskStatus = "under_review"
	// TaskStatusApproved - принята (проверяющий принял результат)
	TaskStatusApproved TaskStatus = "approved"
	// TaskStatusRejected - отклонена (проверяющий отклонил результат, требуется доработка)
	TaskStatusRejected TaskStatus = "rejected"
	// TaskStatusCancelled - отменена
	TaskStatusCancelled TaskStatus = "cancelled"
)

// IsValid проверяет, является ли статус допустимым
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusDraft, TaskStatusAssigned, TaskStatusInProgress,
		TaskStatusCompleted, TaskStatusUnderReview, TaskStatusApproved,
		TaskStatusRejected, TaskStatusCancelled:
		return true
	default:
		return false
	}
}

// IsFinal проверяет, является ли статус финальным
func (s TaskStatus) IsFinal() bool {
	return s == TaskStatusApproved || s == TaskStatusCancelled
}

// TaskAssignment представляет назначение роли пользователю или группе в задаче
type TaskAssignment struct {
	ID           types.ID  `json:"id"`
	TaskID       types.ID  `json:"task_id"`
	Role         TaskRole  `json:"role"`             // Роль: creator, assignee, reviewer
	AssigneeType string    `json:"assignee_type"`    // "user" или "group"
	UserID       *types.ID `json:"user_id,omitempty"`
	GroupID      *types.ID `json:"group_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    types.ID  `json:"created_by"`
	// Для ответственного: кто конкретно из группы выполняет задачу (если назначена группа)
	ActualExecutorID *types.ID `json:"actual_executor_id,omitempty"`
}

// Validate проверяет корректность назначения
func (ta *TaskAssignment) Validate() error {
	if ta.TaskID.String() == "" {
		return errors.New("task_id is required")
	}
	if !ta.Role.IsValid() {
		return errors.New("invalid role")
	}
	if ta.AssigneeType != "user" && ta.AssigneeType != "group" {
		return errors.New("assignee_type must be 'user' or 'group'")
	}
	if ta.AssigneeType == "user" && ta.UserID == nil {
		return errors.New("user_id is required for user assignee")
	}
	if ta.AssigneeType == "group" && ta.GroupID == nil {
		return errors.New("group_id is required for group assignee")
	}
	if ta.CreatedBy.String() == "" {
		return errors.New("created_by is required")
	}
	return nil
}

// GetAssigneeID возвращает ID назначенного (пользователя или группы)
func (ta *TaskAssignment) GetAssigneeID() *types.ID {
	if ta.AssigneeType == "user" {
		return ta.UserID
	}
	return ta.GroupID
}

// IsUserAssignment проверяет, назначена ли задача конкретному пользователю
func (ta *TaskAssignment) IsUserAssignment() bool {
	return ta.AssigneeType == "user"
}

// IsGroupAssignment проверяет, назначена ли задача группе
func (ta *TaskAssignment) IsGroupAssignment() bool {
	return ta.AssigneeType == "group"
}

// TaskComment представляет комментарий к задаче
type TaskComment struct {
	ID        types.ID  `json:"id"`
	TaskID    types.ID  `json:"task_id"`
	AuthorID  types.ID  `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate проверяет корректность комментария
func (tc *TaskComment) Validate() error {
	if tc.TaskID.String() == "" {
		return errors.New("task_id is required")
	}
	if tc.AuthorID.String() == "" {
		return errors.New("author_id is required")
	}
	if len(tc.Content) == 0 || len(tc.Content) > 2000 {
		return errors.New("content is required and must be <= 2000 characters")
	}
	return nil
}

// TaskAttachment представляет вложение к задаче
type TaskAttachment struct {
	ID         types.ID  `json:"id"`
	TaskID     types.ID  `json:"task_id"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	FileType   string    `json:"file_type"`
	UploadedBy types.ID  `json:"uploaded_by"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Validate проверяет корректность вложения
func (ta *TaskAttachment) Validate() error {
	if ta.TaskID.String() == "" {
		return errors.New("task_id is required")
	}
	if len(ta.FileName) == 0 || len(ta.FileName) > 255 {
		return errors.New("file_name is required and must be <= 255 characters")
	}
	if ta.UploadedBy.String() == "" {
		return errors.New("uploaded_by is required")
	}
	return nil
}
