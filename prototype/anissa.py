from flask import Flask, request, current_app
import RPi.GPIO as GPIO

LED_OUT = 11

app = Flask(__name__)

@app.route('/', methods = ['GET'])
def index():
    return current_app.send_static_file("index.html")

@app.route('/turn_on', methods = ['POST'])
def turnOn():
    GPIO.output(LED_OUT, GPIO.HIGH)
    return "success"

@app.route('/turn_off', methods = ['POST'])
def turnOff():
    GPIO.output(LED_OUT, GPIO.LOW)
    return "success"

def init():
    GPIO.setmode(GPIO.BOARD)
    GPIO.setup(LED_OUT, GPIO.OUT)

def main():
    app.run(debug=True, host='0.0.0.0', port=4200)
    
if __name__ == '__main__':
    init()
    main()
    
