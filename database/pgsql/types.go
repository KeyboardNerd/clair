package pgsql

import (
	"time"

	"github.com/coreos/clair/database"
)

// Raw types in the database
type namespace struct {
	id            int64
	name          string
	versionFormat string
}

type feature struct {
	id            int64
	name          string
	version       string
	versionFormat string
}

type layer struct {
	id   int64
	hash string
}

type layerNamespace struct {
	id          int64
	layerID     int64
	namespaceID int64
}

type layerFeature struct {
	id        int64
	layerID   int64
	featureID int64
}

type layerLister struct {
	id      int64
	layerID int64
	lister  string
}

type layerDetector struct {
	id       int64
	layerID  int64
	detector string
}

type ancestry struct {
	id   int64
	name string
}

type ancestryLister struct {
	id         int64
	ancestryID int64
	lister     string
}

type ancestryDetector struct {
	id         int64
	ancestryID int64
	detector   string
}

type ancestryLayer struct {
	id            int64
	ancestryID    int64
	layerID       int64
	ancestryIndex int64
}

type namespacedFeature struct {
	id          int64
	featureID   int64
	namespaceID int64
}

type ancestryFeature struct {
	id                  int64
	ancestryLayerID     int64
	namespacedFeatureID int64
}

type vulnerability struct {
	id          int64
	namespaceID int64
	name        string
	description string
	link        string
	severity    database.Severity
}

type vulnerabilityAffectedFeature struct {
	id              int64
	vulnerabilityID int64
	featureName     string
	affectedVersion string
	fixedIn         string
}

type vulnerabilityAffectedNamespacedFeature struct {
	id                  int64
	vulnerabilityID     int64
	namespacedFeatureID int64
	addedBy             int64 // vulnerabilityAffectedFeature
}

type vulnerabilityNotification struct {
	id                 int64
	name               string
	createdAt          time.Time
	notifiedAt         time.Time
	deletedAt          time.Time
	oldVulnerabilityID int64
	newVulnerabilityID int64
}
