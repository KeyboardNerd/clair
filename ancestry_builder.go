package clair

import (
	"fmt"

	"github.com/deckarep/golang-set"

	"github.com/coreos/clair/database"
)

type ancestryBuilder struct {
	ancestryName string
	layerIndex   int
	layerNames   []string
	detectors    []database.Detector
	namespaces   map[database.Detector]layerIndexedNamespace
	features     map[database.Detector][]layerIndexedFeature
}

type layerIndexedFeature struct {
	feature      database.Feature
	namespaceKey database.Detector
	introducedIn int
}

type layerIndexedNamespace struct {
	namespace    database.Namespace
	introducedIn int
}

// newAncestryBuilder creates a new ancestry builder
//
// ancestry builder takes in the extracted layer information and produce a set of
// namespaces, features, and the relation between features for the whole image.
func newAncestryBuilder(name string) *ancestryBuilder {
	return &ancestryBuilder{
		layerIndex: 0,
		namespaces: make(map[database.Detector]layerIndexedNamespace),
		features:   make(map[database.Detector][]layerIndexedFeature),
	}
}

// TODO(sidac): refactor database.Layer to use the map to store detector and
// feature relationship
func groupFeaturesByDetector(layer *database.Layer) map[database.Detector][]database.Feature {
	detected := map[database.Detector][]database.Feature{}
	for _, feature := range layer.Features {
		if features, ok := detected[feature.By]; ok {
			features = append(features, feature.Feature)
		} else {
			detected[feature.By] = []database.Feature{feature.Feature}
		}
	}

	return detected
}

// AddLeafLayer adds a leaf layer to the ancestry.
//
// The Add function gives each feature a namespace and the layer that introduces
// the feature.
// It requires all layers added to the ancestry to have the same set of
// detectors, otherwise, it'll fail.
func (b *ancestryBuilder) AddLeafLayer(layer database.Layer) error {
	// TODO(sidac): update when layer has mapping from detector to content
	if b.layerIndex == 0 {
		b.detectors = layer.By
	}

	// update the namespace detected by the same detector
	for _, layerNamespace := range layer.Namespaces {
		b.addNamespace(&layerNamespace)
	}

	for detector, features := range groupFeaturesByDetector(&layer) {
		b.addLayerFeatures(detector, features)
	}

	b.layerIndex++
	return nil
}

func (b *ancestryBuilder) addLayerFeatures(detector database.Detector, features []database.Feature) {
	var (
		existingFeatures []layerIndexedFeature
		foundExisting    bool
	)

	if existingFeatures, foundExisting = b.features[detector]; !foundExisting {
		existingFeatures = []layerIndexedFeature{}
		b.features[detector] = existingFeatures
	}

	// Remove features detected by the same detector but does not exist in the
	// current layer.
	i := 0
	isOldFeature := mapset.NewSet()
	for i < len(existingFeatures) {
		foundExisting = false
		for j, feature := range features {
			if feature == existingFeatures[i].feature {
				foundExisting = true
				isOldFeature.Add(j)
				break
			}
		}

		if !foundExisting {
			existingFeatures = append(existingFeatures[:i], existingFeatures[i+1:]...)
		}

		i++
	}

	// Add features that has not introduced by any previous layers.
	for i, feature := range features {
		if !isOldFeature.Contains(i) {
			key, ok := b.lookupNamespaceKey(&detector, &feature)
			if !ok {
				// By best effort, skip if the feature can't be related with a
				// namespace.
				continue
			}

			existingFeatures = append(existingFeatures, layerIndexedFeature{
				feature:      feature,
				namespaceKey: key,
				introducedIn: b.layerIndex,
			})
		}
	}

	b.features[detector] = existingFeatures
}

func (b *ancestryBuilder) addNamespace(layerNamespace *database.LayerNamespace) {
	if previous, ok := b.namespaces[layerNamespace.By]; !ok || previous.namespace != layerNamespace.Namespace {
		b.namespaces[layerNamespace.By] = layerIndexedNamespace{
			layerNamespace.Namespace,
			b.layerIndex,
		}
		// Features refer namespaces by the detector key, and therefore, their
		// referred namespaces will be automatically updated.
	}
}

func (b *ancestryBuilder) lookupNamespaceKey(detector *database.Detector, feature *database.Feature) (database.Detector, bool) {
	// TODO(sidac): Looking up namespace is a best effort function.
	// The look up function first try to determine the namespace by the
	// given feature's raw context, then try to relate the type of the feature
	// to OS namespace, or language level namespace.
	for detector, namespace := range b.namespaces {
		if namespace.namespace.VersionFormat == feature.VersionFormat {
			return detector, true
		}
	}

	return database.Detector{}, false
}

func (b *ancestryBuilder) ancestryFeatures(index int) []database.AncestryFeature {
	ancestryFeatures := []database.AncestryFeature{}
	for detector, features := range b.features {
		for _, feature := range features {
			if feature.introducedIn == index {
				nsDetector := feature.namespaceKey
				if _, ok := b.namespaces[nsDetector]; !ok {
					panic(fmt.Sprintf("namespace detector is missing %s in the ancestry builder", nsDetector.Name))
				}

				namespace := b.namespaces[nsDetector]
				ancestryFeatures = append(ancestryFeatures, database.AncestryFeature{
					NamespacedFeature: database.NamespacedFeature{
						Feature:   feature.feature,
						Namespace: namespace.namespace,
					},
					FeatureBy:   detector,
					NamespaceBy: nsDetector,
				})
			}
		}
	}

	return ancestryFeatures
}

func (b *ancestryBuilder) ancestryLayers() []database.AncestryLayer {
	layers := make([]database.AncestryLayer, 0, b.layerIndex)
	for i := 0; i < b.layerIndex; i++ {
		layers = append(layers, database.AncestryLayer{
			Hash:     b.layerNames[i],
			Features: b.ancestryFeatures(i),
		})
	}

	return layers
}

func (b *ancestryBuilder) Ancestry() *database.Ancestry {
	return &database.Ancestry{
		Name:   b.ancestryName,
		By:     b.detectors,
		Layers: b.ancestryLayers(),
	}
}

// computeAncestryLayers computes ancestry's layers along with what features are
// introduced.
func computeAncestryLayers(layers []database.Layer) ([]database.AncestryLayer, []database.Detector, error) {
	builder := newAncestryBuilder("")
	for _, layer := range layers {
		if err := builder.AddLeafLayer(layer); err != nil {
			return nil, nil, err
		}
	}

	ancestry := builder.Ancestry()
	return ancestry.Layers, ancestry.By, nil
}
