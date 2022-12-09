import socket
from datetime import datetime

HOST = '127.0.0.1'
PORT = 3333
BUF_SIZE = 4096

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.bind((HOST, PORT))
s.listen()

conn, addr = s.accept()
print('Connected by', addr)
count = 0
while 1:
    data = conn.recv(BUF_SIZE)
    count += 1
    if not data:
        break
    else:
        if (count % 500 == 0):
            print(datetime.now().strftime("%H:%M:%S.%f"))
conn.close()