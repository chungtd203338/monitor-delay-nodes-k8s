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
***Lưu ý***
Metrics có dạng
```
node1node2 51.912
node2node3 52.104
.................
```
node1node2 là độ trễn từ node1 -> node2 (node1 node2 là label gắn ở các bước đầu),
51.912 đơn vị là ms

Truy cập vào Grafana bằng Service Grafana thực hiện config đến prometheus và config dashboard mong muốn
```
NAME                 TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
grafana              NodePort   10.107.82.54    <none>        3000:30226/TCP   109m
```
