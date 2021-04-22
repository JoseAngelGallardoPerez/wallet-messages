print("Wallet Messages")

load("ext://restart_process", "docker_build_with_restart")

cfg = read_yaml(
    "tilt.yaml",
    default = read_yaml("tilt.yaml.sample"),
)

local_resource(
    "messages-build-binary",
    "make fast_build",
    deps = ["./cmd", "./internal", "./rpc/pb_server.go", "./rpc/providers.go"],
)
local_resource(
    "messages-generate-protpbuf",
    "make gen-protobuf",
    deps = ["./rpc/messages/messages.proto"],
)

docker_build(
    "velmie/wallet-messages-db-migration",
    ".",
    dockerfile = "Dockerfile.migrations",
    only = "migrations",
)
k8s_resource(
    "wallet-messages-db-migration",
    trigger_mode = TRIGGER_MODE_MANUAL,
    resource_deps = ["wallet-messages-db-init"],
)

wallet_messages_options = dict(
    entrypoint = "/app/service_messages",
    dockerfile = "Dockerfile.prebuild",
    port_forwards = [],
    helm_set = [],
)

if cfg["debug"]:
    wallet_messages_options["entrypoint"] = "$GOPATH/bin/dlv --continue --listen :%s --accept-multiclient --api-version=2 --headless=true exec /app/service_messages" % cfg["debug_port"]
    wallet_messages_options["dockerfile"] = "Dockerfile.debug"
    wallet_messages_options["port_forwards"] = cfg["debug_port"]
    wallet_messages_options["helm_set"] = ["containerLivenessProbe.enabled=false", "containerPorts[0].containerPort=%s" % cfg["debug_port"]]

docker_build_with_restart(
    "velmie/wallet-messages",
    ".",
    dockerfile = wallet_messages_options["dockerfile"],
    entrypoint = wallet_messages_options["entrypoint"],
    only = [
        "./build",
        "zoneinfo.zip",
    ],
    live_update = [
        sync("./build", "/app/"),
    ],
)
k8s_resource(
    "wallet-messages",
    resource_deps = ["wallet-messages-db-migration"],
    port_forwards = wallet_messages_options["port_forwards"],
)

yaml = helm(
    "./helm/wallet-messages",
    # The release name, equivalent to helm --name
    name = "wallet-messages",
    # The values file to substitute into the chart.
    values = ["./helm/values-dev.yaml"],
    set = wallet_messages_options["helm_set"],
)

k8s_yaml(yaml)
