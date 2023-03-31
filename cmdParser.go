package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseGroup(data string) (Group, error) {
	var mygroup Group
	// parse group
	parts := strings.Split(data, "=")
	if len(parts) != 2 {
		return mygroup, fmt.Errorf("invalid group format: %s", data)
	}
	switch parts[0] {
	case "id":
		mygroup.Id = parts[1]
	case "objectId":
		mygroup.ObjectId = parts[1]
	case "objectid":
		mygroup.ObjectId = parts[1]
	case "name":
		mygroup.Name = parts[1]
	default:
		return mygroup, fmt.Errorf("invalid group format: %s", data)
	}
	return mygroup, nil
}

func parseListeners(listeners string) ([]Listener, error) {
	var mylisteners []Listener
	// Listeners (comma separated) - list of listeners in format ListenerPort;Protocol;ForwardPort;ForwardHost;Description
	// example: 80;tcp;8080;myhost.example.com,443;tcp;8443;myhost.example.com
	for _, l := range strings.Split(listeners, ",") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// parse listener
		parts := strings.Split(l, ";")
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid listener format: %s", l)
		}
		var desc string
		if len(parts) > 4 {
			desc = parts[4]
		}
		listenPort, _ := strconv.Atoi(parts[0])
		forwardPort, _ := strconv.Atoi(parts[2])
		if listenPort < 1 || listenPort > 65535 {
			return nil, fmt.Errorf("invalid listener port: %s", parts[0])
		}
		if forwardPort < 1 || forwardPort > 65535 {
			return nil, fmt.Errorf("invalid forward port: %s", parts[2])
		}
		mylistener := Listener{
			ListenPort:  listenPort,
			Protocol:    parts[1],
			ForwardPort: forwardPort,
			ForwardHost: parts[3],
			Description: desc,
		}
		if !regexp.MustCompile(`^(tcp|udp)$`).MatchString(mylistener.Protocol) {
			return nil, fmt.Errorf("invalid protocol: %s", mylistener.Protocol)
		}
		if mylistener.ForwardHost == "" {
			return nil, fmt.Errorf("invalid forward host: %s", mylistener.ForwardHost)
		}
		mylisteners = append(mylisteners, mylistener)
	}
	return mylisteners, nil
}

func parseFirewallRules(groups string) ([]FirewallRule, error) {
	var rules []FirewallRule
	// Array of input firewall rules in format protocol;port;host;group-ids.
	// example: any;any;any,tcp;22;group;id=demo.shieldoo.net:groups:1,udp;53;group;objectId=e7549a43-f3c2-4d0d-9cd1-6811a107cdc4;objectId=a3e4ead5-ffb7-4d94-ba71-0185b5466426

	// parse rule from array
	for _, r := range strings.Split(groups, ",") {
		r := strings.TrimSpace(r)
		if r == "" {
			continue
		}
		// parse rule
		parts := strings.Split(r, ";")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid rule format: %s", r)
		}
		myrule := FirewallRule{
			Protocol: parts[0],
			Port:     parts[1],
			Host:     parts[2],
		}
		// validate data
		if regexp.MustCompile(`^(any|icmp|tcp|udp)$`).MatchString(myrule.Protocol) == false {
			return nil, fmt.Errorf("invalid protocol: %s", myrule.Protocol)
		}
		if regexp.MustCompile(`^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$|^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])-([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$|^any$`).MatchString(myrule.Port) == false {
			return nil, fmt.Errorf("invalid port: %s", myrule.Port)
		}
		if regexp.MustCompile(`^(any|group)$`).MatchString(myrule.Host) == false {
			return nil, fmt.Errorf("invalid host: %s", myrule.Host)
		}
		// parse group
		for _, g := range parts[3:] {
			g := strings.TrimSpace(g)
			if g == "" {
				continue
			}
			// parse group
			mygroup, err := parseGroup(g)
			if err != nil {
				return nil, err
			}
			// add group to rule
			myrule.Groups = append(myrule.Groups, mygroup)
		}
		// add rule to array
		rules = append(rules, myrule)
	}
	return rules, nil
}
