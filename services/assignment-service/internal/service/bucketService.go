package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
)

type BucketService struct {
	salt        string
	bucketCount int64
	// TODO: we need int32 here I think
}

func NewBucketService(salt string, bucketCount int64) *BucketService {
	return &BucketService{
		salt:        salt,
		bucketCount: bucketCount,
	}
}

func (s *BucketService) GetBucketFor(userId string) int32 {
	hash := s.createMD5For(userId)
	hashHex := hex.EncodeToString(hash[:])

	hashInt := new(big.Int)
	hashInt.SetString(hashHex, 16)

	bucket := new(big.Int)
	bucket.Mod(hashInt, big.NewInt(s.bucketCount))

	return int32(bucket.Int64())
}

func (s *BucketService) createMD5For(userId string) [16]byte {
	toHash := fmt.Sprintf("%s:%s", s.salt, userId)
	return md5.Sum([]byte(toHash))
}
