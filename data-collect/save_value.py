import serial
import csv
import datetime

# 串口配置
serial_port = '/dev/cu.usbmodem21301'  # 根据实际情况修改
baud_rate = 9600
timeout = 1

# 打开串口
ser = serial.Serial(serial_port, baud_rate, timeout=timeout)

# 打开CSV文件以便写入
csv_file = str(datetime.datetime.now().strftime('%Y%m%d-%H:%M:%S')) + '.csv'
fieldnames = ['timestamp', 'AX1', 'AY1', 'AZ1', 'AX2', 'AY2', 'AZ2', 'Light', 'Sound1', 'Sound2']

with open(csv_file, mode='w', newline='') as file:
    writer = csv.DictWriter(file, fieldnames=fieldnames)
    writer.writeheader()

    try:
        while True:
            # 读取串口数据
            line = ser.readline().decode('utf-8').strip()
            if line:
                print(f"Received line: {line}")
                try:
                    # 解析数据
                    parts = line.split(' | ')
                    
                    # 提取加速度值
                    ax1 = float(parts[0].split(':')[1])
                    ay1 = float(parts[1].split(':')[1])
                    az1 = float(parts[2].split(':')[1])
                    ax2 = float(parts[3].split(':')[1])
                    ay2 = float(parts[4].split(':')[1])
                    az2 = float(parts[5].split(':')[1])

                    # 获取光传感器和声音传感器的值
                    light = int(parts[6].split(':')[1])
                    sound1 = int(parts[7].split(':')[1])
                    sound2 = int(parts[8].split(':')[1])

                    # 获取当前时间戳
                    timestamp = datetime.datetime.now().strftime('%H:%M:%S.%f')

                    # 写入CSV文件
                    writer.writerow({
                        'timestamp': timestamp,
                        'AX1': ax1,
                        'AY1': ay1,
                        'AZ1': az1,
                        'AX2': ax2,
                        'AY2': ay2,
                        'AZ2': az2,
                        'Light': light,
                        'Sound1': sound1,
                        'Sound2': sound2
                    })
                except Exception as e:
                    print(f"Error parsing line: {line} - {e}")
    except KeyboardInterrupt:
        print("Terminating the script.")
    finally:
        ser.close()
