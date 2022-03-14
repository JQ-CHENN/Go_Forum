package redis

import "webapp/models"

func GetPostIDsInOrder(p *models.ParamPostList) ([]string ,error) {
	// 从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}

	// 确定查询的起始
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(key, start, end).Result()
}

// 按社区根据id列表查找
func GetCommunityPostIDsInOrder(p *models.ParamCommunityList) ([]string ,error) {
	// 从redis获取id
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}

	// 确定查询的起始
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(key, start, end).Result()
}