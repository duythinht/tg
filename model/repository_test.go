package model_test

import (
	"fmt"
	"testing"

	model "github.com/duythinht/tg/model"
)

func TestRepositoryFromHttpURL(t *testing.T) {

	owner := "duythinht"
	repo := "tg"

	httpURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)

	r, err := model.RepositoryFromURL(httpURL)

	if err != nil {
		t.Fatal(err)
	}

	if r.Owner != owner {
		t.Logf("owner should be %s, got %s", owner, r.Owner)
		t.Fail()
	}

	if r.Repository != repo {
		t.Logf("repo should be %s, got: `%s`", repo, r.Repository)
		t.Fail()
	}
}

func TestRepositoryFromSshURL(t *testing.T) {

	owner := "xyz"
	repo := "abc"

	sshURL := fmt.Sprintf("git@github.com:%s/%s.git", owner, repo)

	r, err := model.RepositoryFromURL(sshURL)

	if err != nil {
		t.Fatal(err)
	}

	if r.Owner != owner {
		t.Logf("owner should be %s, got %s", owner, r.Owner)
		t.Fail()
	}

	if r.Repository != repo {
		t.Logf("repo should be %s, got: %s", repo, r.Repository)
		t.Fail()
	}
}

func TestRepositoryFromWrongURLPattern(t *testing.T) {

	owner := "xyz"
	repo := "abc"

	url0 := fmt.Sprintf("git@xyz.com:%s/%s.git", owner, repo)

	_, err := model.RepositoryFromURL(url0)

	if err == nil {
		t.Fatal("error should not be nil")
	}

	url1 := fmt.Sprintf("git@xyz.com:%s%s", owner, repo)

	model.RepositoryFromURL(url1)

	if err == nil {
		t.Fatal("error should not be nil")
	}
}
