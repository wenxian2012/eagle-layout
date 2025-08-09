package service

import (
	"github.com/google/wire"

	"github.com/go-eagle/eagle-layout/internal/repository"
)

// ServiceSet is service providers.
var ServiceSet = wire.NewSet(
	NewUserService,     // for simple http service
	NewUserHTTPService, // for full HTTP service
	repository.RepositorySet,
)
