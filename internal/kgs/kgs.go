package kgs

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

type SnowflakeGen struct {
	node *snowflake.Node
}

func New(nodeID int64) (*SnowflakeGen, error) {
	node, err := snowflake.NewNode(nodeID)

	if err != nil {
		return nil, fmt.Errorf("failed to create new snowflake node: %w", err)
	}

	return &SnowflakeGen{
		node: node,
	}, nil
}

func (s *SnowflakeGen) Generate() uint64 {
	return uint64(s.node.Generate())
}
