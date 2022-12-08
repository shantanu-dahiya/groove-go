package parser

import (
	"math/rand"
	"strconv"
)

type Buddy struct {
	Id             int
	SymmetricKey   []byte
	Circuits       []*Circuit
	RNG            *rand.Rand
	DeadDrop       []byte
	TerminalServer *Server
}

func NewBuddy(id int, symmetricKey []byte) *Buddy {
	var seed int64 = 0
	for _, bt := range symmetricKey {
		seed += int64(bt)
	}
	rng := rand.New(rand.NewSource(seed))

	b := &Buddy{Id: id, SymmetricKey: symmetricKey, RNG: rng}
	b.AddDeadDrop()
	b.CreateCircuits() // TODO: make two circuits

	return b
}

func (b *Buddy) AddDeadDrop() {
	terminalServers := FetchServers(LastLayer)
	address := b.RNG.Int()
	i := address % len(terminalServers)
	b.DeadDrop = []byte(strconv.Itoa(address))
	b.TerminalServer = terminalServers[i]
}
