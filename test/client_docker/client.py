import socket
from datetime import datetime

HOST = '127.0.0.1'
PORT = 3333
PACK_SIZE = 300 * 1024 * 1024

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

s.connect((HOST, PORT))
print('Client connected.')

content = 'a' * PACK_SIZE
packet_num = 15
for packet_id in list(range(packet_num)):
    s.sendall(content.encode())
    print('Packet [' + str(packet_id) + '] has been sent.')
    print(datetime.now().strftime("%H:%M:%S.%f"))
    print()
