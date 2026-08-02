package types

import "net/http"

type Middleware func(handlerFunc http.Handler) http.Handler
