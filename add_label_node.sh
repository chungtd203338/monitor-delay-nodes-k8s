kubectl label node master pod-assign=node1
kubectl label node worker1 pod-assign=node2
kubectl label node worker2 pod-assign=node3

# check label node
kubectl get nodes --show-labels=true