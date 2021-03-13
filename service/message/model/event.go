package model

import (
	"babygrow/event"
	"context"
	"log"
)

// 消息各种各种，通过调方法的方式太麻烦了，
// 改用事件驱动，

func init() {
	event.Ebus.Subscribe("serve:message", func(m *GMessage) {
		log.Printf("==========%+v", m)
		insert(context.Background(), m)
	})
}
