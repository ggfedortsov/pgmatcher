package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"
)

var rarityValues = []string{"Rare", "Epic", "Common", "Legendary", "Mythical"}

const (
	minRank     = 1
	maxRank     = 10
	minPosition = 1
	maxPosition = 100
)

func GenerateAsset() Asset {
	return Asset{
		Title:    gofakeit.Name(),
		Rank:     gofakeit.Number(minRank, maxRank),
		Rarity:   gofakeit.RandomString(rarityValues),
		Position: gofakeit.Number(minPosition, maxPosition),
		Price:    gofakeit.Number(1, 1000),
	}
}

func GenerateRule() Rule {
	return Rule{
		ID:         gofakeit.UUID(),
		Price:      gofakeit.Number(1, 1000),
		Priority:   gofakeit.Number(1, 10),
		Conditions: generateConditions(),
	}
}

func generateConditions() []string {
	size := gofakeit.Number(2, 3)
	randomCondition := []string{genRank(), genRarity(), genPosition()}

	return lo.RepeatBy[string](size, func(i int) string {
		return randomCondition[i]
	})
}

func genRank() string {
	return gofakeit.RandomString([]string{
		genSimpleRank(),
		genRankIn(),
		genRankRange(),
	})
}

func genSimpleRank() string {
	op := gofakeit.RandomString([]string{"==", ">", "<"})
	val := gofakeit.Number(minRank, maxRank)

	return fmt.Sprintf("Rank %s %d", op, val)
}

func genRankRange() string {
	return fmt.Sprintf("Rank <= %d && Rank >= %d", gofakeit.Number(minRank, maxRank/2), gofakeit.Number(maxRank/2+1, maxRank))
}

func genRankIn() string {
	ints := lo.Uniq(lo.RepeatBy(gofakeit.Number(2, 10), func(index int) string {
		return fmt.Sprintf("%d", gofakeit.Number(minRank, maxRank))
	}))
	sort.Strings(ints)

	return fmt.Sprintf("Rank in [%s]", strings.Join(ints, ","))
}

func genRarity() string {
	val := gofakeit.RandomString(rarityValues)

	return fmt.Sprintf(`Rarity == "%s"`, val)
}

func genPosition() string {
	op := gofakeit.RandomString([]string{"==", ">", "<"})
	val := gofakeit.Number(minPosition, maxPosition)

	return fmt.Sprintf("Position %s %d", op, val)
}
