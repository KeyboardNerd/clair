package database

import (
	"fmt"
)

// DetectorType is the type of a detector.
type DetectorType string

const (
	// NamespaceType is a type of detector that extracts the namespaces.
	NamespaceType DetectorType = "Namespace"
	// FeatureType is a type of detector that extracts the features.
	FeatureType DetectorType = "Feature"
)

// DetectorTypes contains all detector types.
var DetectorTypes = []DetectorType{
	NamespaceType,
	FeatureType,
}

// Detector is an extention to scan a layer's content.
type Detector struct {
	// Name of the detector
	Name string
	// Version of the detector
	Version string
	// Type of the detector
	Type DetectorType
}

func (d Detector) String() string {
	return fmt.Sprintf("%sDetector/%s/%s", d.Type, d.Name, d.Version)
}

// NewNamespaceDetector returns a new namespace detector.
func NewNamespaceDetector(name string, version string) Detector {
	return Detector{
		Name:    name,
		Version: version,
		Type:    NamespaceType,
	}
}

// NewFeatureDetector returns a new feature detector.
func NewFeatureDetector(name string, version string) Detector {
	return Detector{
		Name:    name,
		Version: version,
		Type:    FeatureType,
	}
}
