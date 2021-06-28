// Copyright © 2020 Dmitry Stoletov <info@imega.ru>
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

package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/imega/daemon"
)

const defaultTimeout = 60

type health struct {
	hcf     []daemon.HealthCheckFunc
	timeout time.Duration
}

// Handler returns an http.Handler
//
// It returns status 204 if all healthcheckers returns true.
// It returns status 503 (unhealthy) if anyone healthcheckers returns false.
func Handler(opts ...Option) http.Handler {
	h := &health{
		timeout: defaultTimeout * time.Second,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// HandlerFunc returns an http.HandlerFunc.
func HandlerFunc(opts ...Option) http.HandlerFunc {
	return Handler(opts...).ServeHTTP
}

// Option adds optional parameter for the HealthcheckHandler.
type Option func(*health)

// WithHealthCheckFuncs adds the functions healthcheck of daemon
// that needs to be added as part of healthcheck.
func WithHealthCheckFuncs(f ...daemon.HealthCheckFunc) Option {
	return func(h *health) {
		h.hcf = f
	}
}

// WithTimeout sets the global timeout for all healthcheckers.
func WithTimeout(timeout time.Duration) Option {
	return func(h *health) {
		h.timeout = timeout
	}
}

func (h *health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := r.Context(), func() {}
	if h.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
	}

	defer cancel()

	statusCh := make(chan bool, len(h.hcf))
	wg := sync.WaitGroup{}

	wg.Add(len(h.hcf))

	for _, f := range h.hcf {
		go func(f daemon.HealthCheckFunc) {
			statusCh <- f()

			wg.Done()
		}(f)
	}

	go func() {
		wg.Wait()
		close(statusCh)
	}()

	for {
		select {
		case <-ctx.Done():
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)

			return

		case r, ok := <-statusCh:
			if !ok {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			if !r {
				http.Error(w, "unhealthy", http.StatusServiceUnavailable)

				return
			}
		}
	}
}
