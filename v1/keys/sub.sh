#!/bin/sh

if [ $# -ne 2 ]
then
	echo "./sub.sh TOPIC QOS:[0..2]"
	exit
fi

mosquitto_sub -d -h mqtt.sandyuraz.com -p 8883 --tls-use-os-certs --cafile ./ca/ca.crt --cert ./client/client.crt --key ./client/client.key --tls-version tlsv1.2 -I sandissa-secret-me -u sandy -q $2 -t $1
