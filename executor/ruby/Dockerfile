FROM ruby:2.6-alpine

RUN apk update

RUN apk add --no-cache build-base && \
    apk del curl --no-cache

RUN apk --no-cache add graphicsmagick libc6-compat

RUN mkdir -p /openfx/executor
WORKDIR /openfx/executor

RUN gem update --system
RUN bundle init
RUN echo 'gem "grpc", platforms: ["ruby"]' >> Gemfile
RUN echo 'gem "google-protobuf", "3.8.0", platforms: ["ruby"]' >> Gemfile
RUN CFLAGS="-Wno-cast-function-type"
RUN gem install bundler && BUNDLE_FORCE_RUBY_PLATFORM=1 bundle install 

COPY . .
