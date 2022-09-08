# Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
# Use of this source code is governed by a license that can be
# found in the LICENSE file.

FROM exploredev/reviewpad:semanticservice-v1.14 as semanticservice

FROM golang:1.18.5 AS build

ENV LIBGIT2_ZIP v1.2.0.zip
ENV LIBGIT2 libgit2-1.2.0

ARG mixpanelToken

WORKDIR /service

# Install necessary packages
RUN apt-get update && apt-get -y install unzip cmake libssl-dev && apt-get clean

# Install libgit2
RUN curl -OL https://github.com/libgit2/libgit2/archive/refs/tags/${LIBGIT2_ZIP} && \
    unzip -o $LIBGIT2_ZIP -d /tmp && \
    cd /tmp/${LIBGIT2} && \
    mkdir build && \
    cd build && \
    cmake .. && \
    cmake --build . --target install

# Download the dependencies as a separate, cacheable step
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

# Build the project
COPY . .
RUN go build -ldflags "-X main.MixpanelToken=$mixpanelToken"

# Final image
FROM gcr.io/distroless/cc:debug

ENV LIBGIT2 libgit2-1.2.0

SHELL ["/busybox/sh", "-c"]

ENV ATLAS_SEMANTIC_PORT="0.0.0.0:3006"
ENV INPUT_SEMANTIC_SERVICE="0.0.0.0:3006"

WORKDIR /app

# Semantic service
COPY --from=build /tmp/${LIBGIT2}/build/libgit2.so* /lib
COPY --from=semanticservice /semantic-server /app/semantic-server

COPY --from=build /service/action /app/action

COPY ./run.sh .
RUN chmod +x /app/run.sh

ENTRYPOINT ["sh","/app/run.sh"]
