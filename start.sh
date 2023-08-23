#!/bin/bash

echo "Determine current machine running the script."
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    *)          machine="UNKNOWN:${unameOut}"
esac
echo "'$machine' machine."

SUDO=""
ENV_FILE=".env.dev"
if [ $machine = "Linux" ]; then
  SUDO="sudo"
  ENV_FILE=".env"
fi

echo "Stoping and removing previous docker containers."
$SUDO docker stop $($SUDO docker ps -aq)
$SUDO docker rm $($SUDO docker ps -aq)
$SUDO docker volume rm $($SUDO docker volume ls -q)

if [ -z ${VOLUMES_PATH+x} ]; then
  source ./$ENV_FILE
fi
VOL_DIR=${VOLUMES_PATH#*=}

echo "Found '$VOL_DIR' as volumes path in environment file"
if [ -d "$VOL_DIR" ]; then
  echo "'$VOL_DIR' found. Cleaning docker volumes."
  $SUDO rm -rf ./volumes/broker
  $SUDO rm -rf ./volumes/zookeeper
else
  echo "Warning: '$VOL_DIR' NOT found. Creating it now."
  $SUDO mkdir ./volumes
fi

echo "Creating docker volumes on disk."
$SUDO mkdir ./volumes/broker
$SUDO mkdir ./volumes/zookeeper

if [ $machine = "Linux" ]; then
  $SUDO chmod 777 ./volumes/broker
  $SUDO chmod 777 ./volumes/zookeeper
fi

echo "Starting up docker containers necessary for FIT."
CONTAINERS=(postgres)
for CON in "${CONTAINERS[@]}"; do
  $SUDO docker-compose --env-file $ENV_FILE -f ${COMPOSE:-docker-compose.yml} up -d ${CON}
done

echo "Finished"
