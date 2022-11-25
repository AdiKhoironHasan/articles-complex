package integration

type integService struct {
}

func NewService() IntegServices {
	return &integService{}
}

// implement integration service
