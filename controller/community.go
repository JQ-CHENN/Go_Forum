package controller

import (
	"strconv"
	"webapp/logic"
	"webapp/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


func CommunityHandler(c *gin.Context) {
	// 查询到所有社区(community_id, community_name)
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList err", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	// 根据请求的社区id，返回社区详情
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		zap.L().Info("community id is invalid")
		ResponseError(c, CodeInvalidParma)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("community detail is not found", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func CreatePostHandler(c *gin.Context) {
	// 获取参数及参数校验

	p := new(models.Post)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("Post with invild parmas", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}

	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic CreatePost falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 获取请求的社区id, 解析参数
	postID := c.Param("id")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		zap.L().Info("post id is invalid")
		ResponseError(c, CodeInvalidParma)
		return
	}
	// 根据id获取帖子数据
	data, err := logic.GetPostDetailByID(id)

	if err != nil {
		zap.L().Error("logic GetPostDetail falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := GetPageInfo(c)

	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic GetPostList falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}


// 根据前端传来的参数动态获取帖子列表
// 按创建时间，或者按投票分数排序
func GetPostListHandler2(c *gin.Context) {
	// 获取参数   /api/v1/posts2?page=1&size=10&order=time (query string)
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
		Order: models.OrderTime, // 默认按时间排序
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("get url params err", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}

	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic GetPostList falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// 根据社区id去查帖子
func GetPostListByCommunityIDHandler(c *gin.Context) {
	// 获取参数   /api/v1/posts2?page=1&size=10&order=time&cid=1 (query string)
	p := &models.ParamCommunityList{
		ParamPostList: &models.ParamPostList{
			Page: 1,
			Size: 10,
			Order: models.OrderTime, // 默认按时间排序
		},
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("get url params err", zap.Error(err))
		ResponseError(c, CodeInvalidParma)
		return
	}

	data, err := logic.GetPostListByCommunityID(p)
	if err != nil {
		zap.L().Error("logic GetPostListByCommunityID falied", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}