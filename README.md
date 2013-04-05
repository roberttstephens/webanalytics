#analytics.go

##Overview
If you want web analytics, you can use one or more of several third party services. Analytics.go is a simple performant open source application that covers some common use cases.

##Use cases

- How many page views am I getting? (Sometimes it's difficult to tell with varnish)
- On which URIs?
- What percentage of users are still on IE x?
- Which content do users click on?

##About the project

- Uses go (golang) to process requests, via a RESTful API.
- Uses a postgresql database. The database design is purposefully simple in order to be efficient with writes.
- Uses javascript to submit posts.

##How to use
- As mentioned in About the project, analytics.go uses postgresql. Create a user and database for analytics.
- cd /path/to/some/directory
- git clone https://github.com/roberttstephens/analytics.go.git
- cd analytics.go
- nano config.json #Configure with your database credentials.
- go build -o analytics
- ./analytics #Runs on port 8080 by default. Change in config/app.json
