#!/bin/sh

if [ $# -ne 1 ]
then
	echo "./gen_client.sh CLIENT_NAME"
	exit
fi

mkdir -p $1
cd $1
openssl genrsa -out $1.key 4096
openssl req -new -out $1.csr -key $1.key -subj "/C=US/ST=KS/L=Lawrence/O=KU/OU=EECS/CN=mqtt.sandyuraz.com"
openssl x509 -req -sha512 -in $1.csr -CA ../ca/ca.crt -CAkey ../ca/ca.key -CAserial ../ca/ca.srl -out $1.crt -days 100
