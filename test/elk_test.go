// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build elk

package test

import (
	"log"
	"testing"
	"time"

	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/cluster_manager"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
)

var (
	clientConfig = &ClientConfig{
		Host:  "192.168.0.7:9100",
		Debug: true,
	}

	RepoNameForTest    = "incubator"
	RuntimeNameForTest = "qingcloud runtime"
	RepoUrlForTest     = "http://139.198.121.182:8879/"

	ClusterConf = `{
    "cluster": {
        "name": "ELK",
        "description": "The description of the ELK service",
        "subnet": "vxnet-dkdu5u0",
        "es_node": {
            "cpu": 2,
            "memory": 4096,
            "count": 3,
            "instance_class": 1,
            "volume_size": 10
        },
        "kbn_node": {
            "cpu": 2,
            "memory": 4096,
            "count": 1,
            "instance_class": 1,
            "volume_size": 10
        },
        "lst_node": {
            "cpu": 2,
            "memory": 4096,
            "count": 1,
            "instance_class": 1,
            "volume_size": 10
        }
    }
}`
)

func TestK8s(t *testing.T) {
	log.SetPrefix("[ === ELK TEST === ] ")

	client := GetClient(clientConfig)

	// create repo
	var repoID string
	{
		describeParams := repo_manager.NewDescribeReposParams()
		describeParams.SetName([]string{RepoNameForTest})
		describeResp, err := client.RepoManager.DescribeRepos(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		repos := describeResp.Payload.RepoSet

		if len(repos) != 0 {
			repoID = repos[0].RepoID
		} else {
			createParams := repo_manager.NewCreateRepoParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRepoRequest{
					Name:        RepoNameForTest,
					Description: "incubator charts",
					Type:        "http",
					URL:         RepoUrlForTest,
					Credential:  `{}`,
					Visibility:  "public",
					Providers:   []string{"qingcloud"},
				})
			createResp, err := client.RepoManager.CreateRepo(createParams)
			if err != nil {
				t.Fatal(err)
			}
			repoID = createResp.Payload.RepoID
		}
	}
	log.Printf("Got repo [%s]\n", repoID)

	// waiting for apps indexed by repo indexer
	var app *models.OpenpitrixApp
	{
		for {
			describeParams := app_manager.NewDescribeAppsParams()
			describeParams.WithRepoID([]string{repoID})
			describeParams.WithName([]string{"elk"})
			describeResp, err := client.AppManager.DescribeApps(describeParams)
			if err != nil {
				t.Fatal(err)
			}
			apps := describeResp.Payload.AppSet
			if len(apps) != 0 {
				app = apps[0]
				break
			}
			log.Printf("Waiting for app ...")
			time.Sleep(5 * time.Second)
		}
	}
	log.Printf("Got app [%s]\n", app.Name)

	var appVersion *models.OpenpitrixAppVersion
	{
		describeParams := app_manager.NewDescribeAppVersionsParams()
		describeParams.SetAppID([]string{app.AppID})
		describeResp, err := client.AppManager.DescribeAppVersions(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		appVersions := describeResp.Payload.AppVersionSet
		if len(appVersions) == 0 {
			t.Fatal("App has no version released")
		}

		if len(appVersions) != 1 {
			t.Fatal("We need only one version to test")
		}

		appVersion = appVersions[0]
	}
	log.Printf("Got app version [%s]\n", appVersion.Name)

	// create runtime
	var runtimeID string
	{
		describeParams := runtime_manager.NewDescribeRuntimesParams()
		describeParams.SetSearchWord(&RuntimeNameForTest)
		describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		runtimes := describeResp.Payload.RuntimeSet
		if len(runtimes) != 0 {
			runtimeID = runtimes[0].RuntimeID
		} else {
			createParams := runtime_manager.NewCreateRuntimeParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRuntimeRequest{
					Name:              RuntimeNameForTest,
					Description:       "qingcloud runtime",
					Provider:          "qingcloud",
					RuntimeURL:        "https://api.qingcloud.com",
					RuntimeCredential: `{"access_key_id": "xxxxxxxxxx", "secret_access_key": "xxxxxxxxxxxxxxxx"}`,
					Zone:              "ap2a",
				})
			createResp, err := client.RuntimeManager.CreateRuntime(createParams)
			if err != nil {
				t.Fatal(err)
			}
			runtimeID = createResp.Payload.RuntimeID
		}
	}
	log.Printf("Got runtime [%s]\n", runtimeID)

	var clusterId string
	log.Printf("Creating cluster...\n")
	{
		createParams := cluster_manager.NewCreateClusterParams()
		createParams.SetBody(&models.OpenpitrixCreateClusterRequest{
			AdvancedParam: []string{},
			AppID:         app.AppID,
			Conf:          ClusterConf,
			RuntimeID:     runtimeID,
			VersionID:     appVersion.VersionID,
		})

		createResp, err := client.ClusterManager.CreateCluster(createParams)
		if err != nil {
			t.Fatal(err)
		}

		clusterId = createResp.Payload.ClusterID
	}
	log.Printf("Cluster [%s] created \n", clusterId)
}
