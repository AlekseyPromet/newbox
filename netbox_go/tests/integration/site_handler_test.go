package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/entity"
	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/enum"
	"github.com/AlekseyPromet/netbox_go/internal/repository"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSiteRepository - мок репозитория сайтов
type MockSiteRepository struct {
	mock.Mock
}

func (m *MockSiteRepository) GetByID(ctx context.Context, id string) (*entity.Site, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Site), args.Error(1)
}

func (m *MockSiteRepository) GetBySlug(ctx context.Context, slug string) (*entity.Site, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Site), args.Error(1)
}

func (m *MockSiteRepository) List(ctx context.Context, filter repository.SiteFilter) ([]*entity.Site, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Site), args.Get(1).(int64), args.Error(2)
}

func (m *MockSiteRepository) Create(ctx context.Context, site *entity.Site) error {
	args := m.Called(ctx, site)
	return args.Error(0)
}

func (m *MockSiteRepository) Update(ctx context.Context, site *entity.Site) error {
	args := m.Called(ctx, site)
	return args.Error(0)
}

func (m *MockSiteRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSiteRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestSiteHandler_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockSiteRepository)
	handler := NewSiteHandler(mockRepo)

	siteID := "550e8400-e29b-41d4-a716-446655440000"
	expectedSite := &entity.Site{
		ID:     types.ID(types.NewID()),
		Name:   "Test Site",
		Slug:   "test-site",
		Status: enum.SiteStatusActive,
	}

	mockRepo.On("GetByID", mock.Anything, siteID).Return(expectedSite, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dcim/sites/"+siteID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(siteID)

	// Act
	err := handler.GetByID(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestSiteHandler_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockSiteRepository)
	handler := NewSiteHandler(mockRepo)

	siteID := "550e8400-e29b-41d4-a716-446655440000"

	mockRepo.On("GetByID", mock.Anything, siteID).Return(nil, repository.ErrNotFound)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dcim/sites/"+siteID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(siteID)

	// Act
	err := handler.GetByID(c)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 404, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestSiteHandler_List_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockSiteRepository)
	handler := NewSiteHandler(mockRepo)

	sites := []*entity.Site{
		{ID: types.ID(types.NewID()), Name: "Site 1", Slug: "site-1"},
		{ID: types.ID(types.NewID()), Name: "Site 2", Slug: "site-2"},
	}

	mockRepo.On("List", mock.Anything, mock.Anything).Return(sites, int64(2), nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dcim/sites?limit=10&offset=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := handler.List(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestSiteHandler_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockSiteRepository)
	handler := NewSiteHandler(mockRepo)

	newSite := &entity.Site{
		Name:   "New Site",
		Slug:   "new-site",
		Status: enum.SiteStatusActive,
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/dcim/sites", strings.NewReader(`{"name":"New Site","slug":"new-site","status":"active"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := handler.Create(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 201, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestSiteHandler_Update_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockSiteRepository)
	handler := NewSiteHandler(mockRepo)

	siteID := "550e8400-e29b-41d4-a716-446655440000"
	existingSite := &entity.Site{
		ID:     types.ID(types.NewID()),
		Name:   "Old Name",
		Slug:   "old-slug",
		Status: enum.SiteStatusActive,
	}
	updatedSite := &entity.Site{
		ID:     existingSite.ID,
		Name:   "Updated Name",
		Slug:   "updated-slug",
		Status: enum.SiteStatusActive,
	}

	mockRepo.On("GetByID", mock.Anything, siteID).Return(existingSite, nil).Once()
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetByID", mock.Anything, siteID).Return(updatedSite, nil).Twice()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/api/dcim/sites/"+siteID, strings.NewReader(`{"name":"Updated Name","slug":"updated-slug","status":"active"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(siteID)

	// Act
	err := handler.Update(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestSiteHandler_Delete_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockSiteRepository)
	handler := NewSiteHandler(mockRepo)

	siteID := "550e8400-e29b-41d4-a716-446655440000"

	mockRepo.On("Delete", mock.Anything, siteID).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/dcim/sites/"+siteID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(siteID)

	// Act
	err := handler.Delete(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 204, rec.Code)
	mockRepo.AssertExpectations(t)
}
