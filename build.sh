#!/bin/sh

version=1_0_3

cd frontend
yarn build
cd ./build
go-bindata -o=../../src/bindata/bindata.go  -pkg=bindata -fs ./...

cd ../..

mkdir -p build

CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-X 'main.GoVersion=$(go version)' -X main.GitCommit=`git rev-parse HEAD` -X 'main.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`'" ./src/main.go
mv main ./build/gatewayAuth_linux_arm_$version
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.GoVersion=$(go version)' -X main.GitCommit=`git rev-parse HEAD` -X 'main.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`'" ./src/main.go
mv main ./build/gatewayAuth_linux_amd64_$version
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.GoVersion=$(go version)' -X main.GitCommit=`git rev-parse HEAD` -X 'main.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`'" ./src/main.go
mv main ./build/gatewayAuth_darwin_amd64_$version
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.GoVersion=$(go version)' -X main.GitCommit=`git rev-parse HEAD` -X 'main.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`'" ./src/main.go
mv main ./build/gatewayAuth_darwin_arm64_$version
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.GoVersion=$(go version)' -X main.GitCommit=`git rev-parse HEAD` -X 'main.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`'" ./src/main.go
mv main.exe ./build/gatewayAuth_windows_amd64_$version.exe

echo "\n\n======================== build =========================\n\n"

ls ./build/
