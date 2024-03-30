PROJECT_NAME=college
MAIN_PATH=cmd/imi/college/main.go

build:
	GOARCH=amd64 GOOS=linux   go build -o .output/${PROJECT_NAME}-amd64     -ldflags "-s -w" ${MAIN_PATH}
	GOARCH=arm64 GOOS=linux   go build -o .output/${PROJECT_NAME}-arm64     -ldflags "-s -w" ${MAIN_PATH}
	GOARCH=amd64 GOOS=windows go build -o .output/${PROJECT_NAME}-amd64.exe -ldflags "-s -w" ${MAIN_PATH}
	GOARCH=arm64 GOOS=windows go build -o .output/${PROJECT_NAME}-arm64.exe -ldflags "-s -w" ${MAIN_PATH}

