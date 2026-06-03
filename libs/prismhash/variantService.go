package prismhash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Dan-Sones/prismhash/model"
)

func GetVariantForExperiment(experiments model.ExperimentWithVariants, userId string) (string, error) {
	numberLinePosition := GetNumberLinePositionForUserAndExperiment(userId, experiments.ExperimentKey, experiments.UniqueSalt)

	for _, variant := range experiments.Variants {
		if numberLinePosition >= variant.LowerBound && numberLinePosition <= variant.UpperBound {
			return variant.VariantKey, nil
		}
	}
	return "", fmt.Errorf("no variant found for user %s in experiment %s with number line position %d", userId, experiments.ExperimentKey, numberLinePosition)
}

func GetNumberLinePositionForUserAndExperiment(userId, experimentKey, uniqueSalt string) int32 {
	toHash := fmt.Sprintf("%s:%s:%s", userId, experimentKey, uniqueSalt)
	hash := md5.Sum([]byte(toHash))

	hashHex := hex.EncodeToString(hash[:])

	hashInt := new(big.Int)
	hashInt.SetString(hashHex, 16)

	numberLinePosition := new(big.Int)
	return int32(numberLinePosition.Mod(hashInt, big.NewInt(100)).Int64())
}
