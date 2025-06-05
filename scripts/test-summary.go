package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type TestEvent struct {
	Action string `json:"Action"`
}

func main() {
	var passed, failed, skipped int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var event TestEvent
		if err := json.Unmarshal([]byte(line), &event); err == nil {
			switch event.Action {
			case "pass":
				passed++
			case "fail":
				failed++
			case "skip":
				skipped++
			}
		}
	}

	fmt.Printf("\n📊 Unit Test Summary:\n")
	fmt.Printf("✅ Passed: %d\n", passed)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("⚠️ Skipped: %d\n", skipped)
}
