package Paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Learner struct {
	lis net.Listener
	//学习者id
	id int
	//记录接受者已接受的提案
	acceptedMsg map[int]MsgArgs
}

func (l *Learner)majority() int{
	return len(l.acceptedMsg) / 2 + 1
}

//将提案编号更大的提案存储到acceptedMsg中
func (l *Learner) Learn(args *MsgArgs,reply *MsgReply)error{
	a := l.acceptedMsg[args.From]
	if a.NUmber < args.NUmber{
		l.acceptedMsg[args.From] = *args
		reply.Ok = true
	}else{
		reply.Ok = false
	}
	return nil
}


//判断一个提案 是否被半数接受者接收 如果是，返回该批准的提案值
func (l *Learner) chosen() interface{}{
	acceptCounts := make(map[int]int)
	acceptMsg := make(map[int]MsgArgs)

	for _,accepted := range l.acceptedMsg{
		if accepted.NUmber != 0{
			acceptCounts[accepted.NUmber]++
			acceptMsg[accepted.NUmber] = accepted
		}
	}

	for n,count := range acceptCounts{
		if count >= l.majority(){
			return acceptMsg[n].Value
		}
	}

	return nil
}


func (l *Learner) server(id int){
	rpcs := rpc.NewServer()
	rpcs.Register(l)

	addr := fmt.Sprintf(":%d",id)
	lis,e := net.Listen("tcp",addr)
	if e != nil{
		log.Fatal("listen error:",e)
	}

	l.lis = lis
	go func() {
		for{
			conn,err := l.lis.Accept()
			if err != nil{
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()

}

func(l *Learner) close(){
	l.lis.Close()
}