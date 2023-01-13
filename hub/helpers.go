package hub

func getOnlineUsers(hub *Hub) []int {
	var onlineUsers []int
	for key := range hub.Clients {
		onlineUsers = append(onlineUsers, key.user.ID)
	}
	return onlineUsers
}
