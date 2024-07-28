#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <Wire.h>

Adafruit_MPU6050 mpu1;
Adafruit_MPU6050 mpu2;

const int lightPin = A0;
const int soundPin1 = A1;
const int soundPin2 = A2;
const unsigned long soundSampleWindow = 50; // Sample window width in ms (50 ms)

// Function to get sound peak-to-peak value
unsigned int getSoundPeakToPeak(int soundPin) {
  unsigned long startMillis = millis(); // Start of sample window
  unsigned int peakToPeak = 0;   // peak-to-peak level

  unsigned int signalMax = 0;
  unsigned int signalMin = 1024;

  // collect data for 50 ms and then calculate peak-to-peak value
  while (millis() - startMillis < soundSampleWindow) {
    int sample = analogRead(soundPin);
    if (sample < 1024) { // toss out spurious readings
      if (sample > signalMax) {
        signalMax = sample;  // save just the max levels
      } else if (sample < signalMin) {
        signalMin = sample;  // save just the min levels
      }
    }
  }
  peakToPeak = signalMax - signalMin;  // max - min = peak-peak amplitude
  return peakToPeak;
}

void setup(void) {
  Serial.begin(9600);
  while (!Serial) {
    delay(10); // will pause Zero, Leonardo, etc until serial console opens
  }

  // Try to initialize the first MPU6050
  if (!mpu1.begin(0x68)) {
    Serial.println("Failed to find MPU6050 chip at 0x68");
    while (1) {
      delay(10);
    }
  }

  // Try to initialize the second MPU6050
  if (!mpu2.begin(0x69)) {
    Serial.println("Failed to find MPU6050 chip at 0x69");
    while (1) {
      delay(10);
    }
  }

  mpu1.setAccelerometerRange(MPU6050_RANGE_16_G);
  mpu1.setGyroRange(MPU6050_RANGE_250_DEG);
  mpu1.setFilterBandwidth(MPU6050_BAND_184_HZ);

  mpu2.setAccelerometerRange(MPU6050_RANGE_16_G);
  mpu2.setGyroRange(MPU6050_RANGE_250_DEG);
  mpu2.setFilterBandwidth(MPU6050_BAND_184_HZ);

  Serial.println("");
  delay(100);
}

void loop() {
  /* Get new sensor events with the readings */
  sensors_event_t a1, g1, temp1;
  sensors_event_t a2, g2, temp2;
  
  mpu1.getEvent(&a1, &g1, &temp1);
  mpu2.getEvent(&a2, &g2, &temp2);

  /* Read light sensor value */
  int lightValue = analogRead(lightPin);

  /* Get sound peak-to-peak values */
  unsigned int soundPeak1 = getSoundPeakToPeak(soundPin1);
  unsigned int soundPeak2 = getSoundPeakToPeak(soundPin2);

  /* Print out the values from the first sensor */
  Serial.print("AX1:");
  Serial.print(a1.acceleration.x);
  Serial.print(" | ");
  Serial.print("AY1:");
  Serial.print(a1.acceleration.y);
  Serial.print(" | ");
  Serial.print("AZ1:");
  Serial.print(a1.acceleration.z);
  Serial.print(" | ");

  /* Print out the values from the second sensor */
  Serial.print("AX2:");
  Serial.print(a2.acceleration.x);
  Serial.print(" | ");
  Serial.print("AY2:");
  Serial.print(a2.acceleration.y);
  Serial.print(" | ");
  Serial.print("AZ2:");
  Serial.print(a2.acceleration.z);
  Serial.print(" | ");

  /* Print out the light sensor value */
  Serial.print("Light:");
  Serial.print(lightValue);
  Serial.print(" | ");

  /* Print out the sound sensor peak-to-peak values */
  Serial.print("Sound1:");
  Serial.print(soundPeak1);
  Serial.print(" | ");
  Serial.print("Sound2:");
  Serial.print(soundPeak2);
  Serial.println("");

  delay(10);
}