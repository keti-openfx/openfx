var PROTO_PATH = __dirname + '/proto/fxwatcher.proto';

var grpc = require('grpc');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
var fx_proto = grpc.loadPackageDefinition(packageDefinition).pb;

var handler = require('./handler.js');

var fs = require('fs');

/**
 * Implements the Call RPC method.
 */
function Call(call, callback) {
		console.log('[fxwatcher] start service.');
		callback(null, {Output: handler(call.request.input)});
}

/**
 * Starts an RPC server that receives requests for the Greeter service at the
 * sample server port
 */
function main() {
  var server = new grpc.Server();
  server.addService(fx_proto.FxWatcher.service, {Call: Call});
  server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());
  server.start();
}

fs.writeFile('/tmp/.lock', '', function (err) {
  if (err) throw err;
  console.log('Writing lock-file to /tmp/.lock');
});
fs.chmodSync('/tmp/.lock', 0660);
main();
