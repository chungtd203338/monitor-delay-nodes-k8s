#!/bin/bash

# Số lượng node trong cụm k8s
NODES=3
SLEEP=5

# Khởi tạo file YAML rỗng
> pods.yaml

# Tạo các pod dựa trên số lượng node
for i in $(seq 1 $NODES); do
  echo "apiVersion: v1
kind: Pod
metadata:
  name: pod$i
spec:
  containers:
  - name: ping-container
    image: portainer/kubectl-shell:latest
    securityContext:
      runAsUser: 0
      runAsGroup: 0
    command:
    - /bin/sh
    - -c
    - |
      sleep 3
      while true; do
        $(for j in $(seq 1 $NODES); do if [ $i -ne $j ]; then echo "        "IP$j='$(kubectl get pod' pod$j '-o jsonpath='{.status.podIP}')'" && echo \"delay{from="node$i",to="node$j"} \$(ping -c 1 "\$IP$j" | grep 'time=' | awk -F'time=' '{print \$2}' | awk '{print \$1}')\"; "; fi; done)
        sleep $SLEEP;
      done
  nodeSelector:
    pod-assign: node$i
---" >> pods.yaml
done