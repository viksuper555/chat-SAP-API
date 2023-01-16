package hub

func getOnlineUsers(hub *Hub) []int {
	var onlineUsers []int
	for key := range hub.Clients {
		if key.user != nil {
			onlineUsers = append(onlineUsers, key.user.ID)
		}
	}
	return onlineUsers
}
