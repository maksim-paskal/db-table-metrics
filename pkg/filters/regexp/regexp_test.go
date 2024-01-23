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
package regexp_test

import (
	"testing"

	"github.com/maksim-paskal/db-table-metrics/pkg/filters/regexp"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	filter := &regexp.Filter{}
	if err := filter.Init(regexp.Filter{
		Regexp: "^(?P<operation_code>\\d{3}).",
	}); err != nil {
		t.Fatal(err)
	}

	tests := make(map[string]string)
	tests["123SOME_NAME"] = "123"
	tests["123SOME_NAME2"] = "123"
	tests["123SOME_NAME3"] = "123"
	tests["4321SOME_NAME4"] = "432"
	tests["4321SOME_NAME5"] = "432"
	tests["4321SOME_NAME6"] = "432"
	tests["SOMEFAKE"] = "SOMEFAKE"

	for input, expected := range tests {
		value := filter.FormatValue(input)

		if value != expected {
			t.Fatalf("expected %s, got %s", expected, value)
		}
	}
}
