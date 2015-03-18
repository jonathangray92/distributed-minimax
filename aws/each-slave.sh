#!/bin/bash

THIS_DIR="$(dirname $0)"
SLAVE_IP_FILE="$THIS_DIR/slave-ip"
KEYPAIR_FILE="$THIS_DIR/master-keypair.pem"
SSH_ARGS="-i $KEYPAIR_FILE -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -q"

if [[ ! -r $SLAVE_IP_FILE ]]; then
    echo not found: $SLAVE_IP_FILE >&2
    exit 1
fi

if [[ ! -r $KEYPAIR_FILE ]]; then
    echo not found: $KEYPAIR_FILE >&2
    exit 1
fi

if [[ $# != 1 ]]; then
    echo usage: $0 [command] >&2
    exit 1
fi

while read -u10 ip; do
    ssh $SSH_ARGS ec2-user@$ip "$1"
done 10< $SLAVE_IP_FILE
