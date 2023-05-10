package main

import (
	"context"
	"encoding/json"
	"net/http"
	"log"
	"flag"
	"fmt"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Request struct {
    Name string `json:"name"`
}

type Response struct {
    Greeting string `json:"greeting"`
}
type K8Resource struct {
	Pods     int `json:"pods"`
	Services int `json:"services"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
        res := Response{
            Greeting: "Hello backend server is running ...!",
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(res)
	 }

func main()  {

    router := mux.NewRouter()

	configPath := flag.String("configPath", "/Users/deepaknishad/.kube/config", "Path to a kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *configPath)
	if err != nil {
		fmt.Printf("error %s build config from flags\n", err.Error())

		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s failed in fetching incluster config\n", err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error %s in creating clientset\n", err.Error())
	}


    router.HandleFunc("/hello", HealthCheckHandler).Methods("GET")
	router.HandleFunc("/getPods", func (w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()	

	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil{
		fmt.Printf("error %s while listing all the pods from default namespace\n", err.Error())
	}
		w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(pods)

}).Methods("GET")


	router.HandleFunc("/getServices", func (w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()	

	services, err := clientset.CoreV1().Services("default").List(ctx, metav1.ListOptions{})

	if err != nil {
        fmt.Printf("error %s while listing all the services from default namespace\n", err.Error())
        http.Error(w, "Failed to list services", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(services)
    if err != nil {
        fmt.Printf("error %s while encoding response\n", err.Error())
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
    defer w.(http.Flusher).Flush()

}).Methods("GET")

    log.Fatal(http.ListenAndServe(":8080", router))
	
}