package io.grpc.fxwatcher;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;
import java.io.FileWriter;

public class fxwatcherServer {
  private static final Logger logger = Logger.getLogger(fxwatcherServer.class.getName());

  private Server server;

  private void start() throws IOException {
    int port = 50051;
    server = ServerBuilder.forPort(port)
        .addService(new FxWatcherImpl())
        .build()
        .start();
    logger.info("[fxwatcher] start service, listening on " + port);
    Runtime.getRuntime().addShutdownHook(new Thread() {
      public void run() {
        System.err.println("*** shutting down gRPC server since JVM is shutting down");
        try {
          fxwatcherServer.this.stop();
        } catch (InterruptedException e) {
          e.printStackTrace(System.err);
        }
        System.err.println("*** server shut down");
      }
    });
  }

  private void stop() throws InterruptedException {
    if (server != null) {
      server.shutdown().awaitTermination(30, TimeUnit.SECONDS);
    }
  }

  /**
   * Await termination on the main thread since the grpc library uses daemon threads.
   */  

  private void blockUntilShutdown() throws InterruptedException {
    if (server != null) {
      server.awaitTermination();
    }
  }

  /**
   * Main launches the server from the command line.
   */

  public static void main(String[] args) throws IOException, InterruptedException {
    final fxwatcherServer server = new fxwatcherServer();
    FileWriter fw = new FileWriter("/tmp/.lock", true);
    fw.close(); 
    server.start();
    server.blockUntilShutdown();
  }

  static class FxWatcherImpl extends FxWatcherGrpc.FxWatcherImplBase {
    public void call(Request req, StreamObserver<Reply> responseObserver) {
      Handler handler = new Handler();
      Reply reply = Reply.newBuilder().setOutput(handler.reply(req.getInput())).build();
      responseObserver.onNext(reply);
      responseObserver.onCompleted();
    }
  }
}
