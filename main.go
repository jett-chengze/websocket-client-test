package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
	"websocket-client-test/configs"

	pb "websocket-client-test/proto"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var MyConnectId int64
var OtherConnectIds []int64

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	websocketHost := configs.EnvConfig.Websocket.Host
	websocketPort := strconv.Itoa(configs.EnvConfig.Websocket.Port)
	websocketAddr := "ws://" + websocketHost + ":" + websocketPort
	conn, _, err := websocket.DefaultDialer.Dial(websocketAddr+"/ws", nil)
	if err != nil {
		log.Fatalf("websocket connect error: %s", err)
	}
	defer conn.Close()

	go wsSendRedisNewStringHandler(conn)
	go wsSendRedisGetStringHandler(conn)

	go wsReceivedMessageHandler(conn)

	go wsSendGetOthersConnectHandler(conn)
	go wsSendKillOtherConnectHandler(conn)

	for {
		select {
		case <-interrupt:
			return
		}
	}
}

func wsSendRedisNewStringHandler(conn *websocket.Conn) {
	time.Sleep(2 * time.Second)
	payloadType := pb.PayloadType_REDIS_NEW_STRING

	redisNewStringRequest := &pb.RedisNewStringRequest{
		Key:   "string1",
		Value: "string_val1",
	}
	sendPayload, err := proto.Marshal(redisNewStringRequest)
	if err != nil {
		log.Printf("marshal send %v payload error: %s", payloadType, err)
	}
	sendEnvelope := &pb.Envelope{
		PayloadType: payloadType,
		Payload:     sendPayload,
	}
	sendData, err := proto.Marshal(sendEnvelope)
	if err != nil {
		log.Printf("marshal send %v envelope error: %s", payloadType, err)
	}

	err = conn.WriteMessage(websocket.BinaryMessage, sendData)
	if err != nil {
		log.Fatalf("websocket send %v envelope error: %s", payloadType, err)
	}
}

func wsSendRedisGetStringHandler(conn *websocket.Conn) {
	time.Sleep(5 * time.Second)
	payloadType := pb.PayloadType_REDIS_GET_STRING

	redisGetStringRequest := &pb.RedisGetStringRequest{
		Key: "string1",
	}
	sendPayload, err := proto.Marshal(redisGetStringRequest)
	if err != nil {
		log.Printf("marshal send %v payload error: %s", payloadType, err)
	}
	sendEnvelope := &pb.Envelope{
		PayloadType: payloadType,
		Payload:     sendPayload,
	}
	sendData, err := proto.Marshal(sendEnvelope)
	if err != nil {
		log.Printf("marshal send %v envelope error: %s", payloadType, err)
	}

	err = conn.WriteMessage(websocket.BinaryMessage, sendData)
	if err != nil {
		log.Printf("websocket send %v envelope error: %s", payloadType, err)
	}
}

func wsSendGetOthersConnectHandler(conn *websocket.Conn) {
	time.Sleep(6 * time.Second)
	payloadType := pb.PayloadType_GET_OTHER_CONNECT_IDS

	getOtherConnectIdsRequest := &pb.GetOtherConnectIdsRequest{
		MyConnectId: MyConnectId,
	}
	sendPayload, err := proto.Marshal(getOtherConnectIdsRequest)
	if err != nil {
		log.Printf("marshal send %v payload error: %s", payloadType, err)
	}
	sendEnvelope := &pb.Envelope{
		PayloadType: payloadType,
		Payload:     sendPayload,
	}
	sendData, err := proto.Marshal(sendEnvelope)
	if err != nil {
		log.Printf("marshal send %v envelope error: %s", payloadType, err)
	}

	err = conn.WriteMessage(websocket.BinaryMessage, sendData)
	if err != nil {
		log.Fatalf("websocket send %v envelope error: %s", payloadType, err)
	}

}

