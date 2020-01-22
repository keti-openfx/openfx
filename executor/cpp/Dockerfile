FROM gcc:9

RUN apt-get update && apt-get install -y \
    build-essential \
    autoconf \
    libtool \
    pkg-config \
    libgflags-dev \
    libgtest-dev \
    clang \
    automake \
    curl \
    make \
    g++ \ 
    unzip \
    libc++-dev

RUN mkdir -p ${GOPATH}/src/github.com
WORKDIR ${GOPATH}/src/github.com
RUN git clone https://github.com/protocolbuffers/protobuf.git
WORKDIR ${GOPATH}/src/github.com/protobuf
RUN git submodule update --init --recursive
RUN ./autogen.sh
RUN ./configure
RUN make
RUN make check
RUN make install
RUN ldconfig

WORKDIR ${GOPATH}/src/github.com
RUN git clone -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc
WORKDIR ${GOPATH}/src/github.com/grpc
RUN git submodule update --init
RUN sed -i 's/-Werror/ /g' Makefile
RUN make
RUN make install

RUN mkdir -p /openfx/executor
WORKDIR /openfx/executor
COPY server.cc .

COPY Makefile .
COPY proto proto/
RUN protoc -I proto --grpc_out=. --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` proto/fxwatcher.proto
RUN protoc -I proto --cpp_out=. proto/fxwatcher.proto
