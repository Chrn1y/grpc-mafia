package main

import (
	"bufio"
	"context"
	"fmt"
	mafia_proto "github.com/Chrn1y/grpc-mafia/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func vote(text string, variants []string) string {
	for {
		fmt.Println(text)
		for _, v := range variants {
			fmt.Println(v)
		}
		reader := bufio.NewReader(os.Stdin)
		chosen, _ := reader.ReadString(byte('\n'))
		chosen = chosen[:len(chosen) - 1]
		for _, v := range variants {
			if chosen == v {
				return chosen
			}
		}
		fmt.Println("Выбери из доступных вариантов")
	}
}

func main() {
	fmt.Print("Введи ip адрес сервера: ")
	reader := bufio.NewReader(os.Stdin)
	addr, _ := reader.ReadString(byte('\n'))
	addr = addr[:len(addr) - 1]
	fmt.Print("Введи имя пользователя: ")
	name, _ := reader.ReadString(byte('\n'))
	name = name[:len(name) - 1]
	ctx := context.Background()
	conn, err := grpc.Dial(addr + ":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	c := mafia_proto.NewAppClient(conn)
	client, err := c.Play(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.Send(&mafia_proto.Request{
		Data: &mafia_proto.Request_Register_{Register: &mafia_proto.Request_Register{Name: name}},
		Type: mafia_proto.RequestType_register,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		recv, err := client.Recv()
		if err != nil {
			fmt.Println(err)
			return
		}
		if recv.Type == mafia_proto.ResponseType_info{
			fmt.Println(recv.Data.(*mafia_proto.Response_Info_).Info.Text)
			if recv.Data.(*mafia_proto.Response_Info_).Info.End{
				return
			}
		} else { // vote
			chosen := vote(recv.Data.(*mafia_proto.Response_Vote_).Vote.Text,
				recv.Data.(*mafia_proto.Response_Vote_).Vote.Choose)
			err := client.Send(&mafia_proto.Request{
				Data: &mafia_proto.Request_Vote_{Vote: &mafia_proto.Request_Vote{Name: chosen}},
				Type: mafia_proto.RequestType_vote_request,
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
