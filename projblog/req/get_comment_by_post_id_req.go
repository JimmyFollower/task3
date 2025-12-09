package req

type GetCommentByPostIdRequest struct {
	PostId int `json:"postId" bingding:"required"`
}
