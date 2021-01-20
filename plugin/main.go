package main

import (
	"context"

	"github.com/nori-io/common/v4/pkg/domain/registry"

	"github.com/go-pg/pg"
	"github.com/nori-io/common/v4/pkg/domain/config"
	em "github.com/nori-io/common/v4/pkg/domain/enum/meta"
	"github.com/nori-io/common/v4/pkg/domain/logger"
	"github.com/nori-io/common/v4/pkg/domain/meta"
	p "github.com/nori-io/common/v4/pkg/domain/plugin"
	m "github.com/nori-io/common/v4/pkg/meta"
	i "github.com/nori-io/interfaces/database/orm/pg"
	"github.com/nori-plugins/database-orm-pg/internal/hook"
)

func New() p.Plugin {
	return &plugin{}
}

type plugin struct {
	db       *pg.DB
	config   conf
	logger   logger.FieldLogger
	dbLogger *hook.DbLogger
}

type conf struct {
	addr     string
	db       string
	user     string
	password string
}

func (p plugin) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.logger = log
	p.dbLogger = &hook.DbLogger{}
	p.config.addr = config.String("sql.pg.addr", "addr")()
	p.config.db = config.String("sql.pg.db", "database name")()
	p.config.user = config.String("sql.pg.user", "user")()
	p.config.password = config.String("sql.pg.password", "password")()
	return nil
}

func (p plugin) Instance() interface{} {
	return p.db
}

func (p plugin) Meta() meta.Meta {
	return m.Meta{
		ID: m.ID{
			ID:      "sql/pg",
			Version: "8.0.7",
		},
		Author: m.Author{
			Name: "Nori.io",
			URL:  "https://nori.io/",
		},
		Dependencies: []meta.Dependency{},
		Description: m.Description{
			Title:       "",
			Description: "This plugin implements instance of ORM PG",
		},
		Interface: i.PgInterface,
		License:   []meta.License{},
		Links:     []meta.Link{},
		Repository: m.Repository{
			Type: em.Git,
			URL:  "https://github.com/nori-io/sql-pg",
		},
		Tags: []string{"orm", "pg", "sql"},
	}
}

func (p plugin) Start(ctx context.Context, registry registry.Registry) error {
	p.db = pg.Connect(&pg.Options{
		Addr:     p.config.addr,
		User:     p.config.user,
		Password: p.config.password,
		Database: p.config.db,
	})

	var n int
	_, err := p.db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		p.logger.Error(err.Error())
	} else {
		p.db.AddQueryHook(p.dbLogger)
	}

	return err
}

func (p plugin) Stop(ctx context.Context, registry registry.Registry) error {
	err := p.db.Close()
	if err != nil {
		p.logger.Error(err.Error())
	}

	return err
}
