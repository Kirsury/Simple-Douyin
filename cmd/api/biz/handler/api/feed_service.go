// Code generated by hertz generator.

package api

import (
	"context"
	"encoding/json"
	"log"

	"Simple-Douyin/cmd/api/biz/handler/pack"
	api "Simple-Douyin/cmd/api/biz/model/api"
	"Simple-Douyin/cmd/api/biz/mw"
	"Simple-Douyin/cmd/api/rpc"
	"Simple-Douyin/kitex_gen/feed"
	"Simple-Douyin/pkg/constants"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
)

// Feed .
// @router /douyin/feed [GET]
func Feed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.FeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		log.Println("[ypx debug] api BindAndValidate error")
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	log.Println("[ypx debug] api BindAndValidate success and prepare to rpc.Feed")
	uid := int64(0)
	if token, err := mw.JwtMiddleware.ParseTokenString(req.Token); err == nil {
		claims := jwt.ExtractClaimsFromToken(token)
		userid, _ := claims[constants.IdentityKey].(json.Number).Int64()
		uid = userid
	}
	log.Println("[ypx debug] api feed userid", uid)
	next_time, videos, err := rpc.Feed(context.Background(), &feed.FeedRequest{
		LatestTime: req.LatestTime,
		UserId:     uid,
	})
	if err != nil {
		log.Println("[ypx debug] api rpc.Feed fail")
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	log.Println("[ypx debug] api rpc.Feed success")
	resp := new(api.FeedResponse)

	resp.StatusCode = 0
	resp.StatusMsg = "视频流推送成功"
	resp.NextTime = next_time
	resp.VideoList = pack.Videos(videos)

	for _, v := range resp.VideoList {
		log.Println("[ypx debug] api resp.VideoList:", v.ID, " ", v.CommentCount, " ", v.CoverURL)
		log.Println("[ypx debug] api resp.VideoList:", v.FavoriteCount, " ", v.PlayURL)
		log.Println("[ypx debug] api resp.VideoList.Author", v.Author.Name)
	}

	c.JSON(consts.StatusOK, resp)
}
