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

import "fmt"

// StorageError is database error
type StorageError struct {
	reason   string
	original error
}

func (e *StorageError) Error() string {
	internalErr := ""
	if e.original != nil {
		internalErr = e.original.Error()
	}
	return fmt.Sprintf("%s, Internal Error='%s'", e.reason, internalErr)
}

// NewStorageErrorWithInternalError creates a new database error
func NewStorageErrorWithInternalError(reason string, originalError error) *StorageError {
	return &StorageError{reason, originalError}
}

// NewStorageError creates a new database error
func NewStorageError(reason string) *StorageError {
	return &StorageError{reason, nil}
}
