package nex_friends_wiiu

import (
	database_wiiu "github.com/PretendoNetwork/friends/database/wiiu"
	"github.com/PretendoNetwork/friends/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	friends_wiiu "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu"
	friends_wiiu_types "github.com/PretendoNetwork/nex-protocols-go/v2/friends-wiiu/types"
)

func GetRequestBlockSettings(err error, packet nex.PacketInterface, callID uint32, pids *types.List[*types.PrimitiveU32]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.FPD.InvalidArgument, "") // TODO - Add error message
	}

	connection := packet.Sender().(*nex.PRUDPConnection)

	settings := types.NewList[*friends_wiiu_types.PrincipalRequestBlockSetting]()
	settings.Type = friends_wiiu_types.NewPrincipalRequestBlockSetting()

	// TODO - Improve this. Use less database_wiiu reads
	if pids.Each(func(i int, pid *types.PrimitiveU32) bool {
		setting := friends_wiiu_types.NewPrincipalRequestBlockSetting()
		setting.PID = pid

		isBlocked, err := database_wiiu.IsFriendRequestBlocked(connection.PID().LegacyValue(), pid.Value)
		if err != nil {
			globals.Logger.Critical(err.Error())
			return true
		}

		setting.IsBlocked = types.NewPrimitiveBool(isBlocked)

		settings.Append(setting)

		return false
	}) {
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, "") // TODO - Add error message
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureEndpoint.LibraryVersions(), globals.SecureEndpoint.ByteStreamSettings())

	settings.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = friends_wiiu.ProtocolID
	rmcResponse.MethodID = friends_wiiu.MethodGetRequestBlockSettings
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
