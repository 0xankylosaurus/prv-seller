package agents

import (
	"fmt"
	"portalfeeders/utils"

	"github.com/sirupsen/logrus"
)

type AgentAbs struct {
	ID        int
	Name      string
	Frequency int // in sec
	Quit      chan bool
	RPCClient *utils.HttpClient
	Network   string // mainnet, testnet, ...
	Logger    *logrus.Entry
}

type Agent interface {
	Execute()
	GetName() string
	GetFrequency() int
	GetQuitChan() chan bool
	GetNetwork() string
}

func (a *AgentAbs) Execute() {
	fmt.Println("Abstract agent is executing...")
}

func (a *AgentAbs) GetName() string {
	return a.Name
}

func (a *AgentAbs) GetFrequency() int {
	return a.Frequency
}

func (a *AgentAbs) GetQuitChan() chan bool {
	return a.Quit
}

func (a *AgentAbs) GetNetwork() string {
	return a.Network
}
