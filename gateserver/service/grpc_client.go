package service

import (
	"io"
	"log"
	"time"

	pb "cuttleserver/common/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

type Tasks struct {
	AccountID uint32
	MessageID uint32
	Content   []byte
}

type GameServer struct {
	GSStreamer pb.Stream_ForwardClient
	tasks      chan *Tasks
}

func NewGameServer() *GameServer {
	gs := &GameServer{
		tasks: make(chan *Tasks),
	}

	gs.initClient()

	go gs.Recv()
	go gs.Send()

	return gs
}

func (self *GameServer) initClient() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	stream := pb.NewStreamClient(conn)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2)
	defer cancel()

	self.GSStreamer, err = stream.Forward(ctx)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

func (self *GameServer) Send() {
	for {
		select {
		case t := <-self.tasks:
			err := self.GSStreamer.Send(&pb.RPCRequest{AccountID: t.AccountID, MessageID: t.MessageID, MessageBody: t.Content})
			if err != nil {
				if err == io.EOF {
					log.Fatalf("GameServer Exit: %v", err)
				}
				log.Fatalf("send message error: %v", err)
			}
		}
	}
}

func (self *GameServer) Recv() {
	for {
		r, err := self.GSStreamer.Recv()
		if err != nil {
			if err == io.EOF {
				log.Fatalf("GameServer Exit: %v", err)
			}
			log.Fatalf("recv message error: %v", err)
		}

		log.Printf("Resp Message %v, %v, %v", r.AccountID, r.MessageID, r.MessageBody)
		sess, err := sessionManager.GetSession(r.AccountID)
		if err != nil {
			log.Fatalf("could not find session : %v", r.AccountID)
			break
		}
		sess.WriteMessage(uint16(r.MessageID), r.MessageBody)
	}
}

/*
func (self *GameServer) DialGRPC() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	for {
		select {
		case t := <-self.tasks:
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2)
			defer cancel()

			r, err := c.SayHello(ctx, &pb.HelloRequest{AccountID: t.AccountID, MessageID: t.MessageID, Content: t.Content})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			log.Printf("Greeting: %v, %v, %v", r.AccountID, r.MessageID, r.Content)
			sess, err := sessionManager.GetSession(r.AccountID)
			if err != nil {
				log.Fatalf("could not find session : %v", r.AccountID)
				break
			}
			sess.WriteMessage(uint16(r.MessageID), r.Content)
		}
	}
}
*/
