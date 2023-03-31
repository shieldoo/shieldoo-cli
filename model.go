package main

type Group struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ObjectId string `json:"objectId"`
}

type FirewallRule struct {
	Protocol string  `json:"protocol"`
	Port     string  `json:"port"`
	Host     string  `json:"host"`
	Groups   []Group `json:"groups"`
}

type Firewall struct {
	Id       string         `json:"id"`
	Name     string         `json:"name"`
	RulesIn  []FirewallRule `json:"rulesIn"`
	RulesOut []FirewallRule `json:"rulesOut"`
}

type Listener struct {
	ListenPort  int    `json:"listenPort"`
	Protocol    string `json:"protocol"`
	ForwardPort int    `json:"forwardPort"`
	ForwardHost string `json:"forwardHost"`
	Description string `json:"description"`
}

type Server struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Groups      []Group    `json:"groups"`
	Firewall    Firewall   `json:"firewall"`
	Listeners   []Listener `json:"listeners"`
	Autoupdate  bool       `json:"autoupdate"`
	IpAddress   string     `json:"ipAddress"`
	Description string     `json:"description"`
}
