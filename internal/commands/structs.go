package commands

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type locationConfig struct {
	Previous string
	Next     string
}
