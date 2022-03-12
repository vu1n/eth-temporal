FROM golang:1.17 as base

WORKDIR /go/src/eth-temporal/app

COPY go.mod go.sum shared.go structs.go ./

RUN go mod download

FROM base AS build

COPY activities ./activities
COPY api ./api
COPY backfill_worker ./backfill_worker
COPY backfiller ./backfiller
COPY starter ./starter
COPY worker ./worker
COPY workflows ./workflows

RUN go clean

RUN go install -v ./api
RUN go install -v ./backfill_worker
RUN go install -v ./backfiller
RUN go install -v ./starter
RUN go install -v ./worker

FROM golang:1.17 AS app

COPY --from=build /go/bin/api /usr/local/bin/et-api
COPY --from=build /go/bin/backfill_worker /usr/local/bin/et-backfill-worker
COPY --from=build /go/bin/backfiller /usr/local/bin/et-backfill
COPY --from=build /go/bin/starter /usr/local/bin/et-start
COPY --from=build /go/bin/worker /usr/local/bin/et-worker
