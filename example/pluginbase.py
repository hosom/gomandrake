from __future__ import print_function

import threading
import Queue
import sys
import signal
import pyjsonrpc

from datetime import datetime

class Plugin:

	def __init__(self, name='Python_BASE'):

		self.__NAME__ = name
		self.log('Initializing plugin.')
		self.queue = Queue.Queue()
		self.printer_thread = threading.Thread(target=self.printer)
		self.printer_thread.start()

		def signal_handler(signal, frame):
			self.queue.put('kill')
			self.printer_thread.join()
			sys.exit(0)

		signal.signal(signal.SIGINT, signal_handler)

	def log(self, *objs):
		'''A wrapper function that makes it easier to print logs in a more go
		friendly style.
		'''
		ts = datetime.now().strftime("%Y/%m/%d %H:%M:%S")
		print('[%s]' % self.__NAME__, ts, *objs, file=sys.stderr)

	def printer(self):
		'''Output handler. This method will poll the results queue and output
		results as they appear.
		'''
		while True:
			out = self.queue.get()
			if out == 'kill':
				self.log('Kill signal received, stopping threads.')
				return
			sys.stdout.write(out + '\n')
			sys.stdout.flush()

		return

	def worker(self, line):
		'''Worker thread that handles RPC server calls.'''
		out = self.rpc.call(line)
		self.log(out)
		self.queue.put(out)
		return

	def listen(self, method):
		'''Listen for JSON RPC method calls.'''
		self.rpc = pyjsonrpc.JsonRpc(methods = {'%s.Analyze' % self.__NAME__ : method})
		line = sys.stdin.readline()

		while line:
			try:
				this_input = line
				t = threading.Thread(target=self.worker, args=[line])
				t.start()
				line = sys.stdin.readline()
			except Exception, e:
				self.log('Exception occurred: ', e)
				self.queue.put('kill')
				self.printer_thread.join()