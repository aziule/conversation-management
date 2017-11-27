package nlp

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrRepositoryNotFound = errors.New("Repository not found")

	// repositoryBuilders stores the available Repository builders
	repositoryBuilders = make(map[string]RepositoryBuilder)
)

// RepositoryBuilder is the interface describing a builder for Repository
type RepositoryBuilder func() Repository

// RegisterRepositoryBuilder adds a new RepositoryBuilder to the list of available builders
func RegisterRepositoryBuilder(name string, builder RepositoryBuilder) {
	_, registered := repositoryBuilders[name]

	if registered {
		log.WithField("name", name).Warning("RepositoryBuilder already registered, ignoring")
	}

	repositoryBuilders[name] = builder
}

// NewRepository tries to create a Repository using the available builders.
// Returns ErrRepositoryNotFound if the repository builder isn't found.
func NewRepository(name string) (Repository, error) {
	repositoryBuilder, ok := repositoryBuilders[name]

	if !ok {
		return nil, ErrRepositoryNotFound
	}

	return repositoryBuilder(), nil
}

// Repository is the main interface used to get / store NLP data
type Repository interface {
	GetIntents() ([]*Intent, error)
}
