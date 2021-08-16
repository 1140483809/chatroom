package main

import "encoding/json"

var h = hub{
	connections: make(map[*connection]bool),

	broadcast: make(chan []byte),

	register: make(chan *connection),

	unregister: make(chan *connection),
}

func (h *hub) run()  {
	for {
		select {
		case c := <- h.register:
			h.connections[c] = true

			c.data.Ip = c.ws.RemoteAddr().String()

			c.data.Type = "handshake"

			c.data.UserList = userList

			dataB,_ := json.Marshal(c.data)
			c.send <-dataB
		case c:= <- h.unregister:
			if _,ok := h.connections[c];ok{
				delete(h.connections,c)
				close(c.send)
			}
		case data := <- h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- data:
				default:
					delete(h.connections,c)
					close(c.send)
				}
			}
		}
	}
}
