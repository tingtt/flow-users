package github

type Application struct {
	ClientId     string
	ClientSecret string
}

func New(clientId string, clientSecret string) (*Application, error) {
	return &Application{clientId, clientSecret}, nil
}
