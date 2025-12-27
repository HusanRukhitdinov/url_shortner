package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/Go11Group/url_shortner/internal/repo/sqlc"
	"github.com/Go11Group/url_shortner/internal/repo/storage"
)

type UrlUseCaseI interface {
	Shorten(ctx context.Context, originalUrl string) (string, error)
	GetOriginal(ctx context.Context, code string) (string, error)
}
type UrlUseCase struct {
	storage storage.StorageI
}

func NewUrlUseCase(storage storage.StorageI) UrlUseCaseI {
	return &UrlUseCase{storage: storage}
}
func (u *UrlUseCase) Shorten(ctx context.Context, originalUrl string) (string, error) {
	b := make([]byte, 3)
	rand.Read(b)
	code := hex.EncodeToString(b)

	_, err := u.storage.Url().CreateUrl(ctx, &sqlc.CreateUrlParams{
		OriginalUrl: originalUrl,
		ShortCode:   code,
	})
	if err != nil {
		return "", err
	}
	return code, nil
}

func (u *UrlUseCase) GetOriginal(ctx context.Context, code string) (string, error) {
	log.Printf("GetOriginal: looking for code: '%s'", code)
	data, err := u.storage.Url().GetUrlByCode(ctx, code)
	if err != nil {
		log.Printf("GetOriginal: error fetching code '%s': %v", code, err)
		return "", err
	}
	_ = u.storage.Url().IncrementClicks(ctx, data.ID)

	return data.OriginalUrl, nil
}
