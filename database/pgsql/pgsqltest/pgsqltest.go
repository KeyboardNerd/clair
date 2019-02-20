package pgsqltest

import (
	"io/ioutil"

	"github.com/coreos/clair/database"
)

func LoadFixture(store database.Datastore, fixturePath string) error {
	storage := store.(*pgSQL)
	d, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		return err
	}

	_, err = storage.Exec(string(d))
	if err != nil {
		return err
	}

	return nil
}
