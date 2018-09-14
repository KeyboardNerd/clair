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

package pgsql

import (
	"database/sql"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/pkg/commonerr"
	"github.com/sirupsen/logrus"
)

const (
	soiDetector = `
INSERT INTO detector (name, version, type)
SELECT CAST ($1 AS TEXT), CAST ($2 AS TEXT), CAST ($3 AS detector_type )
WHERE NOT EXISTS (SELECT id FROM detector WHERE name = $1 AND version = $2 AND type = $3);`

	selectAncestryDetectors = `
SELECT d.name, d.version, d.type
FROM ancestry_detector, detector AS d
WHERE ancestry_detector.detector_id = d.id AND ancestry.id = $1;`

	selectLayerDetectors = `
	SELECT d.name, d.version, d.type
	FROM layer_detector, detector AS d
	WHERE layer_detector.detector_id = d.id AND layer.id = $1;`

	insertAncestryDetectors = `
INSERT INTO ancestry_detector (ancestry_id, detector_id)
SELECT $1, $2
WHERE NOT EXISTS (SELECT id FROM ancestry_detector WHERE ancestry_id = $1 AND detector_id $2)`

	selectDetector = `SELECT id FROM detector WHERE name = $1 AND version = $2 AND type = $3`
)

func (tx *pgSession) PersistDetector(d database.Detector) error {
	if !d.Valid() {
		return commonerr.NewBadRequestError("invalid detector")
	}

	r, err := tx.Exec(soiDetector, d.Name, d.Version, d.Type)
	if err != nil {
		return handleError("soiDetector", err)
	}

	count, err := r.RowsAffected()
	if err != nil {
		return handleError("soiDetector", err)
	}

	if count == 0 {
		logrus.Debug("detector already exists: ", d)
	}

	return nil
}

func (tx *pgSession) persistAncestryDetectors(id int64, detectors []database.Detector) error {
	// find the index of all the detectors
	detectorIDs, err := tx.getDetectorIDs(detectors)
	if err != nil {
		return err
	}

	// insert all the detector relationship
	for _, detectorID := range detectorIDs {
		if _, err := tx.Exec(insertAncestryDetectors, id, detectorID); err != nil {
			return err
		}
	}

	return nil
}

func (tx *pgSession) getAncestryDetectors(id int64) ([]database.Detector, error) {
	return tx.getDetectors(selectAncestryDetectors, id)
}

func (tx *pgSession) getLayerDetectors(id int64) ([]database.Detector, error) {
	return tx.getDetectors(selectLayerDetectors, id)
}

// getDetectorIDs retrieve ids of the detectors from the database, if any is not
// found, return the error.
func (tx *pgSession) getDetectorIDs(detectors []database.Detector) ([]int64, error) {
	// TODO(sidac): use lru cache.
	ids := []int64{}
	for d := range detectors {
		id := sql.NullInt64{}
		err := tx.QueryRow(selectDetector).Scan(&id)
		if err != nil {
			return nil, handleError("selectDetector", err)
		}

		if !id.Valid {
			return nil, database.ErrInconsistent
		}

		ids = append(ids, id.Int64)
	}

	return ids, nil
}

func (tx *pgSession) getDetectors(query string, id int64) ([]database.Detector, error) {
	rows, err := tx.Query(query, id)
	if err != nil {
		return nil, handleError("getDetectors", err)
	}

	detectors := []database.Detector{}
	for rows.Next() {
		d := database.Detector{}
		err := rows.Scan(&d.Name, &d.Version, &d.Type)
		if err != nil {
			return nil, err
		}

		if !d.Valid() {
			return nil, database.ErrInvalidDetector
		}
	}

	return detectors, nil
}
