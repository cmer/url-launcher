# URL Launcher

A very very simple Golang web server that takes a POST at `/` and opens the URL passed in the request body in the default browser.

## Run it

```bash
./url-launcher 0.0.0.0:8080
```

## Try it

```bash
curl --request POST \
  --url http://127.0.0.1:8080/ \
  --header 'Content-Type: text/plain' \
  --data 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
```
