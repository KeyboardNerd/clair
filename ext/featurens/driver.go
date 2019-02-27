// Copyright 2017 clair authors
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

// Package featurens exposes functions to dynamically register methods for
// determining a namespace for features present in an image layer.
package featurens

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/pkg/tarutil"
)

var (
	detectorsM sync.RWMutex
	detectors  = make(map[string]detector)
)

// Detector represents an ability to detect a namespace used for organizing
// features present in an image layer.
type Detector interface {
	// Detect attempts to determine a Namespace from a FilesMap of an image
	// layer.
	Detect(tarutil.FilesMap) (database.NamespaceDetectResult, error)

	// RequiredFilenames returns the list of files required to be in the FilesMap
	// provided to the Detect method.
	//
	// Filenames must not begin with "/".
	RequiredFilenames() []string
}

type detector struct {
	Detector

	info database.Detector
}

// RegisterDetector makes a detector available by the provided name.
//
// If called twice with the same name, the name is blank, or if the provided
// Detector is nil, this function panics.
func RegisterDetector(name string, version string, d Detector) {
	if name == "" || version == "" {
		panic("namespace: could not register a Detector with an empty name or version")
	}
	if d == nil {
		panic("namespace: could not register a nil Detector")
	}

	detectorsM.Lock()
	defer detectorsM.Unlock()

	if _, ok := detectors[name]; ok {
		panic("namespace: RegisterDetector called twice for " + name)
	}

	detectors[name] = detector{d, database.NewNamespaceDetector(name, version)}
}

// Detect uses detectors specified to retrieve the detect result.
func Detect(files tarutil.FilesMap, toUse []database.Detector) (database.DetectedNamespaces, error) {
	detectorsM.RLock()
	defer detectorsM.RUnlock()

	namespaces := make(database.DetectedNamespaces)
	for _, d := range toUse {
		// Only use the detector with the same type
		if d.DType != database.NamespaceDetectorType {
			continue
		}

		if detector, ok := detectors[d.Name]; ok {
			var err error
			if namespaces[d], err = detector.Detect(files); err != nil {
				return nil, err
			}
		} else {
			panic(fmt.Sprintf("unknown namespace detector: %#v", d))
		}
	}

	return namespaces, nil
}

// RequiredFilenames returns all files required by the give extensions. Any
// extension metadata that has non namespace-detector type will be skipped.
func RequiredFilenames(toUse []database.Detector) (files []string) {
	detectorsM.RLock()
	defer detectorsM.RUnlock()

	for _, d := range toUse {
		if d.DType != database.NamespaceDetectorType {
			continue
		}

		files = append(files, detectors[d.Name].RequiredFilenames()...)
	}

	return
}

// ListDetectors returns the info of all registered namespace detectors.
func ListDetectors() []database.Detector {
	r := make([]database.Detector, 0, len(detectors))
	for _, d := range detectors {
		r = append(r, d.info)
	}
	return r
}

// TestData represents the data used to test an implementation of Detector.
type TestData struct {
	Files             tarutil.FilesMap
	ExpectedNamespace *database.Namespace
}

// TestDetector runs a Detector on each provided instance of TestData and
// asserts the output to be equal to the expected output.
func TestDetector(t *testing.T, d Detector, testData []TestData) {
	for _, td := range testData {
		namespace, err := d.Detect(td.Files)
		assert.Nil(t, err)

		if namespace.Namespace == nil {
			assert.Equal(t, td.ExpectedNamespace, namespace)
		} else {
			assert.Equal(t, td.ExpectedNamespace.Name, namespace.Namespace.Name)
		}
	}
}
