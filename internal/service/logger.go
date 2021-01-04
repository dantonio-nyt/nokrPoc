package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptrace"
)

const LOGGER = "logger"
// getLoggingMiddleware adds a logger to the context and adds request scoped attributes to the
// kit logger so that log entries will contain information about the initiating http request.
//	Example:
//	import http "github.com/go-kit/kit/transport/http"
// 	clientOptions := []http.ClientOption{
// 		http.ClientBefore(logHTTPTrace("your-service-label")),
// 	}
func (h *HermesService) getLoggingMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := kithttp.PopulateRequestContext(r.Context(), r)
			kitlog, logClose, _ := kit.NewLogger(ctx, h.config.ProjectID)
			kitlog = kit.AddLogKeyVals(ctx, kitlog)
			ctx = context.WithValue(r.Context(), LOGGER, &serviceLog{kitlog, logClose})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getLogger retrieves a logger from the context
func getLogger(ctx context.Context) Logger {
	l, ok := ctx.Value(LOGGER).(Logger)
	// fallback gracefully if logger is not in the context
	if !ok {
		kitlog, logClose, _ := kit.NewLogger(ctx, "")
		l = &serviceLog{kitlog, logClose}
	}
	return l
}

// logHTTPTrace can be injected as a RequestFunc through go-kit http ClientBefore hook
// to log detailed http traces.
func logHTTPTrace(label string) kithttp.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		return httptrace.WithClientTrace(req.Context(), getHttpTrace(ctx, label))
	}
}

// getHttpTrace logs out detailed http traces to debug connectivity issues.
func getHttpTrace(ctx context.Context, label string) *httptrace.ClientTrace {
	l := getLogger(ctx)
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] GetConn: %s\n", label, hostPort))
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] DNSStart: %+v\n", label, info))
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] DNSDone: %+v\n", label, dnsInfo))
		},
		ConnectStart: func(network, addr string) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] ConnectStart: network:%s addr:%s", label, network, addr))
		},
		ConnectDone: func(network, addr string, err error) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] ConnectDone: network:%s addr%s err:%+v", label, network, addr, err))
		},
		TLSHandshakeStart: func() {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] TLSHandshakeStart", label))
		},
		TLSHandshakeDone: func(state tls.ConnectionState, e error) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] TLSHandshakeDone state:%+v error:%+v", label, state, e))
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] GotConn: %+v", label, connInfo))
		},
		GotFirstResponseByte: func() {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] GotFirstResponseByte", label))
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] WroteRequest: %+v", label, info))
		},
		WroteHeaderField: func(key string, value []string) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] WroteHeaderField %s:%+v", label, key, value))
		},
		Got100Continue: func() {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] Got100Continue", label))
		},
		PutIdleConn: func(err error) {
			l.logDebug(fmt.Sprintf("HTTP Trace - [%s] PutIdleConn, %+v", label, err))
		},
	}
}

// Custom logger definition to simplify logging patterns.
type Logger interface {
	logDebug(msg string)
	logError(msg string)
	logWarn(msg string)
	logInfo(msg string)
	Log(keyvals ...interface{}) error
}

// serviceLog is a container to wrap the kit logger with some convenience methods.
type serviceLog struct {
	kitlog   log.Logger
	logClose func() error
}

func (s *serviceLog) logDebug(msg string) {
	s.kitlog.Log("level", level.DebugValue(), "message", msg)
}

func (s *serviceLog) logError(msg string) {
	s.kitlog.Log("level", level.ErrorValue(), "message", msg)
}

func (s *serviceLog) logInfo(msg string) {
	s.kitlog.Log("level", level.InfoValue(), "message", msg)
}

func (s *serviceLog) logWarn(msg string) {
	s.kitlog.Log("level", level.WarnValue(), "message", msg)
}

func (s *serviceLog) Log(keyvals ...interface{}) error {
	return s.kitlog.Log(keyvals)
}
