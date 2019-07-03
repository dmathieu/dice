package kubernetes

// FlagValue allows configuring the value for the node flag
func FlagValue(flag string) func() error {
	return func() error {
		flagValue = flag
		return nil
	}
}

// Setup allows setting global configuration on the package
func Setup(opts ...func() error) error {
	for _, option := range opts {
		err := option()
		if err != nil {
			return err
		}
	}

	return nil
}
