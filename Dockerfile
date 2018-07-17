###                        _       _
#__      _____  __ ___   ___  __ _| |_ ___
#\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
# \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
#  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
#
# Copyright © 2016 - 2018 Weaviate. All rights reserved.
# LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
# AUTHOR: Bob van Luijt (bob@kub.design)
# See www.creativesoftwarefdn.org for details
# Contact: @CreativeSofwFdn / bob@kub.design
###

# Build container
FROM golang as BUILDER
RUN go get -u golang.org/x/vgo
WORKDIR /go/src/github.com/creativefotwarefdn/weaviate
COPY . .
RUN  CGO_ENABLED=1 GOOS=linux vgo install -a -ldflags '-extldflags "-static"' ./...
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"'
# vgo install ./...

# Base image for Weaviate
FROM alpine
ENV PATH=$PATH:/opt/weaviate
RUN mkdir -p /etc/weaviate/
COPY --from=BUILDER /go/bin/weaviate-server /opt/weaviate/
COPY --from=BUILDER /go/bin/contextionary-generator /opt/weaviate/

COPY ./weaviate.conf.json /etc/weaviate/config.json
COPY ./test/schema/test-action-schema.json /etc/weaviate/
COPY ./test/schema/test-thing-schema.json  /etc/weaviate/

# Copy script in container
COPY ./weaviate-entrypoint.sh /weaviate-entrypoint.sh

# Set workdir
WORKDIR /var/weaviate/

ENV WEAVIATE_SCHEME="http" \
    WEAVIATE_PORT="80" \
    WEAVIATE_HOST="0.0.0.0" \
    WEAVIATE_CONFIG="cassandra_docker" 

CMD ["sh", "-c", "/opt/weaviate/weaviate-server --scheme=${WEAVIATE_SCHEME} --port=${WEAVIATE_PORT} --host=${WEAVIATE_HOST} --config-file=/etc/weaviate/config.json --config=${WEAVIATE_CONFIG}"]
