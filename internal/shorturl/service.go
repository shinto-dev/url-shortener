package shorturl

type CreateRequest struct {
	OriginalURL string
	CustomAlias string
}

type ShortURL struct {
	OriginalURL   string
	ShortURLToken string
}

//  mockery --name=Service --inpackage
type Service interface {
	Create(CreateRequest) (ShortURL, error)
}
