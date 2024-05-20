package models

type VideoResponse struct {
	Status  bool        `json:"status"`
	Data    Video       `json:"data"`
	Message interface{} `json:"message"`
}
type Video struct {
	VideoID        string         `json:"videoId"`
	Title          string         `json:"title"`
	Thumbnails     []Thumbnail    `json:"thumbnails"`
	Likes          *int           `json:"likes"`
	ViewsCount     string         `json:"viewsCount"`
	UploadDate     string         `json:"uploadDate"`
	VideoCategory  string         `json:"videoCategory"`
	Description    string         `json:"description"`
	Dislikes       *int           `json:"dislikes"`
	IsShort        bool           `json:"isShort"`
	CreatorDetails CreatorDetails `json:"creatorDetails"`
	LastUpdated    string         `json:"lastUpdated"`
	Categories     []string       `json:"categories"`
}

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type CreatorDetails struct {
	Name             string `json:"name"`
	ChannelLink      string `json:"channerlLink"`
	SubscribersCount int    `json:"subscribersCount"`
	ProfilePic       string `json:"profilePic"`
	LastUpdated      string `json:"lastUpdated"`
}
