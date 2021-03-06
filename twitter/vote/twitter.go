package vote

import (
	"context"
	"github.com/gomodule/oauth1/oauth"
	"github.com/joeshaw/envdecode"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var conn net.Conn
var reader io.ReadCloser
var (
	authClient *oauth.Client
	creds *oauth.Credentials
	authSetupOnce sync.Once
	httpclient *http.Client
)

func dial(ctx context.Context,netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}

	conn = netc
	return netc, nil
}

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

func setupOAuth() {
	var ts struct{
		ConsumerKey string `env:"SP_TWITTER_KEY, required"`
		ConsumerSecret string `env:SP_TWITTER_SECRET, required`
		AccessToken string `env:SP_TWITTER_ACCESS_TOKEN,required`
		AccessSecret string `env:SP_TWITTER_ACCESS_SECRET,required`
	}
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatal("setup failure", err)
	}

	creds = &oauth.Credentials{
		Token: ts.AccessToken,
		Secret: ts.AccessSecret,
	}
	authClient = &oauth.Client{Credentials: oauth.Credentials{
		Token:ts.ConsumerKey,
		Secret: ts.ConsumerSecret,
	}}
}

func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
	authSetupOnce.Do(func() {
		setupOAuth()
		httpclient = &http.Client{
			Transport: &http.Transport{DialContext: dial},
		}
	})
	formEnc := params.Encode()
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
	return httpclient.Do(req)
}