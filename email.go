package mailinabox

import (
	"bufio"
	"encoding/json"
	"github.com/dghubble/sling"
	"net/http"
	"strings"
)

// MailService accesses mailinabox mail management endpoints.
type MailService struct {
	sling  *sling.Sling
	client *http.Client
}

// NewMailService returns a new MailService
func NewMailService(sling *sling.Sling, httpClient *http.Client) *MailService {
	return &MailService{
		sling:  sling.Path("mail/users/"),
		client: httpClient,
	}
}

// Params used to modify response output
type Params struct {
	Format string `json:"format"`
}

// Form used to modify email accounts
type Form struct {
	Email    string `json:"email,omitempty" url:"email,omitempty"`
	Password string `json:"password,omitempty" url:"password,omitempty"`
	Quota    string `json:"quota,omitempty" url:"quota,omitempty"`
}

// User is an email user
type User struct {
	Email string `json:"email"`
	Quota string `json:"quota"`
}

const (
	MAIL_USER_ADDED       = "mail user added"
	MAIL_USER_REMOVED     = "mail user removed"
	THATS_NOT_A_USER      = "That's not a user ()."
	INVALID_EMAIL_ADDRESS = "Invalid email address."
	NO_PASSWORD_MESSAGE   = "No password provided"
)

// ListByDomain filters email accounts by domain.
func (o *MailService) ListByDomain(domain string) ([]*User, *http.Response, error) {
	var errMessage json.RawMessage
	resp, err := o.sling.New().Get("").QueryStruct(Params{Format: "json"}).Receive(nil, &errMessage)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	list := []*User{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		unsplit := strings.Split(line, "@")
		if strings.Contains(unsplit[1], domain) {
			unsplit = strings.Split(line, "|")
			list = append(list, &User{unsplit[1], unsplit[0]})
		}
	}
	err = scanner.Err()
	return list, resp, err
}

// List all email accounts.
func (o *MailService) List(email string) ([]*User, *http.Response, error) {
	var errMessage json.RawMessage

	resp, err := o.sling.New().Get("").QueryStruct(Params{Format: "json"}).Receive(nil, &errMessage)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	list := []*User{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		unsplit := strings.Split(line, "|")
		list = append(list, &User{unsplit[1], unsplit[0]})
	}
	err = scanner.Err()
	return list, resp, err
}

// Add a new email account.
func (o *MailService) Add(email, password, quota string) (*[]byte, *http.Response, error) {
	req, err := o.sling.New().Post("add").BodyForm(Form{email, password, quota}).Request()
	if err != nil {
		return nil, nil, nil
	}
	resp, err := o.client.Do(req)
	return responseHandler(resp, err)
}

// Remove an email account
func (o *MailService) Remove(email string) (*[]byte, *http.Response, error) {
	req, err := o.sling.New().Post("remove").BodyForm(Form{Email: email}).Request()
	if err != nil {
		return nil, nil, nil
	}
	resp, err := o.client.Do(req)
	return responseHandler(resp, err)
}

// SetQuota sets the quota of an Email account
func (o *MailService) SetQuota(email, quota string) (*[]byte, *http.Response, error) {

	req, err := o.sling.New().Post("quota").BodyForm(Form{Email: email, Quota: quota}).Request()
	if err != nil {
		return nil, nil, nil
	}
	resp, err := o.client.Do(req)
	return responseHandler(resp, err)

}

func (o *MailService) SetPassword(email, password string) (*[]byte, *http.Response, error) {

	req, err := o.sling.New().Post("password").BodyForm(Form{Email: email, Password: password}).Request()
	if err != nil {
		return nil, nil, nil
	}
	resp, err := o.client.Do(req)
	return responseHandler(resp, err)

}
