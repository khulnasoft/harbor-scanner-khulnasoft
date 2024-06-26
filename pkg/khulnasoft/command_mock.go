package khulnasoft

import "github.com/stretchr/testify/mock"

type MockCommand struct {
	mock.Mock
}

func (c *MockCommand) Scan(imageRef ImageRef) (ScanReport, error) {
	args := c.Called(imageRef)
	return args.Get(0).(ScanReport), args.Error(1)
}
