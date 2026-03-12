package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Zapi-web/url-shortener/internal/domain"
	"github.com/go-chi/chi/v5"
)

type URLGetter interface {
	Get(ctx context.Context, key string) (string, error)
}

func GetNew(data URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := chi.URLParam(r, "short_url")

		if shortUrl == "" {
			slog.Error("short url is empty")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		slog.Debug("short url received", "short_url", shortUrl)

		val, err := data.Get(r.Context(), shortUrl)

		if err != nil {
			if errors.Is(err, domain.ErrKeyisEmpty) {
				slog.Warn(domain.ErrKeyisEmpty.Error())
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			if errors.Is(err, domain.ErrUrlNotFound) {
				slog.Info(domain.ErrUrlNotFound.Error())
				http.Error(w, "not found", http.StatusNotFound)
				return
			}

			slog.Error("failed to get url", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		slog.Info("url received", "url", val)

		http.Redirect(w, r, val, http.StatusFound)
	}
}
