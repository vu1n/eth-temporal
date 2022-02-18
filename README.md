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
worker     Start the worker
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