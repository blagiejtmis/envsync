package env

import "sort"

// CompareResult holds the outcome of comparing two sets of env entries.
type CompareResult struct {
	OnlyInA  []Entry // keys present in A but not B
	OnlyInB  []Entry // keys present in B but not A
	Changed  []EntryPair // keys present in both but with different values
	Identical []Entry // keys present in both with identical values
}

// EntryPair holds corresponding entries from two sources.
type EntryPair struct {
	A Entry
	B Entry
}

// HasDifferences returns true if there is any difference between A and B.
func (r CompareResult) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Changed) > 0
}

// Compare compares two slices of Entry and returns a CompareResult.
// Comparison is key-based; order does not matter.
func Compare(a, b []Entry) CompareResult {
	aMap := make(map[string]string, len(a))
	for _, e := range a {
		aMap[e.Key] = e.Value
	}

	bMap := make(map[string]string, len(b))
	for _, e := range b {
		bMap[e.Key] = e.Value
	}

	var result CompareResult

	for _, e := range a {
		bVal, inB := bMap[e.Key]
		if !inB {
			result.OnlyInA = append(result.OnlyInA, e)
		} else if e.Value == bVal {
			result.Identical = append(result.Identical, e)
		} else {
			result.Changed = append(result.Changed, EntryPair{
				A: e,
				B: Entry{Key: e.Key, Value: bVal},
			})
		}
	}

	for _, e := range b {
		if _, inA := aMap[e.Key]; !inA {
			result.OnlyInB = append(result.OnlyInB, e)
		}
	}

	sort.Slice(result.OnlyInA, func(i, j int) bool { return result.OnlyInA[i].Key < result.OnlyInA[j].Key })
	sort.Slice(result.OnlyInB, func(i, j int) bool { return result.OnlyInB[i].Key < result.OnlyInB[j].Key })
	sort.Slice(result.Changed, func(i, j int) bool { return result.Changed[i].A.Key < result.Changed[j].A.Key })
	sort.Slice(result.Identical, func(i, j int) bool { return result.Identical[i].Key < result.Identical[j].Key })

	return result
}
