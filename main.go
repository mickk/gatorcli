package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mickk/gatorcli/internal/config"
	"github.com/mickk/gatorcli/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.list[cmd.name]
	if !exists {
		return errors.New("command does not exist")
	}

	return f(s, cmd)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username arg is required")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("user has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("name arg is required")
	}

	regParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), regParams)
	if err != nil {
		return fmt.Errorf("unable to create user: %v", err)
	}

	fmt.Printf("user %v has been created\n", user.Name)

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("user has been set")

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	currUser := s.config.CurrentUserName
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == currUser {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	// if len(cmd.args) == 0 {
	// 	return errors.New("feed url is required")
	// }

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %v", err)
	}

	fmt.Printf("%v\n", feed)

	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("unable to read config file: %v\n", err)
		os.Exit(1)
	}

	s := state{
		config: &cfg,
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("cannot connect to db: %v\n", err)
		os.Exit(1)
	}

	s.db = database.New(db)

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("reset", handlerReset)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("command is required")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
