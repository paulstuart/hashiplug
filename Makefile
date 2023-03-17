KV_PLUGIN ?= "./kv-go-grpc"
KV_PROTO ?= "grpc"

KEY ?= "hello"
WORD ?= "everybody"

.phony: put get gen all clean

all: kv kv-plugin

clean:
	rm -f kv kv-plugin
	rm -f proto/*.go

# This builds the main CLI
kv:
	go build -o kv

kv-plugin:
	go build -o kv-plugin ./plugin-go

put:
	KV_PROTO=${KV_PROTO} ./kv put ${KEY} "${WORD}"

get:
	KV_PROTO=${KV_PROTO} ./kv get ${KEY}

gen:
	./regen.sh
