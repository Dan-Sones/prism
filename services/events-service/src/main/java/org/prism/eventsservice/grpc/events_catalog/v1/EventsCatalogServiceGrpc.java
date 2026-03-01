package org.prism.eventsservice.grpc.events_catalog.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@io.grpc.stub.annotations.GrpcGenerated
public final class EventsCatalogServiceGrpc {

  private EventsCatalogServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "events_catalog.v1.EventsCatalogService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest,
      org.prism.eventsservice.grpc.events_catalog.v1.EventType> getGetEventTypeByKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetEventTypeByKey",
      requestType = org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest.class,
      responseType = org.prism.eventsservice.grpc.events_catalog.v1.EventType.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest,
      org.prism.eventsservice.grpc.events_catalog.v1.EventType> getGetEventTypeByKeyMethod() {
    io.grpc.MethodDescriptor<org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest, org.prism.eventsservice.grpc.events_catalog.v1.EventType> getGetEventTypeByKeyMethod;
    if ((getGetEventTypeByKeyMethod = EventsCatalogServiceGrpc.getGetEventTypeByKeyMethod) == null) {
      synchronized (EventsCatalogServiceGrpc.class) {
        if ((getGetEventTypeByKeyMethod = EventsCatalogServiceGrpc.getGetEventTypeByKeyMethod) == null) {
          EventsCatalogServiceGrpc.getGetEventTypeByKeyMethod = getGetEventTypeByKeyMethod =
              io.grpc.MethodDescriptor.<org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest, org.prism.eventsservice.grpc.events_catalog.v1.EventType>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetEventTypeByKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.prism.eventsservice.grpc.events_catalog.v1.EventType.getDefaultInstance()))
              .setSchemaDescriptor(new EventsCatalogServiceMethodDescriptorSupplier("GetEventTypeByKey"))
              .build();
        }
      }
    }
    return getGetEventTypeByKeyMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static EventsCatalogServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceStub>() {
        @java.lang.Override
        public EventsCatalogServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventsCatalogServiceStub(channel, callOptions);
        }
      };
    return EventsCatalogServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static EventsCatalogServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceBlockingV2Stub>() {
        @java.lang.Override
        public EventsCatalogServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventsCatalogServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return EventsCatalogServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static EventsCatalogServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceBlockingStub>() {
        @java.lang.Override
        public EventsCatalogServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventsCatalogServiceBlockingStub(channel, callOptions);
        }
      };
    return EventsCatalogServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static EventsCatalogServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<EventsCatalogServiceFutureStub>() {
        @java.lang.Override
        public EventsCatalogServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new EventsCatalogServiceFutureStub(channel, callOptions);
        }
      };
    return EventsCatalogServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     */
    default void getEventTypeByKey(org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest request,
        io.grpc.stub.StreamObserver<org.prism.eventsservice.grpc.events_catalog.v1.EventType> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetEventTypeByKeyMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service EventsCatalogService.
   */
  public static abstract class EventsCatalogServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return EventsCatalogServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service EventsCatalogService.
   */
  public static final class EventsCatalogServiceStub
      extends io.grpc.stub.AbstractAsyncStub<EventsCatalogServiceStub> {
    private EventsCatalogServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventsCatalogServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventsCatalogServiceStub(channel, callOptions);
    }

    /**
     */
    public void getEventTypeByKey(org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest request,
        io.grpc.stub.StreamObserver<org.prism.eventsservice.grpc.events_catalog.v1.EventType> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetEventTypeByKeyMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service EventsCatalogService.
   */
  public static final class EventsCatalogServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<EventsCatalogServiceBlockingV2Stub> {
    private EventsCatalogServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventsCatalogServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventsCatalogServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     */
    public org.prism.eventsservice.grpc.events_catalog.v1.EventType getEventTypeByKey(org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetEventTypeByKeyMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service EventsCatalogService.
   */
  public static final class EventsCatalogServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<EventsCatalogServiceBlockingStub> {
    private EventsCatalogServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventsCatalogServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventsCatalogServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public org.prism.eventsservice.grpc.events_catalog.v1.EventType getEventTypeByKey(org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetEventTypeByKeyMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service EventsCatalogService.
   */
  public static final class EventsCatalogServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<EventsCatalogServiceFutureStub> {
    private EventsCatalogServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected EventsCatalogServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new EventsCatalogServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.prism.eventsservice.grpc.events_catalog.v1.EventType> getEventTypeByKey(
        org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetEventTypeByKeyMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_EVENT_TYPE_BY_KEY = 0;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_EVENT_TYPE_BY_KEY:
          serviceImpl.getEventTypeByKey((org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest) request,
              (io.grpc.stub.StreamObserver<org.prism.eventsservice.grpc.events_catalog.v1.EventType>) responseObserver);
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

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getGetEventTypeByKeyMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest,
              org.prism.eventsservice.grpc.events_catalog.v1.EventType>(
                service, METHODID_GET_EVENT_TYPE_BY_KEY)))
        .build();
  }

  private static abstract class EventsCatalogServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    EventsCatalogServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.prism.eventsservice.grpc.events_catalog.v1.EventsCatalog.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("EventsCatalogService");
    }
  }

  private static final class EventsCatalogServiceFileDescriptorSupplier
      extends EventsCatalogServiceBaseDescriptorSupplier {
    EventsCatalogServiceFileDescriptorSupplier() {}
  }

  private static final class EventsCatalogServiceMethodDescriptorSupplier
      extends EventsCatalogServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    EventsCatalogServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (EventsCatalogServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new EventsCatalogServiceFileDescriptorSupplier())
              .addMethod(getGetEventTypeByKeyMethod())
              .build();
        }
      }
    }
    return result;
  }
}
