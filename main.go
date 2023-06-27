package main

import (
	"context"
	"log"
	"time"

	"github.com/samber/lo"
	"pgmatcher/internal/matcher"
	"pgmatcher/internal/model"
	"pgmatcher/internal/repository/postgresql"
)

func main() {
	ctx := context.Background()
	log.Println("init db connection")

	rep, err := postgresql.New(ctx, "postgres://user:password@localhost:5432/db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("close connection")
		rep.Close()
	}()

	// init data 100 0000 rules
	const (
		N         = 100
		batchSize = 1000
	)
	for i := 0; i < batchSize; i++ {
		rules := lo.RepeatBy[*model.Rule](N, func(index int) *model.Rule {
			return model.GenerateRule()
		})

		if err := rep.Store(ctx, rules); err != nil {
			log.Fatal(err)
		}
	}

	conditions, err := rep.GetAllConditions(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("all conditions:", len(conditions))

	m, err := matcher.New(ctx, rep)
	if err != nil {
		log.Fatal(err)
	}

	asset := model.GenerateAsset()
	log.Println("asset:", asset)

	now := time.Now()

	rules, err := m.Match(ctx, asset)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rules:", rules)
	log.Println(time.Since(now))
}
