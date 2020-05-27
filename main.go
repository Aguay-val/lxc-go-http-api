// Copyright © 2013, 2014, The Go-LXC Authors. All rights reserved.
// Use of this source code is governed by a LGPLv2.1
// license that can be found in the LICENSE file.

// +build linux,cgo

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/lxc/go-lxc.v2"
)

var (
	lxcpath string
)

func init() {
	flag.StringVar(&lxcpath, "lxcpath", lxc.DefaultConfigPath(), "Use specified container path")
	flag.Parse()
}
func listContainer() []string {
	var ctList []string
	c := lxc.DefinedContainers(lxcpath)
	for i := range c {
		ctList = append(ctList, c[i].Name())
		c[i].Release()
	}

	return ctList
}
func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path	[1:])
	containerName := listContainer()
	fmt.Fprintf(w, "Defined containers:\n")
	fmt.Fprintf(w, "%s", containerName)
}
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
