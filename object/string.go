package object

type String struct {
	Value string
}

func (s *String) ObjectType() ObjectType { return STRING }
func (s *String) Inspect() string        { return s.Value }
