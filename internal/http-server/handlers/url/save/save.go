package save

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Zapi-web/url-shortener/internal/lib/random"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

type URLSetter interface {
	Set(ctx context.Context, key, value string) error
}

func New(data URLSetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		slog.Debug("trying to decode request")
		r.Body = http.MaxBytesReader(w, r.Body, 1024*10)
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			slog.Error("failed to decode request body", "err", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		slog.Debug("request decoded", "url", req.URL)

		if _, err := url.ParseRequestURI(req.URL); err != nil {
			slog.Error("invalid request", "URL", req.URL)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		alias, err := random.RandomKey()
		if err != nil {
			slog.Error("failed to generate a key", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = data.Set(r.Context(), alias, req.URL)
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
