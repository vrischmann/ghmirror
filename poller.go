package main

import (
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
	repos, _, err := p.gh.Repositories.List("", nil)
	if err != nil {
		log.Printf("unable to get user repositories. err=%v", err)
		return
	}

	// Obtain the datastore lock
	p.ds.Lock()
	defer p.ds.Unlock()

	for _, repo := range repos {
		id := int64(*repo.ID)

		ok, err := p.ds.HasRepository(id)
		if err != nil {
			log.Printf("error while checking for repository in the datastore. err=%v", err)
			return
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

				if stringSliceContains(conf.WebHook.ValidOwnerLogins, *repo.Owner.Login) {
					login := *repo.Owner.Login

					log.Printf("check the webhook exist")

					ok, err := p.webHookExist(login, *repo.Name, p.gh)
					if err != nil {
						log.Printf("error while checking the webhook exist. err=%v", err)
						return
					}

					if !ok {
						log.Printf("webhook does not exists for %d, %s", id, *repo.FullName)
						log.Printf("creating webhook for repository %d, %s", id, *repo.FullName)
						hookId, err := p.createWebHook(*repo.Owner.Login, *repo.Name, p.gh)
						if err != nil {
							log.Printf("error while creating webhook. err=%v", err)
							return
						}
						r.HookID = hookId
					} else {
						log.Printf("webhook already exists for %d, %s", id, *repo.FullName)
					}
				}

				if err := p.ds.AddRepository(r); err != nil {
					log.Printf("error while adding repository to the datastore. err=%v", err)
					return
				}
			} else {
				r, err = p.ds.GetByID(id)
				if err != nil {
					log.Printf("error while getting repository from the datastore. err=%v", err)
					return
				}
			}
		}

		log.Printf("updating repo %d, %s", r.ID, *repo.FullName)

		if err := r.Update(); err != nil {
			log.Printf("error while cloning repository. err=%v", err)
			return
		}

		log.Printf("repo %d, %s updated", r.ID, *repo.FullName)
	}

	log.Println("repositories updated")
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

		if v2 == conf.WebHook.Endpoint {
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
			"url":          conf.WebHook.Endpoint,
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
