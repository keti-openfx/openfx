#!/usr/bin/env python
from concurrent import futures
import sys
import imp
import os
import time
import signal
import grpc
import logging

import fxwatcher_pb2
import fxwatcher_pb2_grpc

_ONE_DAY_IN_SECONDS = 60 * 60 * 24

def loadPlugin(moduleName, funcName):
    mod = imp.load_source(funcName, moduleName)
    # load user function from module
    return getattr(mod, funcName)

'''
req = {
    'Input': 'test',
    'Info': {
        'FunctionName':''
    }
}
'''

class FxWatcherServicer(fxwatcher_pb2_grpc.FxWatcherServicer):
    def __init__(self, f):
        self.userfunc = f

    def Call(self, request, context):
        if self.userfunc == None:
            logging.info('Function not found')

        out = self.userfunc(request)
        return fxwatcher_pb2.Reply(Output=out)

def createLockFile():
    path = "/tmp/.lock"
    open(path, "w")
    os.chmod(path, 0o660)
    logging.info('Writing lock-file to: '+ path)

def serve():
    fxwatcherPort = os.getenv('WATCHER_PORT', '50051')
    fxmeshPort = os.getenv('MESH_PORT', '50052')
    moduleName = os.getenv('HANDLER_FILE', '/openfx/handler/handler.py')
    funcName = os.getenv('HANDLER_NAME', 'Handler')

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    watcherServer = FxWatcherServicer(loadPlugin(moduleName, funcName))
    fxwatcher_pb2_grpc.add_FxWatcherServicer_to_server(
        watcherServer, server)

    server.add_insecure_port('[::]:'+fxwatcherPort)
    server.add_insecure_port('[::]:'+fxmeshPort)

    createLockFile()

    logging.info('[fxwatcher] start service.')
    logging.info('[fxmeshserver] start service.')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)
        logging.info("[fxwatcher] received SIGTERM.")

if __name__ == '__main__':
    serve()
