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

package namespace

import (
	"database/sql"
	"sort"

	"github.com/coreos/clair/database"
	"github.com/coreos/clair/database/pgsql/util"
	"github.com/coreos/clair/pkg/commonerr"
)

var (
	cache *namespaceCache
	db    *sql.DB
)

/*
id SERIAL PRIMARY KEY,
name TEXT NULL,
version_format TEXT,
UNIQUE (name, version_format));
*/

func init() {
	cache = newCache(50)
}

func Register(b *sql.DB) {
	db = b
}

func GetID(n database.Namespace) (int, bool, error) {
	if id, ok := cache.Get(n); ok {
		return id, ok, nil
	}

	var id int
	if err := db.QueryRow(
		`SELECT id FROM namespace WHERE name = $1 AND version_format = $2;`,
		n.Name,
		n.VersionFormat,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return -1, false, nil
		}

		return -1, false, err
	}

	cache.Add(n, id)
	return id, true, nil
}

func Get(id int) (database.Namespace, bool, error) {
	if v, ok := cache.GetValue(id); ok {
		return v, ok, nil
	}

	var v database.Namespace
	if err := db.QueryRow(
		`SELECT name, version_format FROM namespace WHERE id = $1`,
		id,
	).Scan(v.Name, v.VersionFormat); err != nil {
		if err == sql.ErrNoRows {
			return v, false, nil
		}

		return v, false, err
	}

	cache.Add(v, id)
	return v, true, nil
}

func Add(n database.Namespace) (int, error) {
	id, ok, err := GetID(n)
	if err != nil || ok {
		return id, err
	}

	err = db.QueryRow(
		`INSERT INTO namespace (name, version_format) VALUES ($1, $2);
		 RETURNING id`,
		n.Name,
		n.VersionFormat,
	).Scan(id)

	if err != nil {
		// race condition.
		id, ok, err = GetID(n)
		if err != nil {
			return id, err
		}

		if !ok {
			panic("wtf")
		}
	}

	return id, err
}

// PersistNamespaces soi namespaces into database.
func PersistNamespaces(tx *sql.Tx, namespaces []database.Namespace) error {
	if len(namespaces) == 0 {
		return nil
	}

	// Sorting is needed before inserting into database to prevent deadlock.
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].Name < namespaces[j].Name &&
			namespaces[i].VersionFormat < namespaces[j].VersionFormat
	})

	keys := make([]interface{}, len(namespaces)*2)
	for i, ns := range namespaces {
		if ns.Name == "" || ns.VersionFormat == "" {
			return commonerr.NewBadRequestError("Empty namespace name or version format is not allowed")
		}
		keys[i*2] = ns.Name
		keys[i*2+1] = ns.VersionFormat
	}

	_, err := tx.Exec(queryPersistNamespace(len(namespaces)), keys...)
	if err != nil {
		return util.HandleError("queryPersistNamespace", err)
	}
	return nil
}

func FindNamespaceIDs(tx *sql.Tx, namespaces []database.Namespace) ([]sql.NullInt64, error) {
	if len(namespaces) == 0 {
		return nil, nil
	}

	keys := make([]interface{}, len(namespaces)*2)
	nsMap := map[database.Namespace]sql.NullInt64{}
	for i, n := range namespaces {
		keys[i*2] = n.Name
		keys[i*2+1] = n.VersionFormat
		nsMap[n] = sql.NullInt64{}
	}

	rows, err := tx.Query(querySearchNamespace(len(namespaces)), keys...)
	if err != nil {
		return nil, util.HandleError("searchNamespace", err)
	}

	defer rows.Close()

	var (
		id sql.NullInt64
		ns database.Namespace
	)
	for rows.Next() {
		err := rows.Scan(&id, &ns.Name, &ns.VersionFormat)
		if err != nil {
			return nil, util.HandleError("searchNamespace", err)
		}
		nsMap[ns] = id
	}

	ids := make([]sql.NullInt64, len(namespaces))
	for i, ns := range namespaces {
		ids[i] = nsMap[ns]
	}

	return ids, nil
}
