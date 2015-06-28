package main

import (
	"fmt"
	"crypto/tls"
)

func main() {
	cert, err := tls.LoadX509KeyPair(
		"src/squirrelchuckle/conf/SquirrelCert.pem.private",
		"src/squirrelchuckle/conf/SquirrelKey.u.pem.private")

	if err != nil {
		fmt.Println(err, "fail")
	} else {
		fmt.Println(cert, "OK")
	}
}
