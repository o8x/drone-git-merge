#!/usr/bin/env bash

docker build -t alextechs/drone-git-merge:latest .

docker run --rm \
    --env=CI_BUILD_NUMBER="61" \
    --env=CI_BUILD_STARTED="1599572953" \
    --env=CI_BUILD_STATUS="success" \
    --env=CI_COMMIT_AUTHOR_EMAIL="im@println.org" \
    --env=CI_COMMIT_AUTHOR_NAME="Alex" \
    --env=CI_COMMIT_BRANCH="branch" \
    --env=CI_COMMIT_MESSAGE="1" \
    --env=DRONE_REPO_LINK="https://host/" \
    --env=PLUGIN_PROJECTS_PATH='/app,/go-app' \
    --env=PLUGIN_SERVER='{"host":"host","port" : "port","user":"user","password":"password"}' \
    --env=PLUGIN_SOURCE_BRANCHS="master,release,physical-server" \
    --env=PLUGIN_TARGET_BRANCH="develops" \
    alextechs/drone-git-merge:latest
