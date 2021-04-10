#!/bin/sh

if [ $# -ne 3 ]
then
	echo "./pub.sh TOPIC QOS:[0..2] 'MESSAGE'"
	exit
fi

mosquitto_pub -d -h mqtt.sandyuraz.com -p 8883 --tls-use-os-certs --cafile ./ca/ca.crt --cert ./client/client.crt --key ./client/client.key --tls-version tlsv1.2 -I sandissa-secret-me -u sandy -m $3 -q $2 -t $1
