Gator is a Blog Aggregator that takes feeds from certain blogs and stores the posts in a Postgresql database!

The database was designed with Goose for easier migrations and SQLC for the type-safe sql to go!

To run this program you'll need Postgres and Go installed.

After you have Go and Postgres installed you have to:
- Create a config file at your Home folder named .gatorconfig.json that has the following:
    {"db_url":"postgres://user:password@ip:port/gator?sslmode=disable"} (edit this to your own Postgres DB)
- Clone the repository wherever you want
- Open the terminal and change to the cloned folder
- In the root folder type "go install"

After that, you can run the program by typing "gator command args"

Here are all the commands available at the moment:
- "login <name>" - Logs into an already registered user (no authentication, just as a way to save who added which feed)
- "register <name>" - Registers a user and sets it as default user, also logs in right away
- "reset" - Resets the entire DB, use only for testing
- "users" - Prints a list of all the users
- "agg <interval (ex: 50s, 5m, 1h)>" - Runs the aggregator every X s/m/h, fetches all the feeds and stores them in the db, also creates posts based on same feeds (this is an infinite loop so be mindful with the timings to not provoke DOS on the blogs)
- "addfeed <url>" - Adds a feed to the user 
- "feeds" - Lists all the feeds and their user
- "follow <feed_name>" - Allows the user to follow a feed created by other users
- "following" - Lists all feeds the user follows
- "unfollow <feed_name>" - Unfollows said feed
- "browse <limit (default: 2)>" - Shows the latest X posts 
