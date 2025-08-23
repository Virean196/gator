package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/virean196/gator/internal/config"
	"github.com/virean196/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {

	var s state
	s.cfg, _ = config.Read()
	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		fmt.Print("error opening db")
		os.Exit(2)
	}
	dbQueries := database.New(db)
	s.db = dbQueries
	commands := commands{handlers: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handleReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	args := os.Args
	if len(args) < 2 {
		fmt.Print("not enough args ")
		os.Exit(1)
	}
	command := command{name: args[1], args: args[1:]}
	commands.run(&s, command)

}
