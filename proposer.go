package Paxos

import "fmt"

type Proposer struct {
	//服务器ID
	id int
	//当前提议者已知最大的轮次
	round int
	//ballot编号
	number int
	//接受者id
	acceptors []int
}

func (p *Proposer)majority() int{
	return len(p.acceptors) / 2 + 1
}

//提案编号(轮次，服务器id)
func (p *Proposer)proposalNumber() int {
	return p.round << 16 | p.id
}


func (p *Proposer) propose(v interface{}) interface{} {
	p.round++
	p.number = p.proposalNumber()

	//第一阶段
	prepareCount := 0
	maxNumber := 0
	for _,aid := range p.acceptors{
		args := MsgArgs{
			NUmber: p.number,
			From: p.id,
			To: aid,
		}

		//给accept发请求
		reply := new(MsgReply)
		//go的rpc call  name这个参数指定响应函数
		err := call(fmt.Sprintf("127.0.0.1:%d",aid),"Acceptor.Prepare",args,reply)
		if !err{
			continue
		}


		if reply.Ok{
			//收到的回复数量，计数
			prepareCount++
			//acceptor发回来的序列号找出最大值
			if reply.Number > maxNumber{
				maxNumber = reply.Number
				v = reply.Value
			}
		}

		if prepareCount == p.majority(){
			//超过最大值了，进入二阶段
			break
		}
	}

	acceptCount := 0
	if  prepareCount >= p.majority(){
		for _,aid := range p.acceptors{
			args := MsgArgs{
				NUmber: p.number,
				Value: v,
				From: p.id,
				To:aid,
			}
			reply := new(MsgReply)
			ok := call(fmt.Sprintf("127.0.0.1:%d",aid),"Acceptor.Accept",args,reply)
			if !ok{
				continue
			}

			if reply.Ok{
				acceptCount++
			}
		}
	}

	if acceptCount >= p.majority(){
		return v
	}

	return nil
}
