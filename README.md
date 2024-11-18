# Simple server

This simple server replies to your requests with a JSON document containing the following information:

```json
{
  "cookie_names": [
    index: "cookie name"
  ],
  "count": <Number of requests for a specific path from a specific source IP>,
  "ip": "Source IP (determined from either last hop or HTTP Heasder x_forwarded_for",
  "lasthop": "IP of last hop",
  "method": "HTTP Method",
  "path": "Requested path",
  "query_string": "everythinf afer the path",
  "x_forwarded_for":	"Proxy IP if there is any",
}
```

# Build and Run

```shell
docker build --no-cache -t goresp .
Step 1/6 : FROM golang:1.20-alpine
...

docker run --rm -p 80:80 --name responder goresp
Starting server on port 80...
yyyy-mm-ddThh:mm:ssZ - lasthop - x_forwarded_for - method - path - count - queryString
...
```
