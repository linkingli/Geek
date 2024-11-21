package handlers

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type PodHandler struct {
}

func (h *PodHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("PodHandler OnAdd: ", obj.(*v1.Pod).Name)
}

func (h *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("PodHandler OnUpdate: ", newObj.(*v1.Pod).Name)
}

func (h *PodHandler) OnDelete(obj interface{}) {
	fmt.Println("PodHandler OnDelete: ", obj.(*v1.Pod).Name)
}

type NewPodHandler struct {
}

func (h *NewPodHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("NewPodHandler OnAdd: ", obj.(*v1.Pod).Name)
}

func (h *NewPodHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("PodHandler OnUpdate: ", newObj.(*v1.Pod).Name)
}

func (h *NewPodHandler) OnDelete(obj interface{}) {
	fmt.Println("PodHandler OnDelete: ", obj.(*v1.Pod).Name)
}
