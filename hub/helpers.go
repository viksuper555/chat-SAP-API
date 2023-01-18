package hub

// Make online users global with Redis
// This function can switch to getHubUsers - to see which users participate in the Room
func getOnlineUsers(hub *Hub) ([]int, error) {
	var onlineUsers []int
	for key := range hub.Clients {
		if key.user != nil {
			onlineUsers = append(onlineUsers, key.user.ID)
		}
	}
	return onlineUsers, nil
}

func getUserRoomIds(userId int) ([]string, error) {
	roomIds := make([]string, 0)
	if err := DB.Table("user_room").Distinct("room_id").Where("user_id = ?", userId).Find(&roomIds).Error; err != nil {
		return nil, err
	}
	return roomIds, nil
}
