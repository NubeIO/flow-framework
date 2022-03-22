#!/bin/bash

# Console colors
DEFAULT="\033[0m"
GREEN="\033[32m"
RED="\033[31m"

PRODUCTION=false

help() {
    echo "Service commands:"
    echo -e "   ${GREEN}--prod | --production: add these suffix to start production"
}

parseCommand() {
    for i in "$@"; do
        case ${i} in
        -h | --help)
            help
            exit 0
            ;;
        --prod | --production)
            PRODUCTION=true
            ;;
        *)
            echo -e "${RED}Unknown options ${i}  (-h, --help for help)${DEFAULT}"
            ;;
        esac
    done
}

parseCommand "$@"

dir=$(pwd)
echo -e "${GREEN}Current working directory is: $dir${DEFAULT}"
pluginDir=$dir/data/plugins

if [ ${PRODUCTION} == true ]; then
    echo -e "${GREEN}We are running in production mode!${DEFAULT}"
    pluginDir=/data/flow-framework/data/plugins
else
    echo -e "${GREEN}We are running in development mode!${DEFAULT}"
fi

echo -e "${GREEN}Creating a plugin directory if does not exist at: ${pluginDir}${DEFAULT}"
rm -rf $pluginDir/* || true
mkdir -p $pluginDir

cd $dir/plugin/nube/system/
go build -buildmode=plugin -o system.so *.go  && cp system.so  $pluginDir

cd $dir/plugin/nube/utils/backup
go build -buildmode=plugin -o backup.so *.go  && cp backup.so  $pluginDir

cd $dir/plugin/nube/utils/mqttapi
go build -buildmode=plugin -o mqttapi.so *.go  && cp mqttapi.so  $pluginDir

cd $dir/plugin/nube/utils/git
go build -buildmode=plugin -o git.so *.go  && cp git.so  $pluginDir

cd $dir/plugin/nube/protocals/edge28
go build -buildmode=plugin -o edge28.so *.go  && cp edge28.so $pluginDir

cd $dir/plugin/nube/protocals/lorawan
go build -buildmode=plugin -o lorawan.so *.go  && cp lorawan.so $pluginDir

cd $dir/plugin/nube/protocals/modbus
go build -buildmode=plugin -o modbus.so *.go  && cp modbus.so $pluginDir

cd $dir/plugin/nube/protocals/lora
go build -buildmode=plugin -o lora.so *.go  && cp lora.so $pluginDir

cd $dir/plugin/nube/protocals/bacnetserver
go build -buildmode=plugin -o bacnetserver.so *.go  && cp bacnetserver.so $pluginDir

cd $dir/plugin/nube/protocals/bacnetmaster
go build -buildmode=plugin -o bacnetmaster.so *.go  && cp bacnetmaster.so $pluginDir

cd $dir/plugin/nube/database/influx
go build -buildmode=plugin -o influx.so *.go  && cp influx.so $pluginDir

cd $dir/plugin/nube/database/history
go build -buildmode=plugin -o history.so *.go  && cp history.so $pluginDir

cd $dir/plugin/nube/protocals/broker
go build -buildmode=plugin -o broker.so *.go  && cp broker.so $pluginDir

cd $dir

if [ ${PRODUCTION} == true ]; then
  go run app.go -g /data/flow-framework  -d data --prod
else
  go run app.go
fi
