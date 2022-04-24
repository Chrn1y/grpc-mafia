package main

import (
	"fmt"
	"github.com/Chrn1y/grpc-mafia/storage"
)

func main() {
	fmt.Println("Введите адрес сервера:")
	addr := ""
	fmt.Scanln(&addr)
	//ctx := context.Background()
	for {
		fmt.Println("Введите номер команды:\nРегистрация пользователя - 1\nУдаление пользователя - 2\nОбновление данных - 3\nПолучение данных пользователя - 4\nПолучить список всех пользователей - 5\nВыход - 0")
		cmd := ""
		fmt.Scanln(&cmd)
		switch cmd {
		case "1":
		case "2":
		case "3":
		case "4":
		case "5":
		case "0":

		}
	}
}

func getUser() *storage.User {
	return nil
}
