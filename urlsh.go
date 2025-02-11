package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Storage interface {
	Get(id string) (ShortURL, error)
	Set(id, url string) error
}

func main() {
	log := slog.Default()

	db, err := NewSqliteDB()
	if err != nil {
		log.Error("Failed to create DB", "err", err)
		return
	}

	store := NewStore(db)

	// urls, err := store.ListShortURL()
	// if err != nil {
	// 	log.Error("failed listing urls", "err", err)
	// 	return
	// }
	//
	// fmt.Println(urls)

	// rec, err := store.Get(
	// 	// "8a928996-8072-4117-b0fe-747d9341133c",
	// 	// "3RtW-6-YuHZ4pZGvQuN7oOkec1Q7ivLwtlo49K_9HpE=",
	// 	// "-B1TgOTGeelvkrcM5hDdAOMlBO5SKzbvEYXrPeUWpu8=",
	// )
	//
	// if err != nil {
	// 	log.Error("error getting rec", "err", err)
	// }
	// log.Info("read rec", "record", rec)

	shortener := NewShortener(&store)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./"))

	mux.Handle("GET /", fs)

	mux.HandleFunc("GET /tst", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("testing"))
	})
	mux.HandleFunc("GET /redirect/{id}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.PathValue("id")
		if shortURL == "" {
			slog.Error("no id in path")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No short URL provided"))
			return
		}

		slog.Info("redirect request received",
			"id", shortURL,
			"path", r.URL.Path,
			"full_url", r.URL.String())

		url, err := shortener.Get(shortURL)
		if err != nil {
			slog.Error("failed getting url", "err", err, "shortURL", shortURL)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Short URL not found"))
			return
		}

		if url == "" {
			slog.Error("empty url returned", "shortURL", shortURL)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL not found"))
			return
		}

		slog.Info("redirecting to", "url", url, "from_short_url", shortURL)

		// Use temporary redirect instead of permanent
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	mux.HandleFunc("POST /redirect", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to parse form"))
			return
		}

		url := r.FormValue("url")
		if url == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("URL is required"))
			return
		}

		origin := r.Header.Get("Origin")

		slog.Info("request", "origin", origin, "url", url)

		encodedUrl, err := shortener.Set(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		slog.Info("encoded url response", "encodedUrl", encodedUrl)

		shortURL := fmt.Sprintf("%s/redirect/%s", origin, encodedUrl)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(shortURL))
	})

	slog.Info("Staring server...", "port", ":8801")
	if err = http.ListenAndServe(":8801", mux); err != nil {
		slog.Error("Stopping server...", "err", err)
	}
}
