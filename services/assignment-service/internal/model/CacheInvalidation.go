package model

type CacheInvalidation struct {
	BucketToInvalidate int32 `json:"bucket_to_invalidate"`
}
