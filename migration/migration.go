package migration

import (
	"regexp"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

var reNotInit = regexp.MustCompile(`relation "[^"]+" does not exist`)

type Migration struct {
	db   *pg.DB
	col  *migrations.Collection
	name string
}

func NewMigration(db *pg.DB, col *migrations.Collection, name string) *Migration {
	return &Migration{db: db, col: col, name: name}
}

func (m *Migration) Run(cmd string) error {
	if _, err := m.col.Version(m.db); err != nil {
		if !reNotInit.MatchString(err.Error()) {
			log.Err(err).Msgf("[%s] [version] failed", m.name)
			return err
		}
		// run migrate init
		if _, _, err := m.col.Run(m.db, "init"); err != nil {
			log.Err(err).Msgf("[%s] [init] failed", m.name)
			return err
		}
	}

	o, n, err := m.col.Run(m.db, cmd)
	if err != nil {
		log.Err(err).Msgf("[%s] [%s] failed", m.name, cmd)
		return err
	} else {
		log.Info().Msgf("[%s] [%s] from version %v to version %v successfully", m.name, cmd, o, n)
		return nil
	}
}

func (m *Migration) Up() error {
	return m.Run("up")
}

func (m *Migration) Down() error {
	return m.Run("down")
}

func (m *Migration) Reset() error {
	return m.Run("reset")
}
