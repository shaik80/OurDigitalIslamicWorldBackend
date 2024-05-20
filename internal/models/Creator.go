package models

type Creator struct {
	CreatorID        string `json:"creatorId"`
	Name             string `json:"name"`
	ChannelLink      string `json:"channerlLink"`
	SubscribersCount int    `json:"subscribersCount"`
	ProfilePic       string `json:"profilePic"`
	LastUpdated      string `json:"lastUpdated"`
}
