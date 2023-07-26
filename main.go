package main

import (
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/araddon/qlbridge/datasource"
	"github.com/araddon/qlbridge/rel"
	"github.com/araddon/qlbridge/vm"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"pgmatcher/internal/model"
)

func main() {
	rules := lo.RepeatBy[model.Rule](1_000_000, func(index int) model.Rule {
		return model.GenerateRule()
	})
	nowMain := time.Now()
	uniq := lo.Uniq(lo.Map(rules, func(it model.Rule, _ int) string {
		return `FILTER AND(` + strings.Join(it.Conditions, `,`) + `)`
	}))
	log.Println("uniq:", len(uniq))

	prepareStmts := lo.Map(uniq, func(it string, _ int) *rel.FilterStatement {
		return rel.MustParseFilter(it)
	})

	asset := model.GenerateAsset()
	obj := datasource.NewContextWrapper(asset)
	log.Println("asset:", asset)

	now := time.Now()
	chunk := lo.Chunk(prepareStmts, 10000)
	log.Println("chunk:", len(chunk))

	var ops uint64
	gr := errgroup.Group{}
	for i := 0; i < len(chunk); i++ {
		ch := chunk[i]
		gr.Go(func() error {
			for _, ql := range ch {

				m, _ := vm.Matches(obj, ql)
				if m {
					atomic.AddUint64(&ops, 1)
				}
			}

			return nil
		})
	}

	gr.Wait()

	log.Println("count:", ops)
	log.Println(time.Since(now))
	log.Println(time.Since(nowMain))
}
