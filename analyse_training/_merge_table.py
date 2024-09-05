import pandas as pd

# 读取数据
table1 = pd.read_csv('./data/08-09_13-13-34_SCADA.csv')
table2 = pd.read_csv('./data/vibration_20240809-13:20:11_normal_EMS.csv')

time_format = '%Y-%m-%d %H:%M:%S.%f'
table1['time'] = pd.to_datetime(table1['time'], format=time_format)
table2['time'] = pd.to_datetime(table2['time'], format=time_format)
table1 = table1['2024-08-09 13:21:00.000' < table1['time']]
table2 = table2['2024-08-09 13:21:00.000' < table2['time']]

table1.set_index('time', inplace=True)
table2.set_index('time', inplace=True)

table1_resampled = table1.resample('200ms').mean().interpolate()
table2_resampled = table2.resample('200ms').mean().interpolate()

merged_table = pd.merge_asof(table1_resampled, table2_resampled, left_index=True,
                             right_index=True, direction='nearest', tolerance=pd.Timedelta('100ms'))

merged_table.reset_index(inplace=True)
merged_table['time'] = merged_table['time'].dt.strftime(time_format)


merged_table = merged_table.dropna()

merged_table.to_csv('./data/merged_table_normal.csv', index=False)
