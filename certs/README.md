# Certificate generation
## Server
```sh
openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 -out cert.pem -keyout key.pem -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost"
```

Stored in
certs/server/cert.pem
certs/server/key.pem

