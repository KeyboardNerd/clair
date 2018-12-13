package featurefmt

import (
	"github.com/deckarep/golang-set"

	"github.com/coreos/clair/database"
)

// Package is some package from a package manager
type Package struct {
	Name          string
	Version       string
	SourceName    string
	SourceVersion string
}

// NewPackage creates a new package
func NewPackage(Name string, Version string, SourceName string, SourceVersion string) *Package {
	return &Package{Name, Version, SourceName, SourceVersion}
}

// PackageManager represents the information about a package manager and its
// associated packages.
type PackageManager struct {
	Packages           mapset.Set
	PackageManagerName string
	VersionFormat      string
}

// NewPackageManager initializes a new package set to contain the scan result.
func NewPackageManager(packageManagerName string, versionFormat string) *PackageManager {
	return &PackageManager{
		Packages:           mapset.NewSet(),
		PackageManagerName: packageManagerName,
		VersionFormat:      versionFormat,
	}
}

// Add a package to package set
func (s *PackageManager) Add(p *Package) {
	s.Add(p)
}

// PackageSet is the scan result of a lister or all listers that can be merged
// together.
type PackageSet struct {
	content map[database.Detector]PackageManager
}

// NewPackageSet creates a package set
func NewPackageSet() *PackageSet {
	return &PackageSet{content: make(map[database.Detector]PackageManager)}
}

// LayerFeatures computes the layer features based on the package information
func (s *PackageSet) LayerFeatures() (features []database.LayerFeature) {
	for detector, content := range s.content {
		for rawPackage := range content.Packages.Iter() {
			p := rawPackage.(Package)
			detectorCopy := detector
			features = append(features, *database.NewLayerFeature(
				p.Name,
				p.Version,
				p.SourceName,
				p.SourceVersion,
				content.VersionFormat,
				&detectorCopy,
			))
		}
	}

	return
}

// Merge two different package sets.
func (s *PackageSet) Merge(p *PackageSet) {
	if p == nil {
		return
	}

	for detector, newContent := range p.content {
		if existingContent, ok := s.content[detector]; ok {
			if existingContent.VersionFormat != newContent.VersionFormat || existingContent.PackageManagerName != newContent.PackageManagerName {
				panic("Detector  detects multiple package managers but should only detect one package manager.")
			}
			existingContent.Packages.Union(newContent.Packages)
		} else {
			s.content[detector] = newContent
		}
	}
}
