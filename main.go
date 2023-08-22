package main

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"pgmatcher/internal/model"
	"strings"
	"sync/atomic"
	"time"
)

//func main() {
//	rules := lo.RepeatBy[model.Rule](1_000_000, func(index int) model.Rule {
//		return model.GenerateRule()
//	})
//	nowMain := time.Now()
//	uniq := lo.Uniq(lo.Map(rules, func(it model.Rule, _ int) string {
//		return `FILTER AND(` + strings.Join(it.Conditions, `,`) + `)`
//	}))
//	log.Println("uniq:", len(uniq))
//
//	prepareStmts := lo.Map(uniq, func(it string, _ int) *rel.FilterStatement {
//		return rel.MustParseFilter(it)
//	})
//
//	asset := model.GenerateAsset()
//	obj := datasource.NewContextWrapper(asset)
//	log.Println("asset:", asset)
//
//	now := time.Now()
//	chunk := lo.Chunk(prepareStmts, 10000)
//	log.Println("chunk:", len(chunk))
//
//	var ops uint64
//	gr := errgroup.Group{}
//	for i := 0; i < len(chunk); i++ {
//		ch := chunk[i]
//		gr.Go(func() error {
//			for _, ql := range ch {
//
//				m, _ := vm.Matches(obj, ql)
//				if m {
//					atomic.AddUint64(&ops, 1)
//				}
//			}
//
//			return nil
//		})
//	}
//
//	gr.Wait()
//
//	log.Println("count:", ops)
//	log.Println(time.Since(now))
//	log.Println(time.Since(nowMain))
//}

func main() {
	asset := model.GenerateAsset()
	assetMap := asset.ToMap()
	slog.Info("asset", slog.Any("k", assetMap))

	rules := lo.RepeatBy[model.Rule](1_000_000, func(index int) model.Rule {
		return model.GenerateRule()
	})
	nowMain := time.Now()
	uniq := lo.Uniq(lo.Map(rules, func(it model.Rule, _ int) string {
		return strings.Join(it.Conditions, ` && `)
	}))

	prepareStmts := lo.Map(uniq, func(it string, _ int) *vm.Program {
		compile, _ := expr.Compile(it, expr.AsBool(), expr.Optimize(true))
		return compile
	})

	now := time.Now()
	chunk := lo.Chunk(prepareStmts, 10000)
	log.Println("chunk:", len(chunk))

	var ops uint64
	gr := errgroup.Group{}
	for i := 0; i < len(chunk); i++ {
		ch := chunk[i]
		gr.Go(func() error {
			for _, ql := range ch {

				m, _ := expr.Run(ql, asset)
				if m.(bool) {
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
