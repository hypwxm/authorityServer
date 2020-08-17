package service

import (
	"babygrowing/service/vote/model/voteM"
	"babygrowing/service/vote/model/voteOptionM"
)

func CreateVote(vote *voteM.WbVote) (string, error) {
	return vote.Insert()
}

func ModifyVote(vote *voteM.WbVote, updateQuery *voteM.UpdateByIDQuery) error {
	return vote.Update(updateQuery)
}

func VoteList(query *voteM.Query) ([]*voteM.WbVote, int64, error) {
	return new(voteM.WbVote).List(query)
}

func VoteInfo(query *voteM.GetQuery) (map[string]interface{}, error) {

	vote, err := new(voteM.WbVote).GetByID(query)
	if err != nil {
		return nil, err
	}

	voteOptions, _, err := VoteOptions(&voteOptionM.Query{
		VoteID: vote.ID,
	})

	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"vote":    vote,
		"options": voteOptions,
	}, nil
}

func CreateVoteOption(vote *voteOptionM.WbVoteOption) (string, error) {
	return vote.Insert()
}

func ModifyVoteOption(vote *voteOptionM.WbVoteOption, updateQuery *voteOptionM.UpdateByIDQuery) error {
	return vote.Update(updateQuery)
}

func VoteOptions(query *voteOptionM.Query) ([]*voteOptionM.WbVoteOption, int64, error) {
	return new(voteOptionM.WbVoteOption).List(query)
}

func DeleteVoteOptions(query *voteOptionM.DeleteQuery) error {
	return new(voteOptionM.WbVoteOption).Delete(query)
}
