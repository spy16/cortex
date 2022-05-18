package store

import (
	log "github.com/sirupsen/logrus"

	"github.com/chunked-app/cortex/chunk"
	"github.com/chunked-app/cortex/store/pgsql"
	"github.com/chunked-app/cortex/user"
)

var (
	_ chunk.Store = (*InMemory)(nil)
	_ chunk.Store = (*pgsql.PostgresQL)(nil)
)

func Open(spec string) (chunk.Store, user.Store, error) {
	if spec == ":memory:" {
		m := &InMemory{}
		log.Warnf("using in-memory database")
		return m, m, nil
	}

	log.Infof("connecting to postgres using '%s'", spec)
	s, err := pgsql.Open(spec)
	if err != nil {
		return nil, nil, err
	}
	return s, s, err
}
