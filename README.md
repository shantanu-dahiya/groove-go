# Groove-Go

Go implementation of the metadata-private messaging system described in the [Groove paper](https://people.csail.mit.edu/nickolai/papers/barman-groove.pdf).

# How to Run

1. Make the client and server scripts executable.
    
     `chmod +x run_clients.sh run_servers.sh`
2. Run the server script followed by the client script.

    `./run_servers.sh`
    
    `./run_clients.sh`

Client 2 is currently hard-coded to send a request to Client 1 to add each other as buddies. Thereafter both clients arrive at a shared deadrop derived from their symmetric key and setup onion circuits to that dead drop.

# Notes
There is a [dockerfile](.devcontainer/Dockerfile) included in the project which can spin up a container with the requirements. If not using docker, make sure that `go` is installed.

[global.json](pkg/global/global.json) contains information about the servers and clients spun up by the system.
