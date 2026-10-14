package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/naufalfmm/gdc-task-management/internal/models/daos"
	"github.com/naufalfmm/gdc-task-management/internal/models/dtos"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/authtoken"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/password"
	"github.com/naufalfmm/gdc-task-management/internal/pkg/txmanager"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/teams"
	"github.com/naufalfmm/gdc-task-management/internal/repositories/users"
)

type Usecase interface {
	RegisterUser(ctx context.Context, r dtos.RegisterRequest) (dtos.RegisterResponse, error)
}

type usecase struct {
	txManager        txmanager.TxManager
	teamRepositories teams.Repository
	userRepositories users.Repository

	passwordHasher password.PasswordHasher
	authToken      authtoken.AuthToken
}

func NewAuthUsecase(txManager txmanager.TxManager, teamRepositories teams.Repository, userRepositories users.Repository, passwordHasher password.PasswordHasher, authToken authtoken.AuthToken) Usecase {
	return &usecase{
		txManager:        txManager,
		teamRepositories: teamRepositories,
		userRepositories: userRepositories,
		passwordHasher:   passwordHasher,
		authToken:        authToken,
	}
}

func (a *usecase) generateAuthToken(user daos.User) (string, string, error) {
	return a.authToken.Generate(&dtos.LoginClaims{
		UserAuth: dtos.UserAuth{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (a *usecase) checkTeamExisting(ctx context.Context, teamName string) (*daos.Team, error) {
	team, err := a.teamRepositories.GetTeamByName(ctx, teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &team, nil
}

func (a *usecase) RegisterUser(ctx context.Context, r dtos.RegisterRequest) (dtos.RegisterResponse, error) {
	tx, ctx, err := a.txManager.Begin(ctx)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	defer tx.Rollback(ctx)

	team, err := a.checkTeamExisting(ctx, r.TeamName)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	if team == nil {
		createdTeam, err := a.teamRepositories.InsertTeam(ctx, daos.Team{
			Name: r.TeamName,
		})
		if err != nil {
			return dtos.RegisterResponse{}, err
		}

		team = &createdTeam
	}

	hashedPass, err := a.passwordHasher.Hash(r.Password)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	user, err := a.userRepositories.InsertUser(ctx, daos.User{
		Email:        r.Email,
		Name:         r.Name,
		PasswordHash: hashedPass,
		TeamID:       team.ID,
	})
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	tx.Commit(ctx)

	accToken, refToken, err := a.generateAuthToken(user)
	if err != nil {
		return dtos.RegisterResponse{}, err
	}

	return dtos.RegisterResponse{
		Email:    user.Email,
		Name:     user.Email,
		TeamName: team.Name,
		Token: dtos.TokenResponse{
			AccessToken:  accToken,
			RefreshToken: refToken,
		},
	}, nil
}
