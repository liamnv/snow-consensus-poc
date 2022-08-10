package consensus

type Preference int

type SnowBall struct {
	preference Preference

	decisionThreshold int

	consecutiveSuccesses map[Preference]int

	decided bool
}

func NewSnowBall(preference Preference, decisionThreshold int) SnowBall {
	return SnowBall{
		preference:           preference,
		consecutiveSuccesses: make(map[Preference]int),
		decided:              false,
		decisionThreshold:    decisionThreshold,
	}
}

func (s *SnowBall) Preference() Preference {
	return s.preference
}

func (s *SnowBall) Decided() bool {
	return s.decided
}
func (s *SnowBall) SuccessPool(preference Preference) {
	if s.decided {
		return
	}
	s.consecutiveSuccesses[preference] = s.consecutiveSuccesses[preference] + 1
	s.preference = preference

	if s.consecutiveSuccesses[preference] > s.decisionThreshold {
		s.decided = true
	}
}
