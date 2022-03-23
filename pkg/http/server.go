package http

const (
	Scheme        = "https"
	Host          = "localhost"
	Port          = "8080"
	TimeoutSecond = 5
)

type Config struct {
	Scheme      string
	Host        string
	Port        string
	URL         string
	AccessToken string
}
