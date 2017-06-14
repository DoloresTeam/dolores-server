package main

import "testing"

func TestGenerateImage(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	generatorAvatar(`id`, `聪灵`)
}
