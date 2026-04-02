package metric

import "fmt"

type ComponentRole string

const (
	ComponentRoleBaseEvent   ComponentRole = "base_event"
	ComponentRoleNumerator   ComponentRole = "numerator"
	ComponentRoleDenominator ComponentRole = "denominator"
)

func (c *ComponentRole) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := ComponentRole(s)
	switch dt {
	case ComponentRoleBaseEvent, ComponentRoleNumerator, ComponentRoleDenominator:
		*c = dt
		return nil
	default:
		return fmt.Errorf("invalid ComponentRole: %s", s)
	}
}
