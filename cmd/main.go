package main

import (
	"bytes"
	"flag"
	"ldaputil"
	"ldaputil/assets"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-routing/content"

	routing "github.com/go-ozzo/ozzo-routing"
)

var (
	username string
	password string
)

func init() {
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.Parse()
}

func main() {
	cfg, err := ldaputil.ParseConfig("./config.yaml")
	if err != nil {
		log.Printf("fail to parse config: %s.", err)
		os.Exit(1)
	}

	log.Println("parsed config.", cfg)

	if username != "" {
		err = ldaputil.UserSetPass(cfg, username, password)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	router := routing.New()
	router.Post("/update", changePass(cfg))
	router.Get("/*", func(c *routing.Context) error {
		name := c.Request.URL.Path
		name = strings.TrimPrefix(name, "/")
		if name == "" {
			name = "index.html"
		}
		b, err := assets.HTML.ReadFile(name)
		if err != nil {
			return err
		}

		http.ServeContent(c.Response, c.Request, filepath.Base(name), time.Now(), bytes.NewReader(b))
		return nil
	})

	log.Printf("listen on %s.", cfg.Listen)

	panic(http.ListenAndServe(cfg.Listen, router))
}

func changePass(cfg ldaputil.Config) routing.Handler {
	return func(c *routing.Context) error {
		c.SetDataWriter(&content.JSONDataWriter{})
		username := c.Form("username")
		password := c.Form("password")
		newPassword := c.Form("new_password")
		if username == "" || password == "" || newPassword == "" {
			return c.Write(map[string]string{"status": "failed"})
		}
		err := ldaputil.UserChangePass(cfg, username, password, newPassword)
		if err != nil {
			return c.Write(map[string]string{"status": "failed"})
		}
		return c.Write(map[string]string{"status": "success"})
	}
}
