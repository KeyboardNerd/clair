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

package database

// Namespace is the contextual information around features.
//
// e.g. Name = Debian, Version = 7.0, VersionFormat = dpkg
type Namespace struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	VersionFormat string `json:"versionFormat"`
}

// NewNamespace creates a new namespace.
func NewNamespace(name string, version string, versionFormat string) *Namespace {
	return &Namespace{name, version, versionFormat}
}

// Valid checks if the namespace is valid.
func (ns *Namespace) Valid() bool {
	if ns == nil || ns.Name == "" || ns.VersionFormat == "" || ns.Version == "" {
		return false
	}

	return true
}
