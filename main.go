package main

import (
	clientset "k8s-crd/pkg/client/clientset/versioned"
	informers "k8s-crd/pkg/client/informers/externalversions"
	"flag"
	"k8s-crd/pkg/signals"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"time"
)

var (
	masterURL string
	kubeconfig string
)

func main()  {
	flag.Parse()
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubecofnig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	networkClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	networkInformerFactory := informers.NewSharedInformerFactory(networkClient, time.Second*30)


	controller := NewController( kubeClient, networkClient,
		networkInformerFactory.Smaplecrd().V1().Networks())

	go networkInformerFactory.Start(stopCh)
	if err = controller.Run(2, stopCh); err != nil{
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init()  {
	flag.StringVar(&kubeconfig, "kubeconfig","","Path to kubeconfig")
	flag.StringVar(&masterURL, "master", "","The address of the kubernetes api server")
}