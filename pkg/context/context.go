/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package context

import (
	informerv1 "k8s.io/client-go/informers/core/v1"
	informerv1beta1 "k8s.io/client-go/informers/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

// ControllerContext holds
type ControllerContext struct {
	IngressInformer  cache.SharedIndexInformer
	ServiceInformer  cache.SharedIndexInformer
	PodInformer      cache.SharedIndexInformer
	NodeInformer     cache.SharedIndexInformer
	EndpointInformer cache.SharedIndexInformer
	// Stop is the stop channel shared among controllers
	StopCh chan struct{}
}

func NewControllerContext(kubeClient kubernetes.Interface, namespace string, resyncPeriod time.Duration, enableEndpointsInformer bool) *ControllerContext {
	context := &ControllerContext{
		IngressInformer: informerv1beta1.NewIngressInformer(kubeClient, namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}),
		ServiceInformer: informerv1.NewServiceInformer(kubeClient, namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}),
		PodInformer:     informerv1.NewPodInformer(kubeClient, namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}),
		NodeInformer:    informerv1.NewNodeInformer(kubeClient, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}),
		StopCh:          make(chan struct{}),
	}
	if enableEndpointsInformer {
		context.EndpointInformer = informerv1.NewEndpointsInformer(kubeClient, namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	}
	return context
}

func (ctx *ControllerContext) Start() {
	go ctx.IngressInformer.Run(ctx.StopCh)
	go ctx.ServiceInformer.Run(ctx.StopCh)
	go ctx.PodInformer.Run(ctx.StopCh)
	go ctx.NodeInformer.Run(ctx.StopCh)
	if ctx.EndpointInformer != nil {
		go ctx.EndpointInformer.Run(ctx.StopCh)
	}
}
