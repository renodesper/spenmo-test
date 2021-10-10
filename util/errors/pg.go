package errors

import (
	"fmt"
	"net/http"

	"gitlab.com/renodesper/spenmo-test/util/error"
)

var (
	FailedNoRows        = error.NewError(http.StatusBadRequest, "PG2000", fmt.Errorf("no rows selected"))
	FailedEmailExist    = error.NewError(http.StatusBadRequest, "PG2001", fmt.Errorf("email is already exist"))
	FailedUsernameExist = error.NewError(http.StatusBadRequest, "PG2002", fmt.Errorf("username is already exist"))
)

var (
	FailedUserExist    = error.NewError(http.StatusBadRequest, "PG2101", fmt.Errorf("user is already exist"))
	FailedUserNotFound = error.NewError(http.StatusNotFound, "PG2102", fmt.Errorf("user cannot be found"))
	FailedUserCreate   = error.NewError(http.StatusInternalServerError, "PG2103", fmt.Errorf("unable to create user"))
	FailedUserUpdate   = error.NewError(http.StatusInternalServerError, "PG2104", fmt.Errorf("unable to update user"))
	FailedUserDelete   = error.NewError(http.StatusInternalServerError, "PG2105", fmt.Errorf("unable to delete user"))
	FailedUserFetch    = error.NewError(http.StatusInternalServerError, "PG2106", fmt.Errorf("unable to fetch user"))
	FailedUsersFetch   = error.NewError(http.StatusInternalServerError, "PG2107", fmt.Errorf("unable to fetch users"))
)

var (
	FailedTeamExist    = error.NewError(http.StatusBadRequest, "PG2201", fmt.Errorf("team is already exist"))
	FailedTeamNotFound = error.NewError(http.StatusNotFound, "PG2202", fmt.Errorf("team cannot be found"))
	FailedTeamCreate   = error.NewError(http.StatusInternalServerError, "PG2203", fmt.Errorf("unable to create team"))
	FailedTeamUpdate   = error.NewError(http.StatusInternalServerError, "PG2204", fmt.Errorf("unable to update team"))
	FailedTeamDelete   = error.NewError(http.StatusInternalServerError, "PG2205", fmt.Errorf("unable to delete team"))
	FailedTeamFetch    = error.NewError(http.StatusInternalServerError, "PG2206", fmt.Errorf("unable to fetch team"))
	FailedTeamsFetch   = error.NewError(http.StatusInternalServerError, "PG2207", fmt.Errorf("unable to fetch teams"))
)

var (
	FailedTeamMemberExist    = error.NewError(http.StatusBadRequest, "PG2301", fmt.Errorf("team member is already exist"))
	FailedTeamMemberNotFound = error.NewError(http.StatusNotFound, "PG2302", fmt.Errorf("team member cannot be found"))
	FailedTeamMemberCreate   = error.NewError(http.StatusInternalServerError, "PG2303", fmt.Errorf("unable to create team member"))
	FailedTeamMemberUpdate   = error.NewError(http.StatusInternalServerError, "PG2304", fmt.Errorf("unable to update team member"))
	FailedTeamMemberDelete   = error.NewError(http.StatusInternalServerError, "PG2305", fmt.Errorf("unable to delete team member"))
	FailedTeamMemberFetch    = error.NewError(http.StatusInternalServerError, "PG2306", fmt.Errorf("unable to fetch team member"))
	FailedTeamMembersFetch   = error.NewError(http.StatusInternalServerError, "PG2307", fmt.Errorf("unable to fetch team members"))
)

var (
	FailedWalletExist    = error.NewError(http.StatusBadRequest, "PG2401", fmt.Errorf("wallet is already exist"))
	FailedWalletNotFound = error.NewError(http.StatusNotFound, "PG2402", fmt.Errorf("wallet cannot be found"))
	FailedWalletCreate   = error.NewError(http.StatusInternalServerError, "PG2403", fmt.Errorf("unable to create wallet"))
	FailedWalletUpdate   = error.NewError(http.StatusInternalServerError, "PG2404", fmt.Errorf("unable to update wallet"))
	FailedWalletDelete   = error.NewError(http.StatusInternalServerError, "PG2405", fmt.Errorf("unable to delete wallet"))
	FailedWalletFetch    = error.NewError(http.StatusInternalServerError, "PG2406", fmt.Errorf("unable to fetch wallet"))
	FailedWalletsFetch   = error.NewError(http.StatusInternalServerError, "PG2407", fmt.Errorf("unable to fetch wallets"))
)

var (
	FailedCardExist    = error.NewError(http.StatusBadRequest, "PG2501", fmt.Errorf("card is already exist"))
	FailedCardNotFound = error.NewError(http.StatusNotFound, "PG2502", fmt.Errorf("card cannot be found"))
	FailedCardCreate   = error.NewError(http.StatusInternalServerError, "PG2503", fmt.Errorf("unable to create card"))
	FailedCardUpdate   = error.NewError(http.StatusInternalServerError, "PG2504", fmt.Errorf("unable to update card"))
	FailedCardDelete   = error.NewError(http.StatusInternalServerError, "PG2505", fmt.Errorf("unable to delete card"))
	FailedCardFetch    = error.NewError(http.StatusInternalServerError, "PG2506", fmt.Errorf("unable to fetch card"))
	FailedCardsFetch   = error.NewError(http.StatusInternalServerError, "PG2507", fmt.Errorf("unable to fetch cards"))
)
