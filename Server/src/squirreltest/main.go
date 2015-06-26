package main

import (
	"fmt"
	"crypto/tls"
)

func main() {
	cert, err := tls.LoadX509KeyPair(
		// "src/squirrelchuckle/conf/cert.p12",
		"src/squirrelchuckle/conf/cert.pem",
		"src/squirrelchuckle/conf/key.pem")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(cert, "OK")
	}
}
