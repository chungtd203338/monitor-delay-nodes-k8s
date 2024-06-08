#!/bin/bash
HUB="chung123abc"
TAG="v2.0"

GREEN="\e[32m"
OPTION=$1
logSuccess() { echo -e "$GREEN-----$message-----";}

buildImage() {
    image=$(docker images | grep metrics-server | awk '{print $1}'):$TAG
    docker rmi $image
    message="Remove old image success" && logSuccess
    docker build -t $HUB/metrics-server:$TAG .
    message="Build Success" && logSuccess
}

pushImage() {
    image=$(docker images | grep metrics-server | awk '{print $1}'):$TAG
    docker push $image
    message="Push Success" && logSuccess
}

if [ $OPTION == "build" ]; then
    buildImage
elif [ $OPTION == "push" ]; then
    pushImage
elif [ $OPTION == "ful" ]; then
    buildImage
    pushImage
fi