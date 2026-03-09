package save

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Zapi-web/url-shortener/internal/lib/random"
	"github.com/Zapi-web/url-shortener/internal/storage/redis"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

func New(db *redis.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		slog.Debug("trying to decode request")
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			slog.Error("failed to decode request body", "err", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		slog.Debug("request decoded", "url", req.URL)

		alias, err := random.RandomKey()
		if err != nil {
			slog.Error("failed to generate a key", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = db.Set(r.Context(), alias, req.URL)
		if err != nil {
			slog.Error("failed to save url", "alias", alias, "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		slog.Info("url saved successfully", "alias", alias)

		res := Response{
			ShortURL: alias,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
