package server

import (
	"errors"
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/constants"
	"github.com/adriein/hastypal/internal/hastypal/shared/helper"
	middleware2 "github.com/adriein/hastypal/internal/hastypal/shared/middleware"
	types2 "github.com/adriein/hastypal/internal/hastypal/shared/types"
	"log"
	"log/slog"
	"net/http"
)

type HastypalApiServer struct {
	address string
	router  *http.ServeMux
}

func New(address string) (*HastypalApiServer, error) {
	router := http.NewServeMux()

	return &HastypalApiServer{
		address: address,
		router:  router,
	}, nil
}

func (s *HastypalApiServer) Start() {
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", s.router))

	MuxMiddleWareChain := middleware2.NewMiddlewareChain(
		middleware2.NewRequestTracingMiddleware,
	)

	server := http.Server{
		Addr:    s.address,
		Handler: MuxMiddleWareChain.ApplyOn(v1),
	}

	slog.Info("Starting the HastypalApiServer at " + s.address)

	err := server.ListenAndServe()

	if err != nil {
		evtErr := types2.ApiError{Msg: err.Error(), Function: "Start", File: "server.go"}

		log.Fatal(evtErr.Error())
	}
}

func (s *HastypalApiServer) Route(url string, handler http.Handler) {
	s.router.Handle(url, handler)
}

func (s *HastypalApiServer) NewHandler(handler types2.HastypalHttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appError types2.ApiErrorInterface

		if err := handler(w, r); errors.As(err, &appError) {
			if appError.IsDomain() {
				response := types2.ServerResponse{
					Ok:    false,
					Error: appError.PresentableError(),
				}

				if encodeErr := helper.Encode[types2.ServerResponse](w, http.StatusOK, response); encodeErr != nil {
					log.Fatal(encodeErr.Error())
				}

				slog.Error(fmt.Sprintf("%s TraceId=%s", appError.Error(), r.Header.Get("traceId")))

				return
			}

			response := types2.ServerResponse{
				Ok:    false,
				Error: constants.ServerGenericError,
			}

			if encodeErr := helper.Encode[types2.ServerResponse](w, http.StatusInternalServerError, response); encodeErr != nil {
				log.Fatal(encodeErr.Error())
			}

			slog.Error(fmt.Sprintf("%s TraceId=%s", appError.Error(), r.Header.Get("traceId")))
		}
	}
}
