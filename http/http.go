package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const enphotoURL = "https://en-photo.net"
const loginURL = enphotoURL + "/login"

type Accessor interface {
	Login(username, password string) error
	GetAlbum(url string) (string, error)
	GetDataSrc(dataSrc string) (string, error)
	GetThumbnail(url string) ([]byte, error)
}

type accessorImpl struct {
	client *http.Client
}

func New() (Accessor, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cookiejar; %w", err)
	}
	client := &http.Client{
		Jar: jar,
	}
	return &accessorImpl{client: client}, nil
}

func (a *accessorImpl) Login(username, password string) error {
	token, err := a.getLogin()
	if err != nil {
		return err
	}
	if err := a.postLogin(username, password, token); err != nil {
		return err
	}
	return nil
}

func (a *accessorImpl) getLogin() (string, error) {
	resp, err := a.client.Get(loginURL)
	if err != nil {
		return "", fmt.Errorf("failed to HTTP GET; %s; %w", loginURL, err)
	}
	defer resp.Body.Close()
	if err := a.isErrorCode(resp); err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML; %w", err)
	}
	selection := doc.Find("[name=\"_token\"]")
	token, ok := selection.Attr("value")
	if !ok {
		return "", fmt.Errorf("failed to find token; %w", err)
	}
	return token, nil
}

func (a *accessorImpl) postLogin(username, password, token string) error {
	values := url.Values{
		"_token":   {token},
		"email":    {username},
		"password": {password},
	}
	resp, err := a.client.PostForm(loginURL, values)
	if err != nil {
		return fmt.Errorf("failed to HTTP POST; %w", err)
	}
	defer resp.Body.Close()
	return a.isErrorCode(resp)
}

func (a *accessorImpl) GetAlbum(url string) (string, error) {
	return a.httpGetString(url)
}

func (a *accessorImpl) GetDataSrc(dataSrc string) (string, error) {
	url := enphotoURL + dataSrc
	return a.httpGetString(url)
}

func (a *accessorImpl) GetThumbnail(url string) ([]byte, error) {
	return a.httpGet(url)
}

func (a *accessorImpl) httpGet(url string) ([]byte, error) {
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to HTTP GET; %s; %w", url, err)
	}
	defer resp.Body.Close()
	if err := a.isErrorCode(resp); err != nil {
		return nil, err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response; %w", err)
	}
	return respBytes, nil
}

func (a *accessorImpl) httpGetString(url string) (string, error) {
	bs, err := a.httpGet(url)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (a *accessorImpl) isErrorCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non 2XX HTTP status code; %d; %s", resp.StatusCode, resp.Status)
	}
	return nil
}
