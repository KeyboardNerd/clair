package dbtest

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/coreos/clair/database/pgsql"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/pkg/pagination"
)

var pgsqlTestConfig = database.RegistrableComponentConfig{
	Type: "pgsql",
	Options: map[string]interface{}{
		"source": 
		"paginationkey": pagination.Must(pagination.NewKey()).String(),
	},
}

func d(){
	if sourceEnv := ; sourceEnv != ""{
		source = sourceEnv
	}


	source := fmt.Sprintf("postgresql://postgres@127.0.0.1:5432/%s?sslmode=disable", dbName)
	if sourceEnv := os.Getenv("CLAIR_TEST_PGSQL"); sourceEnv != "" {
		source = fmt.Sprintf(sourceEnv, dbName)
	}
}

// CreateTestDatabase creates a test database
func CreateTestDatabase(dbType string, loadFixture bool) (database.Datastore, func()) {
	var config database.RegistrableComponentConfig
	switch dbType {
	case "pgsql":
		config = pgsqlTestConfig
	default:
		panic(fmt.Sprintf("Missing test configuration for %s", dbType))
	}

	store, err := database.Open(config)
	if err != nil {
		panic(err)
	}

	if !store.Ping() {
		store.Close()
		panic("Cannot reach test storage")
	}

	if loadFixture {
		switch dbType {
		case "pgsql":
			_, filename, _, _ := runtime.Caller(0)
			fixturePath := filepath.Join(filepath.Dir(filename)) + "/testdata/data.sql"
			if err := pgsql.LoadFixture(store, fixturePath); err != nil {
				store.Close()
				panic("Failed to load pgsql test fixture")
			}
		}
	}

	return store, func() {
		store.Close()
	}
}
