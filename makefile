APP_NAME = ampctl

.PHONY: build clean

build:
	go build -o build/$(APP_NAME)

clean:
	rm -f build/$(APP_NAME)
