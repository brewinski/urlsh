package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
)

type ShortenerService struct {
	store Storage
}

func NewShortener(store Storage) ShortenerService {
	return ShortenerService{
		store,
	}
}

func (ss *ShortenerService) Get(id string) (string, error) {
	shortUrl, err := getLookupTableURL(ss.store, id)
	if err != nil {
		return "", err
	}

	return shortUrl.URL, nil
}

func (ss *ShortenerService) Set(url string) (string, error) {
	encodedUrl, err := encodeURL(url)
	if err != nil {
		return "", err
	}

	err = updateLookupTableURL(ss.store, url, encodedUrl)
	if err != nil {
		return encodedUrl, err
	}

	return encodedUrl, nil
}

func encodeURL(url string) (string, error) {
	hashUrl := sha256.New()

	_, err := hashUrl.Write([]byte(url))
	if err != nil {
		return "", fmt.Errorf("failed to hash url string: %w", err)
	}

	bs := hashUrl.Sum(nil)

	slog.Info("sha", "sha", bs)

	encodedUrl := base64.RawURLEncoding.EncodeToString(bs)
	return encodedUrl, nil
}

func getLookupTableURL(store Storage, encodedPath string) (ShortURL, error) {
	url, err := store.Get(encodedPath)

	if err != nil {
		return url, fmt.Errorf("error looking up encoded path: %w", err)
	}

	if url.ID != encodedPath {
		return url, fmt.Errorf("path not stored in database, %s, %v", encodedPath, url)
	}

	return url, nil
}

func updateLookupTableURL(store Storage, url, encoding string) error {
	storedUrl, err := store.Get(encoding)
	if err != nil {
		return fmt.Errorf("error reading store: %w", err)
	}

	if storedUrl.ID != encoding {
		store.Set(encoding, url)
		return nil
	}

	if storedUrl.URL == url {
		return nil
	}

	for i := 0; true; i++ {
		newKey := encoding + string(i)
		stored, err := store.Get(newKey)
		if err != nil {
			return fmt.Errorf("failed to read storage for key: %w", err)
		}

		if stored.ID != "" {
			err := store.Set(newKey, url)
			if err != nil {
				return fmt.Errorf("failed to read storeage for new key: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("URL Collision, this url's encoding matched a different urls encoding")
}
