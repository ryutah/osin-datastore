# osin-datastore
[![CircleCI](https://circleci.com/gh/ryutah/osin-datastore/tree/master.svg?style=shield&circle-token=:circle-token)](https://circleci.com/gh/ryutah/osin-datastore/tree/master)
[![GoDoc](https://godoc.org/github.com/ryutah/osin-datastore/v1?status.svg)](https://godoc.org/github.com/ryutah/osin-datastore/v1)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE.md)

A Google Cloud Datstore storage for [RangelReale/osin](https://github.com/RangelReale/osin).

This storage can be used on Google App Engine, Google Compute Engine, Googke Kubernetes Engine, On-premises environment and so on.

## Install
```console
$ go get -u github.com/ryutah/osin-datastore/v1
```

## Usage
### Google App Engine (Standard Edition)
```go
package main

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/RangelReale/osin"
	"github.com/ryutah/osin-datastore/v1"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		storage, err := datastore.NewStorageForGAE(ctx)
		if err != nil {
			http.Error(w, "failed to initialize storage client", http.StatusInternalServerError)
			return
		}
		defer storage.Close()

		server := osin.NewServer(osin.NewServerConfig(), storage)

        // do sometihng.
    }
}
```

### Other Platforms
```go
package main

import (
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/ryutah/osin-datastore/v1"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		storage, err := datastore.NewStorage(r.Context())
		if err != nil {
			http.Error(w, "failed to initialize storage client", http.StatusInternalServerError)
			return
		}
		defer storage.Close()

		server := osin.NewServer(osin.NewServerConfig(), storage)

        // do sometihng.
    }
}
```

[Full Examples](example/gaese)
