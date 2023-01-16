package hub

// Make online users global with Redis
// This function can switch to getHubUsers - to see which users participate in the Room
func getOnlineUsers(hub *Hub) []int {
	var onlineUsers []int
	for key := range hub.Clients {
		if key.user != nil {
			onlineUsers = append(onlineUsers, key.user.ID)
		}
	}
	return onlineUsers
}
