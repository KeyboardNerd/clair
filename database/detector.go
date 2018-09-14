// Copyright 2018 clair authors
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

package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// ErrFailedToParseDetectorType is the error returned when a detector type could
// not be parsed from a string.
var ErrFailedToParseDetectorType = errors.New("failed to parse DetectorType from input")

// DetectorType is the type of a detector.
type DetectorType string

// Value implements the database/sql/driver.Valuer interface.
func (s DetectorType) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan implements the database/sql.Scanner interface.
func (s *DetectorType) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New("could not scan a Severity from a non-string input")
	}

	var err error
	*s, err = NewDetectorType(string(val))
	if err != nil {
		return err
	}

	return nil
}

// NewDetectorType attempts to parse a string into a standard DetectorType
// value.
func NewDetectorType(s string) (DetectorType, error) {
	for _, ss := range DetectorTypes {
		if strings.EqualFold(s, string(ss)) {
			return ss, nil
		}
	}

	return "", ErrFailedToParseSeverity
}

// Valid checks if a detector type is defined.
func (s DetectorType) Valid() bool {
	for _, t := range DetectorTypes {
		if s == t {
			return true
		}
	}

	return false
}

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
	if d.Name == "" || d.Version == "" || !d.Type.Valid() {
		return false
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
