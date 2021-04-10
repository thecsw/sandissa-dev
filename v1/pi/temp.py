# https://learn.adafruit.com/adafruits-raspberry-pi-lesson-11-ds18b20-temperature-sensing/software
# https://www.earth.li/~noodles/blog/2018/05/rpi-mqtt-temp.html
import paho.mqtt.publish as publish
import glob
import time
import kitty_config

# Directories and files for w1-bus temperature sensors
W1_DIRECTORY = "/sys/bus/w1/devices/"
W1_TEMP_FILES = glob.glob(W1_DIRECTORY + "28*")[0]
W1_TEMP_SLAVE = W1_TEMP_FILES + "/w1_slave"


def read_temp_raw():
    f = open(W1_TEMP_SLAVE, "r")
    lines = f.readlines()
    f.close()
    return lines


def read_temp_celsius():
    lines = read_temp_raw()
    lines_stripped = lines[0].strip()
    # Check that the slice has at least 3 elements
    while len(lines_stripped) >= 3 and lines_stripped[-3:] != "YES":
        time.sleep(0.2)
        lines = read_temp_raw()
    equals_pos = lines[1].find("t=")
    if equals_pos != -1:
        temp_string = lines[1][equals_pos + 2 :]
        temp_c = float(temp_string) / 1000.0
        # Optional farenheit conversion
        # temp_f = temp_c * 9.0 / 5.0 + 32.0
        return temp_c


def main():
    while True:
        temp = read_temp_celsius()
        print(time.time(), temp)
        try:
            publish.single(
                kitty_config.TOPIC_NAME,
                payload=str(temp),
                hostname=kitty_config.HOSTNAME,
                port=kitty_config.PORT,
                client_id=kitty_config.CLIENT_ID,
                qos=kitty_config.QOS_PUBLISH,
                auth={"username": kitty_config.USERNAME, "password": kitty_config.PASSWORD},
                tls={
                    "ca_certs": kitty_config.CA_CERT,
                    "certfile": kitty_config.CLIENT_CERT,
                    "keyfile": kitty_config.CLIENT_KEY,
                    "tls_version": kitty_config.TLS_VERSION,
                },
            )
        except Exception as e:
            print("MQTT failed to publish temperature:", e)
        time.sleep(1)
