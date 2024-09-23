package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
)

func NewWebDavHandler(prefix, dir string) http.Handler {
	return &webdav.Handler{
		Prefix:     fmt.Sprintf("/%s", prefix),
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(req *http.Request, err error) {
			if err != nil {
				logrus.Error(err)
			}
		},
	}
}
