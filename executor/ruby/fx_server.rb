this_dir = File.expand_path(File.dirname(__FILE__))
lib_dir = File.join(this_dir, '.')
$LOAD_PATH.unshift(lib_dir) unless $LOAD_PATH.include?(lib_dir)

require 'grpc'
require 'fxwatcher_services_pb'

require './handler'
require './fx_mesh'

class FxWatcherServer < Pb::FxWatcher::Service
  def call(req, _unused_call)
    p "[fxwatcher] start service."
    p "[fxmesh] start service."
    tmp = FxWatcher.Handler(req.input)
    Pb::Reply.new(Output: "#{tmp}")
  end
end

def createLockFile
  p "Writing lock-file to /tmp/.lock"
  File.open("/tmp/.lock", "w")
  File.chmod(0660, "/tmp/.lock")
end

createLockFile

# main starts an RpcServer that receives requests to GreeterServer at the sample
# server port.
def main
  s = GRPC::RpcServer.new
  s.add_http2_port('0.0.0.0:50051', :this_port_is_insecure)
  s.add_http2_port('0.0.0.0:50052', :this_port_is_insecure)
  s.handle(FxWatcherServer)
  # Runs the server with SIGHUP, SIGINT and SIGQUIT signal handlers to 
  #   gracefully shutdown.
  # User could also choose to run server via call to run_till_terminated
  s.run_till_terminated_or_interrupted([1, 'int', 'SIGQUIT'])
end

main
