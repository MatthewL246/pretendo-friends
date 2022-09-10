package database

import (
	"encoding/base64"

	"github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"go.mongodb.org/mongo-driver/bson"
)

// Get a users sent friend requests
func GetUserFriendRequestsOut(pid uint32) []*nexproto.FriendRequest {
	var sliceMap []map[string]interface{}
	var err error

	if sliceMap, err = cassandraClusterSession.Query(`SELECT id, recipient_pid, sent_on, expires_on, message, received FROM pretendo_friends.friend_requests WHERE sender_pid=? AND accepted=false AND denied=false ALLOW FILTERING`, pid).Iter().SliceMap(); err != nil {
		logger.Critical(err.Error())

		return make([]*nexproto.FriendRequest, 0)
	}

	friendRequestsOut := make([]*nexproto.FriendRequest, 0)

	for i := 0; i < len(sliceMap); i++ {
		recipientPID := uint32(sliceMap[i]["recipient_pid"].(int))

		recipientUserInforation := GetUserInfoByPID(recipientPID)
		encodedMiiData := recipientUserInforation["mii"].(bson.M)["data"].(string)
		decodedMiiData, _ := base64.StdEncoding.DecodeString(encodedMiiData)

		friendRequest := nexproto.NewFriendRequest()

		friendRequest.PrincipalInfo = nexproto.NewPrincipalBasicInfo()
		friendRequest.PrincipalInfo.PID = recipientPID
		friendRequest.PrincipalInfo.NNID = recipientUserInforation["username"].(string)
		friendRequest.PrincipalInfo.Mii = nexproto.NewMiiV2()
		friendRequest.PrincipalInfo.Mii.Name = recipientUserInforation["mii"].(bson.M)["name"].(string)
		friendRequest.PrincipalInfo.Mii.Unknown1 = 0 // replaying from real server
		friendRequest.PrincipalInfo.Mii.Unknown2 = 0 // replaying from real server
		friendRequest.PrincipalInfo.Mii.Data = decodedMiiData
		friendRequest.PrincipalInfo.Mii.Datetime = nex.NewDateTime(0)
		friendRequest.PrincipalInfo.Unknown = 2 // replaying from real server

		friendRequest.Message = nexproto.NewFriendRequestMessage()
		friendRequest.Message.FriendRequestID = uint64(sliceMap[i]["id"].(int64))
		friendRequest.Message.Received = sliceMap[i]["received"].(bool)
		friendRequest.Message.Unknown2 = 1
		friendRequest.Message.Message = sliceMap[i]["message"].(string)
		friendRequest.Message.Unknown3 = 0
		friendRequest.Message.Unknown4 = ""
		friendRequest.Message.GameKey = nexproto.NewGameKey()
		friendRequest.Message.GameKey.TitleID = 0
		friendRequest.Message.GameKey.TitleVersion = 0
		friendRequest.Message.Unknown5 = nex.NewDateTime(134222053376) // idk what this value means but its always this
		friendRequest.Message.ExpiresOn = nex.NewDateTime(uint64(sliceMap[i]["expires_on"].(int64)))
		friendRequest.SentOn = nex.NewDateTime(uint64(sliceMap[i]["sent_on"].(int64)))

		friendRequestsOut = append(friendRequestsOut, friendRequest)
	}

	return friendRequestsOut
}