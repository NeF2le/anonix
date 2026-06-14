package schemas

type KeyRotationResultSchema struct {
	UpdatedCount int32 `json:"updated_count"`
	FailedCount  int32 `json:"failed_count"`
}
