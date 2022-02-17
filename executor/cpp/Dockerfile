FROM gcc:9 

RUN apt-get update && apt-get install -y \
    build-essential \
    autoconf \
    libtool \
    pkg-config \
    cmake

ENV INSTALL_DIR=$HOME/.local
RUN mkdir -p $INSTALL_DIR
ENV PATH="$INSTALL_DIR/bin:$PATH"

RUN git clone --recurse-submodules -b v1.42.0 https://github.com/grpc/grpc
WORKDIR grpc
RUN mkdir -p cmake/build
WORKDIR cmake/build

RUN cmake -DgRPC_INSTALL=ON \
          -DgRPC_BUILD_TESTS=OFF \
          -DCMAKE_INSTALL_PREFIX=$INSTALL_DIR \
          ../..
RUN make -j 2
RUN make install

RUN mkdir -p /openfx/executor
WORKDIR /openfx/executor

COPY server.cc .
COPY mesh.cc .
COPY proto proto/