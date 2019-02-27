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
id SERIAL PRIMARY KEY
layer_detector_id NOT NULL INT REFERENCES layer_detector ON DELETE CASCADE
namespace_id INT REFERENCES namespace ON DELETE CASCADE
*/
const (
	findLayerNamespaces = `
	SELECT ns.name, ns.version_format, ln.detector_id
	FROM layer_namespace AS ln, namespace AS ns
	WHERE ln.namespace_id = ns.id
		AND ln.layer_id = $1`
)
