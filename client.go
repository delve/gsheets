package gsheets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

// Client is a gsheets client.
type Client struct {
	// credentials, token string
	options []option.ClientOption
	srv     *sheets.Service

	scope string
}

// Type and optionFunc are now only used in NewForCli, which i am not using atm
// ClientOption is an option function.
type ClientOption func(c *Client) *Client

// ClientWritable is an option to change client writable.
func ClientWritable() func(c *Client) *Client {
	return func(c *Client) *Client {
		c.scope = "https://www.googleapis.com/auth/spreadsheets"
		return c
	}
}

// New returns a gsheets client.
func New(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	client := &Client{
		options: opts,
	}
	return new(ctx, client)
}

func new(ctx context.Context, initialClient *Client) (*Client, error) {
	srv, err := sheets.NewService(ctx, initialClient.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %v", err)
	}
	initialClient.srv = srv
	return initialClient, nil
}

// TODO: Convert this to the way New() works with direct google clientoptions
// NewForCLI returns a gsheets client.
// This function is intended for CLI tools.
func NewForCLI(ctx context.Context, authFile string, opts ...ClientOption) (*Client, error) {

	client := &Client{
		scope: "https://www.googleapis.com/auth/spreadsheets.readonly",
	}

	for _, opt := range opts {
		client = opt(client)
	}

	cb, err := os.ReadFile(authFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	tokenFile := "token.json"
	tb, err := os.ReadFile(tokenFile)

	var token string
	if err == nil {
		token = string(tb)
	} else {
		// if there are no token file, get from Web
		config, err := google.ConfigFromJSON(cb, client.scope)
		if err != nil {
			return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
		}

		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		fmt.Printf("Go to the following link in your browser then type the "+
			"authorization code: \n%v\n", authURL)

		var authCode string
		if _, err := fmt.Scan(&authCode); err != nil {
			return nil, fmt.Errorf("unable to read authorization code: %v", err)
		}

		tok, err := config.Exchange(context.Background(), authCode)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
		}

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tok)
		token = b.String()

		// save token
		fmt.Printf("Saving credential file to: %s\n", tokenFile)
		f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		defer f.Close()
		if err != nil {
			return nil, fmt.Errorf("unable to cache oauth token: %v", err)
		}
		fmt.Fprint(f, token)
	}

	// changed to alloiw compile
	// return new(ctx, string(cb), token, client)
	return new(ctx, client)
}
