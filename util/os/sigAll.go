package os

// SIGALL represents all signals and is not an actual os signal
const SIGALL = SigAll("SIGALL")

type SigAll string

func (SigAll) Signal() {}

func (s SigAll) String() string {
	return string(s)
}
