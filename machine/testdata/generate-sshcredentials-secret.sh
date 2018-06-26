#!/usr/bin/env bash
username=$1
private_key_file=$2

cat <<EOF
apiVersion: v1
kind: Secret
data:
  username: $(echo -n "$username" | base64 | tr -d '\n')
  privateKey: $(cat "$private_key_file" | base64 | tr -d '\n')
EOF
