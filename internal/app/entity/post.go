package entity

// Post represents single post in feed
type Post struct {
	UID			string	`json:"uid"`
	Content		string	`json:"content`
	Likes		uint	`json:"likes"`
	Dislikes	uint	`json:"dislikes"`
}
