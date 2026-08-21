package domain

import "context"

type CacheTask struct {
	Key            string
	Value          string
	RequestContext context.Context
}
