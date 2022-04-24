package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Chrn1y/grpc-mafia/storage"
	"io"
	"log"
	"net/http"
	"os"
)

const(
	port = "8001"
	contentType = "application/json"
)

func main() {
	fmt.Println("Введите адрес сервера:")
	addr := ""
	fmt.Scanln(&addr)
	addr += "8001"
	//ctx := context.Background()
	for {
		fmt.Println("Введите номер команды:\nРегистрация пользователя - 1\nУдаление пользователя - 2\nОбновление данных - 3\nПолучение данных пользователя - 4\nПолучить список всех пользователей - 5\nВыход - 0")
		cmd := ""
		fmt.Scanln(&cmd)
		switch cmd {
		case "1":
			body, _ := json.Marshal(getUser())
			_, err := http.Post(addr + "/user", contentType, bytes.NewReader(body))
			if err != nil {
				log.Fatal(err) 
			}
		case "2":
			fmt.Println("Введите айди пользователя")
			id := ""
			fmt.Scanln(&id)
			_, err := http.NewRequest(http.MethodDelete, addr + "/user/" + id, nil)
			if err != nil {
				log.Fatal(err)
			}
		case "3":
			fmt.Println("Введите айди пользователя")
			id := ""
			fmt.Scanln(&id)
			user := getUser()
			user.ID = id
			body, _ := json.Marshal(user)
			_, err := http.NewRequest(http.MethodPut, addr + "/user/" + id, bytes.NewReader(body))
			if err != nil {
				log.Fatal(err)
			}
		case "4":
			fmt.Println("Введите айди пользователя")
			id := ""
			fmt.Scanln(&id)
			get, err := http.Get(addr + "/user/" + id)
			if err != nil {
				log.Fatalln(err)
			}
			b, err := io.ReadAll(get.Body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(b))
		case "5":
			get, err := http.Get(addr + "/users")
			if err != nil {
				log.Fatalln(err)
			}
			b, err := io.ReadAll(get.Body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(b))
		case "0":
			return
		}
	}
}

func getUser() *storage.User {
	res := &storage.User{}
	fmt.Println("Введите имя пользователя")
	name := ""
	fmt.Scanln(&name)
	res.Name = name
	gender := ""
	fmt.Scanln(&gender)
	fmt.Println("Введите пол пользователя: Мужчина - 1; Женщина - 2")
	switch gender {
	case "1":
		res.Gender = storage.Male
	case "2":
		res.Gender = storage.Female
	}
	fmt.Println("Введите адрес электронной почты пользователя")
	email := ""
	fmt.Scanln(&email)
	res.Name = email
	imagePath := ""
	fmt.Scanln(&imagePath)

	fmt.Println("Укажите путь к аватарке")

	file, err := os.Open(imagePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	res.Avatar = bytes
	return res
}
