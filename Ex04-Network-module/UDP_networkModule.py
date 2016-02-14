#haakoneh & erlendvd

#import socket 
#import select
import msgClass

UDP_PORT_IP = 30000
UDP_PORT = 20009
UDP_IP = "129.241.187.156"
UDP_IP_2 = "129.241.187.144"
buffer_size = 1024



def main():
	msg = msgClass.MessageClass(1, UDP_PORT, UDP_IP, buffer_size)
	msg.setMsg('lol')
	msg.printMsg()
		
	msg.sendMsg()
	#msg.retrieveMsg()

main()




#def UDP_init(){
	#find it's own IP address/port number
	#global identifier = that stuff
#}