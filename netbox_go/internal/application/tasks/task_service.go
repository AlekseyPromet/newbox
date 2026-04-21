// Package tasks содержит сервисы для управления задачами, группами и видами работ
package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlekseyPromet/netbox_go/internal/domain/tasks/entity"
	taskrepo "github.com/AlekseyPromet/netbox_go/internal/domain/tasks/repository"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

var (
	// ErrTaskNotFound ошибка когда задача не найдена
	ErrTaskNotFound = errors.New("task not found")
	// ErrWorkTypeNotFound ошибка когда вид работ не найден
	ErrWorkTypeNotFound = errors.New("work type not found")
	// ErrGroupNotFound ошибка когда группа не найдена
	ErrGroupNotFound = errors.New("group not found")
	// ErrInvalidRoleAssignment ошибка недопустимого назначения роли
	ErrInvalidRoleAssignment = errors.New("invalid role assignment")
	// ErrGroupCannotPerformWork ошибка когда группа не может выполнить данный вид работ
	ErrGroupCannotPerformWork = errors.New("group cannot perform this work type")
)

// WorkTypeService сервис для управления видами работ
type WorkTypeService struct {
	workTypeRepo taskrepo.WorkTypeRepository
}

// NewWorkTypeService создаёт новый сервис видов работ
func NewWorkTypeService(workTypeRepo taskrepo.WorkTypeRepository) *WorkTypeService {
	return &WorkTypeService{
		workTypeRepo: workTypeRepo,
	}
}

// CreateWorkTypeParams параметры создания вида работ
type CreateWorkTypeParams struct {
	Name        string
	Description string
}

// CreateWorkType создаёт новый вид работ
func (s *WorkTypeService) CreateWorkType(ctx context.Context, params CreateWorkTypeParams) (*entity.WorkType, error) {
	now := time.Now()
	wt := &entity.WorkType{
		ID:          types.NewID(),
		Name:        params.Name,
		Description: params.Description,
		Created:     now,
		Updated:     now,
	}

	if err := wt.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.workTypeRepo.Create(ctx, wt); err != nil {
		return nil, fmt.Errorf("failed to create work type: %w", err)
	}

	return wt, nil
}

// GetWorkType получает вид работ по ID
func (s *WorkTypeService) GetWorkType(ctx context.Context, id string) (*entity.WorkType, error) {
	wt, err := s.workTypeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get work type: %w", err)
	}
	if wt == nil {
		return nil, ErrWorkTypeNotFound
	}
	return wt, nil
}

// UpdateWorkType обновляет вид работ
func (s *WorkTypeService) UpdateWorkType(ctx context.Context, id string, params CreateWorkTypeParams) (*entity.WorkType, error) {
	wt, err := s.GetWorkType(ctx, id)
	if err != nil {
		return nil, err
	}

	wt.Name = params.Name
	wt.Description = params.Description
	wt.Updated = time.Now()

	if err := wt.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.workTypeRepo.Update(ctx, wt); err != nil {
		return nil, fmt.Errorf("failed to update work type: %w", err)
	}

	return wt, nil
}

// DeleteWorkType удаляет вид работ
func (s *WorkTypeService) DeleteWorkType(ctx context.Context, id string) error {
	if err := s.workTypeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete work type: %w", err)
	}
	return nil
}

// ListWorkTypes получает список видов работ
func (s *WorkTypeService) ListWorkTypes(ctx context.Context, filter taskrepo.WorkTypeFilter) ([]*entity.WorkType, int64, error) {
	return s.workTypeRepo.List(ctx, filter)
}

// GroupService сервис для управления группами
type GroupService struct {
	groupRepo taskrepo.GroupRepository
}

