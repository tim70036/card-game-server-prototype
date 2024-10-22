package server

import (
	commonserver "card-game-server-prototype/pkg/common/server"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/coregrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type Initiator interface{}

func Init(
	connectionServiceServer *commonserver.ConnectionServiceServer,
	emoteRpcServiceServer *commonserver.EmoteRpcServiceServer,
	actionServiceServer *ActionServiceServer,
	messageServiceServer *MessageServiceServer,
	grpcServer *commonserver.GrpcServer,
) Initiator {
	coregrpc.RegisterConnectionServiceServer(
		grpcServer,
		connectionServiceServer,
	)

	commongrpc.RegisterEmoteRpcServiceServer(
		grpcServer,
		emoteRpcServiceServer,
	)

	txpokergrpc.RegisterActionServiceServer(
		grpcServer,
		actionServiceServer,
	)

	txpokergrpc.RegisterMessageServiceServer(
		grpcServer,
		messageServiceServer,
	)
	go grpcServer.Run()
	return nil
}
