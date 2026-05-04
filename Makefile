# ====== CONFIG ======
APP_NAME=gopher-gateway

# ====== APP ======

run:
	go run cmd/server/main.go

build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

# ====== CLEAN ======

clean:
	rm -rf bin/