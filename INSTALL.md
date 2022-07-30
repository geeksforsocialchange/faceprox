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
2. Set up the firewall as below:

```
sudo iptables -A INPUT -p tcp -m multiport --dports 80,443 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
sudo iptables -A OUTPUT -p tcp -m multiport --dports 80,443 -m conntrack --ctstate ESTABLISHED -j ACCEPT
```

3. `git clone https://github.com/geeksforsocialchange/faceprox.git`
4. `cd faceprox`
5. Set your hostname and email in the `.env` file
6. Run `docker-compose up` and check everything works as it should

When you're satisfied everything is working, quit the server with `ctrl-c` and run the following to load in daemon mode:

```
docker-compose up -d
```

You can now exit ssh and your server is complete.
