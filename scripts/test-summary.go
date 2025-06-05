package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type TestEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
}

func main() {
	label := "Unit Test Summary"
	if len(os.Args) > 1 {
		label = os.Args[1]
	}

	passed := 0
	failed := 0
	skipped := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var event TestEvent
		err := json.Unmarshal(scanner.Bytes(), &event)
		if err != nil {
			continue
		}

		switch event.Action {
		case "pass":
			if event.Test != "" {
				passed++
			}
		case "fail":
			if event.Test != "" {
				failed++
			}
		case "skip":
			if event.Test != "" {
				skipped++
			}
		}
	}

	fmt.Printf("📊 %s:\n", label)
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("⚠️ Skipped: %d\n", skipped)
}
