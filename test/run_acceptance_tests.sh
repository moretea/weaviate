# Running Cassandra deamon
docker run -it --name weaviate_db_1_travis -e CASSANDRA_BROADCAST_ADDRESS=127.0.0.1 -p 7000:7000 -p 9042:9042 -d cassandra:3
sleep 45
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
set -e
# Test!
go test -v ./test -args -api-key=${ROOTKEY} -api-token=${ROOTTOKEN} -server-host=127.0.0.1 -server-port=8080 -server-scheme=http
