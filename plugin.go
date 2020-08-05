package main

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/nori-io/common/v3/config"
	"github.com/nori-io/common/v3/logger"
	"github.com/nori-io/common/v3/meta"
	"github.com/nori-io/common/v3/plugin"
	i"github.com/nori-io/interfaces/public/sql/pg"
)

type service struct {
	db *pg.DB
	config   *pluginConfig

}

//dialect - mysql или postgres.
//one parameter with format  host:port
type pluginConfig struct {
	addr string
	db string
	user string
	password string
	dialect string
}

var (
	Plugin plugin.Plugin = &service{}
)

func (p *service) Init(ctx context.Context, config config.Config, log logger.FieldLogger) error {
	p.config.addr=config.String("host", "host")()
	p.config.db=config.String("db", "db")()
	p.config.user=config.String("user", "user")()
	p.config.password=config.String("password", "password")()
	p.config.dialect=config.String("dialect", "dialect")()
	return nil
}

func (p *service) Instance() interface{} {
	return p.db
}

func (p *service) Meta() meta.Meta {
	return &meta.Data{
		ID: meta.ID{
			ID:      "public/sql/pg",
			Version: "8.0.7",
		},
		Author: meta.Author{
			Name: "Nori.io",
			URI:  "https://nori.io/",
		},
		Core: meta.Core{
			VersionConstraint: "=0.2.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name:        "Nori: ORM PG",
			Description: "This plugin implements instance of ORG PG",
		},
		Interface: i.PgInterface,
		License: []meta.License{
		{
				Title: "GPLv3",
				Type:  "GPLv3",
				URI:   "https://www.gnu.org/licenses/"},
		},
		Tags: []string{"orm", "pg", "sql"},
	}

}

func (p *service) Start(ctx context.Context, registry plugin.Registry) error {
	if p.db == nil {
			p.db= &pg.DB{}
	}
	p.db= pg.Connect(&pg.Options{
		Addr:                  p.config.addr,
		User:                  p.config.user,
		Password:              p.config.password,
		Database:              p.config.db,

	})

	return nil
}

func (p *service) Stop(ctx context.Context, registry plugin.Registry) error {
	p.db.Close()
	p.db = nil
	return nil
}
