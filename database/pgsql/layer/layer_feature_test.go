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

package layer_test

import (
	"testing"

	"github.com/coreos/clair/database"

	"github.com/stretchr/testify/require"

	"github.com/coreos/clair/database/pgsql/detector"

	"github.com/coreos/clair/database/pgsql/layer"
	"github.com/coreos/clair/database/pgsql/testutil"
)

func TestFindLayerFeatures(t *testing.T) {
	tx, cleanup := testutil.CreateTestTxWithFixtures(t, "TestFindLayerFeatures")
	defer cleanup()

	m, err := detector.FindAllDetectors(tx)
	require.Nil(t, err)
	features, err := layer.FindLayerFeatures(tx, 6, m)
	require.Nil(t, err)
	database.AssertLayerFeaturesEqual(t, testutil.RealLayers[6].Features, features)
}
