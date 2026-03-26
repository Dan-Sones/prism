package model

type ActionPerformer interface {
	PerformAction(request EventRequest)
}
