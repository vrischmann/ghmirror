package internal

type Repository struct {
	ID        int64
	Name      string
	LocalPath string
	CloneURL  string
	HookID    int64
}

func NewRepository(id int64, name, localPath, cloneURL string) *Repository {
	return &Repository{
		ID:        id,
		Name:      name,
		LocalPath: localPath,
		CloneURL:  cloneURL,
	}
}

type Repositories []*Repository

type BlacklistedOwner struct {
	ID   int64
	Name string
}

type OwnersBlacklist []*BlacklistedOwner

type BlacklistedRepository struct {
	ID           int64
	Organization string
	Name         string
}

type RepositoriesBlacklist []*BlacklistedRepository
