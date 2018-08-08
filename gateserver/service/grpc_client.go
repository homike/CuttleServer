package service

import (
	"log"
	"os"
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
	tasks chan *Tasks
}

func NewGameServer() *GameServer {
	gs := &GameServer{
		tasks: make(chan *Tasks),
	}
	go gs.DialGRPC()

	return gs
}

func (self *GameServer) DialGRPC() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	_ = name

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
