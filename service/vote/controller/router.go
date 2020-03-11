package controller

import "github.com/hypwxm/rider"

func Router() *rider.Router {
	route := rider.NewRouter()

	// 创建一个投票
	route.POST("/createVote", createVote)
	route.POST("/modifyVote", modifyVote)
	route.POST("/list", list)
	route.POST("voteInfo", voteInfo)


	route.POST("/createVoteOption", createVoteOption)
	route.POST("/modifyVoteOption", modifyVoteOption)
	route.POST("/voteOptions", voteOptions)
	route.POST("/deleteVoteOptions", deleteVoteOptions)


	return route
}
