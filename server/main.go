package main

import (
	"errors"
	"fmt"
	"github.com/Chrn1y/grpc-mafia/mafia"
	mafia_proto "github.com/Chrn1y/grpc-mafia/proto"
	"google.golang.org/grpc"
	"io"
	"net"
)

type Impl struct {
	mafia_proto.UnimplementedAppServer
	m *mafia.Mafia
}

func (i *Impl) Play(stream mafia_proto.App_PlayServer) error {
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if in.GetType() != mafia_proto.RequestType_register{
		return errors.New("need to register first")
	}
	input, output, err := i.m.Join(in.GetData().(*mafia_proto.Request_Register_).Register.Name)
	if err != nil {
		return err
	}
	for {
		data := <-output
		//println("server sending", data.String())
		err = stream.Send(data)
		if err != nil {
			return err
		}
		if data.GetType() == mafia_proto.ResponseType_vote_response{
			data, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			input <- data
		}
	}
}

const num = 5

func main() {
	s := grpc.NewServer()
	mafia_proto.RegisterAppServer(s, &Impl{
		m:                      mafia.New(5),
	})
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		return
	}
	fmt.Println("Starting server...")
	if err = s.Serve(l); err != nil {
		return
	}
}
