package mailinabox

import (
	"net/http"

	"github.com/dghubble/sling"
	"github.com/ernesto-jimenez/httplogger"
)

// Client is a Mailinabox client for accessing the admin api.
type Client struct {
	sling *sling.Sling

	// All Mailinabox services
	Mail *MailService
}

// NewClient returns a new Client.
func NewClient(url, username, password string, httpClients ...*http.Client) *Client {
	var httpClient *http.Client

	if len(httpClients) > 0 {
		httpClient = httpClients[0]
	} else {
		// httpClient = http.DefaultClient
		httpClient = &http.Client{
			Transport: httplogger.NewLoggedTransport(http.DefaultTransport, newLogger()),
		}
	}
	base := sling.New().Client(httpClient).SetBasicAuth(username, password).Base(url + "/admin")

	return &Client{
		Mail: NewMailService(base.New(), httpClient),
	}
}
