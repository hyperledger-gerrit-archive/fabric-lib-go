#!/bin/bash -e
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

Parse_Arguments() {
    while [ $# -gt 0 ]; do
          case $1 in
              --env_Info)
                 env_Info
                 ;;
              --clean_Environment)
                 clean_Environment
                 ;;
          esac
          shift
    done
}

clean_Environment() {

  echo "-----------> Clean Docker Containers & Images, unused/lefover build artifacts"
  clearContainers () {
        CONTAINER_IDS=$(docker ps -aq)
        if [ -z "$CONTAINER_IDS" ] || [ "$CONTAINER_IDS" = " " ]; then
                echo "---- No containers available for deletion ----"
        else
                docker rm -f $CONTAINER_IDS || true
        fi
  }

  removeUnwantedImages() {
        # Delete <none> images
        docker images | grep none | awk '{ print $3; }' | xargs docker rmi || true
        # Get the latest baseimage version from Makefile of fabric master branch
        curl -L https://raw.githubusercontent.com/hyperledger/fabric/master/Makefile > Makefile
        # Fetch baseimage release version
        BASE_IMAGE=$(cat Makefile | grep "BASEIMAGE_RELEASE =" | awk '{print $3}')
        # Deleete Makefile
        rm -rf Makefile
        # Delete all docker images except the latest one fetched from fabric master Makefile
        IMAGE_IDS=$(docker images | grep -v "$BASE_IMAGE" | awk 'NR>1 {print $3}')
        if [[ -z ${IMAGE_IDS// } ]]; then
             echo "---- No Images available for deletion ----"
        else
             # Delete all list docker images
             docker rmi -f $IMAGE_IDS
             echo -e "\033[32m Docker Images List \033[0m"
             docker images
        fi
  }

  # Delete nvm prefix & then delete nvm
  rm -rf $HOME/.node-gyp/ $HOME/.npm/ $HOME/.npmrc  || true

  # remove tmp/hfc and hfc-key-store data
  rm -rf /home/jenkins/npm /tmp/fabric-shim /tmp/hfc* /tmp/npm* /home/jenkins/kvsTemp /home/jenkins/.hfc-key-store

  rm -rf /var/hyperledger/*

clearContainers
removeUnwantedImages

}

env_Info() {
        # This function prints system info

        echo -e "\033[32m -----------> Build Env INFO" "\033[0m"
        # Output all information about the Jenkins environment
        uname -a
        env
        gcc --version
        docker version
        docker info
        docker-compose version
}
