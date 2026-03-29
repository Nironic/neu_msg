import socket

sock = socket.socket()

sock.connect(('localhost', 8080))

key = b"11112\n"

class Client:
    def connect(self, key):
        sock.send(key)
        sock.recv(1024)
    
    def run(self):
        while True:
            sock.send(b"Hello\n")
            print(sock.recv(1024).decode().strip())


cl = Client()
cl.connect(key)
cl.run()