package httpwareclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/improbable-eng/go-httpwares"
	http_retry "github.com/improbable-eng/go-httpwares/retry"
	"github.com/sirupsen/logrus"
)

var (
	httpClient HTTPClientDo
	logger     *logrus.Entry
	rmw        []RequestFunc
)

// SendIn returns a new request given a method, URL, and optional body,
// and optional coder, and RequestFunc for modify of request before send
type SendIn struct {
	Method      string
	URL         string
	Headers     map[string]string
	BodySend    interface{}
	BodyRecv    interface{}
	Coder       Coder
	RequestFunc RequestFunc
}

// HTTPClientDo is interface http.Client
type HTTPClientDo interface {
	Do(req *http.Request) (*http.Response, error)
}

// WithLogger append logger
func WithLogger(l *logrus.Entry) {
	logger = l
}

// RequestFunc is a signature for all http request middleware
type RequestFunc func(req *http.Request)

// WithTripperware append tripperwares
func WithTripperware(httpClient *http.Client, tripperwares ...httpwares.Tripperware) {
	httpClient = httpwares.WrapClient(httpClient, tripperwares...)
}

// WithRequestWares append ware in request
func WithRequestWares(wares ...RequestFunc) {
	rmw = append(rmw, wares...)
}

// RetryTriceTripperwares will retry three times to send request
func RetryTriceTripperwares() []httpwares.Tripperware {
	var wares []httpwares.Tripperware

	wares = append(wares, http_retry.Tripperware(
		http_retry.WithMax(3),
		http_retry.WithBackoff(func(attempt uint) time.Duration {
			return time.Duration(time.Duration(int64(attempt)) * 100 * time.Millisecond)
		}),
	))

	return wares
}

// DefaultHTTPClient returns default http.Client with set timeouts
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
}

// Send request with context
func Send(ctx context.Context, in *SendIn) error {
	if in.Coder == nil {
		in.Coder = &nullCoder{}
	}
	reader, err := in.Coder.Encode(in.BodySend)
	if err != nil {
		return fmt.Errorf("failed to encode request body, %s", err)
	}

	req, err := http.NewRequest(in.Method, in.URL, reader)
	if err != nil {
		return err
	}

	if len(rmw) > 0 {
		for _, w := range rmw {
			w(req)
		}
	}

	for k, v := range in.Headers {
		req.Header.Set(k, v)
	}

	if in.RequestFunc != nil {
		in.RequestFunc(req)
	}

	if logger != nil && logger.Logger.Level == logrus.DebugLevel {
		dump, _ := httputil.DumpRequestOut(req, true)
		logger.Debugf("http-client send request, %s", string(dump))
	}

	if httpClient == nil {
		httpClient = httpwares.WrapClient(DefaultHTTPClient(), RetryTriceTripperwares()...)
	}

	r, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	if logger != nil && logger.Logger.Level == logrus.DebugLevel {
		dump, _ := httputil.DumpResponse(r, true)
		logger.Debugf("http-client getting response, %s", string(dump))
	}

	if in.BodyRecv == nil {
		return nil
	}

	if err := in.Coder.Decode(r.Body, in.BodyRecv); err != nil {
		return fmt.Errorf("failed to decode response body, %s", err)
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	return nil
}
