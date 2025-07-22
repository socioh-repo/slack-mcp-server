package server

import (
	"net/http"
	"time"

	"github.com/korotovsky/slack-mcp-server/pkg/provider"
	"go.uber.org/zap"
)

type HealthHandler struct {
	provider *provider.ApiProvider
	logger   *zap.Logger
}

func NewHealthHandler(provider *provider.ApiProvider, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		provider: provider,
		logger:   logger,
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}

	// Check if provider is ready
	ready, err := h.provider.IsReady()
	if err != nil {
		h.logger.Error("Health check failed - provider error", zap.Error(err))
		http.Error(w, "Service not ready", http.StatusServiceUnavailable)
		return
	}

	if !ready {
		h.logger.Warn("Health check failed - provider not ready")
		http.Error(w, "Service warming up", http.StatusServiceUnavailable)
		return
	}

	// Check Slack authentication
	_, err = h.provider.Slack().AuthTest()
	if err != nil {
		h.logger.Error("Health check failed - Slack auth error", zap.Error(err))
		http.Error(w, "Slack authentication failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
} 