/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"context"
	"database/sql"
	"time"

	"github.com/maksim-paskal/db-table-metrics/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

const initTimeout = 10 * time.Second

func Init(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, initTimeout)
	defer cancel()

	var err error

	db, err = sql.Open(config.Get().Driver, config.Get().DB)
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}

	if err = db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping database")
	}

	return nil
}

type GetQueryResult struct {
	Code  string
	Count int
}

func GetQuery(ctx context.Context, sql string) ([]*GetQueryResult, error) {
	log.Debug(sql)

	results, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query database")
	}

	if results.Err() != nil {
		return nil, errors.Wrap(results.Err(), "failed to query database")
	}

	defer results.Close()

	result := make([]*GetQueryResult, 0)

	for results.Next() {
		item := GetQueryResult{}

		err = results.Scan(&item.Code, &item.Count)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan results")
		}

		log.Debugf("%+v", item)

		result = append(result, &item)
	}

	return result, nil
}
