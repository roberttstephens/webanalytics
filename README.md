#webanalytics

##Overview
If you want web analytics, you can use one or more of several third party services. Webanalytics is a simple performant open source application that covers some common use cases.

- How many page views am I getting? (Sometimes it's difficult to tell with varnish)
- On which URLs?
- What percentage of users are still on IE x?
- Which content do users click on?

##About the project

- This is not production ready yet, though feel free to try it and report any bugs.
- Uses go (golang) to process requests.
- Uses a postgresql database. The database design is purposefully simple in order to be efficient with writes.
- Inserts into the database in a goroutine so requests are handled concurrently.
- Uses javascript to submit posts.

##How to use

Webanalytics is broken into two parts. The server side application and the javascript.

###Server side application
If you have already set up your $GOPATH and added $GOPATH/bin to your $PATH you should:
- Create a postgres user and database for webanalytics.
- run "go get github.com/roberttstephens/webanalytics" without quotation marks.
- Copy $GOPATH/src/github.com/roberttstephens/webanalytics/config.json to somewhere of your choice.
- Edit config.json to reflect your new database connection and desired port.
- Run "webanalytics --config path/to/config.json" without quotation marks.


###Javascript
The javascript is in poor shape right now. However, you should be able to copy docs/webanalytics.js to your site, change your domain (and possibly port) and start receiving POSTs.  Please reach out to me if something doesn't work, so I can fix it.

## How to contribute
Right here on github.com is easiest.
 - Fork the project.
 - Make a commit or two.
 - Perform a [pull request](http://help.github.com/pull-requests/).
