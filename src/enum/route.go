package enum

type RouteStatus string

const (
	RouteStatusActive   RouteStatus = "active"
	RouteStatusInactive RouteStatus = "inactive"
)

type RoutePathType string

const (
	RoutePathTypeExact  RoutePathType = "exact"
	RoutePathTypePrefix RoutePathType = "prefix"
)
