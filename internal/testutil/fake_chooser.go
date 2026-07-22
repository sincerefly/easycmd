package testutil

type FakeChooser struct {
	Fixed string
	Err   error
}

func (f FakeChooser) Choice(_ []string) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.Fixed, nil
}
