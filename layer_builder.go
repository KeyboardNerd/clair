package clair

import "github.com/coreos/clair/database"

// For building a layer from the scanned result and the database.
type LayerBuilder struct {
	diffsha256 string
	features   map[database.Detector][]database.Feature
	namespaces map[database.Detector]*database.Namespace
}

// NewLayerBuilder returns a new layer builder
func NewLayerBuilder(sha string, existingLayer *database.Layer) *LayerBuilder {
	builder := &LayerBuilder{sha, make(map[database.Detector][]database.Feature), make(map[database.Detector]*database.Namespace)}
	if existingLayer != nil {
		if existingLayer.Hash != sha {
			panic("existing Layer sha must match the given sha")
		}

		builder.init(existingLayer)
	}

	return builder
}

func (b *LayerBuilder) layerFeatures() []database.LayerFeature{
	layerFeatures := []database.LayerFeature{}
	for detector, features := range b.features{
		layerFeatures = append(database.LayerFeature{
			
		})
	}
}

// Layer generates a layer based on the result of the builder
func (b *LayerBuilder) Layer() *database.Layer {
	detectors := make([]database.Detector, 0, len(b.features) + len(b.namespaces))
	for _, d := range b.features{

	}

	layer := &database.Layer{
		Hash: b.diffsha256,
		By: 
	}
}

// init initializes the builder with content in the existing layer.
func (b *LayerBuilder) init(layer *database.Layer) {
	for _, feature := range layer.Features {
		b.AddFeature(feature.By, feature.Feature)
	}

	for _, namespace := range layer.Namespaces {
		b.AddNamespace(namespace.By, namespace.Namespace)
	}

	for _, detector := range layer.By {
		if _, ok := b.features[detector]; detector.DType == database.FeatureDetectorType && !ok {
			b.features[detector] = []database.Feature{}
		}

		if _, ok := b.namespaces[detector]; detector.DType == database.NamespaceDetectorType && !ok {
			b.namespaces[detector] = nil
		}
	}
}

// AddFeature adds a feature to the layer
func (b *LayerBuilder) AddFeature(detector database.Detector, feature database.Feature) *LayerBuilder {
	if detector.DType != database.FeatureDetectorType {
		panic("Invalid Feature detector must be feature detector type.")
	}

	for _, existingFeature := range b.features[detector] {
		if existingFeature == feature {
			return b
		}
	}

	b.features[detector] = append(b.features[detector], feature)
	return b
}

func (b *LayerBuilder) AddNamespace(detector database.Detector, namespace database.Namespace) *LayerBuilder {
	// TODO: Differentiate the two types of detectors during compile time.
	if detector.DType != database.NamespaceDetectorType {
		panic("Namespace expects namespace detector type.")
	}

	if _, ok := b.namespaces[detector]; ok {
		panic("One detector can only detect one namespace.")
	}

	b.namespaces[detector] = &namespace
	return b
}
