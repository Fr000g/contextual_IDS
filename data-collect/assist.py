import pandas as pd
import numpy as np
import matplotlib.pyplot as plt

# 读取CSV文件
def read_data(file_path):
    data = pd.read_csv(file_path)
    return data

# 计算振动幅度
def calculate_magnitude(data):
    ax1, ay1, az1 = -0.597253,-0.265987,10.883562
    ax2, ay2, az2 = -0.857146,0.047639,9.733455

    data['magnitude1'] = np.sqrt((data['AX1'] - ax1) ** 2 + (data['AY1'] - ay1) ** 2 + (data['AZ1'] - az1) ** 2)
    data['magnitude2'] = np.sqrt((data['AX2'] - ax2) ** 2 + (data['AY2'] - ay2) ** 2 + (data['AZ2'] - az2) ** 2)
    return data

# 使用滑动窗口平滑数据
def smooth_data(data, window_size):
    data['magnitude1'] = data['magnitude1'].rolling(window=window_size, center=True).mean()
    data['magnitude2'] = data['magnitude2'].rolling(window=window_size, center=True).mean()
    return data

# 主函数
def main(file_path, window_size):
    data = read_data(file_path)
    data = calculate_magnitude(data)
    data = smooth_data(data, window_size)
    
    # 取前1/3数据
    n = len(data) // 3
    
    
    return data

def display(data):
    columns = data.columns
    plt.figure(figsize=(20, 8))
    columns = data.columns

    for i, column in enumerate(columns[1:], 1):
        n = len(columns)
        plt.subplot(int(n/2)+1, 2, i)
        plt.plot(data['time'], data[column])
        # plt.title(column)
        plt.xlabel('Time')
        # plt.ylabel(column)
        plt.grid(True)

    plt.tight_layout()
    plt.show()

# 使用示例
file_path = '20240809-10:12:23_attack.csv'
window_size = 5  # 滑动窗口大小，可根据需要调整
smoothed_data = main(file_path, window_size)
smoothed_data['magnitude1'] = smoothed_data['magnitude1'].round(6)
smoothed_data['magnitude2'] = smoothed_data['magnitude2'].round(6)
# smoothed_data.to_csv('vibration.csv', index=False)
# display(smoothed_data)
smoothed_data.to_csv('vibration.csv',index=False)
