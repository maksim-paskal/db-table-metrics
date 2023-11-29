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
package collect

import (
	"context"
	"fmt"
	"time"

	"github.com/maksim-paskal/db-table-metrics/pkg/api"
	"github.com/maksim-paskal/db-table-metrics/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

func NewCollector() *Collector {
	return &Collector{
		Namespace: "db_metrics",
	}
}

type Collector struct {
	Namespace string
}

func (c *Collector) Start(ctx context.Context, input *config.CollectMetric) {
	log.Infof("Starting collector %s", input.String())

	gauge := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Name:      input.Name,
		Help:      input.Help,
	}, input.Labels)

	errorCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: c.Namespace,
		Help:      "Errors count",
		Name:      fmt.Sprintf("%s_errors", input.Name),
	})

	executionTimeSeconds := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: c.Namespace,
		Help:      "Execution time in seconds",
		Name:      fmt.Sprintf("%s_duration", input.Name),
	})

	tiker := time.NewTicker(input.GetInterval())

	for ctx.Err() == nil {
		parentCtx := ctx

		go func() {
			start := time.Now()
			defer func() {
				executionTimeSeconds.Observe(time.Since(start).Seconds())
			}()

			ctx, cancel := context.WithTimeout(parentCtx, input.GetTimeout())
			defer cancel()

			operationCodes, err := api.GetQuery(ctx, input.SQL)
			if err != nil {
				errorCounter.Inc()
				log.WithError(err).Error("error while getting results")
			}

			for _, operationCode := range operationCodes {
				for _, filter := range input.Filters {
					operationCode.Code = filter.GetFilter().FormatValue(operationCode.Code)
				}

				gauge.WithLabelValues(operationCode.Code).Set(float64(operationCode.Count))
			}
		}()

		select {
		case <-tiker.C:
		case <-parentCtx.Done():
		}
	}
}
