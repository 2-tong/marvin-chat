package handler

import (
	_ "regexp"
	"strings"
)

var handlerMap = map[msgMatcher]MsgHandler{}

type msgMatcher interface {
	//IsMatch 是否处理当前
	IsMatch(msg string) bool
}

type MsgHandler func(msg string) string

type SimpleMsgMatcher struct {
	command string
}

func (r *SimpleMsgMatcher) IsMatch(msg string) bool {
	return strings.Contains(msg, r.command)
}

func (r *SimpleMsgMatcher) InitRegexp(command string) {
	r.command = command
}

func RegisterSimpleMsgHandler(command string, handler MsgHandler) {
	matcher := SimpleMsgMatcher{command}
	RegisterMsgHandler(&matcher, handler)
}

func RegisterMsgHandler(matcher msgMatcher, handler MsgHandler) {
	handlerMap[matcher] = handler
}

func HandleTextMsg(in string, replyFunc func(reply string)) {
	matchTimes := 0
	for matcher, handler := range handlerMap {
		if matcher.IsMatch(in) {
			matchTimes++
			str := handler(in)
			if str != "" {
				replyFunc(str)
			}
		}
	}
	if matchTimes == 0 {
		replyFunc("听不懂思密达 😅")
	}
}
