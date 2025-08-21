package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/virean196/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	value, exists := c.handlers[cmd.name]
	if exists {
		err := value(s, cmd)
		if err != nil {
			return fmt.Errorf("error running the function")
		}
		return err
	} else {
		return fmt.Errorf("command doesnt exist")
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) <= 1 {
		fmt.Print("username is required\n")
		os.Exit(1)
	}
	userData, err := s.db.GetUser(context.Background(), cmd.args[1])
	if err != nil {
		fmt.Print("You can't login to an account that doesn't exist!\n")
		os.Exit(1)
	}
	s.cfg.SetUser(userData.Name)
	fmt.Printf("User: %s has been set\n", userData.Name)
	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Print("failed reseting users table")
		os.Exit(1)
	}
	fmt.Print("db reset successfully\n")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	userData, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error obtaining list of users")
		os.Exit(1)
	}
	for i := range userData {
		if userData[i].Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", userData[i].Name)
		} else {
			fmt.Printf("* %s\n", userData[i].Name)
		}
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) <= 1 {
		fmt.Print("username is required\n")
		os.Exit(1)
	}
	userData, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: int32(uuid.New().ID()), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[1]})
	if err != nil {
		fmt.Print("error creating the user")
		os.Exit(1)
	}
	s.cfg.SetUser(userData.Name)
	fmt.Printf("User %s registered!\n", userData.Name)
	fmt.Printf("User data: %v\n", userData)
	return nil
}

func handlerAgg(s *state, cmd command) error {
	var rssfeed *RSSFeed
	rssfeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed")
	}
	cleanFeed(rssfeed)
	fmt.Print(rssfeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 3 {
		fmt.Print("feed name and link are required\n")
		os.Exit(1)
	}
	currentUser := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("user doesn't exist")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID: int32(uuid.New().ID()), Name: cmd.args[1], Url: cmd.args[2], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("error creating feed")
	}
	fmt.Printf("Feed ID: %v\nFeed Name: %s\nFeed URL: %v\n User ID: %v", feed.ID, feed.Name, feed.Url, feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		fmt.Print("no comamnd found\n")
		os.Exit(1)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Print("error obtaining feeds")
		os.Exit(1)
	}
	for i := range feeds {
		user, err := s.db.GetFeedUser(context.Background(), feeds[i].UserID)
		if err != nil {
			fmt.Print("invalid user id")
			os.Exit(1)
		}
		fmt.Printf("Feed ID: %v\nFeed Name: %s\nFeed Url: %s\nFeed User Name: %s\n\n", feeds[i].ID, feeds[i].Name, feeds[i].Url, user.Name)
	}
	return nil
}
