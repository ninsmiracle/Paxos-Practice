package Paxos

import "testing"

func start(acceptorIds []int,learnerIds[]int)([]*Acceptor,[]*Learner){
	acceptors := make([]*Acceptor,0)
	for _,aid := range acceptorIds{
		a:= newAcceptor(aid,learnerIds)
		acceptors = append(acceptors,a)
	}

	learners := make([]*Learner,0)
	for _,lid := range learnerIds{
		l := newLearner(lid,acceptorIds)
		learners = append(learners,l)
	}

	return acceptors,learners
}

func cleanup(acceptors []*Acceptor,learners []*Learner){
	for _,a := range acceptors{
		a.close()
	}

	for _,l := range learners{
		l.close()
	}
}

func TestSinglePropose(t *testing.T){
	acceptorIds := []int{1001,1002,1003}

	learnerIds := []int{2001}
	acceptors,learners := start(acceptorIds,learnerIds)

	defer cleanup(acceptors,learners)

	p := &Proposer{
		id:1,
		acceptors:acceptorIds,
	}

	value := p.propose("hello world")
	if value != "hello world"{
		t.Errorf("value = %s,excepted %s",value,"hello world")
	}

	learnValue := learners[0].chosen()
	if learnValue != value{
		t.Errorf("learnValue = %s,excepted %s",learnValue,"hello world")
	}
}

func TestTwoProposers(t *testing.T){
	acceptorIds := []int{1001,1002,1003}

	learnerIds := []int{2001}
	acceptors,learners := start(acceptorIds,learnerIds)
	defer cleanup(acceptors,learners)

	//ID是用来区分propose的ID
	p1 := &Proposer{
		id:1,
		acceptors:acceptorIds,
	}
	value1 := p1.propose("hello world")

	p2 := &Proposer{
		id:2,
		acceptors:acceptorIds,
	}
	value2 := p2.propose("bad world")

	if value1 != value2{
		t.Errorf("value1 = %s,value2 = %s",value1,value2)
	}

	learnValue := learners[0].chosen()
	if learnValue != value1{
		t.Errorf("learnValue = %s,excepted %s",learnValue,value1)
	}



}
