package jsparser

type stateList []LexicalState

func (l *stateList) pushMode(t *Tokeniser, s LexicalState) {
	t.state = s
	*l = append(*l, s)
}

func (l *stateList) popMode(t *Tokeniser, d LexicalState) LexicalState {
	p := t.state

	s := d
	if len(*l) > 0 {
		*l = (*l)[0 : len(*l)-1]

		if len(*l) > 0 {
			s = (*l)[len(*l)-1]
		}
	}

	t.state = s

	return p
}
