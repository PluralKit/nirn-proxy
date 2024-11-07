package lib

import (
	"github.com/hashicorp/memberlist"
	"os"
	"time"
	"net"
)

func InitMemberList(advertiseAddr string, knownMembers []string, port int, proxyPort string, manager *QueueManager) *memberlist.Memberlist {
	config := memberlist.DefaultLANConfig()
	config.BindPort = port
	if advertiseAddr != "" {
		// i'm pretty sure being able to pass a dns name into BIND_IP is a quirk of go's http server
		// but it's useful for us, so let's handle that for memberlist as well
		if net.ParseIP(advertiseAddr) == nil {
			addrs, err := net.LookupHost(advertiseAddr)
			if err != nil {
				panic(err)
			}
			if len(addrs) == 0 {
				panic("memberlist init: could not find ip address for advertiseAddr")
			}
			advertiseAddr = addrs[0]
		}
		config.AdvertiseAddr = advertiseAddr
	}
	config.AdvertisePort = port
	config.Delegate = NirnDelegate{
		proxyPort: proxyPort,
	}

	config.Events = manager.GetEventDelegate()

	//DEBUG CODE
	if os.Getenv("NODE_NAME") != "" {
		config.Name = os.Getenv("NODE_NAME")
		config.DeadNodeReclaimTime = 1 * time.Nanosecond
	}

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	manager.SetCluster(list, proxyPort)

	_, err = list.Join(knownMembers)
	if err != nil {
		logger.Info("Failed to join existing cluster, ok if this is the first node")
		logger.Error(err)
	}

	var members string
	for _, member := range list.Members() {
		members += member.Name + " "
	}

	logger.Info("Connected to cluster nodes: [ " + members + "]")
	return list
}