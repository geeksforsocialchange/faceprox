# Faceprox

A privacy-respecting proxy for Facebook content. Currently only handles event pages.

Listens on port 8000 and responds to the following paths:

* /events/:eventid:
* /events/:eventid:.json
* /events/:eventid:.ics
* /page/:pagename:
* /page/:pagename:.json
* /page/:pagename:.ics