func wsSendKillOtherConnectHandler(conn *websocket.Conn) {
	time.Sleep(8 * time.Second)

	if len(OtherConnectIds) == 0 {
		return
	}
	payloadType := pb.PayloadType_KILL_OTHER_CONNECT

	killOtherConnectRequest := &pb.KillOtherConnectRequest{
		ConnectId: OtherConnectIds[0],
	}
	sendPayload, err := proto.Marshal(killOtherConnectRequest)
	if err != nil {
		log.Fatalf("marshal send %v payload error: %s", payloadType, err)
	}
	sendEnvelope := &pb.Envelope{
		PayloadType: payloadType,
		Payload:     sendPayload,
	}
	sendData, err := proto.Marshal(sendEnvelope)
	if err != nil {
		log.Fatalf("marshal send %v error: %s", payloadType, err)
	}

	err = conn.WriteMessage(websocket.BinaryMessage, sendData)
	if err != nil {
		log.Fatalf("websocket send %v envelope error: %s", payloadType, err)
	}
}

func wsReceivedMessageHandler(conn *websocket.Conn) {
	for {
		_, receivedMessage, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("websocket received message error: %s", err)
		}

		var receivedEnvelope pb.Envelope
		if err := proto.Unmarshal(receivedMessage, &receivedEnvelope); err != nil {
			log.Fatalf("unmarshal received envelope error: %s", err)
		}
		payloadType := receivedEnvelope.PayloadType
		//log.Printf("received envelope payloadType: %v\n", payloadType)

		switch receivedEnvelope.PayloadType {
		case pb.PayloadType_FIRST_CONNECT:
			var firstConnectResponse pb.FirstConnectResponse
			if err := proto.Unmarshal(receivedEnvelope.Payload, &firstConnectResponse); err != nil {
				log.Fatalf("unmarshal received %v payload error: %s", payloadType, err)
			}
			log.Printf("received FirstConnectResponse: %v", firstConnectResponse.ConnectId)
			MyConnectId = firstConnectResponse.ConnectId

		case pb.PayloadType_SERVER_TIMING_BROADCAST:
			var serverTimingBroadCastRequest pb.ServerTimingBroadCastRequest
			if err := proto.Unmarshal(receivedEnvelope.Payload, &serverTimingBroadCastRequest); err != nil {
				log.Fatalf("unmarshal received %v payload error: %s", payloadType, err)
			}
			log.Printf("received ServerTimingBroadCastRequest: %v", serverTimingBroadCastRequest.Msg)

		case pb.PayloadType_GET_OTHER_CONNECT_IDS:
			var getOtherConnectIdsResponse pb.GetOtherConnectIdsResponse
			if err := proto.Unmarshal(receivedEnvelope.Payload, &getOtherConnectIdsResponse); err != nil {
				log.Fatalf("unmarshal received %v payload error: %s", payloadType, err)
			}
			log.Printf("received GetOtherConnectIdsResponse: %v", getOtherConnectIdsResponse.OtherConnectIds)
			OtherConnectIds = getOtherConnectIdsResponse.OtherConnectIds

		case pb.PayloadType_KILL_OTHER_CONNECT:
			var killOtherConnectResponse pb.KillOtherConnectResponse
			if err := proto.Unmarshal(receivedEnvelope.Payload, &killOtherConnectResponse); err != nil {
				log.Fatalf("unmarshal received %v payload error: %s", payloadType, err)
			}
			log.Printf("received KillOtherConnectResponse: %v", killOtherConnectResponse.IsSuccess)

		case pb.PayloadType_REDIS_NEW_STRING:
			var redisNewStringResponse pb.RedisNewStringResponse
			if err := proto.Unmarshal(receivedEnvelope.Payload, &redisNewStringResponse); err != nil {
				log.Fatalf("unmarshal received %v payload error: %s", payloadType, err)
			}
			log.Printf("received RedisNewStringResponse: %v %v", redisNewStringResponse.IsSuccess, redisNewStringResponse.Result)

		case pb.PayloadType_REDIS_GET_STRING:
			var redisGetStringResponse pb.RedisGetStringResponse
			if err := proto.Unmarshal(receivedEnvelope.Payload, &redisGetStringResponse); err != nil {
				log.Fatalf("unmarshal received %v payload error: %s", payloadType, err)
			}
			log.Printf("received RedisGetStringReponse: %v %v", redisGetStringResponse.IsSuccess, redisGetStringResponse.Result)
		}

	}
}
