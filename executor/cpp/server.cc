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
  std::string mesh_address("0.0.0.0:50052");
  FxWatcherServiceImpl service;

  ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.AddListeningPort(mesh_address, grpc::InsecureServerCredentials());

  builder.RegisterService(&service);

  std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "[fxwatcher] start service." << server_address << std::endl;
  std::cout << "[fxmesh] start service." << mesh_address << std::endl;

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
