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

// Package aptsources implements a featurens.Detector for apt based container
// image layers.
//
// This detector is necessary to determine the precise Debian version when it
// is an unstable version for instance.
package aptsources

import (
	"bufio"
	"strings"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/ext/featurens"
	"github.com/coreos/clair/ext/versionfmt/dpkg"
	"github.com/coreos/clair/pkg/tarutil"
)

type detector struct{}

func init() {
	featurens.RegisterDetector("apt-sources", "1.0", &detector{})
}

func (d detector) Detect(files tarutil.FilesMap) (database.NamespaceDetectResult, error) {
	detectResult := database.NamespaceDetectResult{Status: database.NotFound}
	f, hasFile := files["etc/apt/sources.list"]
	if !hasFile {
		return detectResult, nil
	}

	detectResult.Status = database.Changed
	var OS, version string

	scanner := bufio.NewScanner(strings.NewReader(string(f)))
	for scanner.Scan() {
		// Format: man sources.list | https://wiki.debian.org/SourcesList)
		// deb uri distribution component1 component2 component3
		// deb-src uri distribution component1 component2 component3
		line := strings.Split(scanner.Text(), " ")
		if len(line) > 3 {
			// Only consider main component
			isMainComponent := false
			for _, component := range line[3:] {
				if component == "main" {
					isMainComponent = true
					break
				}
			}
			if !isMainComponent {
				continue
			}

			var found bool
			version, found = database.DebianReleasesMapping[line[2]]
			if found {
				OS = "debian"
				break
			}

			line[2] = strings.Split(line[2], "/")[0]
			version, found = database.UbuntuReleasesMapping[line[2]]
			if found {
				OS = "ubuntu"
				break
			}
		}
	}

	if OS != "" && version != "" {
		detectResult.Namespace = &database.Namespace{
			Name:          OS + ":" + version,
			VersionFormat: dpkg.ParserName,
		}
	}

	return detectResult, nil
}

func (d detector) RequiredFilenames() []string {
	return []string{"etc/apt/sources.list"}
}
