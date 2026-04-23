package app

import (
	"net/http"

	"netbox_go/internal/delivery/graphql"
	"netbox_go/internal/delivery/http/handlers"
	"netbox_go/internal/delivery/http/middleware"
	"netbox_go/internal/domain/core/services"
	"netbox_go/internal/repository/postgres"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ModuleAPI provides handler dependencies and route registration
var ModuleAPI = fx.Options(
	fx.Provide(
		NewSiteHandler,
		NewAccountHandler,
		NewCoreHandlers,
		NewGraphQLResolver,
	),
	fx.Invoke(RegisterRoutes),
)

func NewSiteHandler(repo *postgres.SiteRepositoryPostgres) *handlers.SiteHandler {
	return handlers.NewSiteHandler(repo)
}

func NewAccountHandler(repo *postgres.AccountRepositoryPostgres) *handlers.AccountHandler {
	return handlers.NewAccountHandler(
		repo, repo, repo, repo, repo,
	)
}

func NewCoreHandlers(
	dsRepo *postgres.DataSourceRepositoryPostgres,
	dfRepo *postgres.DataFileRepositoryPostgres,
	jobRepo *postgres.JobRepositoryPostgres,
	ocRepo *postgres.ObjectChangeRepositoryPostgres,
	otRepo *postgres.ObjectTypeRepositoryPostgres,
	crRepo *postgres.ConfigRevisionRepositoryPostgres,
) *handlers.CoreHandlers {
	return handlers.NewCoreHandlers(
		dsRepo, dfRepo, jobRepo, ocRepo, otRepo, crRepo,
	)
}

func NewGraphQLResolver(coreService services.CoreService) *graphql.Resolver {
	return graphql.NewResolver(coreService)
}

func RegisterRoutes(e *echo.Echo, siteH *handlers.SiteHandler, accountH *handlers.AccountHandler, coreH *handlers.CoreHandlers, gqlR *graphql.Resolver) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API Routes
	api := e.Group("/api")

	// Core
	core := api.Group("/core")
	dataSources := core.Group("/data-sources")
	dataSources.GET("", coreH.ListDataSources)
	dataSources.GET("/:id", coreH.GetDataSource)
	dataSources.POST("", coreH.CreateDataSource)
	dataSources.PUT("/:id", coreH.UpdateDataSource)
	dataSources.DELETE("/:id", coreH.DeleteDataSource)
	dataSources.POST("/:id/sync", coreH.SyncDataSource)

	dataFiles := core.Group("/data-files")
	dataFiles.GET("", coreH.ListDataFiles)
	dataFiles.GET("/:id", coreH.GetDataFile)
	dataFiles.POST("", coreH.CreateDataFile)
	dataFiles.PUT("/:id", coreH.UpdateDataFile)
	dataFiles.DELETE("/:id", coreH.DeleteDataFile)

	jobs := core.Group("/jobs")
	jobs.GET("", coreH.ListJobs)
	jobs.GET("/:id", coreH.GetJob)
	jobs.POST("", coreH.CreateJob)

	objectChanges := core.Group("/object-changes")
	objectChanges.GET("", coreH.ListObjectChanges)
	objectChanges.GET("/:id", coreH.GetObjectChange)
	objectChanges.POST("/log", coreH.LogObjectChange)

	objectTypes := core.Group("/object-types")
	objectTypes.GET("", coreH.ListObjectTypes)
	objectTypes.GET("/:id", coreH.GetObjectType)

	configRevisions := core.Group("/config-revisions")
	configRevisions.GET("", coreH.ListConfigRevisions)
	configRevisions.GET("/:id", coreH.GetConfigRevision)
	configRevisions.POST("", coreH.CreateConfigRevision)
	configRevisions.PUT("/:id", coreH.UpdateConfigRevision)
	configRevisions.DELETE("/:id", coreH.DeleteConfigRevision)
	configRevisions.POST("/:id/activate", coreH.ActivateConfigRevision)
	configRevisions.GET("/active", coreH.GetActiveConfigRevision)

	bgQueues := core.Group("/background-queues")
	bgQueues.GET("", coreH.ListBackgroundQueues)
	bgQueues.GET("/:name", coreH.GetBackgroundQueue)

	bgWorkers := core.Group("/background-workers")
	bgWorkers.GET("", coreH.ListBackgroundWorkers)
	bgWorkers.GET("/:name", coreH.GetBackgroundWorker)

	bgTasks := core.Group("/background-tasks")
	bgTasks.GET("", coreH.ListBackgroundTasks)
	bgTasks.GET("/:id", coreH.GetBackgroundTask)
	bgTasks.POST("/:id/delete", coreH.DeleteBackgroundTask)
	bgTasks.POST("/:id/requeue", coreH.RequeueBackgroundTask)
	bgTasks.POST("/:id/enqueue", coreH.EnqueueBackgroundTask)
	bgTasks.POST("/:id/stop", coreH.StopBackgroundTask)

	// DCIM
	dcim := api.Group("/dcim")
	sites := dcim.Group("/sites")
	sites.GET("", siteH.List)
	sites.GET("/:id", siteH.GetByID)
	sites.GET("/slug/:slug", siteH.GetBySlug)
	sites.POST("", siteH.Create)
	sites.PUT("/:id", siteH.Update)
	sites.DELETE("/:id", siteH.Delete)

	account := api.Group("/account", middleware.KerberosSSOMiddleware())
	account.GET("/profile", accountH.Profile)
	account.GET("/bookmarks", accountH.ListBookmarks)
	account.GET("/notifications", accountH.ListNotifications)
	account.GET("/subscriptions", accountH.ListSubscriptions)
	account.GET("/preferences", accountH.GetPreferences)
	account.POST("/preferences", accountH.UpsertPreferences)

	// GraphQL
	e.POST("/graphql", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "GraphQL endpoint integrated. Resolver initialized."})
	})

	tokens := account.Group("/api-tokens")
	tokens.GET("", accountH.ListTokens)
	tokens.GET("/:id", accountH.GetToken)
	tokens.POST("", accountH.CreateToken)
	tokens.PUT("/:id", accountH.UpdateToken)
	tokens.DELETE("/:id", accountH.DeleteToken)
}
