package hub

// Make online users global with Redis
// This function can switch to getHubUsers - to see which users participate in the Room
func getOnlineUsers(hub *Hub) ([]int, error) {
	var onlineUsers []int
	for _, c := range hub.Clients {
		if c.user != nil {
			onlineUsers = append(onlineUsers, c.user.ID)
		}
	}
	return onlineUsers, nil
}
