package mcapi

type Client struct {
	baseURL    string
	privateKey string
	publicKey  string
}

func NewClient(url, privateKey, publicKey string) *Client {
	return &Client{
		baseURL:    url,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}
