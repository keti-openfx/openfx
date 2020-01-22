#include <iostream>
#include <memory>
#include <string>
#include <fstream>

#include <grpcpp/grpcpp.h>
#include "fxwatcher.grpc.pb.h"
#include "handler.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using pb::FxWatcher;
using pb::Request;
using pb::Reply;

class FxWatcherServiceImpl final : public FxWatcher::Service {
  Status Call(ServerContext* context, const Request* request, 
              Reply* reply) override {
    reply->set_output(Handler(request->input()));
    return Status::OK;
  }
};

void RunServer() {
  std::string server_address("0.0.0.0:50051");
  FxWatcherServiceImpl service;

  ServerBuilder builder;
  // Listen on the given address without any authentication mechanism.
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  
  // Register "service" as the instance through which we'll communicate with
  // clients. In this case it corresponds to an *synchronous* service.
  builder.RegisterService(&service);
  
  // Finally assemble the server.
  std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "[fxwatcher] start service." << server_address << std::endl;

  // Wait for the server to shutdown. Note that some other thread must be
  // responsible for shutting down the server for this call to ever return.
  server->Wait();
}

int main(int argc, char** argv) {
  std::string filePath = "/tmp/.lock";
  std::ofstream writeFile (filePath);
  writeFile << " " << std::endl;
  writeFile.close(); 
 
  RunServer();

  return 0;
}
