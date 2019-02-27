package database

// Layer is a layer with all the detected features and namespaces.
type Layer struct {
	// Hash is the sha-256 tarsum on the layer's blob content.
	Hash string `json:"hash"`
	// Features contain
	Features   DetectedFeatures   `json:"features"`
	Namespaces DetectedNamespaces `json:"namespaces"`
}

// DetectedFeatures is a map from the detector to the result of the detector.
type DetectedFeatures map[Detector]FeatureDetectResult

// DetectedNamespaces is a map from the detector to the result of the detector.
type DetectedNamespaces map[Detector]NamespaceDetectResult

// NamespaceDetectResult contains the status of detection and the shitted out result
type NamespaceDetectResult struct {
	Status    DetectorStatus `json:"status"`
	Namespace *Namespace     `json:"namespace"`
}

// FeatureDetectResult contains the status of detection and the shitted out result
type FeatureDetectResult struct {
	Status   DetectorStatus `json:"status"`
	Features []Feature      `json:"features"`
}

// AllUniqueFeatures returns a list of unique features detected by all detectors
func (d DetectedFeatures) AllUniqueFeatures() []Feature {
	features := make([]Feature, 0)
	for _, featureMap := range d {
		for _, feature := range featureMap.Features {
			existing := false
			for _, existingFeature := range features {
				if existingFeature == feature {
					existing = true
					break
				}
			}

			if existing {
				continue
			}

			copiedFeature := feature
			features = append(features, copiedFeature)
		}
	}

	return features
}

// AllUniqueNamespaces returns all unique namespaces in the detection result
func (d DetectedNamespaces) AllUniqueNamespaces() []Namespace {
	namespaces := make([]Namespace, 0)
	for _, detectResult := range d {
		existing := false
		for _, existingNamespace := range namespaces {
			if detectResult.Namespace != nil && existingNamespace == *detectResult.Namespace {
				existing = true
				break
			}
		}

		if existing {
			continue
		}

		namespaces = append(namespaces, *detectResult.Namespace)
	}

	return namespaces
}
