# Trịnh Đức Chung - K65 - HUST <3
# Hệ thống giám sát độ trễ giữa các node trong K8s #
**Phương pháp sử dụng: Triên khai mỗi node 1 pod xong ping lẫn nhau, thu thập độ trễ từ các pod ở các node khác nhau (Gần đúng so với độ trễ đứng từ các node ping nhau @@@)**
# Cách sử dụng #
Clone project về máy
```console
git clone https://github.com/chungtd203338/monitor-delay-nodes-k8s
cd monitor-delay-nodes-k8s
```
Đánh label để phân biệt và triển khai pod giữa các node trong cụm K8s (config sao cho phù hợp với cụm k8s)
```console
chmod +x add_label_node.sh
./add_label_node.sh
```
Xác định số node trong cụm k8s và thời gian giữa các lần đo, config vào file render_yaml.sh
```console
NODES=3
SLEEP=5
```
Thực thi file render_yaml.sh để render ra file pods.yaml
```console
chmod +x render_yaml.sh
./render_yaml.sh
```
Triển khai các pod lên các node
```console
kubectl apply -f pods.yaml
```
```
NAME                            READY   STATUS    RESTARTS   AGE    IP              NODE      NOMINATED NODE   READINESS GATES
pod1                            1/1     Running   0          122m   10.10.219.107   master    <none>           <none>
pod2                            1/1     Running   0          122m   10.10.235.157   worker1   <none>           <none>
pod3                            1/1     Running   0          122m   10.10.189.87    worker2   <none>           <none>
```
**Triển khai metrics-app để thu thập metrics**
```console
cd app
chmod +x build.sh
./build.sh build
docker images
```
```
REPOSITORY                   TAG       IMAGE ID       CREATED       SIZE
chung123abc/metrics-server   v2.0      0a50d56b7b51   2 hours ago   1.01GB
```
Có thể config tên images và push lên dockerhub trong file build.sh

Triển khai metrics-app lên cụm k8s (thay đúng tên image đã build trong file metrics-server.yaml)
```console
cd ..
cd manifest
kubectl apply -f metrics-server.yaml
```

Triển khai Prometheus để collect metrics của metrics-app và Grafana để hiện thị metrics
```console
kubectl apply -f prometheus.yaml
kubectl apply -f grafana.yaml
```

Thực hiện config prometheus trỏ đúng vào ip và pod của metrics-app ```metrics-app(x.x.x.x):1323``` trong configmap in file prometheus.yaml
```
NAME                            READY   STATUS    RESTARTS   AGE    IP              NODE      NOMINATED NODE   READINESS GATES
metrics-server                  1/1     Running   0          170m   10.10.235.151   worker1   <none>           <none>

```
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-server-conf
  namespace: prometheus
data:
  prometheus.yml: |
    global:
      scrape_interval: 5s
    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
          - targets: ['metrics-app(x.x.x.x):1323']
```

Truy cập vào Prometheus bằng Service Prometheus thực hiện truy vấn metrics
```
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
prometheus-service   NodePort   10.100.89.224   <none>        80:30165/TCP     130m
```
**Lưu ý** </br>
Metrics có dạng
```
delay{from=node1,to=node2} 52.351
delay{from=node1,to=node3} 52.069
delay{from=node2,to=node3} 2.117
.................
```
label chỉ rõ độ trễ giữa 2 node (ví dụ node1 node2 là label gắn ở các bước đầu) </br>
51.912 đơn vị là ms

Truy cập vào Grafana bằng Service Grafana thực hiện config đến prometheus và config dashboard mong muốn
```
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
grafana              NodePort   10.107.82.54    <none>        3000:30226/TCP   109m
```
