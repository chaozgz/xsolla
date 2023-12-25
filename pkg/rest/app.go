package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sxolla-rest-api/config"
	"sxolla-rest-api/ent"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func RunGinService(port string, entClient *ent.Client) {

	config.SetClient(entClient)

	router, err := inject()

	if err != nil {
		log.Err(err).Msg("Failure to inject data sources")
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("Failed to initialize server")
		}
	}()

	log.Info().Msgf("gin listening on port %v", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	log.Warn().Msg("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Err(err).Msg("Server forced to shutdown:")
	}
}
