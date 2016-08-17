#!/usr/bin/python

from __future__ import print_function

import pyjsonrpc
import sys
import Queue
import threading
import signal

class HelloWorldPython(pyjsonrpc.JsonRpc):
	'''
	JsonRpc server example. It has one method Analyze().
	'''

	@pyjsonrpc.rpcmethod
	def Analyze(self, fmeta):
		'''
		Perform analysis on file based on fmeta.
		'''

		fm = json.loads(fmeta)
		return fm.Path

def analyze(fmeta):
	fm = json.loads(fmeta)
	return '{"Hello":"%s"}' % fm['Filepath']

def worker(line, q, rpc_client):
	'''
	Worker thread that handles the RPC server calls when they come in
	from stdin.
	'''
	out = rpc_client.call(line)
	q.put(out)

def printer(q):

	while True:
		out = q.get()
		if out == 'kill':
			return
		sys.stdout.write(out + '\n')
		sys.stdout.flush()

def init(printer_thread, q):
	'''
	Initialize the printer thread and exit signal handler so that we
	kill long running threads on exit.
	'''

	printer_thread.start()

	def signal_handler(signal, frame):
		q.put('kill')
		printer_thread.join()
		sys.exit(0)

	signal.signal(signal.SIGINT, signal_handler)

def main():
	print("Starting HelloWorldPython", file=sys.stderr)
	q = Queue.Queue()
	printer_thread = threading.Thread(target=printer, args=[q])

	init(printer_thread, q)


	rpc = pyjsonrpc.JsonRpc(methods = {"HelloWorldPython.Analyze" : analyze})
	#rpc = HelloWorldPython()
	print(rpc)
	line = sys.stdin.readline()

	while line:
		print(line)
		try:
			current_input = line 
			t = threading.Thread(target=worker, args=[line, q, rpc])
			t.start()
			line = sys.stdin.readline()
		except Exception as e:
			warning('Exception occrred: ', e)
			q.put('kill')
			printer_thread.join()

if __name__ == '__main__':
	print ('Testing 1234', file=sys.stderr)
	main()