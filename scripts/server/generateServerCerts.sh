#!/bin/bash
CERT_PATH="scripts"

openssl req -new -subj "/C=US/ST=California/CN=localhost" \
	-newkey rsa:2048 -nodes -keyout "$CERT_PATH/serverCert.key" -out "$CERT_PATH/serverCert.csr"
pwd
openssl x509 -req -days 365 -in "$CERT_PATH/serverCert.csr" \
	-signkey "$CERT_PATH/serverCert.key" -out "$CERT_PATH/certificates/serverCert.crt" \
	-extfile "$CERT_PATH/selfSignedCert.ext"
