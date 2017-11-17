package service

// remote settings
const (
	ClbtModeNormal        Mode = iota
	ClbtModeOnlyRegister       // collaborator only registers service, service is not accessible until it has been changed to ClbtModeNormal
	ClbtModeOnlySubscribe      // subscriber only subscribes to collaborator service at coordinator, no service redirection will be provided
	LBModeRandom               // assign tasks as per weighted probability
	LBModeLeastActive          // assign tasks to active responders
	LBModeRoundRobin           // assign tasks sequentially based on the order of collaborator
	LBModeTokenHash            // assign tasks based on the hash value of subscriber token
)

type Mode int

func (m *Mode) GetMode() Mode {
	return *m
}
