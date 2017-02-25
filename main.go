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
)

func HandleProxy(target string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		p := &httputil.ReverseProxy{Director: director}
		p.ServeHTTP(w, r)
	}
}

func main() {

	Config := map[string]string{
		"/": "127.0.0.1:8787",
	}

	for Path, Target := range Config {
		// avoid add comments as route
		http.HandleFunc(Path, HandleProxy(Target))
		log.Printf("%s > %s", Path, Target)
	}

	Address := fmt.Sprintf("%s:%s", "", "8080")
	log.Print("start listening on " + Address)
	http.ListenAndServe(Address, nil)
}
