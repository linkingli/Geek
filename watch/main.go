package main

import (
	"fmt"

	"github.com/xingyunyang01/watch/pkg/config"
	"github.com/xingyunyang01/watch/pkg/handlers"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func informer(lw *cache.ListWatch) {
	options := cache.InformerOptions{
		ListerWatcher: lw,
		ObjectType:    &v1.Pod{},
		ResyncPeriod:  0,
		Handler:       &handlers.PodHandler{},
	}

	_, informer := cache.NewInformerWithOptions(options)

	informer.Run(wait.NeverStop)
}

func sharedInformer(lw *cache.ListWatch) {
	sharedInformer := cache.NewSharedInformer(lw, &v1.Pod{}, 0)
	sharedInformer.AddEventHandler(&handlers.PodHandler{})
	sharedInformer.AddEventHandler(&handlers.NewPodHandler{})
	sharedInformer.Run(wait.NeverStop)
}

func sharedInformerFactory(client *kubernetes.Clientset, lw *cache.ListWatch) {
	fact := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))

	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&handlers.PodHandler{})

	svcInformer := fact.Core().V1().Services()
	svcInformer.Informer().AddEventHandler(&handlers.ServiceHandler{})

	fact.Start(wait.NeverStop)
}

func sharedInformerFactoryLister(client *kubernetes.Clientset, lw *cache.ListWatch) {
	fact := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))

	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&handlers.PodHandler{})

	ch := make(chan struct{})
	fact.Start(ch)
	fact.WaitForCacheSync(ch)

	podlist, _ := podInformer.Lister().List(labels.Everything())

	fmt.Println(podlist)
}

func sharedInformerFactoryForResource(client *kubernetes.Clientset, lw *cache.ListWatch) {
	fact := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("default"))

	grv := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	informer, err := fact.ForResource(grv)
	if err != nil {
		panic(err)
	}
	informer.Informer().AddEventHandler(&cache.ResourceEventHandlerFuncs{})

	ch := make(chan struct{})
	fact.Start(ch)
	fact.WaitForCacheSync(ch)

	list, _ := informer.Lister().List(labels.Everything())

	fmt.Println(list)
}

func main() {
	client := config.NewK8sConfig().InitRestConfig().InitClientSet()

	lw := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", "default", fields.Everything())

	//informer(lw)

	//sharedInformer(lw)

	//sharedInformerFactory(client, lw)

	//sharedInformerFactoryLister(client, lw)

	sharedInformerFactoryForResource(client, lw)

	select {}
}
