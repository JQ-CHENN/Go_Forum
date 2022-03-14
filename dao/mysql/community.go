package mysql

import (
	"database/sql"
	"errors"
	"strconv"
	"webapp/models"

	"go.uber.org/zap"
)

var (
	ErrorCommunityInvalidID = errors.New("无效的community_id")
	ErrorPostInvalidID = errors.New("无效的post_id")
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	sqlStr := "select community_id, community_name, introduction,create_time from community where community_id = ?"
	detail := new(models.CommunityDetail)
	err := db.Get(detail, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = ErrorCommunityInvalidID
		}
	}
	return detail, err
}

func InsertPost(post *models.Post) (err error) {
	sqlStr := "insert into post(post_id, title, content, author_id, community_id) values(?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, post.ID, post.Title, post.Content, post.AuthorID, post.CommunityID)

	return
}

func GetPostDetailByID(id int64) (*models.Post, error) {
	sqlStr := "select post_id,title,content,author_id,community_id,status,create_time from post where post_id = ?"
	detail := new(models.Post)
	err := db.Get(detail, sqlStr, id)

	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("the post not found")
			err = ErrorPostInvalidID
		}
	}

	return detail, err
}

func GetPostList(page, size int64) (postlist []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,status,create_time 
			   from post 
			   order by create_time
			   desc
			   limit ?,?`
	postlist = make([]*models.Post, 0, 10)
	err = db.Select(&postlist, sqlStr, (page - 1) * size, size)

	return
}

// 根据给定id列表查询数据
func GetPostsByIds(ids []string) (postlist []*models.Post, err error) {
	postlist = make([]*models.Post, 0)
	for _, id := range ids {
		postID, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		post, err := GetPostDetailByID(int64(postID))
		if err != nil {
			return postlist, err
		}
		postlist = append(postlist, post)
	}

	return postlist, err
}