#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>
#include "fxwatcher.grpc.pb.h"
#include "handler.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

using pb::FxWatcher;
using pb::Request;
using pb::Reply;

class FxMeshClient {
    public:
        FxMeshClient(std::shared_ptr<Channel> channel) 
            : stub_(FxWatcher::NewStub(channel)) {}

    string Call(string input) {
        Request request;
        request.set_input(input);

        Reply reply;
        ClientContext context;
        Status status = stub_->Call(&context, request, &reply);

        if(status.ok()){
            return reply.output();
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
            return 0;
        }
    }

    private:
        std::unique_ptr<FxWatcher::Stub> stub_;
};

string MeshCall(string functionName, string input) {
    std::string address(functionName + ".openfx-fn:50052");
    FxMeshClient client(
        grpc::CreateChannel(
            address, 
            grpc::InsecureChannelCredentials()
        )
    );

    string response = client.Call(input);
    return response;
}
