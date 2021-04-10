import paho.mqtt.client as mqtt
import RPi.GPIO as GPIO
import ssl

import kitty_config

GPIO.setmode(GPIO.BOARD)
GPIO.setwarnings(False)
GPIO.setup(11, GPIO.OUT)
GPIO.output(11, GPIO.LOW)


def on_connect(client, userdata, flag, rc):
    client.subscribe("controls/LED")


def on_message(client, userdata, msg):
    if "on" in msg.payload:
        GPIO.output(11, GPIO.HIGH)

    if "off" in msg.payload:
        GPIO.output(11, GPIO.LOW)


client = mqtt.Client(client_id=kitty_config.CLIENT_ID)
client.on_connect = on_connect
client.on_message = on_message

# update cert paths
client.tls_set(
    ca_certs=kitty_config.CA_CERT,
    certfile=kitty_config.CLIENT_CERT,
    keyfile=kitty_config.CLIENT_KEY,
    tls_version=kitty_config.TLS_VERSION,
)

# set the username and password
client.username_pw_set(username=kitty_config.USERNAME, password=kitty_config.PASSWORD)

# update host to be broker host
client.connect(kitty_config.HOSTNAME, kitty_config.PORT, 60)

client.loop_forever()
GPIO.cleanup()
