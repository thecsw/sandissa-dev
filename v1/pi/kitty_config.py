
# MQTT certificates
CA_CERT = "/home/pi/shell/ca/ca.crt"
CLIENT_CERT = "/home/pi/shell/new_client/new_client.crt"
CLIENT_KEY = "/home/pi/shell/new_client/new_client.key"
TLS_VERSION = 2

# MQTT broker options
TOPIC_NAME = "sensors/Temp"
HOSTNAME = "mqtt.sandyuraz.com"
PORT = 8883

# MQTT user options
CLIENT_ID = "sandissa-secret-me"
USERNAME = "kitty"
PASSWORD = None

# MQTT publish/subscribe options
QOS_PUBLISH = 2
QOS_SUBSCRIBE = 2
