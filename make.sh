#!/bin/bash -e

function latest_version()
{
    aws s3 ls s3://nxcloud-devtools/linux/config_helper/  | grep PRE | grep -v latest | sed 's%.*PRE \([^/]*\).*$%\1%' | sort -V | tail -1
}

function _increment_build_number()
{
    local version=$1

    echo $version | awk -F. '{ split($0, a); printf("%d.%d.%d", a[1], a[2], a[3] + 1)}'
}

function build()
{
    local goos=$1
    local version=$2

    [ -z "$goos" ] && goos=$(uname -s | tr A-Z a-z)
    [ -z "$version" ] && version=dev

    mkdir -p bin/$goos

    local binary=$(_binary_name $goos)

    mkdir -p bin/$goos/{amd64,arm64}
    GOOS=$goos GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o bin/$goos/amd64/$binary
    GOOS=$goos GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o bin/$goos/arm64/$binary
}

function _goos_to_os()
{
    local goos=$1

    [ "$goos" = "darwin" ] && { echo mac; return; }

    echo $goos
}

function _binary_name()
{
    local goos=$1

    [ "$goos" = "windows" ] && { echo config_helper.exe; return; }

    echo config_helper
}

function publish()
{
    local gooses="darwin linux windows"
    local version=$(_increment_build_number $(latest_version))

    for goos in $gooses
    do
        local binary=$(_binary_name $goos)

        build $goos $version
    done

    #TODO Add unit tests here

    for goos in $gooses
    do
        local os=$(_goos_to_os $goos)
        local binary=$(_binary_name $goos)

        aws s3 cp --acl public-read bin/$goos/amd64/$binary s3://nxcloud-devtools/$os/config_helper/$version/amd64/$binary
        aws s3 cp --acl public-read bin/$goos/arm64/$binary s3://nxcloud-devtools/$os/config_helper/$version/arm64/$binary
    done
}

numargs=$#
for ((n=1;n <= numargs; n++))
do
    func=$1; shift
    $func
done
