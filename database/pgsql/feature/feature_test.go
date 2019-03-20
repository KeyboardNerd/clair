// Copyright 2016 clair authors
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/database/pgsql/feature"
	"github.com/coreos/clair/database/pgsql/testutil"
)

func TestPersistFeatures(t *testing.T) {
	tx, cleanup := testutil.CreateTestTx(t, "TestPersistFeatures")
	defer cleanup()

	invalid := database.Feature{}
	valid := *database.NewBinaryPackage("mount", "2.31.1-0.4ubuntu3.1", "dpkg")

	// invalid
	require.NotNil(t, feature.PersistFeatures(tx, []database.Feature{invalid}))
	// existing
	require.Nil(t, feature.PersistFeatures(tx, []database.Feature{valid}))
	require.Nil(t, feature.PersistFeatures(tx, []database.Feature{valid}))

	features := testutil.SelectAllFeatures(t, tx)
	assert.Equal(t, []database.Feature{valid}, features)
}

func TestPersistNamespacedFeatures(t *testing.T) {
	tx, cleanup := testutil.CreateTestTxWithFixtures(t, "TestPersistNamespacedFeatures")
	defer cleanup()

	// existing features
	f1 := database.NewSourcePackage("ourchat", "0.5", "dpkg")
	// non-existing features
	f2 := database.NewSourcePackage("fake!", "", "")
	// exising namespace
	n1 := database.NewNamespace("debian", "7", "dpkg")
	// non-existing namespace
	n2 := database.NewNamespace("debian", "non", "dpkg")
	// existing namespaced feature
	nf1 := database.NewNamespacedFeature(n1, f1)
	// invalid namespaced feature
	nf2 := database.NewNamespacedFeature(n2, f2)
	// namespaced features with namespaces or features not in the database will
	// generate error.
	assert.Nil(t, feature.PersistNamespacedFeatures(tx, []database.NamespacedFeature{}))
	assert.NotNil(t, feature.PersistNamespacedFeatures(tx, []database.NamespacedFeature{*nf1, *nf2}))
	// valid case: insert nf3
	assert.Nil(t, feature.PersistNamespacedFeatures(tx, []database.NamespacedFeature{*nf1}))

	all := testutil.ListNamespacedFeatures(t, tx)
	assert.Contains(t, all, *nf1)
}
