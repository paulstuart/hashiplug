KV_PLUGIN ?= "./kv-go-grpc"
KV_PROTO ?= "grpc"

KEY ?= "hello"
WORD ?= "everybody"

.phony: put get gen all clean

all: kv kv-plugin kv-go-grpc

clean:
	rm -f kv kv-plugin kv-go-grpc

# This builds the main CLI
kv:
	go build -o kv

kv-plugin:
	go build -o kv-plugin ./plugin-go

# This builds the plugin written in Go
kv-go-grpc:
	go build -o kv-go-grpc ./plugin-go-grpc

# This tells the KV binary to use the "kv-go-grpc" binary
put:
	@#KV_PLUGIN="${KV_PLUGIN}" KV_PROTO=${KV_PROTO} ./kv put ${KEY} ${WORD}
	KV_PROTO=${KV_PROTO} ./kv put ${KEY} "${WORD}"

get:
	@#KV_PLUGIN="${KV_PLUGIN}" KV_PROTO=${KV_PROTO} ./kv get ${KEY}
	KV_PROTO=${KV_PROTO} ./kv get ${KEY}

gen:
	./regen.sh
