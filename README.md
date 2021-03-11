# rworker

## Building
```sh
go build ./...
```

To compile the server and client binaries:
```sh
go build cmd/client/client.go
go build cmd/server/server.go
```

As of now, the paths used for the server certificates are hardcoded so that
running from the root of the repository works.

## Testing

```sh
go test ./...
```

## Manual testing

### Setup

From the root of the repository:
```sh
./server --sk test/server/certs/server/key.pem --sc test/server/certs/server/cert.pem --addr localhost:8443 --cc test/server/certs/server/known_clients/client.pem
```

### CLI

The client has to be compiled first:
```sh
go build cmd/client/client.go
```


We can add a job:
```sh
./client --cc test/server/certs/client/cert.pem --ck test/server/certs/client/key.pem --sc test/server/certs/server/cert.pem --addr localhost:8443 add ls /

```

The server answers with the job ID, e.g. with:
```
{
  "id": "72a4b8b0-78ef-4085-91fd-324617849728"
}
```


We can fetch the log and status:
```sh
./client --cc test/server/certs/client/cert.pem --ck test/server/certs/client/key.pem --sc test/server/certs/server/cert.pem --addr localhost:8443 log "72a4b8b0-78ef-4085-91fd-324617849728"
```


Response:
```
{
  "stdout": "bin\nboot\nconfig\ndev\ndownloads\netc\nhome\ninitrd.img\ninitrd.img.old\nlib\nlib32\nlib64\nlibx32\nlost+found\nmedia\nmnt\nopt\nproc\nroot\nrun\nsbin\nsrv\nsys\ntmp\nusr\nvar\nvmlinuz\nvmlinuz.old",
  "stderr": ""
}
```



```sh
./client --cc test/server/certs/client/cert.pem --ck test/server/certs/client/key.pem --sc test/server/certs/server/cert.pem --addr localhost:8443 status "72a4b8b0-78ef-4085-91fd-324617849728"
```

Response:
```
{
    "id": "72a4b8b0-78ef-4085-91fd-324617849728",
        "command": "ls",
        "arguments": [
            "/"
        ],
        "status": "done",
        "exitCode": 0
}
```

Or stop a long running job:
```sh
./client --cc test/server/certs/client/cert.pem --ck test/server/certs/client/key.pem --sc test/server/certs/server/cert.pem --addr localhost:8443 stop "72a4b8b0-78ef-4085-91fd-324617849728"
```
