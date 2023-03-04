#!/bin/bash
sleep 2
cp /certificates/serverCert.crt /usr/local/share/ca-certificates/serverCert.crt
update-ca-certificates
