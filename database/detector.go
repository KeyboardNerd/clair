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
	// Name of a detector should be non-empty and uniquely identifies the
	// detector.
	Name string
	// Version of a detector should be non-empty.
	Version string
	// Type of a detector should be one of the types in DetectorTypes.
	Type DetectorType
}

// Valid checks if all fields in the detector satisfies the spec.
func (d Detector) Valid() bool {
	if d.Name == "" && d.Version == "" {
		return false
	}

	for _, t := range DetectorTypes {
		if d.Type == t {
			return true
		}
	}

	return false
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
