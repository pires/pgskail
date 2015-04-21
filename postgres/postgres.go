package postgres

import (
	"log"

	"github.com/pires/pgskail/util"
)

const (
	PG_DIR = "/var/lib/data/postgres"
)

func Start() {
	// is data dir empty?
	if isDirEmpty, err := util.IsDirEmpty(PG_DIR); err != nil {
		log.Fatal("There was an error while validating PostgreSQL data directory", err)
	} else {
		if isDirEmpty {
			// TODO race to set initialization key
			hasInitializationKey := false
			// has initilization key?
			if hasInitializationKey {
				// initialize database
			} else {
				// TODO wait for leader to have initialization key
				// TODO pg_basebackup from leader
			}
		} // else?

		// TODO start PostgreSQL
	}
}
