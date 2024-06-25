package globals

import (
	"context"
	"strconv"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"google.golang.org/grpc/metadata"
)

func AccountDetailsByPID(pid *types.PID) (*nex.Account, *nex.Error) {
	if pid.Equals(AuthenticationEndpoint.ServerAccount.PID) {
		return AuthenticationEndpoint.ServerAccount, nil
	}

	if pid.Equals(SecureEndpoint.ServerAccount.PID) {
		return SecureEndpoint.ServerAccount, nil
	}

	if pid.Equals(GuestAccount.PID) {
		return GuestAccount, nil
	}

	ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

	response, err := GRPCAccountClient.GetNEXPassword(ctx, &pb.GetNEXPasswordRequest{Pid: pid.LegacyValue()})
	if err != nil {
		Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidPID, "Invalid PID")
	}

	username := strconv.Itoa(int(pid.Value()))
	account := nex.NewAccount(pid, username, response.Password)

	return account, nil
}
