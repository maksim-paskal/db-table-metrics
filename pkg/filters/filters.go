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
package filters

import (
	"github.com/maksim-paskal/db-table-metrics/pkg/filters/regexp"
	"github.com/pkg/errors"
)

type Filter interface {
	Init(config interface{}) error
	FormatValue(value string) string
}

func NewFilter(name string, config interface{}) (Filter, error) { //nolint:ireturn
	var filter Filter

	switch name {
	case "regexp":
		filter = &regexp.Filter{}
	default:
		return nil, errors.New("unknown filter " + name)
	}

	if err := filter.Init(config); err != nil {
		return nil, errors.Wrap(err, "filter init failed")
	}

	return filter, nil
}
