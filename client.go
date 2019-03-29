package gsheets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

// Client is a gsheets client.
type Client struct {
	credentials, token string
	srv                *sheets.Service
}

// New returns a gsheets client.
func New(ctx context.Context, credentials, token string) (*Client, error) {

	config, err := google.ConfigFromJSON([]byte(credentials), "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fmt.Errorf("Unable to parse json to config: %v", err)
	}
	tok := &oauth2.Token{}
	if err := json.NewDecoder(strings.NewReader(token)).Decode(tok); err != nil {
		return nil, fmt.Errorf("Unable to parse json to token: %v", err)
	}
	srv, err := sheets.New(config.Client(ctx, tok))
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	return &Client{
		credentials: credentials,
		token:       token,
		srv:         srv,
	}, nil
}

// NewForCLI returns a gsheets client.
// This function is intended for CLI tools.
func NewForCLI(ctx context.Context, authFile string) (*Client, error) {

	cb, err := ioutil.ReadFile(authFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	tokenFile := "token.json"
	tb, err := ioutil.ReadFile(tokenFile)

	var token string
	if err == nil {
		token = string(tb)
	} else {
		// if there are no token file, get from Web
		config, err := google.ConfigFromJSON(cb, "https://www.googleapis.com/auth/spreadsheets.readonly")
		if err != nil {
			return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
		}

		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)

		var authCode string
		if _, err := fmt.Scan(&authCode); err != nil {
			return nil, fmt.Errorf("Unable to read authorization code: %v", err)
		}

		tok, err := config.Exchange(oauth2.NoContext, authCode)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve token from web: %v", err)
		}

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tok)
		token = b.String()

		// save token
		fmt.Printf("Saving credential file to: %s\n", tokenFile)
		f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		defer f.Close()
		if err != nil {
			return nil, fmt.Errorf("Unable to cache oauth token: %v", err)
		}
		fmt.Fprint(f, token)
	}

	return New(ctx, string(cb), token)
}