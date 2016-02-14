#haakoneh & erlendvd

#import socket 
#import select
import msgClass
import time

UDP_PORT_IP = 30000
UDP_PORT = 20009
UDP_IP = "129.241.187.156"
UDP_IP_2 = "129.241.187.159"
buffer_size = 1024


def main():
	msg = msgClass.MessageClass(1, UDP_PORT, UDP_IP_2, buffer_size)
	
	#msg.setMsg('lol')
	#msg.printMsg()
	#msg.sendMsg()
	#msg.retrieveMsg()


	while True:
		time.sleep(1)
		msg.retrieveMsg()


main()
