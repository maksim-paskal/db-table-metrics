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
package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/maksim-paskal/db-table-metrics/pkg/filters"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type collectType string

func (c collectType) Validate() error {
	if c != TypeCounter && c != TypeDictionary {
		return errors.New("unknown type")
	}

	return nil
}

const (
	TypeCounter    collectType = "counter"
	TypeDictionary collectType = "dictionary"
)

type Filters struct {
	Name   string
	Config interface{}
	filter filters.Filter
}

func (f *Filters) Init() error {
	filter, err := filters.NewFilter(f.Name, f.Config)
	if err != nil {
		return errors.Wrap(err, "error in filters.NewFilter")
	}

	f.filter = filter
	f.Config = &f.filter

	return nil
}

func (f *Filters) GetFilter() filters.Filter { //nolint:ireturn
	return f.filter
}

type CollectMetric struct {
	Type            collectType
	SQL             string
	Name            string
	Help            string
	Labels          []string
	TimeoutSeconds  int `yaml:"timeoutSeconds"`
	IntervalSeconds int `yaml:"intervalSeconds"`
	Filters         []*Filters
}

func (c *CollectMetric) Normalize() {
	if c.Type == "" {
		c.Type = TypeCounter
	}

	if len(c.Labels) == 0 {
		c.Labels = []string{"operation_code"}
	}

	if c.IntervalSeconds == 0 {
		c.IntervalSeconds = Get().IntervalSeconds
	}

	if c.TimeoutSeconds == 0 {
		c.TimeoutSeconds = Get().TimeoutSeconds
	}
}

func (c *CollectMetric) Validate() error {
	if c.Type == "" {
		return errors.New("type is required")
	}

	if err := c.Type.Validate(); err != nil {
		return errors.Wrap(err, "error in c.Type.Validate")
	}

	if c.SQL == "" {
		return errors.New("sql is required")
	}

	if c.Name == "" {
		return errors.New("name is required")
	}

	for _, filter := range c.Filters {
		if err := filter.Init(); err != nil {
			return errors.Wrap(err, "error in filters.NewFilter")
		}
	}

	return nil
}

func (c *CollectMetric) String() string {
	return fmt.Sprintf("%s, interval=%s, timeout=%s",
		c.Name,
		c.GetInterval(),
		c.GetTimeout(),
	)
}

func (c *CollectMetric) GetInterval() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}

func (c *CollectMetric) GetTimeout() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

type Type struct {
	DB              string           `yaml:"db"`
	Driver          string           `yaml:"driver"`
	TimeoutSeconds  int              `yaml:"timeoutSeconds"`
	IntervalSeconds int              `yaml:"intervalSeconds"`
	Metrics         []*CollectMetric `yaml:"metrics"`
}

func (c *Type) Normalize() {
	if c.Driver == "" {
		c.Driver = "mysql"
	}

	if c.IntervalSeconds == 0 {
		c.IntervalSeconds = 60
	}

	if c.TimeoutSeconds == 0 {
		c.TimeoutSeconds = c.IntervalSeconds
	}

	if len(c.DB) == 0 {
		c.DB = os.Getenv("DB")
	}
}

func (c *Type) Validate() error {
	if len(c.Driver) == 0 {
		return errors.New("driver is required")
	}

	if len(c.DB) == 0 {
		return errors.New("db is required")
	}

	return nil
}

func (c *Type) String() string {
	jsomBytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Sprintf("error in yaml.Marshal: %s", err.Error())
	}

	return string(jsomBytes)
}

var config = &Type{
	Metrics: []*CollectMetric{},
}

func Get() *Type {
	return config
}

var configPath = flag.String("config", "config.yaml", "Path to config file")

func Load() error {
	configByte, err := os.ReadFile(*configPath)
	if err != nil {
		return errors.Wrap(err, "error in os.ReadFile")
	}

	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		return errors.Wrap(err, "error in yaml.Unmarshal")
	}

	if len(config.Metrics) == 0 {
		return errors.New("no metrics found in config")
	}

	config.Normalize()

	if err := config.Validate(); err != nil {
		return errors.Wrap(err, "error in config.Validate")
	}

	for _, metric := range config.Metrics {
		metric.Normalize()

		if err := metric.Validate(); err != nil {
			return errors.Wrap(err, "error in metric.Validate")
		}
	}

	return nil
}
