package commentsdtostructs

type CommentDTO struct {
	ID         int
	UserID     int
	PictureURL string
	UserName   string
	Content    string
	CreatedAt  string
	Likes      int
	Dislikes   int
	Replies    []CommentDTO
}
