#!/bin/bash

# if directory community-operators does not exist, checkout from git@github.com:k8s-operatorhub/community-operators.git 
if [ ! -d "community-operators" ]; then
    git clone git@github.com:k8s-operatorhub/community-operators.git
fi

# ask if we want to start a new branch from main
read -p "Do you want to start a new branch from main? (y/n): " START_NEW_BRANCH
if [ "$START_NEW_BRANCH" == "y" ]; then
    git fetch --all
    read -p "Enter the name of the new branch: " NEW_BRANCH
    git checkout remotes/origin/main -B $NEW_BRANCH
fi

cd community-operators

VERSION=$(grep "^  version:" ../../bundle/manifests/jumpstarter-operator.clusterserviceversion.yaml | awk '{print $2}')

echo "Updating community-operators to version ${VERSION}"

# make sure that the operators/jumpstarter-operator/${VERSION} directory exists
mkdir -p operators/jumpstarter-operator/${VERSION}

cp -v -r -f ../../bundle/* operators/jumpstarter-operator/${VERSION}

echo You can now review the changes and commit them.
