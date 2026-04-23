package app

import (
	"database/sql"

	"netbox_go/internal/repository/postgres"

	"go.uber.org/fx"
)

// ModuleRepository provides repository dependencies
var ModuleRepository = fx.Options(
	fx.Provide(
		NewSiteRepository,
		NewAccountRepository,
		NewDataSourceRepository,
		NewDataFileRepository,
		NewJobRepository,
		NewObjectChangeRepository,
		NewObjectTypeRepository,
		NewConfigRevisionRepository,
	),
)

func NewSiteRepository(db *sql.DB) *postgres.SiteRepositoryPostgres {
	return postgres.NewSiteRepositoryPostgres(db)
}

func NewAccountRepository(db *sql.DB) *postgres.AccountRepositoryPostgres {
	return postgres.NewAccountRepositoryPostgres(db)
}

func NewDataSourceRepository(db *sql.DB) *postgres.DataSourceRepositoryPostgres {
	return postgres.NewDataSourceRepositoryPostgres(db)
}

func NewDataFileRepository(db *sql.DB) *postgres.DataFileRepositoryPostgres {
	return postgres.NewDataFileRepositoryPostgres(db)
}

func NewJobRepository(db *sql.DB) *postgres.JobRepositoryPostgres {
	return postgres.NewJobRepositoryPostgres(db)
}

func NewObjectChangeRepository(db *sql.DB) *postgres.ObjectChangeRepositoryPostgres {
	return postgres.NewObjectChangeRepositoryPostgres(db)
}

func NewObjectTypeRepository(db *sql.DB) *postgres.ObjectTypeRepositoryPostgres {
	return postgres.NewObjectTypeRepositoryPostgres(db)
}

func NewConfigRevisionRepository(db *sql.DB) *postgres.ConfigRevisionRepositoryPostgres {
	return postgres.NewConfigRevisionRepositoryPostgres(db)
}
