# build
docker build -t websocket-client-test .  

# run
docker run --name=websocket-client-test -p 8081:8081 websocket-client-test  
docker run --name=websocket-client-test -p 8081:8081 -e WEBSOCKET_PORT=8083 websocket-client-test  

# use
go get github.com/gorilla/websocket  
go get github.com/spf13/viper  
go get google.golang.org/protobuf  