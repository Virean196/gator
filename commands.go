package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	if len(cmd.args) < 1 {
		fmt.Print("no comamnd found\n")
		os.Exit(1)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[1])
	if err != nil {
		return err
	}
	fmt.Printf("Colleting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 3 {
		fmt.Print("feed name and link are required\n")
		os.Exit(1)
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID: int32(uuid.New().ID()), Name: cmd.args[1], Url: cmd.args[2], UserID: user.ID})
	if err != nil {
		return fmt.Errorf("error creating feed")
	}
	fmt.Printf("Feed ID: %v\nFeed Name: %s\nFeed URL: %v\nUser ID: %v\n", feed.ID, feed.Name, feed.Url, feed.UserID)
	fmt.Print("\nFollowing the feed...\n")
	res, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: int32(uuid.New().ID()), CreatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("error creating feed follow")
	}
	fmt.Printf("Feed Name: %s\nCurrent User: %s\nFeed Followed!\n", res.FeedName, res.UserName)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Print("no comamnd found\n")
		os.Exit(1)
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.args[1])
	if err != nil {
		return fmt.Errorf("invalid url, not a followable feed")
	}
	res, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: int32(uuid.New().ID()), CreatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("error creating feed follow")
	}
	fmt.Printf("Feed Name: %s\nCurrent User: %s\n", res.FeedName, res.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Print("no comamnd found\n")
		os.Exit(1)
	}
	followingFeeds, err := s.db.GetFollowing(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error getting the following for the current user")
	}
	for _, feed := range followingFeeds {
		fmt.Printf("- %s\n", feed)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Print("no comamnd found\n")
		os.Exit(1)
	}
	if len(cmd.args) < 2 {
		fmt.Print("no url found\n")
		os.Exit(1)
	}
	s.db.Unfollow(context.Background(), database.UnfollowParams{Name: user.Name, Url: cmd.args[1]})
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) < 1 {
		fmt.Print("no comamnd found\n")
		os.Exit(1)
	}
	if len(cmd.args) > 1 {
		limit, _ = strconv.Atoi(cmd.args[1])
	}
	posts, err := s.db.GetPostsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("Title: %s\nLink: %s\nDescription: %s\nPublish Date: %v\n\n", posts[i].Title, posts[i].Url, posts[i].Description, posts[i].PublishedAt)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedTofetch(context.Background())
	if err != nil {
		return err
	}
	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}, ID: feedToFetch.ID})
	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	cleanFeed(feed)
	if err != nil {
		return err
	}
	fmt.Printf("%s", feed.Channel.Item[1].Description)
	for _, item := range feed.Channel.Item {
		parsedTime, err := parsePubDate(item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing the time: %w", err)
		}
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{ID: int32(uuid.New().ID()), CreatedAt: time.Now(), Title: item.Title,
			Url: item.Link, Description: item.Description, PublishedAt: parsedTime, FeedID: feedToFetch.ID})
		if (err != nil) && (!strings.Contains(err.Error(), "duplicate key value violates unique constraint")) {
			return fmt.Errorf("error creating post: %w", err)
		}
	}
	return nil
}

func parsePubDate(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,                     // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,                      // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC3339,                      // "2006-01-02T15:04:05Z07:00"
		time.RFC822Z,                      // "02 Jan 06 15:04 -0700"
		"Mon, 02 Jan 2006 15:04:05 -0700", // Sometimes necessary as string literal
		"Mon, 02 Jan 2006 15:04:05 MST",   // Covers GMT/UTC
	}
	dateStr = strings.TrimSpace(dateStr)
	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", dateStr)
}
