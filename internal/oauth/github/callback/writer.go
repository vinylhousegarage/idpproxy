package callback

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func writeJSON(w http.ResponseWriter, status int, v any, logger *zap.Logger) {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error("writeJSON marshal error", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		logger.Error("writeJSON write error", zap.Error(err))
	}
}
