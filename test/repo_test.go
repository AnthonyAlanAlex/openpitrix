// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/topic"
	apiclient "openpitrix.io/openpitrix/test/client"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/models"
)

//var clientConfig = &ClientConfig{}
//
//func init() {
//	clientConfig = GetClientConfig()
//	log.Printf("Got Client Config: %+v", clientConfig)
//}

//func TestMain(m *testing.M) {
//	os.Exit(m.Run())
//}

var (
	repoUrl = "https://helm-chart-repo.pek3a.qingstor.com/svc-catalog-charts/"
	//repoUrl = "https://helm-chart-repo.pek3a.qingstor.com/kubernetes-charts/"
)

func deleteRepo(t *testing.T, client *apiclient.Openpitrix, testRepoName string) {
	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.SetName([]string{testRepoName})
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.RepoManager.DescribeRepos(describeParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	repos := describeResp.Payload.RepoSet
	for _, repo := range repos {
		deleteParams := repo_manager.NewDeleteReposParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteReposRequest{
				RepoID: []string{repo.RepoID},
			})
		_, err := client.RepoManager.DeleteRepos(deleteParams, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestRepo(t *testing.T) {
	client := GetClient(clientConfig)

	// delete old repo
	testRepoName := "e2e_test_repo"
	deleteRepo(t, client, testRepoName)
	ioClient := GetIoClient(clientConfig)

	// test validate repo
	repoType := "https"
	credential := "{}"
	validateParams := repo_manager.NewValidateRepoParams()
	validateParams.SetType(&repoType)
	validateParams.SetURL(&repoUrl)
	validateParams.SetCredential(&credential)
	validateResp, err := client.RepoManager.ValidateRepo(validateParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	if validateResp.Payload.Ok != true {
		t.Fatal("validate repo failed")
	}

	// create repo
	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        testRepoName,
			Description: "description",
			Type:        "https",
			URL:         repoUrl,
			Credential:  `{}`,
			Visibility:  "public",
			Providers:   []string{constants.ProviderKubernetes},
			CategoryID:  "xx,yy,zz",
		},
	)
	createResp, err := client.RepoManager.CreateRepo(createParams, nil)
	require.NoError(t, err)

	repoId := createResp.Payload.RepoID

	msg := ioClient.ReadMessage()
	require.Equal(t, "repo", msg.Resource.ResourceType)
	require.Equal(t, topic.Create, msg.Type)
	require.Equal(t, "active", msg.Resource.Values["status"])
	require.Equal(t, repoId, msg.Resource.ResourceId)

	// repo-event pending
	msg = ioClient.ReadMessage()
	repoEventId := msg.Resource.ResourceId
	require.Equal(t, "repo_event", msg.Resource.ResourceType)
	require.Equal(t, topic.Create, msg.Type)
	require.Equal(t, "pending", msg.Resource.Values["status"])

	// repo-event success
	msg = ioClient.ReadMessage()
	require.Equal(t, "repo_event", msg.Resource.ResourceType)
	require.Equal(t, topic.Update, msg.Type)
	require.Equal(t, "successful", msg.Resource.Values["status"])
	require.Equal(t, repoEventId, msg.Resource.ResourceId)

	// modify repo
	modifyParams := repo_manager.NewModifyRepoParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyRepoRequest{
			RepoID:      repoId,
			Description: "cc",
			Type:        "https",
			URL:         repoUrl,
			Credential:  `{}`,
			Visibility:  "private",
			Providers:   []string{constants.ProviderKubernetes},
			CategoryID:  "aa,bb,cc,xx",
		})
	modifyResp, err := client.RepoManager.ModifyRepo(modifyParams, nil)
	if err != nil {
		t.Fatal(err)
	}

	msg = ioClient.ReadMessage()
	require.Equal(t, "repo", msg.Resource.ResourceType)
	require.Equal(t, topic.Update, msg.Type)
	require.Equal(t, repoId, msg.Resource.ResourceId)

	t.Log(modifyResp)
	// describe repo
	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.SetName([]string{testRepoName})
	describeParams.SetStatus([]string{constants.StatusActive})
	describeParams.WithRepoID([]string{repoId})
	describeResp, err := client.RepoManager.DescribeRepos(describeParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	repos := describeResp.Payload.RepoSet
	if len(repos) != 1 {
		t.Fatalf("failed to describe repos with params [%+v]", describeParams)
	}
	repo := repos[0]
	require.Equal(t, testRepoName, repo.Name)
	require.Equal(t, "cc", repo.Description)
	require.Equal(t, repoUrl, repo.URL)

	var enabledCategoryIds []string
	var disabledCategoryIds []string
	for _, a := range repo.CategorySet {
		if a.Status == constants.StatusEnabled {
			enabledCategoryIds = append(enabledCategoryIds, a.CategoryID)
		}
		if a.Status == constants.StatusDisabled {
			disabledCategoryIds = append(disabledCategoryIds, a.CategoryID)
		}
	}

	require.Equal(t, "aa,bb,cc,xx", getSortedString(enabledCategoryIds))
	require.Equal(t, "yy,zz", getSortedString(disabledCategoryIds))
	// delete repo
	deleteParams := repo_manager.NewDeleteReposParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteReposRequest{
		RepoID: []string{repoId},
	})
	deleteResp, err := client.RepoManager.DeleteRepos(deleteParams, nil)
	if err != nil {
		t.Fatal(err)
	}

	msg = ioClient.ReadMessage()
	require.Equal(t, "repo", msg.Resource.ResourceType)
	require.Equal(t, topic.Update, msg.Type)
	require.Equal(t, repoId, msg.Resource.ResourceId)

	t.Log(deleteResp)
	// describe deleted repo
	describeParams.WithRepoID([]string{repoId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithName(nil)
	describeResp, err = client.RepoManager.DescribeRepos(describeParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	repos = describeResp.Payload.RepoSet
	if len(repos) != 1 {
		t.Fatalf("failed to describe repos with params [%+v]", describeParams)
	}
	repo = repos[0]
	if repo.RepoID != repoId {
		t.Fatalf("failed to describe repo")
	}
	if repo.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete repo, got repo status [%s]", repo.Status)
	}

	t.Log("test repo finish, all test is ok")
}

func testDescribeReposWithLabelSelector(t *testing.T,
	repoId string,
	labels string,
	selectors string) {
	client := GetClient(clientConfig)

	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.SetLabel(&labels)
	describeParams.SetSelector(&selectors)
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.RepoManager.DescribeRepos(describeParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	if describeResp.Payload.RepoSet[0].RepoID != repoId {
		t.Fatalf("describe repo with filter failed")
	}
	//repo := describeResp.Payload.RepoSet[0]
	//for i, label := range repo.Labels {
	//	if label.LabelKey != labels[i].LabelKey {
	//		t.Fatalf("repo label key not matched")
	//	}
	//	if label.LabelValue != labels[i].LabelValue {
	//		t.Fatalf("repo label value not matched")
	//	}
	//}
	//for i, selector := range repo.Selectors {
	//	if selector.SelectorKey != selectors[i].SelectorKey {
	//		t.Fatalf("repo selector key not matched")
	//	}
	//	if selector.SelectorValue != selectors[i].SelectorValue {
	//		t.Fatalf("repo selector value not matched")
	//	}
	//}
}

func TestRepoLabelSelector(t *testing.T) {
	client := GetClient(clientConfig)
	// Create a test repo that can attach label and selector on it
	testRepoName := "e2e_test_repo"
	labels := generateLabels()
	selectors := generateLabels()
	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        testRepoName,
			Description: "description",
			Type:        "https",
			URL:         repoUrl,
			Credential:  `{}`,
			Visibility:  "public",
			Providers:   []string{constants.ProviderKubernetes},
			Labels:      labels,
			Selectors:   selectors,
		})
	createResp, err := client.RepoManager.CreateRepo(createParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	repoId := createResp.Payload.RepoID
	testDescribeReposWithLabelSelector(t, repoId, labels, selectors)

	i := 0
	for i < 10 {
		i++
		newLabels := generateLabels()
		newSelectors := generateLabels()
		modifyParams := repo_manager.NewModifyRepoParams()
		modifyParams.SetBody(
			&models.OpenpitrixModifyRepoRequest{
				RepoID:    repoId,
				Providers: []string{constants.ProviderKubernetes},
				Labels:    newLabels,
				Selectors: newSelectors,
			},
		)
		_, err := client.RepoManager.ModifyRepo(modifyParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		testDescribeReposWithLabelSelector(t, repoId, newLabels, newSelectors)
	}

	// delete repo
	deleteParams := repo_manager.NewDeleteReposParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteReposRequest{
		RepoID: []string{repoId},
	})
	deleteResp, err := client.RepoManager.DeleteRepos(deleteParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)

	t.Log("test repo label and selector finish, all test is ok")
}

func TestDeleteInternalRepo(t *testing.T) {
	client := GetClient(clientConfig)

	// test delete internal repo, should be failed
	for _, repoId := range constants.InternalRepos {
		deleteParams := repo_manager.NewDeleteReposParams()
		deleteParams.WithBody(&models.OpenpitrixDeleteReposRequest{
			RepoID: []string{repoId},
		})
		_, err := client.RepoManager.DeleteRepos(deleteParams, nil)
		require.Error(t, err)
	}
}
