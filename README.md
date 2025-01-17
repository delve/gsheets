# gsheets

[![CircleCI](https://circleci.com/gh/delve/gsheets/tree/master.svg?style=shield&circle-token=30f6e95108024e7a0562f630c69209783e5086ec)](https://circleci.com/gh/delve/gsheets/tree/master)
[![codecov](https://codecov.io/gh/delve/gsheets/branch/master/graph/badge.svg)](https://codecov.io/gh/delve/gsheets)
[![GoDoc](https://godoc.org/github.com/delve/gsheets?status.svg)](https://godoc.org/github.com/delve/gsheets)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

A golang wrapper package for `golang.org/x/oauth2` and `google.golang.org/api/sheets/v4`.
You can easily manipulate spreadsheets.

**!!! Only for personal use !!!**

## Installation

```bash
go get github.com/delve/gsheets
```

## Requirement

This package uses Google OAuth2.0. So before executing tool, you have to prepare credentials.json.
See [Go Quickstart](https://developers.google.com/sheets/api/quickstart/go), or [Blog (Japanese)](https://medium.com/veltra-engineering/how-to-use-google-sheets-api-with-golang-9e50ee9e0abc) for the details.

## Usage

### Create Cache

If you want to use the cache, initialize the context.
If you are updating sheets, you should not use Cache.

```go
ctx := gsheets.WithCache(ctx)
```

### Create New Client

```go
client, err := gsheets.New(ctx, `{"credentials": "json"}`, `{"token": "json"}`)
```

```go
client, err := gsheets.NewForCLI(ctx, "credentials.json")
```

If you are updating sheets, create a client with `ClientWritable` option.

```go
client, err := gsheets.New(ctx, `{"credentials": "json"}`, `{"token": "json"}`, gsheets.ClientWritable())
```

### Get Sheet Information

```go
func (*Client) GetTitle(ctx context.Context, spreadsheetID string) (string, error)
```

```go
func (*Client) GetSheetNames(ctx context.Context, spreadsheetID string) ([]string, error)
```

```go
func (*Client) GetSheet(ctx context.Context, spreadsheetID, sheetName string) (Sheet, error)
```

### Update Sheet Values

```go
func (c *Client) Update(ctx context.Context, spreadsheetID, sheetName string, rowNo int, values []interface{}) error
```

```go
func (c *Client) BatchUpdate(ctx context.Context, spreadsheetID string, updateValues ...UpdateValue) error
```

### Manipulate Sheet Values

If the index is out of range, `Value` method returns empty string.

```go
s, err := client.GetSheet(ctx, "spreadsheetID", "sheetName")
if err != nil {
  return err
}

fmt.Println(s.Value(row, clm))

for _, r := range s.Rows() {
  fmt.Println(r.Value(clm))
}
```
