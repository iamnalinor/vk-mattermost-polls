package handler

import (
	"github.com/iamnalinor/vk-mattermost-polls/internal/polling"
)

func RegisterHandlers(dispatcher *polling.Dispatcher) {
	dispatcher.Register("!newpoll", newPollHandler)
	dispatcher.Register("", pollTextHandler)
	dispatcher.Register("!done", pollDoneHandler)
	dispatcher.Register("!cancel", pollCancelHandler)
	dispatcher.Register("!poll", getPollHandler)
	dispatcher.Register("!vote", voteHandler)
	dispatcher.Register("!closepoll", closePollHandler)
	dispatcher.Register("!deletepoll", deletePollHandler)
}
