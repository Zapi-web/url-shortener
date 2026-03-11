package get

import (
	"log/slog"
	"net/http"

	"github.com/Zapi-web/url-shortener/internal/storage/db"
	"github.com/go-chi/chi/v5"
)

func GetNew(d *db.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := chi.URLParam(r, "short_url")

		if shortUrl == "" {
			slog.Error("short url is empty")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		slog.Debug("short url received", "short_url", shortUrl)

		val, err := d.Get(r.Context(), shortUrl)

		if err != nil {
			slog.Error("failed to get url", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		slog.Info("url received", "url", val)

		http.Redirect(w, r, val, http.StatusFound)
	}
}
