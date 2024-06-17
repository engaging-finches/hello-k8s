# Ubuntu base image
FROM ubuntu:24.04

# update and install stuff 
RUN apt-get update && \
    apt-get upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y curl jq sudo bash

# creates user GHA with no password and adds to sudoers
RUN useradd -m GHA && \
    mkdir -p /etc/sudoers.d && \
    echo "GHA ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/GHA && \
    chmod 0440 /etc/sudoers.d/GHA

# switch to user GHA and set the working directory
USER GHA
WORKDIR /home/GHA/actions-runner
    
# create a directory for the runner and install
RUN sudo mkdir -p /home/GHA/actions-runner && \
    cd /home/GHA/actions-runner && \
    sudo curl -O -L https://github.com/actions/runner/releases/download/v2.312.0/actions-runner-linux-arm64-2.312.0.tar.gz && \
    sudo tar xzf ./actions-runner-linux-arm64-2.312.0.tar.gz && \
    sudo rm -f ./actions-runner-linux-arm64-2.312.0.tar.gz && \
    sudo ./bin/installdependencies.sh  

# copy any necessary files
COPY ./start.sh .

# make start.sh executable and change ownership
RUN sudo chmod +x /home/GHA/actions-runner/start.sh && \
    sudo chown -R GHA:GHA /home/GHA/actions-runner 
    
# start the runner
# ENTRYPOINT ["/bin/bash", "-c", "while true; do sleep 1000; done"]
ENTRYPOINT ["./start.sh"]