package types

import "net/http"

type HastypalHttpHandler func(w http.ResponseWriter, r *http.Request) error
