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

import (
	"database/sql/driver"
	"fmt"
)

// DetectorStatus is the status if a file is changed, or removed, or not found
// in the layer blob.
type DetectorStatus string

const (
	// Whiteout means the required feature file is removed from the base layer.
	// https://github.com/moby/moby/blob/master/pkg/archive/whiteouts.go#L9
	Whiteout DetectorStatus = "file_removed"
	// Changed means the File is changed and exists in the current layer
	Changed DetectorStatus = "file_changed"
	// NotFound means the file does not exists in the current layer and
	// therefore, either it doesn't exists or it's not changed.
	NotFound DetectorStatus = "file_not_found"
	// Corrupted means the file to scan is corrupted.
	Corrupted DetectorStatus = "file_corrupted"
	// BlackListed means the detector encountered some file that's not allowed
	// to exist.
	BlackListed DetectorStatus = "file_blacklisted"
)

var detectorStatuses = []DetectorStatus{
	Whiteout,
	Changed,
	NotFound,
	Corrupted,
}

// Scan implements the database/sql.Scanner interface.
func (d *DetectorStatus) Scan(value interface{}) error {
	val := value.(string)
	for _, ft := range detectorStatuses {
		if string(ft) == val {
			*d = ft
			return nil
		}
	}

	panic(fmt.Sprintf("invalid feature type received from database: '%s'", val))
}

// Value implements the database/sql/driver.Valuer interface.
func (d *DetectorStatus) Value() (driver.Value, error) {
	return string(*d), nil
}
