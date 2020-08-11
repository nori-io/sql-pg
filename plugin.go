package main

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/nori-io/common/v3/config"
	"github.com/nori-io/common/v3/logger"
	"github.com/nori-io/common/v3/meta"
	"github.com/nori-io/common/v3/plugin"
	i "github.com/nori-io/interfaces/public/sql/pg"
	"github.com/nori-io/sql-pg/internal/hook"
)

type service struct {
	db       *pg.DB
	config   *pluginConfig
	logger   logger.FieldLogger
	dbLogger *hook.DbLogger
}

type pluginConfig struct {
	addr     string
	db       string
	user     string
	password string
}

var (
	Plugin plugin.Plugin = &service{}
)

func (p *service) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.logger = log
	p.dbLogger = &hook.DbLogger{}
	p.config.addr = config.String("sql.pg.addr", "addr")()
	p.config.db = config.String("sql.pg.db", "database name")()
	p.config.user = config.String("sql.pg.user", "user")()
	p.config.password = config.String("sql.pg.password", "password")()
	return nil
}

func (p *service) Instance() interface{} {
	return p.db
}

func (p *service) Meta() meta.Meta {
	return meta.Data{
		ID: meta.ID{
			ID:      "sql/pg",
			Version: "8.0.7",
		},
		Author: meta.Author{
			Name: "Nori.io",
			URI:  "https://nori.io/",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name:        "Nori: ORM PG",
			Description: "This plugin implements instance of ORM PG",
		},
		Core: meta.Core{
			VersionConstraint: "^0.2.0",
		},
		Interface: i.PgInterface,
		License: []meta.License{
			{
				Title: "GPLv3",
				Type:  "GPLv3",
				URI:   "https://www.gnu.org/licenses/"},
		},
		Links: []meta.Link{},
		Repository: meta.Repository{
			Type: "git",
			URI:  "https://github.com/nori-io/sql-pg",
		},
		Tags: []string{"orm", "pg", "sql"},
	}

}

func (p *service) Start(ctx context.Context, registry plugin.Registry) error {

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

func (p *service) Stop(ctx context.Context, registry plugin.Registry) error {
	err := p.db.Close()

	if err != nil {
		p.logger.Error(err.Error())
	}

	return err
}
