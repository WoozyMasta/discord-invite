package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// newHandler creates a new HTTP handler with logging middleware.
// It handles incoming requests by generating a Discord invite and redirecting the user.
func newHandler(cfg *Config, logger zerolog.Logger) http.Handler {
	// Setting up middleware for logging
	handler := hlog.NewHandler(logger)(
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			log.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Str("remote", r.RemoteAddr).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("handled request")
		})(
			hlog.RequestIDHandler("req_id", "Request-Id")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				inviteCode, err := cfg.makeInvite()
				if err != nil {
					log.Error().Err(err).Msg("Failed to generate Discord invite")
					http.Error(w, "Failed to generate Discord invite", http.StatusInternalServerError)
					return
				}
				log.Debug().Msgf("Created invite code %s", inviteCode)
				redirectURL := fmt.Sprintf("https://discord.gg/%s", inviteCode)
				http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			})),
		),
	)

	return handler
}
