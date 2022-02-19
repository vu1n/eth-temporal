# Ethereum Temporal Project

This project will spin up a [Temporal](https://temporal.io/) cluster. Results are initially stored in a Postgres instance.


```
❯ make
up         Spin up Temporal cluster
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

### Worker Program
To start the worker program. We can scale horizontally by running more programs.

`$ make worker`

### Start fetching blocks
The included go program will spawn worker processes to fetch new blocks, starting from the `LATEST` block and incrementing from there.

`$ make fetch`

### Backfilling
To backfill specific blocks, run the included go program to populate the task queue. Currently set to use the same task queue as the fetch task.

Start the backfill worker:

`$ make backfill-worker`

Add the backfill tasks:

`$ go run ./backfiller/ -start 200 -end 210 -size 3`

### API
There is a basic API to query block numbers. If the block number does not exist in the postgres database, it will start a workflow to fetch the block from the web3 provider and then return the results.

Start the api server:

`$ make api`

Make an HTTP GET request to:

`http://localhost:8081/blockNumber/{blockNumber}[0-9]+`