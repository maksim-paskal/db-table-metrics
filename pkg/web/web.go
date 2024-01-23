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
package web

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/maksim-paskal/db-table-metrics/pkg/metrics"
	log "github.com/sirupsen/logrus"
)

var address = flag.String("web.listen-address", ":8080", "Address to listen on for web interface and telemetry.")

func NewServer() *Server {
	return &Server{
		serverReadTimeout: 10 * time.Second, //nolint:gomnd
	}
}

type Server struct {
	serverReadTimeout time.Duration
}

func (s *Server) GetHandler() *http.ServeMux {
	mux := http.NewServeMux()

	// metrics
	mux.Handle("/metrics", metrics.GetHandler())

	return mux
}

func (s *Server) Start(ctx context.Context) {
	log.Infof("Starting web server %s", *address)

	server := &http.Server{
		Addr:              *address,
		Handler:           s.GetHandler(),
		ReadHeaderTimeout: s.serverReadTimeout,
	}

	if err := server.ListenAndServe(); err != nil && ctx.Err() == nil {
		log.WithError(err).Fatal("error while starting web server")
	}
}
