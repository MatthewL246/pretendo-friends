package database_wiiu

import (
	"github.com/PretendoNetwork/friends/database"
)

// UnsetUserBlocked removes a block from a user
func UnsetUserBlocked(user1_pid uint32, user2_pid uint32) error {
	result, err := database.Manager.Exec(`
		DELETE FROM wiiu.blocks WHERE blocker_pid=$1 AND blocked_pid=$2`, user1_pid, user2_pid)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrPIDNotFound
	}

	return nil
}
