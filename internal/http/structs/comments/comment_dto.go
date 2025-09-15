package commentsdtostructs

type CommentDTO struct {
	ID           int
	UserID       int
	PictureURL   string
	UserName     string
	Content      string
	LastUpdated  string
	Likes        int
	Dislikes     int
	Replies      []CommentDTO
	UserReaction *string // "like", "dislike" or nil if no reaction
}
