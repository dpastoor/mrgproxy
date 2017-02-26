// simple http reverse proxy
// Copyright (C) 2017  geosoft1  geosoft1@gmail.com
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/pressly/chi"
)

// CreateProxyHandler creates a new reverse proxy instance while rewriting to strip the leading prefix
func CreateProxyHandler(target string, prefix string) http.Handler {
	fmt.Println("got here")
	director := func(req *http.Request) {
		req.Host = target // for cors
		req.URL.Host = target
		req.URL.Scheme = "http"
		req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, prefix)
	}
	p := &httputil.ReverseProxy{Director: director}
	return p
}

func main() {
	// not that this must be /<versiontag>/
	// given /<versiontag> the relative rstudio requests to /rstudio go to the original
	// domain and therefore fail to be managed properly
	Config := map[string]string{
		"/latest/":  "127.0.0.1:8787",
		"/v0.7.10/": "127.0.0.1:8788",
	}

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("future landing page"))
	})

	for Path, Target := range Config {
		// Mount expects a http.Handler, namely something that satisfies the interface with method ServeHTTP
		// given the reverseProxy can be called such that p.ServeHTTP(w, r)
		r.Mount(Path, CreateProxyHandler(Target, Path))
	}

	r.Get("/:newRoute", func(w http.ResponseWriter, req *http.Request) {
		path := "/" + chi.URLParam(req, "newRoute") + "/"
		fmt.Println("about to mount to path: ", path)
		r.Mount(path, CreateProxyHandler("127.0.0.1:8787", path))
	})

	Address := fmt.Sprintf("%s:%s", "", "8080")
	log.Print("start listening on " + Address)
	http.ListenAndServe(":8080", r)
}
