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
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/maksim-paskal/db-table-metrics/pkg/api"
	"github.com/maksim-paskal/db-table-metrics/pkg/collect"
	"github.com/maksim-paskal/db-table-metrics/pkg/config"
	"github.com/maksim-paskal/db-table-metrics/pkg/web"
	log "github.com/sirupsen/logrus"
)

const defaultGraceInterval = 5 * time.Second

var (
	logLevel      = flag.String("log.level", "info", "Log level")
	graceInterval = flag.Duration("grace.interval", defaultGraceInterval, "Grace interval")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChanInterrupt := make(chan os.Signal, 1)
	signal.Notify(signalChanInterrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signalChanInterrupt:
			log.Error("Got interruption signal...")
			cancel()
		case <-ctx.Done():
		}
		<-signalChanInterrupt
		os.Exit(1)
	}()

	logLevel, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.WithError(err).Fatal("error parsing log level")
	}

	log.SetLevel(logLevel)

	if _, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if err := config.Load(); err != nil {
		log.WithError(err).Fatal("error while loading config")
	}

	log.Infof("config loaded:\n %s", config.Get().String())

	if err := api.Init(ctx); err != nil {
		log.WithError(err).Fatal("error while initializing api")
	}

	// start all collectors
	for _, collectMetric := range config.Get().Metrics {
		go collect.NewCollector().Start(ctx, collectMetric)
	}

	go web.NewServer().Start(ctx)

	<-ctx.Done()

	if *graceInterval > 0 {
		log.Infof("Waiting grace interval %s", *graceInterval)
		time.Sleep(*graceInterval)
	}
}
