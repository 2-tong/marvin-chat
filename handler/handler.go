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

// MsgReply 定义消息回复的结构
type MsgReply struct {
	Type    string // "text" 或 "image"
	Content string // 文本内容或图片URL
}

// MsgHandler 修改函数签名
type MsgHandler func(msg string) MsgReply

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

func HandleTextMsg(in string, replyFunc func(reply MsgReply)) {
	matchTimes := 0
	for matcher, handler := range handlerMap {
		if matcher.IsMatch(in) {
			matchTimes++
			reply := handler(in)
			if reply.Content != "" {
				replyFunc(reply)
			}
		}
	}
	if matchTimes == 0 {
		replyFunc(MsgReply{Type: "text", Content: "听不懂思密达 😅"})
	}
}
