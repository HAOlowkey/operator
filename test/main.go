package main

import (
	"context"
	"fmt"
	clientset "github.com/HAOlowkey/operator/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	unitClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	unitList, err := unitClient.ApplicationV1().Units("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, unit := range unitList.Items {
		fmt.Println(unit.Name, unit.Spec.IpAddr, unit.OwnerReferences, unit.Spec.SideCarContainer[0].Name, unit.Status)
	}
}
