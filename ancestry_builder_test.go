package clair

import (
	"testing"

	"github.com/coreos/clair/database"
	"github.com/stretchr/testify/require"
)

// Testing
// 1. Adding feature, removing feature based on existence of feature in current
// layer.
// 2. Adding Namespace, Removing Namespace.
// 3. Looking up Feature with namespace in current layer, in
// parent layers, handling of case namespace cannot be determined.
// 4. Relation: Relating Feature to current layer's, Relating Feature to parent layer's.
// 5. Relate Features with correct Ancestry Layer.
//
// Not supported Yet:
// 1. Explicitly removing packages,
// 2. Package indicated namespace

var testDetectors = []database.Detector{
	database.NewFeatureDetector("FeatureDetector1", "1.0"),
	database.NewNamespaceDetector("NSDetector1", "1.0"),
}

var testFeatures = []*database.Feature{
	NewFeature("feature1", "1.0", "dpkg", database.SourcePackage),
}

func NewFeature(name string, version string, versionfmt string, featureType database.FeatureType) *database.Feature {
	return &database.Feature{name, version, versionfmt, featureType}
}

func TestAddLayer(t *testing.T) {
	cases := []struct {
		title    string
		baseCase *ancestryBuilder
		layer    database.Layer

		expected          *ancestryBuilder
		expectedErrString string
	}{
		{
			title:    "empty ancestry, empty layer",
			baseCase: newAncestryBuilder("test"),
			expected: newAncestryBuilder("test"),
		},
		{
			"empty ancestry, layer with 1 feature",
			newAncestryBuilder("test"),
			NewLayerBuilder("test", nil).AddFeature(testDetectors[0], testFeatures[0]).Layer(),
		},
		{
			"Multiple Detectors: Ubuntu installs 2 features",
		},
		{
			"Multiple Detectors: Ubuntu installs 2 features, and RPM installs 2 features, both under the same namespace",
		},
		{
			"Upgrade Namespace from CentOS:7 to CentOS:8",
		},
		{
			"Downgrade Namespace from CentOS:8 to CentOS:7",
		},
		{
			"Multiple Namespace: Python + Ubuntu package",
		},
		{
			"Multiple Namespace: CentOS installs RHEL package",
		},
		{
			"Multiple Namespace: Javascript + Python + Ubuntu packages",
		},
		{
			"Multiple Namespace: A feature is under Python + Ubuntu namespace",
		},
		{
			"Relational Feature: Multiple packages with same source package",
		},
		{
			"Relational Feature: Package with source package in parent layer",
		},
	}

	for _, test := range cases {
		t.Run(title, func(t *testing.T) {
			err := test.baseCase.Add(test.layer)
			if expectedErrString != "" {
				require.EqualError(t, err, expectedErrString)
				return
			}

			require.Empty(t, err)
			assertAncestryBuilderEqual(t, baseCase, expected)
		})
	}
}

func assertAncestryBuilderEqual(t *testing.T, actual, expected ancestryBuilder) {

}
