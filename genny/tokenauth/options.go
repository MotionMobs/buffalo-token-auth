package tokenauth

// Options for generating token-auth
type Options struct {
	// add your stuff here
	Prefix     string
	UserFields []string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	return nil
}
