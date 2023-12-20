BuildVersion = v1
Registry = docker.io
LDFlags = "-X 'main.BuildVersion=$(BuildVersion)'"
Image = $(Registry)/hunter2019/sota:$(BuildVersion)

all:
	go build -ldflags $(LDFlags) -o sota-mesh ./cmd/sota/*.go

build:
	docker build --build-arg LDFLAGS=$(LDFlags) -t $(Image) .

sample:
	go build -o ./cmd/sample/httpserver ./cmd/sample/httpserver/*.go
	#go build -o ./cmd/sample/zmqclient ./cmd/sample/zmqclient/*.go

push:
	docker push $(Image)
