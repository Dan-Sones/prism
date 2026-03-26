package model

import (
	"fmt"
	"log"
	"math/rand/v2"
)

type UserIds []int
type VariantUserIds map[string]UserIds

func (vuId *VariantUserIds) SelectRandomUserIdForVariant(variantKey string) int {
	userIdsForVariant := (*vuId)[variantKey]
	randomIndex := rand.IntN(len(userIdsForVariant))
	return userIdsForVariant[randomIndex]
}

func (vuId *VariantUserIds) GetFirstXUserIdsForVariant(variantKey string, x int) []int {
	userIdsForVariant := (*vuId)[variantKey]
	if x > len(userIdsForVariant) {
		log.Fatal(fmt.Sprintf("Not enough user ids for variant %s to get the first %d user ids!!!", variantKey, x))
	}
	return userIdsForVariant[:x]
}
