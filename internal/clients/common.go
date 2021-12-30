package clients

import (
	"fmt"
)

var (
	LeaderNotFoundError = fmt.Errorf("service leader not found")
)
