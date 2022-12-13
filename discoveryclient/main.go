package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. 加载配置文件，生成config对象
	config, err := clientcmd.BuildConfigFromFlags("", "/home/zhangzeng/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// 2. 实例化客户端对象，这里是实例化 discorvery客户端对象
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 3.发送请求，获取GVR数据
	_, apiResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err.Error())
	}

	// ServerGroups负责获取GV数据，然后调用fetchGroupVersionResources 且这给这个方法传递GV参数，
	// 然后通过调用ServerResourcesForGroupVersion（resetClient）方法获取GV对应的 Resource数据，也就是资源数据
	// 同时返回一个map[gv]resourceList的数据格式，最后处理map > slice 然后返回GVR slice
	for _, list := range apiResources {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			panic(err.Error())
		}
		for _, resource := range list.APIResources {
			fmt.Printf("name: %v, gorup: %v,version: %v\n", resource.Name, gv.Group, gv.Version)
		}
	}
}
