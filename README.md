# REST

Any path except for `/stream` will result in an HTTP response containing the `Bearer` authorization content.
```
$ curl http://localhost:8080/rest
Authorization: Bearer GeneratedJWTTokenForPath:/rest
```

# Streaming

The same for streaming:
```
$ websocat ws://localhost:8080/stream
Authorization: Bearer GeneratedJWTTokenForPath:/stream
```