// NewGroupService создаёт новый сервис групп
func NewGroupService(groupRepo taskrepo.GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

// CreateGroupParams параметры создания группы
type CreateGroupParams struct {
	Name        string
	Type        entity.GroupType
	Description string
	ShiftStart  *time.Time
	ShiftEnd    *time.Time
	WorkDays    []string
}

// CreateGroup создаёт новую группу
func (s *GroupService) CreateGroup(ctx context.Context, params CreateGroupParams, createdBy types.ID) (*entity.Group, error) {
	now := time.Now()
	g := &entity.Group{
		ID:          types.NewID(),
		Name:        params.Name,
		Type:        params.Type,
		Description: params.Description,
		ShiftStart:  params.ShiftStart,
		ShiftEnd:    params.ShiftEnd,
		WorkDays:    params.WorkDays,
		Created:     now,
		Updated:     now,
	}

	if err := g.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.groupRepo.Create(ctx, g); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return g, nil
}

// GetGroup получает группу по ID
func (s *GroupService) GetGroup(ctx context.Context, id string) (*entity.Group, error) {
	g, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if g == nil {
		return nil, ErrGroupNotFound
	}
	return g, nil
}

// UpdateGroup обновляет группу
func (s *GroupService) UpdateGroup(ctx context.Context, id string, params CreateGroupParams) (*entity.Group, error) {
	g, err := s.GetGroup(ctx, id)
	if err != nil {
		return nil, err
	}

	g.Name = params.Name
	g.Type = params.Type
	g.Description = params.Description
	g.ShiftStart = params.ShiftStart
	g.ShiftEnd = params.ShiftEnd
	g.WorkDays = params.WorkDays
	g.Updated = time.Now()

	if err := g.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.groupRepo.Update(ctx, g); err != nil {
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	return g, nil
}

// DeleteGroup удаляет группу
func (s *GroupService) DeleteGroup(ctx context.Context, id string) error {
	if err := s.groupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

// ListGroups получает список групп
func (s *GroupService) ListGroups(ctx context.Context, filter taskrepo.GroupFilter) ([]*entity.Group, int64, error) {
	return s.groupRepo.List(ctx, filter)
}

// AddGroupMember добавляет участника в группу
func (s *GroupService) AddGroupMember(ctx context.Context, groupID, userID string, addedBy types.ID) error {
	gID, err := types.ParseID(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}
	uID, err := types.ParseID(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	member := &entity.GroupMember{
		ID:      types.NewID(),
		GroupID: gID,
		UserID:  uID,
		Created: time.Now(),
		AddedBy: addedBy,
	}

	if err := member.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.groupRepo.AddMember(ctx, member); err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}

	return nil
}

// RemoveGroupMember удаляет участника из группы
func (s *GroupService) RemoveGroupMember(ctx context.Context, groupID, userID string) error {
	if err := s.groupRepo.RemoveMember(ctx, groupID, userID); err != nil {
		return fmt.Errorf("failed to remove group member: %w", err)
	}
	return nil
}

// GetGroupMembers получает список участников группы
func (s *GroupService) GetGroupMembers(ctx context.Context, groupID string) ([]*entity.GroupMember, error) {
	return s.groupRepo.GetMembers(ctx, groupID)
}

// AddGroupWorkType добавляет компетенцию группе
func (s *GroupService) AddGroupWorkType(ctx context.Context, groupID, workTypeID string, addedBy types.ID) error {
	// Проверяем существование группы и вида работ
	if _, err := s.GetGroup(ctx, groupID); err != nil {
		return err
	}

	gID, err := types.ParseID(groupID)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}
	wtID, err := types.ParseID(workTypeID)
	if err != nil {
		return fmt.Errorf("invalid work type ID: %w", err)
	}

	gwt := &entity.GroupWorkType{
		ID:         types.NewID(),
		GroupID:    gID,
		WorkTypeID: wtID,
		Created:    time.Now(),
		AddedBy:    addedBy,
	}

	if err := gwt.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.groupRepo.AddWorkType(ctx, gwt); err != nil {
		return fmt.Errorf("failed to add group work type: %w", err)
	}

	return nil
}

// RemoveGroupWorkType удаляет компетенцию у группы
func (s *GroupService) RemoveGroupWorkType(ctx context.Context, groupID, workTypeID string) error {
	if err := s.groupRepo.RemoveWorkType(ctx, groupID, workTypeID); err != nil {
		return fmt.Errorf("failed to remove group work type: %w", err)
	}
	return nil
}

// GetGroupWorkTypes получает список компетенций группы
func (s *GroupService) GetGroupWorkTypes(ctx context.Context, groupID string) ([]*entity.GroupWorkType, error) {
	return s.groupRepo.GetWorkTypes(ctx, groupID)
}

// CanGroupPerformWork проверяет, может ли группа выполнить данный вид работ
func (s *GroupService) CanGroupPerformWork(ctx context.Context, groupID, workTypeID string) (bool, error) {
	return s.groupRepo.CanPerformWork(ctx, groupID, workTypeID)
}

// TaskService сервис для управления задачами
type TaskService struct {
	taskRepo           taskrepo.TaskRepository
	taskAssignmentRepo taskrepo.TaskAssignmentRepository
	groupRepo          taskrepo.GroupRepository
	workTypeRepo       taskrepo.WorkTypeRepository
}

// NewTaskService создаёт новый сервис задач
func NewTaskService(
	taskRepo taskrepo.TaskRepository,
	taskAssignmentRepo taskrepo.TaskAssignmentRepository,
	groupRepo taskrepo.GroupRepository,
	workTypeRepo taskrepo.WorkTypeRepository,
) *TaskService {
	return &TaskService{
		taskRepo:           taskRepo,
		taskAssignmentRepo: taskAssignmentRepo,
		groupRepo:          groupRepo,
		workTypeRepo:       workTypeRepo,
	}
}

// CreateTaskParams параметры создания задачи
type CreateTaskParams struct {
	Title       string
	Description string
	WorkTypeID  string
	Priority    int
	DueDate     *time.Time
}

// CreateTask создаёт новую задачу
func (s *TaskService) CreateTask(ctx context.Context, params CreateTaskParams, createdBy types.ID) (*entity.Task, error) {
	now := time.Now()
	
	// Проверяем существование вида работ
	if _, err := s.workTypeRepo.GetByID(ctx, params.WorkTypeID); err != nil {
		return nil, fmt.Errorf("work type not found: %w", err)
	}

	wtID, err := types.ParseID(params.WorkTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid work type ID: %w", err)
	}

	task := &entity.Task{
		ID:          types.NewID(),
		Title:       params.Title,
		Description: params.Description,
		WorkTypeID:  wtID,
		Status:      entity.TaskStatusDraft,
		Priority:    params.Priority,
		CreatedByID: createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		DueDate:     params.DueDate,
	}

	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Создаём назначение создателя как creator
	creatorAssignment := &entity.TaskAssignment{
		ID:           types.NewID(),
		TaskID:       task.ID,
		Role:         entity.TaskRoleCreator,
		AssigneeType: "user",
		UserID:       &createdBy,
		CreatedAt:    now,
		CreatedBy:    createdBy,
	}

	if err := s.taskAssignmentRepo.Create(ctx, creatorAssignment); err != nil {
		return nil, fmt.Errorf("failed to create creator assignment: %w", err)
	}

	return task, nil
}

// AssignTaskParams параметры назначения ответственного/проверяющего
type AssignTaskParams struct {
	TaskID       string
	Role         entity.TaskRole
	AssigneeType string // "user" или "group"
	UserID       string
	GroupID      string
}

// AssignTask назначает ответственного или проверяющего на задачу
func (s *TaskService) AssignTask(ctx context.Context, params AssignTaskParams, assignedBy types.ID) (*entity.TaskAssignment, error) {
	now := time.Now()

	// Проверяем существование задачи
	task, err := s.taskRepo.GetByID(ctx, params.TaskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Нельзя назначать роли если задача в финальном статусе
	if task.Status.IsFinal() {
		return nil, errors.New("cannot assign roles to a task in final status")
	}

	// Проверяем роль
	if !params.Role.IsValid() || params.Role == entity.TaskRoleCreator {
		return nil, ErrInvalidRoleAssignment
	}

	// Если назначаем группу ответственных, проверяем компетенцию
	if params.Role == entity.TaskRoleAssignee && params.AssigneeType == "group" {
		canPerform, err := s.groupRepo.CanPerformWork(ctx, params.GroupID, task.WorkTypeID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to check group competence: %w", err)
		}
		if !canPerform {
			return nil, ErrGroupCannotPerformWork
		}
	}

	var userID *types.ID
	var groupID *types.ID

	if params.AssigneeType == "user" {
		uid, err := types.ParseID(params.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
		userID = &uid
	} else if params.AssigneeType == "group" {
		gid, err := types.ParseID(params.GroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid group ID: %w", err)
		}
		groupID = &gid
	} else {
		return nil, errors.New("assignee_type must be 'user' or 'group'")
	}

	assignment := &entity.TaskAssignment{
		ID:           types.NewID(),
		TaskID:       task.ID,
		Role:         params.Role,
		AssigneeType: params.AssigneeType,
		UserID:       userID,
		GroupID:      groupID,
		CreatedAt:    now,
		CreatedBy:    assignedBy,
	}

	if err := assignment.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.taskAssignmentRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	// Обновляем статус задачи если назначен ответственный
	if params.Role == entity.TaskRoleAssignee {
		task.Status = entity.TaskStatusAssigned
		task.UpdatedAt = now
		if err := s.taskRepo.Update(ctx, task); err != nil {
			return nil, fmt.Errorf("failed to update task status: %w", err)
		}
	}

	return assignment, nil
}

// GetTask получает задачу по ID
func (s *TaskService) GetTask(ctx context.Context, id string) (*entity.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// UpdateTask обновляет задачу
func (s *TaskService) UpdateTask(ctx context.Context, id string, params CreateTaskParams) (*entity.Task, error) {
	task, err := s.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Title = params.Title
	task.Description = params.Description
	task.WorkTypeID = types.ParseID(params.WorkTypeID)
	task.Priority = params.Priority
	task.DueDate = params.DueDate
	task.UpdatedAt = time.Now()

	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// DeleteTask удаляет задачу
func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	if err := s.taskRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// ListTasks получает список задач
func (s *TaskService) ListTasks(ctx context.Context, filter taskrepo.TaskFilter) ([]*entity.Task, int64, error) {
	return s.taskRepo.List(ctx, filter)
}

// StartTask начинает выполнение задачи (меняет статус на in_progress)
func (s *TaskService) StartTask(ctx context.Context, taskID, userID string) (*entity.Task, error) {
	task, err := s.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != entity.TaskStatusAssigned {
		return nil, errors.New("task is not in assigned status")
	}

	// Проверяем, что пользователь имеет право начать эту задачу
	hasPermission, err := s.userCanStartTask(ctx, taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return nil, errors.New("user does not have permission to start this task")
	}

	task.Status = entity.TaskStatusInProgress
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}

// CompleteTask завершает выполнение задачи (меняет статус на completed)
func (s *TaskService) CompleteTask(ctx context.Context, taskID, userID, comment string) (*entity.Task, error) {
	task, err := s.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != entity.TaskStatusInProgress {
		return nil, errors.New("task is not in progress")
	}

	// Проверяем, что пользователь имеет право завершить эту задачу
	hasPermission, err := s.userCanCompleteTask(ctx, taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return nil, errors.New("user does not have permission to complete this task")
	}

	now := time.Now()
	task.Status = entity.TaskStatusCompleted
	task.ReviewComment = comment
	task.CompletedAt = &now
	task.UpdatedAt = now

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}

// ReviewTaskParams параметры проверки задачи
type ReviewTaskParams struct {
	TaskID      string
	Approved    bool
	Comment     string
	ReviewerID  string
}

// ReviewTask проверяет задачу (принимает или отклоняет)
func (s *TaskService) ReviewTask(ctx context.Context, params ReviewTaskParams) (*entity.Task, error) {
	task, err := s.GetTask(ctx, params.TaskID)
	if err != nil {
		return nil, err
	}

	if task.Status != entity.TaskStatusCompleted {
		return nil, errors.New("task is not completed")
	}

	// Проверяем, что пользователь имеет право проверить эту задачу
	hasPermission, err := s.userCanReviewTask(ctx, params.TaskID, params.ReviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return nil, errors.New("user does not have permission to review this task")
	}

	now := time.Now()
	if params.Approved {
		task.Status = entity.TaskStatusApproved
	} else {
		task.Status = entity.TaskStatusRejected
	}
	task.ReviewComment = params.Comment
	task.ReviewedAt = &now
	task.UpdatedAt = now

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}

// userCanStartTask проверяет, может ли пользователь начать задачу
func (s *TaskService) userCanStartTask(ctx context.Context, taskID, userID string) (bool, error) {
	assignment, err := s.taskAssignmentRepo.GetByTaskAndRole(ctx, taskID, entity.TaskRoleAssignee)
	if err != nil {
		return false, err
	}
	if assignment == nil {
		return false, nil
	}

	if assignment.AssigneeType == "user" && assignment.UserID != nil {
		return assignment.UserID.String() == userID, nil
	}

	if assignment.AssigneeType == "group" && assignment.GroupID != nil {
		// Проверяем, является ли пользователь участником группы
		members, err := s.groupRepo.GetMembers(ctx, assignment.GroupID.String())
		if err != nil {
			return false, err
		}
		for _, member := range members {
			if member.UserID.String() == userID {
				return true, nil
			}
		}
	}

	return false, nil
}

// userCanCompleteTask проверяет, может ли пользователь завершить задачу
func (s *TaskService) userCanCompleteTask(ctx context.Context, taskID, userID string) (bool, error) {
	return s.userCanStartTask(ctx, taskID, userID)
}

// userCanReviewTask проверяет, может ли пользователь проверить задачу
func (s *TaskService) userCanReviewTask(ctx context.Context, taskID, reviewerID string) (bool, error) {
	assignment, err := s.taskAssignmentRepo.GetByTaskAndRole(ctx, taskID, entity.TaskRoleReviewer)
	if err != nil {
		return false, err
	}
	if assignment == nil {
		return false, nil
	}

	if assignment.AssigneeType == "user" && assignment.UserID != nil {
		return assignment.UserID.String() == reviewerID, nil
	}

	if assignment.AssigneeType == "group" && assignment.GroupID != nil {
		// Проверяем, является ли пользователь участником группы
		members, err := s.groupRepo.GetMembers(ctx, assignment.GroupID.String())
		if err != nil {
			return false, err
		}
		for _, member := range members {
			if member.UserID.String() == reviewerID {
				return true, nil
			}
		}
	}

	return false, nil
}

// SetActualExecutor устанавливает фактического исполнителя задачи (когда задача назначена группе)
func (s *TaskService) SetActualExecutor(ctx context.Context, taskID, assignmentID, executorUserID string) error {
	assignment, err := s.taskAssignmentRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return fmt.Errorf("assignment not found: %w", err)
	}

	if assignment.TaskID.String() != taskID {
		return errors.New("assignment does not belong to this task")
	}

	if !assignment.IsGroupAssignment() {
		return errors.New("can only set actual executor for group assignments")
	}

	if err := s.taskAssignmentRepo.SetActualExecutor(ctx, assignmentID, executorUserID); err != nil {
		return fmt.Errorf("failed to set actual executor: %w", err)
	}

	return nil
}

// GetTaskAssignments получает все назначения для задачи
func (s *TaskService) GetTaskAssignments(ctx context.Context, taskID string) ([]*entity.TaskAssignment, error) {
	return s.taskAssignmentRepo.GetByTaskID(ctx, taskID)
}

// GetUserTasks получает задачи пользователя по роли
func (s *TaskService) GetUserTasks(ctx context.Context, userID string, role entity.TaskRole) ([]*entity.TaskAssignment, error) {
	return s.taskAssignmentRepo.GetByUserAndRole(ctx, userID, role)
}

// GetGroupTasks получает задачи группы по роли
func (s *TaskService) GetGroupTasks(ctx context.Context, groupID string, role entity.TaskRole) ([]*entity.TaskAssignment, error) {
	return s.taskAssignmentRepo.GetByGroupAndRole(ctx, groupID, role)
}
