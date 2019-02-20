package dbtest_test

import (
	"testing"

	"github.com/coreos/clair/database/dbtest"
	_ "github.com/coreos/clair/database/pgsql"
	"github.com/stretchr/testify/require"
)

func TestCreateTestSession(t *testing.T) {
	storage, cleanup := dbtest.CreateTestDatabase("pgsql", false)
	defer cleanup()

	require.True(t, storage.Ping())
}
