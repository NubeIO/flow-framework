#!/bin/bash

# Console colors
DEFAULT="\033[0m"
GREEN="\033[32m"
RED="\033[31m"

PRODUCTION=false
BUILD_ONLY=false
SYSTEM=false
EDGE28=false
MODBUS=false
LORA=false
BACNET=false
LORAWAN=false
BACNET_MASTER=false
HISTORY=false
INFLUX=false
RUBIXIO=false
MODBUSSERVER=false
POSTGRES=false

help() {
  echo "Service commands:"
  echo -e "   --prod | --production: add these suffix to start production"
  echo -e "   --build-only : don't run program"
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
    --build-only)
      BUILD_ONLY=true
      ;;
    --system)
      SYSTEM=true
      ;;
    --edge28)
      EDGE28=true
      ;;
    --modbus)
      MODBUS=true
      ;;
    --lora)
      LORA=true
      ;;
    --bacnet)
      BACNET=true
      ;;
    --lorawan)
      LORAWAN=true
      ;;
    --bacnetmaster)
      BACNET_MASTER=true
      ;;
    --history)
      HISTORY=true
      ;;
    --influx)
      INFLUX=true
      ;;
    --rubixio)
      RUBIXIO=true
      ;;
    --modbusserver)
      MODBUSSERVER=true
      ;;
    --postgres)
      POSTGRES=true
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

BUILD_ERROR=false

function buildPlugin {
  go build -buildmode=plugin -o system.so $2/*.go && cp system.so $pluginDir
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}BUILD $1"
  else
    echo -e "${RED}ERROR BUILD $1"
    BUILD_ERROR=true
  fi
}

pushd $dir > /dev/null

if [ ${SYSTEM} == true ]; then
  buildPlugin "SYSTEM" plugin/nube/system
fi
if [ ${EDGE28} == true ]; then
  buildPlugin "EDGE28" plugin/nube/protocals/edge28
fi
if [ ${MODBUS} == true ]; then
  buildPlugin "MODBUS" plugin/nube/protocals/modbus
fi
if [ ${LORA} == true ]; then
  buildPlugin "LORA" plugin/nube/protocals/lora
fi
if [ ${BACNET} == true ]; then
  buildPlugin "BACNET" plugin/nube/protocals/bacnetserver
fi
if [ ${LORAWAN} == true ]; then
  buildPlugin "LORAWAN" plugin/nube/protocals/lorawan
fi
if [ ${BACNET_MASTER} == true ]; then
  buildPlugin "BACNET_MASTER" plugin/nube/protocals/bacnetmaster
fi
if [ ${HISTORY} == true ]; then
  buildPlugin "HISTORY" plugin/nube/database/history
fi
if [ ${INFLUX} == true ]; then
  buildPlugin "INFLUX" plugin/nube/database/influx
fi
if [ ${RUBIXIO} == true ]; then
  buildPlugin "RUBIXIO" plugin/nube/protocals/rubixio
fi
if [ ${MODBUSSERVER} == true ]; then
  buildPlugin "MODBUSSERVER" plugin/nube/protocals/modbusserver
fi
if [ ${POSTGRES} == true ]; then
  buildPlugin "POSTGRES" plugin/nube/database/postgres
fi

if [ ${BUILD_ERROR} == true ]; then
    exit -1
fi

popd > /dev/null

if [ ${BUILD_ONLY} == true ]; then
    exit 0
fi

if [ ${PRODUCTION} == true ]; then
  go run app.go -g /data/flow-framework -d data --prod
else
  go run app.go
fi
