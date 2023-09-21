package Paxos

import "net"

type Learner struct {
	lis net.Listener
	//学习者id
	id int
	//记录接受者已接受的提案
	acceptedMsg map[int]MsgArgs
}
