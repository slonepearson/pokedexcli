package commands

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	Previous string
	Next     string
}
