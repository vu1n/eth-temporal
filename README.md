# Ethereum Temporal Project

This project will spin up a [Temporal](https://temporal.io/) cluster. Results are initially stored in a Postgres instance.


```
❯ make
up         Spin up Temporal cluster
up-b       Spin up and force build Temporal cluster
down       Destroy the Temporal cluster
stop       Stop the Temporal cluster
ps         Check the status of Temporal services
shell      Start a shell with the Temporal CLI
db-init    Initialize the database
api        Start up API server
worker     Start the worker
bworker    Start the backfill worker
fetch      Fetch latest after worker has started
test       Run tests
```

## How to start

### Temporal cluster
The project uses a (slightly) modified docker-compose from [Temporal](https://github.com/temporalio/docker-compose)

To spin up the Temporal cluster

`$ make up`

To destroy the cluster

`$ make down`

To stop the cluster

`$ make stop`

### Start fetching blocks
The included go program will spawn worker processes to fetch new blocks, starting from the `LATEST` block and incrementing from there.

`$ make fetch`

### Backfilling
Add the backfill tasks:

`$ go run ./backfiller/ -start 200 -end 210 -size 3`

### API
Make an HTTP GET request to:

`http://localhost:8081/blockNumber/{blockNumber}[0-9]+`

`http://localhost:8081/traceBlock/{blockNumber}[0-9]+`