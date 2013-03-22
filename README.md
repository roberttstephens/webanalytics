#In development, do not use.

If you want web analytics, you can use one or more of several third party services. I wanted to come up with a simple performant open source application that covers some common use cases.

Some example use cases:

How many page views am I getting? (Sometimes it's difficult to tell with varnish)
What percentage of users are still on IE x?
Are users clicking on anything in the related articles block on my page?
Should I display 20 blog posts, or 10? (Is the content on the bottom being clicked on?)

The database design is purposefully simple. It needs to be performant on writes. Very little processing is done on the way out. We're expecting to deal with a large number of records. It's probably best to set up a cron that moves records to back up tables. I'm going to play around with consolodating into a more relational design.


Tables

page_view
id, ip_address, url, timestamp, user_agent, screen size, 

href_click
pvid, href, href_location, timestamp

url = document.URL
browser = http://stackoverflow.com/questions/2400935/browser-detection-in-javascript
href_location = getBoundingClientRect
  
top: -337
right: 391
bottom: -337
left: 164

rectangle will be a postgres type 'box'
it will be inserted as
top, right, bottom, left

##TODO

###Domains
create a table for domains
add columns to page_view and href_click for did (domain id)
on the start of analytics, load all domains and primary id into memory.
on a request, map the domain to the id. 
insert did into table.

###Batch insert
Read from app.json and wrap multiple inserts into a transaction.
