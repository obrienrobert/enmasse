/*
 * Copyright 2018-2019, EnMasse authors.
 * License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/enmasseproject/enmasse/pkg/logs"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/openshift/api"

	enmassescheme "github.com/enmasseproject/enmasse/pkg/client/clientset/versioned/scheme"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/enmasseproject/enmasse/pkg/cache"
	"github.com/enmasseproject/enmasse/pkg/controller"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var log = logf.Log.WithName("cmd")

func main() {
	flag.Parse()

	logs.InitLog()
	logs.PrintVersions(log)

	namespace := os.Getenv("NAMESPACE")

	log.Info("Watching on namespace", "namespace", namespace)

	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Failed to get configuration")
		os.Exit(1)
	}

	ctx := context.TODO()
	serverClient, err := client.New(cfg, client.Options{})

	monitoring := os.Getenv("ENABLE_MONITORING")
	if monitoring == "true" {
		err = installMonitoring(ctx, serverClient)
		if err != nil {
			fmt.Print(err)
		}
	}

	mgr, err := manager.New(cfg, manager.Options{
		Namespace: namespace,
		NewCache: cache.NewDelegateCacheBuilder(namespace,
			schema.GroupVersionKind{
				Group:   "user.enmasse.io",
				Version: "v1beta1",
				Kind:    "MessagingUser",
			},
			schema.GroupVersionKind{
				Group:   "user.enmasse.io",
				Version: "v1beta1",
				Kind:    "MessagingUserList",
			},
			schema.GroupVersionKind{
				Group:   "enmasse.io",
				Version: "v1beta1",
				Kind:    "AddressSpace",
			},
			schema.GroupVersionKind{
				Group:   "enmasse.io",
				Version: "v1beta1",
				Kind:    "AddressSpaceList",
			},
			schema.GroupVersionKind{
				Group:   "enmasse.io",
				Version: "v1beta1",
				Kind:    "Address",
			},
			schema.GroupVersionKind{
				Group:   "enmasse.io",
				Version: "v1beta1",
				Kind:    "AddressList",
			},
		),
	})
	if err != nil {
		log.Error(err, "Failed to create manager")
		os.Exit(1)
	}

	log.Info("Registering components...")

	// register APIs

	if err := api.Install(scheme.Scheme); err != nil {
		log.Error(err, "Failed to register OpenShift schema")
		os.Exit(1)
	}

	if err := enmassescheme.AddToScheme(scheme.Scheme); err != nil {
		log.Error(err, "Failed to register EnMasse schema")
		os.Exit(1)
	}

	if err := controller.CheckUpgrade(mgr); err != nil {
		log.Error(err, "Failed to upgrade")
		os.Exit(1)
	}

	// register controller

	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "Failed to register controller")
		os.Exit(1)
	}

	// starting ...

	log.Info("Starting the operator")

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "manager exited non-zero")
		os.Exit(1)
	}
}

func installMonitoring(ctx context.Context, client client.Client) error {
	params := map[string]string{}
    params["Namespace"] = os.Getenv("NAMESPACE")

	templateHelper := NewTemplateHelper(params)

	for _, template := range templateHelper.TemplateList {
		resource, err := templateHelper.CreateResource(template)
		if err != nil {
			return fmt.Errorf("createResource failed: %w", err)
		}
		err = client.Create(ctx, resource)
		if err != nil {
			if !kerrors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating resource: %w", err)
			}
		}
	}

	return nil
}
