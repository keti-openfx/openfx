package io.grpc.fxwatcher;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.26.0)",
    comments = "Source: fxwatcher.proto")
public final class FxWatcherGrpc {

  private FxWatcherGrpc() {}

  public static final String SERVICE_NAME = "pb.FxWatcher";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<io.grpc.fxwatcher.Request,
      io.grpc.fxwatcher.Reply> getCallMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Call",
      requestType = io.grpc.fxwatcher.Request.class,
      responseType = io.grpc.fxwatcher.Reply.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<io.grpc.fxwatcher.Request,
      io.grpc.fxwatcher.Reply> getCallMethod() {
    io.grpc.MethodDescriptor<io.grpc.fxwatcher.Request, io.grpc.fxwatcher.Reply> getCallMethod;
    if ((getCallMethod = FxWatcherGrpc.getCallMethod) == null) {
      synchronized (FxWatcherGrpc.class) {
        if ((getCallMethod = FxWatcherGrpc.getCallMethod) == null) {
          FxWatcherGrpc.getCallMethod = getCallMethod =
              io.grpc.MethodDescriptor.<io.grpc.fxwatcher.Request, io.grpc.fxwatcher.Reply>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Call"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  io.grpc.fxwatcher.Request.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  io.grpc.fxwatcher.Reply.getDefaultInstance()))
              .setSchemaDescriptor(new FxWatcherMethodDescriptorSupplier("Call"))
              .build();
        }
      }
    }
    return getCallMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FxWatcherStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FxWatcherStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FxWatcherStub>() {
        @java.lang.Override
        public FxWatcherStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FxWatcherStub(channel, callOptions);
        }
      };
    return FxWatcherStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FxWatcherBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FxWatcherBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FxWatcherBlockingStub>() {
        @java.lang.Override
        public FxWatcherBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FxWatcherBlockingStub(channel, callOptions);
        }
      };
    return FxWatcherBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FxWatcherFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FxWatcherFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FxWatcherFutureStub>() {
        @java.lang.Override
        public FxWatcherFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FxWatcherFutureStub(channel, callOptions);
        }
      };
    return FxWatcherFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class FxWatcherImplBase implements io.grpc.BindableService {

    /**
     */
    public void call(io.grpc.fxwatcher.Request request,
        io.grpc.stub.StreamObserver<io.grpc.fxwatcher.Reply> responseObserver) {
      asyncUnimplementedUnaryCall(getCallMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCallMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                io.grpc.fxwatcher.Request,
                io.grpc.fxwatcher.Reply>(
                  this, METHODID_CALL)))
          .build();
    }
  }

  /**
   */
  public static final class FxWatcherStub extends io.grpc.stub.AbstractAsyncStub<FxWatcherStub> {
    private FxWatcherStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FxWatcherStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FxWatcherStub(channel, callOptions);
    }

    /**
     */
    public void call(io.grpc.fxwatcher.Request request,
        io.grpc.stub.StreamObserver<io.grpc.fxwatcher.Reply> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCallMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class FxWatcherBlockingStub extends io.grpc.stub.AbstractBlockingStub<FxWatcherBlockingStub> {
    private FxWatcherBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FxWatcherBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FxWatcherBlockingStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.fxwatcher.Reply call(io.grpc.fxwatcher.Request request) {
      return blockingUnaryCall(
          getChannel(), getCallMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class FxWatcherFutureStub extends io.grpc.stub.AbstractFutureStub<FxWatcherFutureStub> {
    private FxWatcherFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FxWatcherFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FxWatcherFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<io.grpc.fxwatcher.Reply> call(
        io.grpc.fxwatcher.Request request) {
      return futureUnaryCall(
          getChannel().newCall(getCallMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CALL = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final FxWatcherImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(FxWatcherImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CALL:
          serviceImpl.call((io.grpc.fxwatcher.Request) request,
              (io.grpc.stub.StreamObserver<io.grpc.fxwatcher.Reply>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class FxWatcherBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FxWatcherBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return io.grpc.fxwatcher.FxWatcherProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FxWatcher");
    }
  }

  private static final class FxWatcherFileDescriptorSupplier
      extends FxWatcherBaseDescriptorSupplier {
    FxWatcherFileDescriptorSupplier() {}
  }

  private static final class FxWatcherMethodDescriptorSupplier
      extends FxWatcherBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    FxWatcherMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (FxWatcherGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FxWatcherFileDescriptorSupplier())
              .addMethod(getCallMethod())
              .build();
        }
      }
    }
    return result;
  }
}
