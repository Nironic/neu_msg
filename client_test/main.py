import socket
from threading import Thread


def send(text):
    sock.send((text + "\n").encode())

class Client:
    def key(self, key):
        send(key)
        data = sock.recv(1024).decode()
        print(data)
    
    def login(self, login, password, ):
        print("Login")
        send("login")
        data = sock.recv(1024).decode()
        print(data)
        send(login)
        print(sock.recv(1024).decode())
        send(password)
        print(sock.recv(1024).decode())
        send("check")
        data = sock.recv(1024).decode()
        print(data)
    
    def register(self, login, password, username):
        print("Reg")
        send("reg")
        print(sock.recv(1024).decode())
        send(login)
        print(sock.recv(1024).decode())
        send(password)
        print(sock.recv(1024).decode())
        send(username)
        print(sock.recv(1024).decode())


sock = socket.socket()

sock.connect(('localhost', 8080))

key = "9S2oPsZJ1ipUxKlbyJvr"


cl = Client()
cl.key(key)
cl.login("Roman", "12345")