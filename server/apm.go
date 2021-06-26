package server

import (
	"net/http"
	"strings"

	newrelictracing "github.com/anant-sharma/go-utils/new-relic/tracing"
)

func apm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn, ctx := newrelictracing.NewTransaction(r.Context(), strings.Join([]string{r.Method, r.RequestURI}, "__"))
		txn.SetWebRequestHTTP(r)
		defer txn.End()

		next.ServeHTTP(txn.SetWebResponse(w), r.WithContext(ctx))
	})
}
