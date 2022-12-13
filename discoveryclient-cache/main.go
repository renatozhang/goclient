package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. 加载配置文件，生成config对象
	config, err := clientcmd.BuildConfigFromFlags("", "/home/zhangzeng/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// 2. 实例化客户端对象 本地客户端负责将GVR数据缓存到本地文件中
	cacheDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(config, "./cache/discovery", "./cache/http", time.Minute*60)
	if err != nil {
		panic(err.Error())
	}

	// 3.发送请求，获取GVR数据
	_, apiResources, err := cacheDiscoveryClient.ServerGroupsAndResources()
	// 1. 先从缓存文件中找GVR数据，有则直接返返回，没有则需要调用APIServer
	// 2.调用APIServer 获取GVR数据
	// 3. 将获取的GVR数据保存到本地，然后返回给客户端
	if err != nil {
		panic(err.Error())
	}
	for _, list := range apiResources {
		gv, err := schema.ParseGroupVersion(list.APIVersion)
		if err != nil {
			panic(err.Error())
		}

		for _, resource := range list.APIResources {
			fmt.Printf("name:%v, group:%v, version:%v\n", resource.Name, gv.Group, gv.Version)
		}
	}

}
