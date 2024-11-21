package handlers

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type ServiceHandler struct {
}

func (h *ServiceHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("ServiceHandler OnAdd: ", obj.(*v1.Service).Name)
}

func (h *ServiceHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("ServiceHandler OnUpdate: ", newObj.(*v1.Service).Name)
}

func (h *ServiceHandler) OnDelete(obj interface{}) {
	fmt.Println("ServiceHandler OnDelete: ", obj.(*v1.Service).Name)
}
