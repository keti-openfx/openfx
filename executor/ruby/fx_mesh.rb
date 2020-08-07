require 'grpc'
require 'fxwatcher_services_pb'

module Fxwatcher
    def FxWatcher.mesh_call(functionName, input)
        address = "#{functionName}.openfx-fn:50052"
        client_stub =  Pb::FxWatcher::Stub.new(address, :this_channel_is_insecure)
        request = Pb::Request.new(input: "#{input}")
        response = client_stub.call(request)
        return response.Output
    end
end
