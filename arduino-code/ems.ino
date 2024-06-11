#include <Ethernet.h>
#include <MQTT.h>
#include <Wire.h>
#include <Adafruit_BMP085.h>
#include <DHT11.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <ArduinoJson.h>

// Ethernet and MQTT setup
byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
byte ip[] = { 192, 168, 3, 2 };  // Your network IP
EthernetClient net;
MQTTClient client;

// Sensor setup
int lightPin = A1;
int soundPin = A0;
int PIRpin = 4;
int vibPin = 6;
Adafruit_BMP085 bmp;
Adafruit_MPU6050 mpu;
DHT11 dht11(5);

// Variables to store sensor data
int lightValue = 0;
int soundValue = 0;
int movement = 0;
int vibration = 0;
int temperatureDHT = 0;
int humidity = 0;
float temperatureBMP = 0.0;
float AccX;
float AccY;
float AccZ;

void setup() {
  Serial.begin(9600);  // Start serial communication
  Ethernet.begin(mac, ip);  // Initialize Ethernet connection
  client.begin("192.168.3.1", 1883, net);  // Initialize MQTT client with broker IP and port
  client.onMessage(messageReceived);  // Set callback function for incoming messages

  // Initialize pins
  pinMode(PIRpin, INPUT);
  pinMode(vibPin, INPUT);

  // Initialize BMP085 sensor
  if (!bmp.begin()) {
    Serial.println("Could not find a valid BMP085 sensor, check wiring!");
    while (1) {}  // Halt if sensor initialization fails
  }

  // Initialize MPU6050 sensor
  if (!mpu.begin()) {
    Serial.println("Failed to find MPU6050 chip");
    while (1) {  // Halt if sensor initialization fails
      delay(10);
    }
  }

  // Set MPU6050 parameters
  mpu.setAccelerometerRange(MPU6050_RANGE_8_G);
  mpu.setGyroRange(MPU6050_RANGE_500_DEG);
  mpu.setFilterBandwidth(MPU6050_BAND_21_HZ);
  delay(20);

  connect();  // Connect to MQTT broker
}

void loop() {
  client.loop();  // Keep MQTT client loop running

  if (!client.connected()) {
    connect();  // Reconnect if connection is lost
  }

  // Read data from MPU6050
  sensors_event_t a, g, temp;
  mpu.getEvent(&a, &g, &temp);
  AccX = a.acceleration.x;
  AccY = a.acceleration.y;
  AccZ = a.acceleration.z;

  // Read data from all sensors
  dht11.readTemperatureHumidity(temperatureDHT, humidity);
  temperatureBMP = bmp.readTemperature();
  lightValue = analogRead(lightPin);
  soundValue = analogRead(soundPin);
  movement = digitalRead(PIRpin);
  vibration = digitalRead(vibPin);

  // Prepare JSON data
  StaticJsonDocument<256> doc;
  doc["temperature_dht"] = temperatureDHT;
  doc["temperature_bmp"] = temperatureBMP;
  doc["humidity"] = humidity;
  doc["light"] = lightValue;
  doc["sound"] = soundValue;
  doc["pressure"] = bmp.readPressure();
  doc["movement"] = movement;
  doc["vibration"] = vibration;
  doc["accX"] = AccX;
  doc["accY"] = AccY;
  doc["accZ"] = AccZ;

  // Serialize JSON data to a buffer
  char jsonBuffer[256];
  serializeJson(doc, jsonBuffer);

  // Publish sensor data to MQTT broker
  client.publish("/sensors", jsonBuffer);

  delay(1000);  // Wait for 1 second before next reading
}

// Callback function for incoming MQTT messages
void messageReceived(String &topic, String &payload) {
  Serial.println("incoming: " + topic + " - " + payload);
}

// Connect to MQTT broker
void connect() {
  Serial.print("connecting...");
  while (!client.connect("arduino", "", "")) {  // Replace with your MQTT broker credentials
    Serial.print(".");
    delay(1000);
  }
  Serial.println("\nconnected!");
}
