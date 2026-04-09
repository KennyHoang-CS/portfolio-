package gorules

type Toolchain struct {
	Version   string
	Installed bool
}

var DefaultToolchain = &Toolchain{
	Version:   "go1.22",
	Installed: true,
}

func (t *Toolchain) Check() error {
	if !t.Installed {
		return ErrToolchainMissing
	}
	return nil
}