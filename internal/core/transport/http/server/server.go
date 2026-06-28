package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/shitaiv1ck/test-effective-mobile/internal/core/logger"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	handler http.Handler
	config  Config
	logger  *logger.Logger
}

func NewHTTPServer(handler http.Handler, config Config, logger *logger.Logger) *HTTPServer {
	return &HTTPServer{
		handler: handler,
		config:  config,
		logger:  logger,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	handler := middleware.ChainMiddleware(
		s.handler,
		middleware.CORS(),
		middleware.RequestID(),
		middleware.Logger(s.logger),
		middleware.Panic(),
		middleware.Trace(),
	)

	httpServer := &http.Server{
		Addr:    s.config.Addr,
		Handler: handler,
	}

	errChan := make(chan error)
	go func() {
		defer close(errChan)

		s.logger.Debug("run http server...", zap.String("addr", s.config.Addr))
		err := httpServer.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		s.logger.Debug("stop http server")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			httpServer.Close()
		}
	case err := <-errChan:
		s.logger.Error("failed to run http server", zap.Error(err))
		return err
	}

	return nil
}
