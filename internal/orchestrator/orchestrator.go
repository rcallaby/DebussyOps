package orchestrator

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/example/orchestrator-template/internal/nlu"
	agentclient "github.com/example/orchestrator-template/internal/agent/client"
)

type Orchestrator struct {
	nlu    nlu.NLU
	log    *zap.SugaredLogger
	agents map[string]*agentclient.AgentClient
}

func NewOrchestrator(n nlu.NLU, logger *zap.SugaredLogger) *Orchestrator {
	return &Orchestrator{nlu: n, log: logger, agents: make(map[string]*agentclient.AgentClient)}
}

func (o *Orchestrator) RegisterAgent(name string, c *agentclient.AgentClient) {
	o.agents[name] = c
	o.log.Infof("agent registered: %s", name)
}

// ParseInput: call NLU to convert free text -> AgentRequest
func (o *Orchestrator) ParseInput(input string) (*nlu.AgentRequest, error) {
	ar, err := o.nlu.Parse(input)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

// RouteAndExecute: route to registered agent and return response
func (o *Orchestrator) RouteAndExecute(ar *nlu.AgentRequest) (map[string]interface{}, error) {
	c, ok := o.agents[ar.Agent]
	if !ok {
		return nil, fmt.Errorf("agent %s not registered", ar.Agent)
	}
	msg := &agentclient.AgentMessage{ID: ar.ID, Action: ar.Action, Payload: ar.Payload}
	res, err := c.Handle(msg)
	if err != nil {
		return nil, err
	}
	return res.Response, nil
}
