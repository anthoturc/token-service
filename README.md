## Token Service

### What?

This is a test application to demonstrate some of the approaches I use to wards building web services in Golang. 

### Build?

This is super simple. No Makefile needed. 
```
go build .
```

### Run

Start the Jaeger all-in-one tracing server via a Docker image. *Note*: You will need to `docker rm jaeger` if you want to run this script more than once.
```
./scripts/init_dependencies.sh
```

Run the binary directly

```
./token-service
```

*Or*

Build a Docker image and run it

```
docker build -t token-service .
```

```
docker run --name token-service --rm token-service:latest
```

I didn't run in "detached" mode but you could easily do that!

### Deploy

Deployments are done to Digital Ocean's app platform. You can deploy anywhere you'd like. I chose DO because it was simple to get something setup and
the pricing is super easy to understand. 

Once you have setup your DO account -- sign up does require a CC but you can avoid being charged at all if you are tearing down your infrastructure after getting the results you want. 

I am deploying directly from this GitHub repository. This required me to install the Digital Ocean github app. You should do the same by going through the settings of 
your repo.

```
doctl apps create --spec spec.yml
```
