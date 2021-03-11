# Certificates

## Certificate generation
### Server

mkcert is used to easily fix the warnings that come with Go 1.5 and 1.6 (see
https://github.com/golang/go/issues/39568).

```sh
mkcert -cert-file certs/server/cert.pem -key-file certs/server/key.pem localhost
```

Stored in
certs/server/cert.pem
certs/server/key.pem


### Client

With the UUID in the CN:

```sh
UUID="405a10fa-083f-49c6-ba79-f6afa3db9bb7" openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 -out certs/client/cert.pem -keyout certs/client/key.pem -subj "/CN=$UUID"
```
