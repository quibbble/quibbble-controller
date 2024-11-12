package parse

type ICard interface {
	GetID() string
	GetCost() int
	GetEnabled() bool
}

type Card struct {
	Enabled bool // only used in parsing to determine if a card is ready for prime-time

	ID          string
	Name        string
	Description string

	Cost int

	Conditions []Condition
	Targets    []Choose

	Hooks  []Hook
	Events []Event

	Traits []Trait
}

func (c *Card) GetID() string {
	return c.ID
}

func (c *Card) GetCost() int {
	return c.Cost
}

func (c *Card) GetEnabled() bool {
	return c.Enabled
}

type ItemCard struct {
	Card `yaml:",inline" mapstructure:",squash"`

	HeldTraits []Trait
}

type SpellCard struct {
	Card `yaml:",inline" mapstructure:",squash"`
}

type UnitCard struct {
	Card `yaml:",inline" mapstructure:",squash"`

	Type       string
	DamageType string
	Attack     int
	Health     int
	Cooldown   int
	Movement   int
	Codex      string
}

type Condition struct {
	Type string
	Not  bool
	Args interface{}
}

type Choose struct {
	Type string
	Args interface{}
}

type Hook struct {
	When            string
	Types           []string
	Conditions      []Condition
	Events          []Event
	ReuseConditions []Condition
}

type Event struct {
	Type string
	Args interface{}
}

type Trait struct {
	Type string
	Args interface{}
}
