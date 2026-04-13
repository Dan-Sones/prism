package model

import (
	"errors"
	"fmt"
)

type UserIds []string
type VariantUserIds map[VariantKey]UserIds

// TODO: string at the moment, but when we add in Bucket finder logic to this, this can be more flexible
func (vuId *VariantUserIds) getUserIdsForVariant(variantKey VariantKey, x int) ([]string, error) {
	userIdsForVariant := (*vuId)[variantKey]
	if x > len(userIdsForVariant) {
		return nil, errors.New(fmt.Sprintf("Not enough user ids for variant %s to get the first %d user ids!!!", variantKey, x))
	}
	return userIdsForVariant[:x], nil
}
