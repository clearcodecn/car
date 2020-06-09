package cargo

import (
	"errors"
	"reflect"
	"sync"
)

var (
	ErrIDNotFound = errors.New("id not found")
)

type PushServer interface {
	AddMember(agent *Agent)
	DelMember(agent *Agent)
	Push(agent *Agent, message interface{}) error
	PushId(id string, message interface{}) error
	AddGroup(id string)
	JoinGroup(gid string, agent *Agent)
	DelGroupAgent(agent *Agent)
	PushGroupId(id string, message interface{})
}

var (
	defaultPushServer = newInnerPushServer()
)

func newInnerPushServer() PushServer {
	ips := new(innerPushServer)
	ips.members = make(map[string]*Agent)
	ips.groups = make(map[string]map[string]*Agent)
	return ips
}

type innerPushServer struct {
	memberMutex sync.Mutex
	members     map[string]*Agent

	groupMutex sync.Mutex
	groups     map[string]map[string]*Agent
}

func (i *innerPushServer) PushId(id string, message interface{}) error {
	i.memberMutex.Lock()
	member, ok := i.members[id]
	i.memberMutex.Unlock()
	if !ok {
		return ErrIDNotFound
	}
	msgName := reflect.TypeOf(message).Name()
	return member.WriteMessage(msgName, message)
}

func (i *innerPushServer) AddGroup(id string) {
	i.groupMutex.Lock()
	defer i.groupMutex.Unlock()

	if _, ok := i.groups[id]; !ok {
		i.groups[id] = make(map[string]*Agent)
	}
}

func (i *innerPushServer) JoinGroup(gid string, agent *Agent) {
	i.groupMutex.Lock()
	defer i.groupMutex.Unlock()
	if _, ok := i.groups[gid]; ok {
		i.groups[gid][agent.ID()] = agent
	}
}

func (i *innerPushServer) PushGroupId(id string, message interface{}) {
	var agents []*Agent
	i.groupMutex.Lock()
	for _, v := range i.groups[id] {
		agents = append(agents, v)
	}
	i.groupMutex.Unlock()
	for _, a := range agents {
		i.PushId(a.ID(), message)
	}
}

func (i *innerPushServer) AddMember(agent *Agent) {
	i.memberMutex.Lock()
	defer i.memberMutex.Unlock()

	i.members[agent.ID()] = agent
}

func (i *innerPushServer) DelMember(agent *Agent) {
	i.memberMutex.Lock()
	defer i.memberMutex.Unlock()

	delete(i.members, agent.ID())
}

func (i *innerPushServer) Push(agent *Agent, message interface{}) error {
	return i.PushId(agent.ID(), message)
}

func (i *innerPushServer) DelGroupAgent(agent *Agent) {
	i.groupMutex.Lock()
	defer i.groupMutex.Unlock()
	for gid, group := range i.groups {
		for id, a := range group {
			if a == agent {
				delete(group, id)
			}
		}
		if len(group) == 0 {
			delete(i.groups, gid)
		}
	}
}

func (s *Server) AddMember(agent *Agent) {
	s.opts.PushServer.AddMember(agent)
}

func (s *Server) PushId(id string, message interface{}) error {
	return s.opts.PushServer.PushId(id, message)
}

func (s *Server) AddGroup(id string) {
	s.opts.PushServer.AddGroup(id)
}

func (s *Server) PushGroupId(id string, message interface{}) {
	s.opts.PushServer.PushGroupId(id, message)
}

func (s *Server) DelMember(agent *Agent) {
	s.opts.PushServer.DelMember(agent)
}

func (s *Server) Push(agent *Agent, message interface{}) error {
	return s.opts.PushServer.Push(agent, message)
}

func (s *Server) JoinGroup(id string, agent *Agent) {
	s.opts.PushServer.JoinGroup(id, agent)
}

func (s *Server) DelGroupAgent(agent *Agent) {
	s.opts.PushServer.DelGroupAgent(agent)
}
