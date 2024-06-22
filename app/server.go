package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// kubeconfig := os.Getenv("KUBECONFIG")
	// if kubeconfig == "" {
	// 	kubeconfig = os.Getenv("HOME") + "/.kube/config"
	// }
	// flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "path to kubeconfig file")
	// flag.Parse()

	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	nodes, _ := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	podLog := make([]*rest.Request, len(nodes.Items))
	log := make([]string, len(nodes.Items))
	podName := make([]string, len(nodes.Items))
	line := int64(2)
	for i := 1; i <= len(nodes.Items); i++ {
		podName[i-1] = "pod" + strconv.Itoa(i)
		podLog[i-1] = clientset.CoreV1().Pods("default").GetLogs(podName[i-1], &corev1.PodLogOptions{
			TailLines: &line,
		})
	}
	e := echo.New()
	e.GET("/metrics", func(c echo.Context) error {
		start := time.Now()
		for i := 0; i < len(nodes.Items); i++ {
			podLogs, err := podLog[i].Stream(context.TODO())
			if err != nil {
				panic(err.Error())
			}
			defer podLogs.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, podLogs)
			if err != nil {
				panic(err.Error())
			}
			str := buf.String()
			lines := strings.Split(str, "\n")
			log[i] = strings.Join(lines, "\n")
		}

		laslog := strings.Join(log, "")
		end := time.Now()
		elapsed := end.Sub(start)
		fmt.Printf("Thời gian thực hiện: %s\n", elapsed)
		return c.String(http.StatusOK, laslog)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
