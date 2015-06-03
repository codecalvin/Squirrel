#!/bin/bash

# go list -f {{.Dir}} pkg-url

# websocket
go get github.com/gorilla/websocket

# mux
go get github.com/gorilla/mux

# scheme
go get github.com/gorilla/schema

# mongodb
go get gopkg.in/mgo.v2

# beego
go get github.com/astaxie.beego
go get github.com/beggo/bee
go get github.com/beego/i18n

# tls
go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1
