package user

import (
	"context"

	"github.com/coopersmall/subswag/apm"
	"github.com/coopersmall/subswag/cache"
	"github.com/coopersmall/subswag/domain"
	"github.com/coopersmall/subswag/domain/user"
	"github.com/coopersmall/subswag/repos"
	servicesdomain "github.com/coopersmall/subswag/services/domain"
	"github.com/coopersmall/subswag/streams/publishers"
	"github.com/coopersmall/subswag/utils"
)

type UsersService struct {
	standardService *servicesdomain.StandardService[user.UserID, *user.User]
	usersCache      cache.IUsersCache
	usersRepo       repos.IUsersRepo
	usersPublisher  publishers.IUsersPublisher
}

func NewUsersService(
	logger utils.ILogger,
	tracer apm.ITracer,
	usersCache cache.IUsersCache,
	usersRepo repos.IUsersRepo,
	usersPublisher publishers.IUsersPublisher,
) *UsersService {
	standardService := servicesdomain.NewStandardService[user.UserID, *user.User](
		"users",
		logger,
		tracer,
		struct {
			Get    func(context.Context, user.UserID) (*user.User, error)
			All    func(context.Context) ([]*user.User, error)
			Create func(context.Context, *user.User) error
			Update func(context.Context, *user.User) error
			Delete func(context.Context, user.UserID) error
		}{
			Get:    usersRepo.Get,
			All:    usersRepo.All,
			Create: usersRepo.Create,
			Update: usersRepo.Update,
			Delete: usersRepo.Delete,
		},
		&struct {
			Get    func(context.Context, user.UserID, func(context.Context, user.UserID) (*user.User, error)) (*user.User, error)
			Set    func(context.Context, user.UserID, *user.User) error
			Delete func(context.Context, user.UserID) error
		}{
			Get:    usersCache.Get,
			Set:    usersCache.Set,
			Delete: usersCache.Delete,
		},
		&struct {
			Create func(context.Context, *user.User) error
			Update func(context.Context, *user.User) error
			Delete func(context.Context, *user.User) error
		}{
			Update: usersPublisher.PublishUpdate,
		},
	)

	return &UsersService{
		standardService: standardService,
		usersRepo:       usersRepo,
		usersCache:      usersCache,
	}
}

func (s *UsersService) GetUser(ctx context.Context, userId user.UserID) (*user.User, error) {
	return s.standardService.Get(ctx, userId)
}

func (s *UsersService) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	return s.standardService.All(ctx)
}

func (s *UsersService) CreateUser(ctx context.Context, data user.UserData) (*user.User, error) {
	user := user.NewUser()
	user.UserData = data
	return user, s.createUser(ctx, user)
}

func (s *UsersService) CreateUserWithId(ctx context.Context, id user.UserID, data user.UserData) (*user.User, error) {
	user := &user.User{
		ID:       id,
		UserData: data,
		Metadata: domain.NewMetadata(),
	}
	return user, s.createUser(ctx, user)
}

func (s *UsersService) createUser(ctx context.Context, user *user.User) error {
	return s.standardService.Create(ctx, user.ID, user)
}

func (s *UsersService) UpdateUser(ctx context.Context, user *user.User) (*user.User, error) {
	return user, s.standardService.Update(ctx, user.ID, user)
}

func (s *UsersService) DeleteUser(ctx context.Context, userId user.UserID) error {
	return s.standardService.Delete(ctx, userId)
}
