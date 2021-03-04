# rworker

## Manual testing (WIP)

To be able to quickly test this implementation, I have hardcoded a job id. With
a post.json containing the following:


```json
{
  "command": "/some_script.sh",
  "arguments": [
    "arg1",
    "arg2",
    "arg3"
  ]
}

```


These commands submit the corresponding job and fetch the status and the log:
```sh
curl -k -X POST -H 'Content-Type: application/json' -d @post.json https://127.0.0.1:8443/job
curl -k -X GET https://127.0.0.1:8443/job/01000000000000000000000000000000
curl -k -X GET https://127.0.0.1:8443/log/01000000000000000000000000000000
```
