package gogogo

type PushServer interface {
	AddMember(agent *Agent)
	DelMember(agent *Agent)
	Push(agent *Agent, message interface{})
}
