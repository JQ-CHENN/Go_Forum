package models

import "time"

// 定义用户请求参数结构体

// 用户注册
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"required,eqfield=Password"`
}

// 用户登录
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 社区及详情
type Community struct {
	ID int64	`json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	Community
	Introduction string `json:"introduction,omitempty" db:"introduction"`
	CreateTime time.Time`json:"create_time,omitempty" db:"create_time"`
}

// 帖子及详情
type Post struct {
	AuthorID int64 `json:"author_id" db:"author_id"`
	ID int64 `json:"id" db:"post_id"`
	Title string `json:"title" db:"title" binding:"required"`
	Content string `json:"content" db:"content" binding:"required"`
	CommunityID int64 `json:"community_id" db:"community_id" binding:"required"`
	Status int32 `json:"status" db:"status"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

type PostDetail struct {
	AuthorName string `json:"author_name"`
	Ups int64 `json:"ups"` // 赞数
	*Post
	*CommunityDetail `json:"community_detail"`
}

// 帖子投票数据
type ParamVoteData struct {
	// UserID 从请求中获取当前用户
	PostID int64 `json:"post_id" binding:"required"`
	UpLow int8 `json:"up_low" binding:"oneof=-1 0 1"` // 1赞成  -1反对  0取消投票
}

const (
	OrderTime = "time"
	OrderScore = "score"
)

// url中帖子列表参数
type ParamPostList struct {
	Page int64 `json:"page" form:"page"`
	Size int64 `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

type ParamCommunityList struct {
	*ParamPostList
	CommunityID int64 `json:"cid" form:"cid"`
}

