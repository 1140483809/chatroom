package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type connection struct {
	ws *websocket.Conn
	send chan []byte
	data *Data
}

type hub struct {
	connections map[*connection]bool
	broadcast chan []byte
	register chan *connection
	unregister chan *connection
}

func (c *connection)writer()  {
	 for message := range c.send{
	 	c.ws.WriteMessage(websocket.TextMessage, message)
	 }
	 c.ws.Close()
}

var userList = []string{}

func (c *connection)reader()  {
	for {
		_ ,message , err := c.ws.ReadMessage()
		if err != nil {
			h.unregister <- c
			break
		}
		json.Unmarshal(message,&c.data)
		switch c.data.Type {
		case "login":
			c.data.User = c.data.Content
			c.data.From = c.data.User

			userList = append(userList,c.data.User)
			c.data.UserList = userList

			dataB ,_ := json.Marshal(c.data)
			h.broadcast <- dataB
		case "user":
			c.data.Type = "user"
			dataB ,_ := json.Marshal(c.data)
			h.broadcast <- dataB
		case "logout":
			c.data.Type = "logout"
			userList = remove(userList,c.data.User)
			c.data.UserList = userList
			c.data.Content = c.data.User

			dataB,_ := json.Marshal(c.data)
			h.broadcast <- dataB
			h.unregister <- c
		default:
			fmt.Println("其他")
		}
	}
}

func remove(slice []string,user string)[]string  {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0]==user{
		return []string{}
	}

	var myslice = []string{}

	for i := range slice{
		if slice[i] == user && i == count{
			return slice[:count]
		}else if slice[i]==user {
			myslice = append(slice[:i],slice[i+1:]...)
			break
		}
	}
	return myslice
}

var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandle(w http.ResponseWriter,r *http.Request)  {
	ws ,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		return
	}

	c := &connection{send: make(chan[]byte,128),ws: ws,data: &Data{}}

	h.register <- c

	go c.writer()
	c.reader()
	defer func() {
		c.data.Type = "logout"
		userList = remove(userList,c.data.User)
		c.data.UserList = userList
		c.data.Content = c.data.User

		dataB,_ := json.Marshal(c.data)
		h.broadcast <- dataB
		h.unregister <- c
	}()
}
