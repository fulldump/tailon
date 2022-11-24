package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fulldump/box"
)

func RecoverFromPanic(next box.H) box.H {
	return func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()

		next(ctx)
	}
}

func AccessLog(l *log.Logger) box.I {
	return func(next box.H) box.H {
		return func(ctx context.Context) {
			r := box.GetRequest(ctx)
			action := ""
			if boxAction := box.GetBoxContext(ctx).Action; boxAction != nil {
				action = boxAction.Name
			}
			now := time.Now()
			defer func() {
				l.Println(formatRemoteAddr(r), r.Method, r.URL.String(), time.Since(now), action)
			}()

			next(ctx)
		}
	}
}

func formatRemoteAddr(r *http.Request) string {
	xorigin := strings.TrimSpace(strings.Split(
		r.Header.Get("X-Forwarded-For"), ",")[0])
	if xorigin != "" {
		return xorigin
	}

	return r.RemoteAddr[0:strings.LastIndex(r.RemoteAddr, ":")]
}

func PrettyErrorInterceptor(next box.H) box.H {
	return func(ctx context.Context) {

		next(ctx)

		err := box.GetError(ctx)
		if err == nil {
			return
		}
		w := box.GetResponse(ctx)

		if err == box.ErrResourceNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message":     err.Error(),
					"description": fmt.Sprintf("resource '%s' not found", box.GetRequest(ctx).URL.String()),
				},
			})
			return
		}

		if err == box.ErrMethodNotAllowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message":     err.Error(),
					"description": fmt.Sprintf("method '%s' not allowed", box.GetRequest(ctx).Method),
				},
			})
			return
		}

		if _, ok := err.(*json.SyntaxError); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"message":     err.Error(),
					"description": "Malformed JSON",
				},
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message":     err.Error(),
				"description": "Unexpected error",
			},
		})

	}
}
