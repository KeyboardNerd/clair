// v3_test is the integration tests for the gRPC API
package v3_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coreos/clair/api/v3"

	"github.com/coreos/clair"
	pb "github.com/coreos/clair/api/v3/clairpb"
	"github.com/coreos/clair/database/dbtest"

	_ "github.com/coreos/clair/database/pgsql"
	_ "github.com/coreos/clair/ext/featurefmt/apk"
	_ "github.com/coreos/clair/ext/featurefmt/dpkg"
	_ "github.com/coreos/clair/ext/featurefmt/rpm"
	_ "github.com/coreos/clair/ext/featurens/alpinerelease"
	_ "github.com/coreos/clair/ext/featurens/aptsources"
	_ "github.com/coreos/clair/ext/featurens/lsbrelease"
	_ "github.com/coreos/clair/ext/featurens/osrelease"
	_ "github.com/coreos/clair/ext/featurens/redhatrelease"
	_ "github.com/coreos/clair/ext/imagefmt/aci"
	_ "github.com/coreos/clair/ext/imagefmt/docker"
	_ "github.com/coreos/clair/ext/notification/webhook"
	_ "github.com/coreos/clair/ext/vulnmdsrc/nvd"
	_ "github.com/coreos/clair/ext/vulnsrc/alpine"
	_ "github.com/coreos/clair/ext/vulnsrc/debian"
	_ "github.com/coreos/clair/ext/vulnsrc/oracle"
	_ "github.com/coreos/clair/ext/vulnsrc/rhel"
	_ "github.com/coreos/clair/ext/vulnsrc/suse"
	_ "github.com/coreos/clair/ext/vulnsrc/ubuntu"
)

var (
	status       *v3.StatusServer
	ancestry     *v3.AncestryServer
	notification *v3.NotificationServer
)

func TestMain(m *testing.M) {
	store, cleanup := dbtest.CreateTestDatabase("pgsql", true)
	defer cleanup()
	clair.InitWorker(store)
	status = &v3.StatusServer{store}
	ancestry = &v3.AncestryServer{store}
	notification = &v3.NotificationServer{store}
	m.Run()
}

func TestGetStatus(t *testing.T) {
	_, err := status.GetStatus(context.Background(), &pb.GetStatusRequest{})
	require.Nil(t, err)
}

func TestPostAncestry(t *testing.T) {

}

func TestGetAncestry(t *testing.T) {

}

func TestGetNotification(t *testing.T) {

}

func TestMarkNotificationAsRead(t *testing.T) {

}
