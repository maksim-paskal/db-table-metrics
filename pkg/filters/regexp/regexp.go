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
package regexp

import (
	"regexp"
	"slices"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Filter struct {
	Regexp                string `yaml:"regexp"`
	Group                 string `yaml:"group"`
	regexpReplace         *regexp.Regexp
	regexpReplacePosition int
}

func (f *Filter) Init(config interface{}) error {
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "yaml.Marshal failed")
	}

	if err := yaml.Unmarshal(configBytes, f); err != nil {
		return errors.Wrap(err, "yaml.Unmarshal failed")
	}

	if len(f.Group) == 0 {
		f.Group = "operation_code"
	}

	if len(f.Regexp) == 0 {
		return errors.New("regexp is empty")
	}

	regexpReplace, err := regexp.Compile(f.Regexp)
	if err != nil {
		return errors.Wrap(err, "error in regexp.Compile")
	}

	pos := slices.Index(regexpReplace.SubexpNames(), f.Group)
	if pos == -1 {
		return errors.New(f.Group + " not found in regexp")
	}

	f.regexpReplace = regexpReplace
	f.regexpReplacePosition = pos

	return nil
}

func (f *Filter) FormatValue(value string) string {
	result := f.regexpReplace.FindStringSubmatch(value)
	if len(result) == 0 || len(result) <= f.regexpReplacePosition {
		return value
	}

	return result[f.regexpReplacePosition]
}
