#!/bin/bash

BUILDS=( \
    'linux;arm' \
    'linux;amd64' \
    'windows;amd64' \
    'darwin;amd64' \
)

BUILDNAME="slms"

DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

if [ ! -d builds ]; then
    mkdir builds
fi

FILELIST="assets/"

for BUILD in ${BUILDS[*]}; do

    IFS=';' read -ra SPLIT <<< "$BUILD"
    OS=${SPLIT[0]}
    ARCH=${SPLIT[1]}

    echo "Building ${OS}_$ARCH..."
    (env GOOS=$OS GOARCH=$ARCH \
        go build -o builds/${BUILDNAME}_${OS}_${ARCH} \
        -ldflags "-X main.appVersion=$TAG -X main.appDaten=$DATE -X main.appCommit=$COMMIT")

    if [ "$OS" = "windows" ]; then
        mv builds/${BUILDNAME}_windows_${ARCH} builds/${BUILDNAME}_windows_${ARCH}.exe
        FILELIST="$FILELIST builds/${BUILDNAME}_windows_${ARCH}.exe"
    else
        FILELIST="$FILELIST builds/${BUILDNAME}_${OS}_${ARCH}"
    fi

done

tar -czvf builds_and_assets.tar.gz $FILELIST

wait