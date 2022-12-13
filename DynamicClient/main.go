package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. 加载配置文件，生成config对象
	config, err := clientcmd.BuildConfigFromFlags("", "/home/zhangzeng/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// 2. 实例化客户端对象，这里是实例化 动态客户端对象
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 3. 配置我们需要的GVR
	gvr := schema.GroupVersionResource{
		Group:    "", //不需要写的， 应为时无名资源组，也就是core资源组
		Version:  "/v1",
		Resource: "pods",
	}

	// 发送请求且得到返回结果
	unStructData, err := dynamicClient.Resource(gvr).Namespace("kube-system").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 5. 转化为结构化的数据
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructData.UnstructuredContent(), podList)
	if err != nil {
		panic(err.Error())
	}

	// Resource 基于gvr生成了一个针对于资源的客户端
	// Namespace 指定一个可操作的命名空间 同时它是dynamicResouceClient的方法
	// List 首先通过RESTClient 调用APIServer的接口返回了Pod的数据，返回的数据格式时二进制的Json格式，然后通过一系列的解析方法，转换成unstructured.UnstructuredList

	for _, item := range podList.Items {
		fmt.Printf("nameapace: %v, name: %v\n", item.Namespace, item.Name)
	}
}
