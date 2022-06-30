# Setting up a faceprox droplet on Digital Ocean

## Create a new droplet

I used:

1. Marketplace distribution: Docker 19.x on Ubuntu
1. Basic plan, $5/mo
1. London datacentre (pick one near you)
1. Add a memorable hostname (I used `faceprox`)

## Set up DNS

Add an 'A' record for your domain of choice pointing to the new server's IP address.

## Set up the droplet

1. `ssh root@my-droplet-ip`
1. `git clone https://github.com/geeksforsocialchange/faceprox.git`
1. Set your hostname and email in the `.env` file
1. Run `docker-compose up` and check everything works as it should
1. When you're ready, run `docker-compose up -d` to run in daemon mode
