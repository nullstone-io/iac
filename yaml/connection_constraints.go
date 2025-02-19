package yaml

import (
	"fmt"
	"strings"
)

type ConnectionConstraints map[string]ConnectionConstraint

func (s ConnectionConstraints) String() string {
	result := make([]string, 0)
	for name, c := range s {
		result = append(result, fmt.Sprintf("%s=%s", name, c))
	}
	return strings.Join(result, ",")
}
