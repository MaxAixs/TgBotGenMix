package userState

type UserState struct {
	State       string
	TobaccoName string
	Strength    string
	FlavorName  string
}

func (state *UserState) Reset() {
	state.State = ""
	state.TobaccoName = ""
	state.FlavorName = ""
}

func (state *UserState) SetState(newState string) {
	state.State = newState
}

func (state *UserState) IsActive() bool {
	return state.State != ""
}
