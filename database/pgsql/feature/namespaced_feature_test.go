// Copyright 2019 clair authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package feature_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/database/pgsql/feature"
	"github.com/coreos/clair/database/pgsql/testutil"
)

func TestFindNamespacedFeatureIDs(t *testing.T) {
	tx, cleanup := testutil.CreateTestTxWithFixtures(t, "TestFindNamespacedFeatureIDs")
	defer cleanup()

	ns := []database.NamespacedFeature{
		testutil.RealNamespacedFeatures[1],
		testutil.RealNamespacedFeatures[2],
		testutil.RealNamespacedFeatures[5],
		testutil.FakeNamespacedFeatures[1],
	}

	ids, err := feature.FindNamespacedFeatureIDs(tx, ns)
	require.Nil(t, err)
	require.True(t, ids[0].Valid && ids[0].Int64 == 1)
	require.True(t, ids[1].Valid && ids[1].Int64 == 2)
	require.True(t, ids[2].Valid)
	require.True(t, ids[2].Int64 == 5)
	require.True(t, !ids[3].Valid)
}
