package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"gitlab.com/renodesper/spenmo-test/service"
)

type Set struct {
	GetHealthCheckEndpoint            endpoint.Endpoint
	GetAllUsersEndpoint               endpoint.Endpoint
	GetUserByIDEndpoint               endpoint.Endpoint
	CreateUserEndpoint                endpoint.Endpoint
	UpdateUserEndpoint                endpoint.Endpoint
	DeleteUserByIDEndpoint            endpoint.Endpoint
	GetAllTeamsEndpoint               endpoint.Endpoint
	GetTeamEndpoint                   endpoint.Endpoint
	CreateTeamEndpoint                endpoint.Endpoint
	UpdateTeamEndpoint                endpoint.Endpoint
	DeleteTeamEndpoint                endpoint.Endpoint
	CreateTeamMemberEndpoint          endpoint.Endpoint
	GetTeamMembersEndpoint            endpoint.Endpoint
	DeleteTeamMemberEndpoint          endpoint.Endpoint
	DeleteTeamMembersByTeamIDEndpoint endpoint.Endpoint
	DeleteTeamMembersByUserIDEndpoint endpoint.Endpoint
	GetAllWalletsEndpoint             endpoint.Endpoint
	GetWalletEndpoint                 endpoint.Endpoint
	CreateWalletEndpoint              endpoint.Endpoint
	UpdateWalletEndpoint              endpoint.Endpoint
	DeleteWalletByIDEndpoint          endpoint.Endpoint
	DeleteWalletsByTeamIDEndpoint     endpoint.Endpoint
	DeleteWalletsByUserIDEndpoint     endpoint.Endpoint
	GetAllCardsEndpoint               endpoint.Endpoint
	GetCardEndpoint                   endpoint.Endpoint
	CreateCardEndpoint                endpoint.Endpoint
	UpdateCardEndpoint                endpoint.Endpoint
	DeleteCardByIDEndpoint            endpoint.Endpoint
	DeleteCardsByWalletIDEndpoint     endpoint.Endpoint
}

// New ...
func New(
	env string,
	healthSvc service.HealthService,
	userSvc service.UserService,
	teamSvc service.TeamService,
	teamMemberSvc service.TeamMemberService,
	walletSvc service.WalletService,
	cardSvc service.CardService,
) Set {
	return Set{
		GetHealthCheckEndpoint:            MakeHealthCheckEndpoint(healthSvc),
		GetAllUsersEndpoint:               MakeGetAllUsersEndpoint(userSvc),
		GetUserByIDEndpoint:               MakeGetUserByIDEndpoint(userSvc),
		CreateUserEndpoint:                MakeCreateUserEndpoint(userSvc),
		UpdateUserEndpoint:                MakeUpdateUserEndpoint(userSvc),
		DeleteUserByIDEndpoint:            MakeDeleteUserByIDEndpoint(userSvc),
		GetAllTeamsEndpoint:               MakeGetAllTeamsEndpoint(teamSvc),
		GetTeamEndpoint:                   MakeGetTeamEndpoint(teamSvc),
		CreateTeamEndpoint:                MakeCreateTeamEndpoint(teamSvc),
		UpdateTeamEndpoint:                MakeUpdateTeamEndpoint(teamSvc),
		DeleteTeamEndpoint:                MakeDeleteTeamEndpoint(teamSvc),
		CreateTeamMemberEndpoint:          MakeCreateTeamMemberEndpoint(teamMemberSvc),
		GetTeamMembersEndpoint:            MakeGetTeamMembersEndpoint(teamMemberSvc),
		DeleteTeamMemberEndpoint:          MakeDeleteTeamMemberEndpoint(teamMemberSvc),
		DeleteTeamMembersByTeamIDEndpoint: MakeDeleteTeamMembersByTeamIDEndpoint(teamMemberSvc),
		DeleteTeamMembersByUserIDEndpoint: MakeDeleteTeamMembersByUserIDEndpoint(teamMemberSvc),
		GetAllWalletsEndpoint:             MakeGetAllWalletsEndpoint(walletSvc),
		GetWalletEndpoint:                 MakeGetWalletEndpoint(walletSvc),
		CreateWalletEndpoint:              MakeCreateWalletEndpoint(walletSvc),
		UpdateWalletEndpoint:              MakeUpdateWalletEndpoint(walletSvc),
		DeleteWalletByIDEndpoint:          MakeDeleteWalletByIDEndpoint(walletSvc),
		DeleteWalletsByTeamIDEndpoint:     MakeDeleteWalletsByTeamIDEndpoint(walletSvc),
		DeleteWalletsByUserIDEndpoint:     MakeDeleteWalletsByUserIDEndpoint(walletSvc),
		GetAllCardsEndpoint:               MakeGetAllCardsEndpoint(cardSvc),
		GetCardEndpoint:                   MakeGetCardEndpoint(cardSvc),
		CreateCardEndpoint:                MakeCreateCardEndpoint(cardSvc),
		UpdateCardEndpoint:                MakeUpdateCardEndpoint(cardSvc),
		DeleteCardByIDEndpoint:            MakeDeleteCardByIDEndpoint(cardSvc),
		DeleteCardsByWalletIDEndpoint:     MakeDeleteCardsByWalletIDEndpoint(cardSvc),
	}
}
