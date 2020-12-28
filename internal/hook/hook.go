package hook

import (
	"github.com/go-pg/pg"
	"github.com/nori-io/common/v4/pkg/domain/logger"
)

type DbLogger struct {
	origin logger.FieldLogger
}

func (d DbLogger) BeforeQuery(q *pg.QueryEvent) {
	r, err := q.FormattedQuery()
	if err != nil {
		d.origin.Error(err.Error())
	} else {
		d.origin.Debug("%s", r)
	}
}

func (d DbLogger) AfterQuery(q *pg.QueryEvent) {
	r, err := q.FormattedQuery()
	if err != nil {
		d.origin.Error(err.Error())
	} else {
		d.origin.Debug("%s", r)
	}
}
