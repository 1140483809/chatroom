package main

type Data struct {
	Ip string `json:"ip"`
	Type string `json:"type"`
	From string `json:"from"`
	Content string `json:"content"`
	User string `json:"user"`
	UserList []string `json:"user_list"`
}
