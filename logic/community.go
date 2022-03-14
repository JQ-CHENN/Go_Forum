package logic

import (
	"webapp/dao/mysql"
	"webapp/dao/redis"
	"webapp/models"
	"webapp/pkg/snowflake"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (detail *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetail(id)
}

func CreatePost(post *models.Post) (err error) {
	post.ID = snowflake.GenID()
	if err = mysql.InsertPost(post); err != nil {
		return
	}

	return redis.CreatePost(post.ID)
}

func GetPostDetailByID(id int64) (*models.PostDetail, error) {
	post, err := mysql.GetPostDetailByID(id)
	if err != nil {
		return nil, err
	}
	authorName, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, err
	}

	communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
	if err != nil {
		return nil, err
	}

	detail := &models.PostDetail{
		AuthorName: authorName,
		Post: post,
		CommunityDetail: communityDetail,
	}

	return detail, err
}

func GetPostList(page, size int64) (posts []*models.PostDetail, err error) {
	postlist, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	posts = make([]*models.PostDetail, 0, len(postlist))
	for _, post := range postlist {
		authorName, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			return nil, err
		}
		detail := &models.PostDetail{
			AuthorName: authorName,
			Post: post,
			CommunityDetail: communityDetail,
		}
		posts = append(posts, detail)
	}

	return 
}

func GetPostList2(p *models.ParamPostList) (posts []*models.PostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	// 根据列表去mysql查询帖子详情
	postlist, err := mysql.GetPostsByIds(ids)
	if err != nil {
		return
	}

	// 查询每篇帖子的赞数
	upslist, err := redis.GetPostVoteUps(ids)
	if err != nil {
		return
	}

	// 组装post,以返回更多详情
	posts = make([]*models.PostDetail, 0, len(postlist))
	for idx, post := range postlist {
		authorName, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			return nil, err
		}
		detail := &models.PostDetail{
			AuthorName: authorName,
			Ups: upslist[idx],
			Post: post,
			CommunityDetail: communityDetail,
		}
		posts = append(posts, detail)
	}

	return
}

func GetPostListByCommunityID(p *models.ParamCommunityList) (posts []*models.PostDetail, err error) {
	// 使用zinterstore 把保存社区下帖子id的set与保存帖子分数的zset 生成一个新的zset
	// 按新的zset按之前逻辑取数据
	// 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	// 根据列表去mysql查询帖子详情
	postlist, err := mysql.GetPostsByIds(ids)
	if err != nil {
		return
	}

	// 查询每篇帖子的赞数
	upslist, err := redis.GetPostVoteUps(ids)
	if err != nil {
		return
	}

	// 组装post,以返回更多详情
	posts = make([]*models.PostDetail, 0, len(postlist))
	for idx, post := range postlist {
		authorName, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			return nil, err
		}
		detail := &models.PostDetail{
			AuthorName: authorName,
			Ups: upslist[idx],
			Post: post,
			CommunityDetail: communityDetail,
		}
		posts = append(posts, detail)
	}

	return
}