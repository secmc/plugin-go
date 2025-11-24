package player

type Player struct {
	uuid string
	name string
}

func New(name, uuid string) *Player {
	return &Player{
		name: name,
		uuid: uuid,
	}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) UUID() string {
	return p.uuid
}
