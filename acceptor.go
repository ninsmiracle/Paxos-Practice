package Paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

//朴素的理解就是pegasus的primary的角色
type Acceptor struct {
	lis net.Listener

	//服务器id
	id int

	//接受者承诺的提案编号，如果为0，说明接受者没有收到过任何prepare消息
	minProposal int
	//接受者已接受的提案编号，如果为0，则表示没有接受任何提案
	acceptedNumber int
	//接受者已接受的提案值，如果没有接受任何提案，则为nil
	acceptedValue interface{}

	//学习者id
	learners []int
}

//用于一阶段处理 准备阶段
func (a *Acceptor)Prepare(args *MsgArgs,reply *MsgReply) error{
	if args.NUmber > a.minProposal{
		a.minProposal = args.NUmber
		reply.Number = a.acceptedNumber
		reply.Value = a.acceptedValue
		reply.Ok = true
	}else{
		reply.Ok = false
	}

	return nil
}

//二阶段处理 接收函数
func (a *Acceptor)Accept(args *MsgArgs,reply *MsgReply) error{
	//新消息编号大于接受者之前的最小承诺值
	if args.NUmber >= a.minProposal{
		a.minProposal = args.NUmber
		a.acceptedNumber = args.NUmber
		a.acceptedValue = args.Value
		reply.Ok = true

		//发给全部的学习者(其实就是副本)
		for _,lid := range a.learners{
			go func(learner int) {
				addr  := fmt.Sprintf("127.0.0.1:%d",learner)
				args.From = a.id
				args.To	 = learner
				resp := new(MsgReply)
				ok:=call(addr,"Learner.Learn",args,resp)

				if !ok{
					return
				}
			}(lid)
		}
	}else {
		reply.Ok = false
	}

	return nil
}

//建立连接
func (a* Acceptor)server()  {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	//创建监听端口
	addr := fmt.Sprintf(":%d",a.id)
	l,e	 := net.Listen("tcp",addr)
	if e != nil{
		log.Fatal("Listen error:",e)
	}

	a.lis = l
	go func() {
		for{
			//监听到有事件
			conn,err := a.lis.Accept()
			if err != nil{
				continue
			}
			//创建连接
			go rpcs.ServeConn(conn)
		}

	}()
}

func(a *Acceptor) close(){
	//停止监听
	a.lis.Close()
}


func newAcceptor(id int,learners []int)*Acceptor{
	acceptor := &Acceptor{
		id : id,
		learners: learners,
	}
	acceptor.server()

	return acceptor
}




