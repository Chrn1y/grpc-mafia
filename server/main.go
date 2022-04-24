package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Chrn1y/grpc-mafia/mafia"
	mafia_proto "github.com/Chrn1y/grpc-mafia/proto"
	"github.com/Chrn1y/grpc-mafia/storage"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Impl struct {
	mafia_proto.UnimplementedAppServer
	m *mafia.Mafia
	s *storage.Storage
}

func (i *Impl) Play(stream mafia_proto.App_PlayServer) error {
	in, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if in.GetType() != mafia_proto.RequestType_register {
		return errors.New("need to register first")
	}
	name := in.GetData().(*mafia_proto.Request_Register_).Register.Name
	if !i.s.Validate(name) {
		return errors.New("need to register first")
	}
	input, output, err := i.m.Join(name)
	if err != nil {
		return err
	}
	for {
		data := <-output
		err = stream.Send(data)
		if err != nil {
			return err
		}
		if data.GetType() == mafia_proto.ResponseType_vote_response {
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

func main() {

	s := storage.New()

	r := mux.NewRouter()
	insert := func(w http.ResponseWriter, r *http.Request) {
		user := &storage.User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id := s.Insert(user)
		//vars := mux.Vars(r)
		_, _ = w.Write([]byte(id))
		w.WriteHeader(http.StatusOK)
	}
	del := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		s.Delete(id)
		w.WriteHeader(http.StatusOK)
	}
	get := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		user, ok := s.Get(id)
		if !ok {
			http.Error(w, "unknown user", http.StatusBadRequest)
			return
		}
		data, _ := json.Marshal(user)
		_, err := w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	list := func(w http.ResponseWriter, r *http.Request) {
		user := s.List()
		data, _ := json.Marshal(user)
		_, err := w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	update := func(w http.ResponseWriter, r *http.Request) {
		user := &storage.User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		vars := mux.Vars(r)
		id := vars["id"]
		ok := s.Update(id, user)
		if !ok {
			http.Error(w, "unknown user", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	r.HandleFunc("/user/{id}", get).Methods(http.MethodGet)
	r.HandleFunc("/users", list).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", del).Methods(http.MethodDelete)
	r.HandleFunc("/user/{id}", update).Methods(http.MethodPut)
	r.HandleFunc("/user", insert).Methods(http.MethodPost)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	// -----------------------
	serv := grpc.NewServer()
	mafia_proto.RegisterAppServer(serv, &Impl{
		m: mafia.New(5),
		s: s,
	})
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		return
	}
	fmt.Println("Starting server...")
	if err = serv.Serve(l); err != nil {
		return
	}
}
