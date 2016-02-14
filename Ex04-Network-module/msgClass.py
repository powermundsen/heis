#haakoneh & erlendvd

import socket 
import select

class MessageClass:
	def __init__(self, id, port, ip, buffer_size):
		self.id = id
		self.port = port
		self.ip = ip
		self.msg = ""
		self.buffer_size = buffer_size
		self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) 
		#self.socket.bind((self.ip, self.port))
	
	def setMsg(self, msg):
		self.msg = msg

	def sendMsg(self):
		self.socket.sendto(self.msg, (self.ip, self.port))

	def retrieveMsg(self):
		self.msg, self.addr = self.socket.recvfrom(self.buffer_size)
		print self.msg
		print self.addr

	def printMsg(self):
		print 'ID: ' + str(self.id)
		print 'Port: ' + str(self.port)
		print 'IP: ' + str(self.ip)
		print 'Msg: ' + str(self.msg)
		print ''


