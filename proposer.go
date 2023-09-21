package Paxos

type Proposer struct {
	//服务器ID
	id int
	//当前提议者已知最大的轮次
	round int
	//ballot编号
	number int
	//接受者id
	acceptor []int
}
