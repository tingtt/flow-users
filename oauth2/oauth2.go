package oauth2

type Provider int

const (
	ProviderGitHub Provider = iota + 1
	ProviderGoogle
	ProviderTwitter
)

func (p Provider) String() string {
	switch p {
	case ProviderGitHub:
		return "github"
	case ProviderGoogle:
		return "google"
	case ProviderTwitter:
		return "twitter"
	default:
		return ""
	}
}
