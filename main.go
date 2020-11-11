package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"portalfeeders/agents"
	"portalfeeders/utils"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Server struct {
	quit   chan os.Signal
	finish chan bool
	agents []agents.Agent
}

func instantiateLogger(workerName string) (*logrus.Entry, error) {
	var log = logrus.New()
	logsPath := filepath.Join(".", "logs")
	os.MkdirAll(logsPath, os.ModePerm)
	file, err := os.OpenFile(fmt.Sprintf("%s/%s.log", logsPath, workerName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Infof("Failed to log to file - with error: %v", err)
		return nil, err
	}
	log.Out = file
	logger := log.WithFields(logrus.Fields{
		"worker": workerName,
	})
	return logger, nil
}

func registerPRVSeller(
	agentsList []agents.Agent,
) []agents.Agent {
	prvSeller := &agents.PRVSeller{}
	prvSeller.ID = 1
	prvSeller.Name = "prv-seller"
	prvSeller.Frequency = 600
	prvSeller.Quit = make(chan bool)
	prvSeller.RPCClient = utils.NewHttpClient("", os.Getenv("INCOGNITO_PROTOCOL"), os.Getenv("INCOGNITO_HOST"), os.Getenv("INCOGNITO_PORT")) // incognito chain rpc endpoint
	prvSeller.Network = "test"

	logger, err := instantiateLogger(prvSeller.Name)
	if err != nil {
		panic("Could instantiate a logger for prv seller")
	}
	prvSeller.Logger = logger

	prvSeller.Counter = 0
	prvSeller.SellerPrivKey = os.Getenv("SELLER_PRIVKEY")
	prvSeller.SellerAddress = os.Getenv("SELLER_ADDRESS")

	return append(agentsList, prvSeller)
}

func NewServer() *Server {
	agents := []agents.Agent{}
	agents = registerPRVSeller(agents)

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGTERM)
	signal.Notify(quitChan, syscall.SIGINT)
	return &Server{
		quit:   quitChan,
		finish: make(chan bool, len(agents)),
		agents: agents,
	}
}

func (s *Server) NotifyQuitSignal(agents []agents.Agent) {
	sig := <-s.quit
	fmt.Printf("Caught sig: %+v \n", sig)
	// notify all agents about quit signal
	for _, a := range agents {
		a.GetQuitChan() <- true
	}
}

func (s *Server) Run() {
	agents := s.agents
	go s.NotifyQuitSignal(agents)
	for _, a := range agents {
		go executeAgent(s.finish, a)
	}
}

func executeAgent(
	finish chan bool,
	agent agents.Agent,
) {
	agent.Execute() // execute as soon as starting up
	for {
		select {
		case <-agent.GetQuitChan():
			fmt.Printf("Finishing task for %s ...\n", agent.GetName())
			time.Sleep(time.Second * 1)
			fmt.Printf("Task for %s done! \n", agent.GetName())
			finish <- true
			break
		case <-time.After(time.Duration(agent.GetFrequency()) * time.Second):
			agent.Execute()
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	var myEnv map[string]string
	myEnv, _ = godotenv.Read()
	fmt.Println("=========Config============")
	for key, value := range myEnv {
		fmt.Println(key + ": " + value)
	}
	fmt.Println("=========End============")

	runtime.GOMAXPROCS(runtime.NumCPU())
	s := NewServer()
	s.Run()
	for range s.agents {
		<-s.finish
	}
	fmt.Println("Server stopped gracefully!")
}
