from time import sleep

workload = 300 * 1024 * 1024
unit = True
unit_size = unit.__sizeof__()
data = [unit] * (workload // unit_size)

while True:
    i = 0
    while i < len(data):
        data[i] = not data[i]
        i+=1
        sleep(0.1)
