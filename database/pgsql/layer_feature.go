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

package pgsql

/*
	schema:
		id SERIAL PRIMARY KEY,
		layer_detector_id NOT NULL INT REFERENCES layer_detector ON DELETE CASCADE,
		feature_id INT REFERENCES feature ON DELETE CASCADE
*/
const (
	findLayerFeatures = `
	SELECT f.name, f.version, f.version_format, t.name, lf.detector_id
	FROM layer_feature AS lf, feature AS f, feature_type AS t
	WHERE lf.feature_id = f.id
		AND t.id = f.type
		AND lf.layer_id = $1`
)

type layerFeature struct {
	featureID int
}

func (tx *pgSession) findLayerFeatures()
