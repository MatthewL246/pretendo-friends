package database_wiiu

import (
	"database/sql"

	"github.com/PretendoNetwork/friends/database"
)

// SaveFriendRequest registers a new friend request
func SaveFriendRequest(senderPID uint32, recipientPID uint32, sentTime uint64, expireTime uint64, message string) (uint64, error) {
	var id uint64

	friendRequestBlocked, err := IsFriendRequestBlocked(recipientPID, senderPID)
	if err != nil {
		return 0, err
	}

	// Make sure we don't already have that friend request! If we do, give them the one we already have.
	row, err := database.Manager.QueryRow(`SELECT id FROM wiiu.friend_requests WHERE sender_pid=$1 AND recipient_pid=$2`, senderPID, recipientPID)
	if err != nil {
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	} else if id != 0 {
		// If they aren't blocked, we want to unset the denied status on the previous request we have so that it appears again.
		if friendRequestBlocked {
			return id, nil
		} else {
			err = UnsetFriendRequestDenied(id)
			if err != nil {
				return 0, err
			}

			return id, nil
		}
	}

	row, err = database.Manager.QueryRow(`
		INSERT INTO wiiu.friend_requests (sender_pid, recipient_pid, sent_on, expires_on, message, received, accepted, denied)
		VALUES ($1, $2, $3, $4, $5, false, false, $6) RETURNING id`, senderPID, recipientPID, sentTime, expireTime, message, friendRequestBlocked)
	if err != nil {
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
