package core

import (
	"context"
	"encoding/json"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
	"github.com/google/uuid"
)

// ChangeLogService сервис для логирования изменений объектов
type ChangeLogService struct {
	objectChangeRepo repository.ObjectChangeRepository
	objectTypeRepo   repository.ObjectTypeRepository
}

func NewChangeLogService(
	ocRepo repository.ObjectChangeRepository,
	otRepo repository.ObjectTypeRepository,
) *ChangeLogService {
	return &ChangeLogService{
		objectChangeRepo: ocRepo,
		objectTypeRepo:   otRepo,
	}
}

// LogChangeParams параметры для логирования изменения
type LogChangeParams struct {
	Action          types.Status
	ObjectType      string
	ObjectID        string
	ObjectRepr      string
	PreChangeData   interface{}
	PostChangeData  interface{}
	UserID          *types.ID
	UserName        string
	RequestID       *string
	RelatedObject   *RelatedObjectInfo
	Message         string
}

// RelatedObjectInfo информация о связанном объекте
type RelatedObjectInfo struct {
	ObjectType string
	ObjectID   string
	ObjectRepr string
}

// LogChange логирует изменение объекта
func (s *ChangeLogService) LogChange(ctx context.Context, params LogChangeParams) error {
	// Сериализация данных изменений
	var objectDataJSON json.RawMessage
	var err error

	if params.PostChangeData != nil {
		objectDataJSON, err = json.Marshal(params.PostChangeData)
		if err != nil {
			return err
		}
	}

	// Генерация request_id если не предоставлен
	requestID := params.RequestID
	if requestID == nil {
		id := uuid.New().String()
		requestID = &id
	}

	oc := &entity.ObjectChange{
		Time:              time.Now(),
		UserID:            params.UserID,
		RequestID:         requestID,
		Action:            params.Action,
		ChangedObjectType: params.ObjectType,
		ChangedObjectID:   params.ObjectID,
		ObjectRepr:        params.ObjectRepr,
		ObjectData:        objectDataJSON,
		RelatedObjectType: nil,
		RelatedObjectID:   nil,
		RelatedObjectRepr: nil,
	}

	if params.RelatedObject != nil {
		oc.RelatedObjectType = &params.RelatedObject.ObjectType
		oc.RelatedObjectID = &params.RelatedObject.ObjectID
		oc.RelatedObjectRepr = &params.RelatedObject.ObjectRepr
	}

	return s.objectChangeRepo.Create(ctx, oc)
}

// GetObjectHistory возвращает историю изменений объекта
func (s *ChangeLogService) GetObjectHistory(
	ctx context.Context,
	objectType string,
	objectID string,
	limit int,
) ([]*entity.ObjectChange, error) {
	filter := repository.ObjectChangeFilter{
		ChangedObjectType: &objectType,
		ChangedObjectID:   &objectID,
		Limit:             limit,
	}

	changes, _, err := s.objectChangeRepo.List(ctx, filter)
	return changes, err
}

// GetRecentChanges возвращает последние изменения
func (s *ChangeLogService) GetRecentChanges(
	ctx context.Context,
	limit int,
	since *time.Time,
) ([]*entity.ObjectChange, error) {
	filter := repository.ObjectChangeFilter{
		Since: since,
		Limit: limit,
	}

	changes, _, err := s.objectChangeRepo.List(ctx, filter)
	return changes, err
}
