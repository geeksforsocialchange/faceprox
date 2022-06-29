# Faceprox

Faceprox is an alternative front-end for Facebook event pages. It has no ads or tracking, and doesn't urge you to sign up.

We're not aware of anybody running public instances of Faceprox yet, so you will need to run it yourself.

## Installation

Download the [latest release](https://github.com/geeksforsocialchange/faceprox/releases/latest) for your system and run it. It will listen on port 8000

Alternatively you can build a docker image by cloning this repository, building an image with `docker build -t faceprox .`, and then run that with `docker run -ti -p 8000:8000 faceprox`

### Docker Compose

If you just want to run this with an SSL certificate using docker-compose:

- Edit the [.env](.env) file to set the hostname and an email
- Check everything is working: `docker-compose up`

## Usage

Faceprox listens for a few different paths. These will render html, json and ical formats for both individual events as well as the next five events from a single page

Open your web browser to http://127.0.0.1:8000 to see the welcome screen, then try out the following paths:

* /events/:eventid:
* /events/:eventid:.json
* /events/:eventid:.ics
* /page/:pagename:
* /page/:pagename:.json
* /page/:pagename:.ics

where `:eventid:` is a numeric event ID for a single event, and `:pagename:` is a page name for a page that might have multiple events. You can find both of these in the original Facebook URLs

## Known Issues

* The html version is really just for basic testing. It has more javascript than mbasic.facebook.com, and the images may get blocked by Firefox containers or other tracking blockers
