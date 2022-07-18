package entity

// Post represents single post in feed
type Post struct {
	UUID		string	`json:"uuid"`
	Content		string	`json:"content`
	Likes		uint	`json:"likes"`
	Dislikes	uint	`json:"dislikes"`
}
