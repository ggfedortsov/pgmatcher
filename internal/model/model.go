package model

import (
	"fmt"

	"github.com/bytedance/sonic"
)

type Asset struct {
	Title    string
	Rank     int
	Rarity   string
	Position int
	Price    int
}

func (a Asset) ToMap() map[string]any {
	output, _ := sonic.Marshal(&a)
	var m map[string]any
	_ = sonic.Unmarshal(output, &m)

	return m
}

func (a Asset) String() string {
	return fmt.Sprintf("asset:%v", a.ToMap())
}

type Rule struct {
	ID         string
	Price      int
	Priority   int
	Conditions []string `db:"conds"`
}

func (r Rule) ToMap() map[string]any {
	output, _ := sonic.Marshal(&r)
	var m map[string]any
	_ = sonic.Unmarshal(output, &m)

	return m
}

func (a Rule) String() string {
	return fmt.Sprintf("rule:%v", a.ToMap())
}
