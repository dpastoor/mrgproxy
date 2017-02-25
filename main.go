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
)

//HandleProxy handles the rpoxy
func HandleProxy(target string, prefix string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("URI:", r.RequestURI)
		director := func(req *http.Request) {
			req = r
			req.Host = target // for cors
			req.URL.Host = target
			req.URL.Scheme = "http"
			fmt.Println("original URI:", r.RequestURI)
			req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, prefix)
			fmt.Println("request now on", req.URL.Path)
		}
		p := &httputil.ReverseProxy{Director: director}
		fmt.Println("r path:", r.URL.Path)
		p.ServeHTTP(w, r)
	}
}

func main() {
	// not that this must be /<versiontag>/
	// given /<versiontag> the relative rstudio requests to /rstudio go to the original
	// domain and therefore fail to be managed properly
	Config := map[string]string{
		"/latest/":  "127.0.0.1:8787",
		"/v0.7.10/": "127.0.0.1:8788",
	}

	for Path, Target := range Config {
		// avoid add comments as route
		http.HandleFunc(Path, HandleProxy(Target, Path))
		log.Printf("%s > %s", Path, Target)
	}

	Address := fmt.Sprintf("%s:%s", "", "8080")
	log.Print("start listening on " + Address)
	http.ListenAndServe(Address, nil)
}
