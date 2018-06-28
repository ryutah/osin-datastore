package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/RangelReale/osin"
	"github.com/RangelReale/osin/example"
	"github.com/ryutah/osin-datastore/v1"
)

func init() {
	cfg := osin.NewServerConfig()
	cfg.AllowGetAccessRequest = true
	cfg.AllowClientSecretInParams = true

	http.HandleFunc("/initialize", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		cstorage, err := datastore.NewClientStorageForGAE(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to create client storage: %v", err)
			http.Error(w, "failed to create client storage", http.StatusInternalServerError)
			return
		}
		client := &datastore.Client{
			ID:          "1234",
			Secret:      "aabbccdd",
			RedirectUri: "http://localhost:8080/appauth/code",
		}
		if err := cstorage.Put(ctx, client); err != nil {
			log.Errorf(ctx, "failed to put client: %v", err)
			http.Error(w, "failed to put client", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/app", http.StatusFound)
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		storage, err := datastore.NewStorageForGAE(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to put initialize storage client: %v", err)
			http.Error(w, "failed to initialize storage client", http.StatusInternalServerError)
			return
		}
		defer storage.Close()

		server := osin.NewServer(cfg, storage)

		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			if !example.HandleLoginPage(ar, w, r) {
				return
			}
			ar.Authorized = true
			server.FinishAuthorizeRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		storage, err := datastore.NewStorageForGAE(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to put initialize storage client: %v", err)
			http.Error(w, "failed to initialize storage client", http.StatusInternalServerError)
			return
		}
		defer storage.Close()

		server := osin.NewServer(cfg, storage)

		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			fmt.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		storage, err := datastore.NewStorageForGAE(ctx)
		if err != nil {
			log.Errorf(ctx, "failed to put initialize storage client: %v", err)
			http.Error(w, "failed to initialize storage client", http.StatusInternalServerError)
			return
		}
		defer storage.Close()

		server := osin.NewServer(cfg, storage)

		resp := server.NewResponse()
		defer resp.Close()

		if ir := server.HandleInfoRequest(resp, r); ir != nil {
			server.FinishInfoRequest(resp, r, ir)
		}
		osin.OutputJSON(resp, w, r)
	})

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>"))
		w.Write([]byte(fmt.Sprintf(
			"<a href=\"/authorize?response_type=code&client_id=1234&state=xyz&scope=everything&redirect_uri=%s\">Login</a><br/>",
			url.QueryEscape("http://localhost:8080/appauth/code"),
		)))
		w.Write([]byte("</body></html>"))
	})

	http.HandleFunc("/appauth/code", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		r.ParseForm()

		code := r.FormValue("code")

		w.Write([]byte("<html><body>"))
		w.Write([]byte("APP AUTH - CODE<br/>"))
		defer w.Write([]byte("</body></html>"))

		if code == "" {
			w.Write([]byte("Nothing to do"))
			return
		}

		jr := make(map[string]interface{})

		// build access code url
		aurl := fmt.Sprintf(
			"/token?grant_type=authorization_code&client_id=1234&client_secret=aabbccdd&state=xyz&redirect_uri=%s&code=%s",
			url.QueryEscape("http://localhost:8080/appauth/code"),
			url.QueryEscape(code),
		)

		// if parse, download and parse json
		if r.FormValue("doparse") == "1" {
			url := fmt.Sprintf("http://localhost:8080%s", aurl)
			preq, err := http.NewRequest("POST", url, nil)
			preq.SetBasicAuth("1234", "aabbccdd")

			fclient := urlfetch.Client(ctx)
			presp, err := fclient.Do(preq)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.Write([]byte("<br/>"))
			}
			if err := json.NewDecoder(presp.Body).Decode(&jr); err != nil {
				w.Write([]byte(err.Error()))
				w.Write([]byte("<br/>"))
			}
		}

		// show json error
		if erd, ok := jr["error"]; ok {
			w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
		}

		// show json access token
		if at, ok := jr["access_token"]; ok {
			w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
		}

		w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

		// output links
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Goto Token URL</a><br/>", aurl)))

		cururl := *r.URL
		curq := cururl.Query()
		curq.Add("doparse", "1")
		cururl.RawQuery = curq.Encode()
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))
	})
}
