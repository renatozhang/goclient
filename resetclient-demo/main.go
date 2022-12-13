package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Pod struct {
	Namespace string
	Status    string
	Name      string
	IP        string
}

func main() {
	var kubeconfig *string
	// home是家目录，如果能取得家目录的值，就可以用来做默认值
	if home := homedir.HomeDir(); home != "" {
		// 如果输入了kubeconfig参数，该参数的值就是kubeconfig文件的绝对路径，
		// 如果没有输入kubeconfig参数，就用默认路径~/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig filess")
	} else {
		// 如果取不到当前用户的家目录，就没办法设置kubeconfig的默认目录了，只能从入参中取
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 从本机加载kubeconfig配置文件，因此第一个参数为空字符串
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// kubeconfig加载失败就直接退出了
	if err != nil {
		panic(err.Error())
	}
	// 2.配置api路径 参考path : /api/v1/namespaces/{namespace}/pods
	config.APIPath = "api"
	// 3.配置分组版本 pod的group是空字符串
	config.GroupVersion = &corev1.SchemeGroupVersion
	// 4.配置数据编解码方式
	config.NegotiatedSerializer = scheme.Codecs
	// 5.实例化restClient对象
	restClient, err := rest.RESTClientFor(config)

	if err != nil {
		panic(err.Error())
	}

	// 6.定义RESTClient对象
	result := &corev1.PodList{}

	// 指定namespace
	// namespace := "kube-system"
	namespace := "kube-system"
	// 7.跟APIServer交互
	// 设置请求参数，然后发起请求
	// GET请求
	//  指定namespace，参考path : /api/v1/namespaces/{namespace}/pods
	err = restClient.Get().
		//  指定namespace，参考path : /api/v1/namespaces/{namespace}/pods
		Namespace(namespace).
		// 查找多个pod，参考path : /api/v1/namespaces/{namespace}/pods
		Resource("pods").
		// 指定大小限制和序列化工具
		VersionedParams(&metav1.ListOptions{Limit: 100}, scheme.ParameterCodec).
		// 触发请求
		Do(context.TODO()).
		// 结果存入result
		Into(result)
	if err != nil {
		panic(err.Error())
	}

	/*
		Get,定义请求方式， 返回一个Request 结构体对象，这个request结构体对象就是构建访问APIServer请求的
		依次执行了 Namespace，Resource，VersionedParams，构建与 APIServer交互的参数
		Do方法通过request发起请求，然后通过transformResponse 解析请求返回，并绑定对应资源对象的结构体对象，这里的话，就表示是corev1.PodList{}
		request显示检查了有没有可用的Client，在这里开始调用net/http包的功能
	*/
	// 写入返回结果表头
	fmt.Printf("namespace\t status\t\t name\n")

	// 每个pod都打印namespace、status.Phase、name三个字段
	for _, item := range result.Items {
		fmt.Printf("%v\t %v\t %v\n",
			item.Namespace,
			item.Status.Phase,
			item.Name)
	}
}
