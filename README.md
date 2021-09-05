
# getting started
rename the `config-example.yml` file to `config.yml`

# default port
1660

# plugins
## See plugin docs
/docs/plugins


## Build plugin
add into /data/plugins

```
go build -buildmode=plugin -o ehco.so echo.go
go build -buildmode=plugin -o echo.so echo.go  && cp echo.so  ../../../data/plugins/  && go run ../../../app.go
```
example to build and run the apps
```
cd plugin/example/system
go build -buildmode=plugin -o system.so system.go  && cp system.so  ../../../data/plugins/ && (cd ~/code/go/flow-framework  && go run app.go)
```

## Logging
```
debug: when we want to show information on debugging issue (we activate this mode on just debugging so will not be that much un-necessary logs)
info: when we want to show meaningful information for user
warn: when we want to give a warning for user for some operations
error: while error happens, show it on red alert  
```