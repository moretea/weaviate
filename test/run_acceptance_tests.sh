set -e
set -m

# Running Cassandra deamon
if (docker ps --format '{{ .Names }}' | grep  -s weaviate_db_1_travis); then
  echo "Deleting weaviate_db_1_travis"
  docker rm -f weaviate_db_1_travis > /dev/null
fi

echo "Starting Cassandra..."
docker run --rm -it --name weaviate_db_1_travis -e CASSANDRA_BROADCAST_ADDRESS=127.0.0.1 -p 7000:7000 -p 9042:9042 -d cassandra:3 > /dev/null
docker exec weaviate_db_1_travis bash -c "echo 'waiting until Cassandra is up...'; while ! (cqlsh -e 'describe cluster' 1>&2 2> /dev/null) ; do sleep 1; done; echo 'it is up!'"

# Running Weaviate as deamon
nohup go run ./cmd/weaviate-server/main.go --scheme=http --port=8080 --host=127.0.0.1 --config="cassandra_docker" --config-file="./weaviate.conf.json" &
# Sleep to make sure all is up and running
sleep 45
# cat nohup for debugging
cat nohup.out
# Find the correct ip, key and token
CASSANDRAIP=$(sudo docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' weaviate_db_1_travis)
ROOTTOKEN=$(cat nohup.out | grep -s ROOTTOKEN|sed 's/.*ROOTTOKEN=//')
ROOTKEY=$(cat nohup.out | grep -s ROOTKEY|sed 's/.*ROOTKEY=//')

# Run the integration tests
# On error, quit
# Test!
go test -v ./test -args -api-key=${ROOTKEY} -api-token=${ROOTTOKEN} -server-host=127.0.0.1 -server-port=8080 -server-scheme=http
exit_status=$?

# Kill all background tasks
kill -9 $(jobs -p)

# Kill cassandra
docker rm -f weaviate_db_1_travis

exit $exit_status
