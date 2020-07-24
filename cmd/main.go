package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"github.com/dkpeakbil/taskserver/api"
	"github.com/dkpeakbil/taskserver/domain"
	"github.com/dkpeakbil/taskserver/game"
	"github.com/dkpeakbil/taskserver/repository"
	"github.com/dkpeakbil/taskserver/repository/mem"
	usecase "github.com/dkpeakbil/taskserver/usecase/impl"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var (
	r repository.Repository
)

func main() {
	var (
		apiAddr = flag.String("api-addr", ":8081", "api listen address")
		wsAddr  = flag.String("ws-addr", ":8080", "websocket listen address")
		debug   = flag.Bool("debug", false, "debug mode")
	)
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.Printf("api addr: %s, ws addr: %s\n", *apiAddr, *wsAddr)

	repo, _ := mem.NewInMemoryRepository()

	if *debug {
		seedUsers(repo)
	}

	log.Printf("repository %s", repo)

	validation := validator.New()
	ucase, _ := usecase.NewUseCase(repo, validation)

	log.Printf("ucase %s", ucase)

	gameApi, _ := api.NewApi(*apiAddr, ucase)
	go gameApi.Run()

	gameServer, _ := game.NewGame(*wsAddr, ucase)
	go gameServer.Run()

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, os.Kill)
	select {
	case <-sg:
		log.Printf("got interrupted %s", sg)
	}
}

func seedUsers(repo repository.Repository) {
	log.Debugf("seeding users")

	pass := md5.Sum([]byte("123456"))

	user1 := &domain.User{
		Username: "enes",
		Password: hex.EncodeToString(pass[:]),
	}
	_, _ = repo.Save(user1)

	user2 := &domain.User{
		Username: "kursad",
		Password: hex.EncodeToString(pass[:]),
	}
	_, _ = repo.Save(user2)
}
