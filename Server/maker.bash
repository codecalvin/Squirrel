# make pem
openssl x509 -in aps_development.cer -inform der -out SquirrelCert.pem

# make private key pem
openssl pkcs12 -nocerts -out SquirrelKey.pem -in SquirrelKey.p12

# test
openssl s_client -connect gateway.sandbox.push.apple.com:2195 -cert SquirrelCert.pem -key SquirrelKey.pem
