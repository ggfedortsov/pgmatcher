package matcher

import (
	"context"
	"log"
	"time"

	"github.com/araddon/qlbridge/datasource"
	"github.com/araddon/qlbridge/rel"
	"github.com/araddon/qlbridge/vm"
	"pgmatcher/internal/model"
)

type Storage interface {
	GetAllConditions(ctx context.Context) ([]string, error)
	GetAllowRules(ctx context.Context, price int, allowConds []string) ([]model.Rule, error)
}

type Matcher struct {
	prepareQlStatement map[string]*rel.FilterStatement
	storage            Storage
}

func New(ctx context.Context, s Storage) (*Matcher, error) {
	conds, err := s.GetAllConditions(ctx)
	if err != nil {
		return nil, err
	}

	pQlStatements := make(map[string]*rel.FilterStatement)
	for _, cond := range conds {
		pQlStatements[cond] = rel.MustParseFilter("FILTER " + cond)
	}

	return &Matcher{
		storage:            s,
		prepareQlStatement: pQlStatements,
	}, nil
}

func (m *Matcher) Match(ctx context.Context, a model.Asset) ([]model.Rule, error) {
	obj := datasource.NewContextWrapper(a)
	t := time.Now()
	var allowConds []string
	for s, ql := range m.prepareQlStatement {
		if matches, _ := vm.Matches(obj, ql); matches {
			allowConds = append(allowConds, s)
		}
	}
	log.Println(time.Since(t))
	rules, err := m.storage.GetAllowRules(ctx, a.Price, allowConds)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
