package controller

import (
	"encoding/json"
	"github.com/hypwxm/rider"
	"worldbar/service/vote/model/voteM"
	"worldbar/service/vote/model/voteOptionM"
	"worldbar/service/vote/service"
	"worldbar/util/response"
)

func createVote(c rider.Context) {
	sender := response.NewSender()
	(func() {
		vote := new(voteM.WbVote)
		err := json.Unmarshal(c.Body(), &vote)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		id, err := service.CreateVote(vote)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
	})()
	c.SendJson(200, sender)
}

func modifyVote(c rider.Context) {
	sender := response.NewSender()
	(func() {
		vote := new(voteM.WbVote)
		err := json.Unmarshal(c.Body(), &vote)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.ModifyVote(vote, &voteM.UpdateByIDQuery{
			ID:      vote.ID,
			Title:   vote.Title,
			Comment: vote.Comment,
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
	})()
	c.SendJson(200, sender)
}

func list(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(voteM.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		list, total, err := service.VoteList(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(list, int(total))
	})()
	c.SendJson(200, sender)
}

func voteInfo(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(voteM.GetQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		data, err := service.VoteInfo(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(data)
	})()
	c.SendJson(200, sender)
}

func createVoteOption(c rider.Context) {
	sender := response.NewSender()
	(func() {
		vote := new(voteOptionM.WbVoteOption)
		err := json.Unmarshal(c.Body(), &vote)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		id, err := service.CreateVoteOption(vote)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success(id)
	})()
	c.SendJson(200, sender)
}

func modifyVoteOption(c rider.Context) {
	sender := response.NewSender()
	(func() {
		vote := new(voteOptionM.WbVoteOption)
		err := json.Unmarshal(c.Body(), &vote)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.ModifyVoteOption(vote, &voteOptionM.UpdateByIDQuery{
			ID:      vote.ID,
			Title:   vote.Title,
			Comment: vote.Comment,
		})
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("操作成功")
	})()
	c.SendJson(200, sender)
}

func voteOptions(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(voteOptionM.Query)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		data, total, err := service.VoteOptions(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.SuccessList(data, int(total))
	})()
	c.SendJson(200, sender)
}

func deleteVoteOptions(c rider.Context) {
	sender := response.NewSender()
	(func() {
		query := new(voteOptionM.DeleteQuery)
		err := json.Unmarshal(c.Body(), &query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		err = service.DeleteVoteOptions(query)
		if err != nil {
			sender.Fail(err.Error())
			return
		}
		sender.Success("删除成功")
	})()
	c.SendJson(200, sender)
}
