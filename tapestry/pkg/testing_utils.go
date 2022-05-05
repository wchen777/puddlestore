package pkg

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	//t "tapestry/tapestry"
)

/*
   Parse an ID from String
*/
func MakeID(stringID string) ID {
	var id ID

	for i := 0; i < DIGITS && i < len(stringID); i++ {
		d, err := strconv.ParseInt(stringID[i:i+1], 16, 0)
		if err != nil {
			return id
		}
		id[i] = Digit(d)
	}
	for i := len(stringID); i < DIGITS; i++ {
		id[i] = Digit(0)
	}

	return id
}

var tapestriesByAddress map[string]*Node = make(map[string]*Node)
var tapestryMapMutex *sync.Mutex = &sync.Mutex{}

func registerCachedTapestry(tapestry ...*Node) {
	tapestryMapMutex.Lock()
	defer tapestryMapMutex.Unlock()
	for _, t := range tapestry {
		tapestriesByAddress[t.Node.Address] = t
	}
}

func unregisterCachedTapestry(tapestry ...*Node) {
	tapestryMapMutex.Lock()
	defer tapestryMapMutex.Unlock()
	for _, t := range tapestry {
		delete(tapestriesByAddress, t.Node.Address)
	}
}

func AddOne(ida string, addr string, tap []*Node) (t1 *Node, tapNew []*Node, err error) {
	t1, err = Start(MakeID(ida), 0, addr)
	if err != nil {
		return nil, tap, err
	}
	registerCachedTapestry(t1)
	tapNew = append(tap, t1)
	time.Sleep(1000 * time.Millisecond) //Wait for availability
	return
}

func MakeTapestries(connectThem bool, ids ...string) ([]*Node, error) {
	tapestries := make([]*Node, 0, len(ids))
	for i := 0; i < len(ids); i++ {
		connectTo := ""
		if i > 0 && connectThem {
			connectTo = tapestries[0].Node.Address
		}
		t, err := Start(MakeID(ids[i]), 0, connectTo)
		if err != nil {
			return tapestries, err
		}
		registerCachedTapestry(t)
		tapestries = append(tapestries, t)
		time.Sleep(10 * time.Millisecond)
	}
	return tapestries, nil
}

func KillTapestries(ts ...*Node) {
	fmt.Println("killing")
	unregisterCachedTapestry(ts...)
	for _, t := range ts {
		t.Kill()
	}
	fmt.Println("finished killing")
}

func MakeRandomTapestries(seed int64, count int) ([]*Node, error) {
	r := rand.New(rand.NewSource(seed))

	ts := make([]*Node, 0, count)

	for i := 0; i < count; i++ {
		connectTo := ""
		if i > 0 {
			connectTo = ts[0].Node.Address
		}
		t, err := Start(IntToID(r.Int()), 0, connectTo)
		if err != nil {
			return ts, err
		}
		registerCachedTapestry(t)
		ts = append(ts, t)
		time.Sleep(10 * time.Millisecond)
	}

	return ts, nil
}

func IntToID(x int) ID {
	var id ID
	for i := range id {
		id[i] = Digit(x % BASE)
		x = x / BASE
	}
	return id
}
