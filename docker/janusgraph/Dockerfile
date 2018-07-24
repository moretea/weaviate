# JanusGraph docker container
FROM alpine as FETCHER
ENV JANUS_VERSION=0.2.1
ARG JANUS_SHA256=0c832568706261f8466e8717bac4a3b1629b52626cf1ecf9490c10cca8a8dc77
ENV JANUS_URL=https://github.com/JanusGraph/janusgraph/releases/download/v${JANUS_VERSION}/janusgraph-${JANUS_VERSION}-hadoop2.zip

WORKDIR /janus
RUN wget -S -O /janus/janus.zip ${JANUS_URL}
RUN (echo -e "${JANUS_SHA256}  /janus/janus.zip" | sha256sum -c || (echo -e "shasum failed!\nExpected ${JANUS_SHA256},\nbut got $(sha256sum /janus/janus.zip)" >&2; exit 1))

# Unzip janus, and put it in /janus.
RUN unzip -q janus.zip && \
    rm /janus/janus.zip && \
    mv /janus/janusgraph-${JANUS_VERSION}-hadoop2/* /janus && \
    rmdir /janus/janusgraph-${JANUS_VERSION}-hadoop2/

FROM openjdk:11
COPY --from=FETCHER /janus /opt/janus
RUN echo -e "storage.backend=embeddedcassandra\n \
    index.search.backend=lucene\n \
    index.search.directory=/data/searchindex \
    " > /opt/janus/conf/gremlin-server/janusgraph-cql-es-server.properties
CMD ["/opt/janus/bin/gremlin-server.sh"]
