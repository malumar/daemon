A package for creating applications and implementing plug-ins (e.g. supported by other system processes or remote services) for it, without directly using the application API.

Forking is not supported in Windows OS

## Example

An example application located in the **example/** folder that implements graceful reload and graceful stop http server

#### Run example

```go run example/example.go```

#### Reloading the application without losing connected clients

```
curl http://localhost:8080/reload
```

#### Graceful stop 

```
curl http://localhost:8080/stop
```
