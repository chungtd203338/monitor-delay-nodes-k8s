#!/bin/bash

# Lấy danh sách các node từ Kubernetes
nodes=$(kubectl get nodes -o jsonpath='{.items[*].metadata.name}')
node_array=($nodes)

# Khởi tạo biến đếm
counter=1

for node in "${node_array[@]}"
do
  label="node$counter"
  echo "Labeling $node with pod-assign=$label"
  kubectl label nodes $node pod-assign=$label --overwrite
  counter=$((counter+1))
done

# check label node
kubectl get nodes --show-labels=true