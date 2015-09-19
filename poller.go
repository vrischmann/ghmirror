package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
)

// poller poll regularly the GitHub API for new repositories
type poller struct {
	ds   DataStore
	gh   *github.Client
	freq time.Duration
}

func (p *poller) run() {
	ticker := time.NewTicker(p.freq)

	force := make(chan struct{}, 1)
	force <- struct{}{}

	for {
		select {
		case <-ticker.C:
			p.updateRepositories()
		case <-force:
			p.updateRepositories()
		}
	}
}

func (p *poller) updateRepositories() {
	count := 0
	for page := 0; ; {
		repos, nextPage, err := p.updateRepositoriesForPage(page)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		count += repos

		if nextPage == 0 {
			break
		}

		page = nextPage
	}

	log.Printf("%d repositories updated", count)
}

func (p *poller) updateRepositoriesForPage(page int) (int, int, error) {
	var opts github.RepositoryListOptions
	opts.ListOptions.Page = page

	repos, resp, err := p.gh.Repositories.List("", &opts)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to get user repositories. err=%v", err)
	}

	// Obtain the datastore lock
	p.ds.Lock()
	defer p.ds.Unlock()

	log.Printf("got %d repositories for page %d", len(repos), page)

	nextPage := resp.NextPage

	if len(repos) == 0 {
		return 0, 0, nil
	}

	count := 0
	for _, repo := range repos {
		id := int64(*repo.ID)

		if repo.Private != nil && *repo.Private {
			log.Printf("ignoring repo %s at %s because it's private", *repo.Name, *repo.CloneURL)
			continue
		}

		count++

		ok, err := p.ds.HasRepository(id)
		if err != nil {
			return 0, 0, fmt.Errorf("error while checking for repository in the datastore. err=%v", err)
		}

		var r *Repository
		{
			// Let's add the new repository if it does not exist
			if !ok {
				log.Printf("repository %d does not exist yet, adding it", id)

				localPath := filepath.Join(conf.RepositoriesPath, *repo.FullName)
				r = NewRepository(
					id,
					*repo.Name,
					localPath,
					*repo.CloneURL,
				)

				if stringSliceContains(conf.Webhook.ValidOwnerLogins, *repo.Owner.Login) {
					login := *repo.Owner.Login

					log.Printf("check the webhook exist")

					ok, err := p.webHookExist(login, *repo.Name, p.gh)
					if err != nil {
						return 0, 0, fmt.Errorf("error while checking the webhook exist. err=%v", err)
					}

					if !ok {
						log.Printf("webhook does not exists for %d, %s", id, *repo.FullName)
						log.Printf("creating webhook for repository %d, %s", id, *repo.FullName)

						hookID, err := p.createWebHook(*repo.Owner.Login, *repo.Name, p.gh)
						if err != nil {
							return 0, 0, fmt.Errorf("error while creating webhook. err=%v", err)
						}
						r.HookID = hookID
					} else {
						log.Printf("webhook already exists for %d, %s", id, *repo.FullName)
					}
				}

				if err := p.ds.AddRepository(r); err != nil {
					return 0, 0, fmt.Errorf("error while adding repository to the datastore. err=%v", err)
				}
			} else {
				r, err = p.ds.GetByID(id)
				if err != nil {
					return 0, 0, fmt.Errorf("error while getting repository from the datastore. err=%v", err)
				}
			}
		}

		log.Printf("updating repo %d, %s", r.ID, *repo.FullName)

		if err := r.Update(); err != nil {
			return 0, 0, fmt.Errorf("error while cloning repository. err=%v", err)
		}

		log.Printf("repo %d, %s updated", r.ID, *repo.FullName)
	}

	return count, nextPage, nil
}

func (p *poller) webHookExist(owner, repo string, gh *github.Client) (bool, error) {
	hooks, _, err := gh.Repositories.ListHooks(owner, repo, nil)
	if err != nil {
		return false, err
	}

	exist := false
	for _, hook := range hooks {
		v, ok := hook.Config["url"]
		if !ok {
			continue
		}

		v2, ok := v.(string)
		if !ok {
			continue
		}

		if v2 == conf.Webhook.Endpoint {
			exist = true
		}
	}

	return exist, nil
}

func (p *poller) createWebHook(owner, repo string, gh *github.Client) (int64, error) {
	name := "web"
	active := true

	hook := &github.Hook{
		Name:   &name,
		Events: []string{"push"},
		Config: map[string]interface{}{
			"url":          conf.Webhook.Endpoint,
			"content_type": "json",
			"secret":       conf.Secret,
		},
		Active: &active,
	}

	hook, _, err := gh.Repositories.CreateHook(owner, repo, hook)
	if err != nil {
		return -1, err
	}

	return int64(*hook.ID), nil
}

func stringSliceContains(sl []string, s string) bool {
	for _, el := range sl {
		if el == s {
			return true
		}
	}

	return false
}
