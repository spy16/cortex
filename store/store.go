package store

import (
	log "github.com/sirupsen/logrus"

	"github.com/chunked-app/cortex/store/pgsql"
	"github.com/chunked-app/cortex/user"
)

var (
	_ block.Store = (*InMemory)(nil)
	_ block.Store = (*pgsql.PostgresQL)(nil)
)

func Open(spec string) (block.Store, user.Store, error) {
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
