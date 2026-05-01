package chat

import (
	"strings"
	"testing"
)

func TestSystemPromptIncludesPersona(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "Rafael Abreu's AI assistant") {
		t.Error("SystemPrompt should include the persona marker phrase")
	}
}

func TestSystemPromptIncludesRules(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "Only answer questions about Rafael") {
		t.Error("SystemPrompt should include the rules marker phrase")
	}
}

func TestSystemPromptStructure(t *testing.T) {
	p := SystemPrompt()
	personaIdx := strings.Index(p, "Rafael Abreu's AI assistant")
	rulesIdx := strings.Index(p, "Only answer questions about Rafael")
	if personaIdx == -1 || rulesIdx == -1 {
		t.Fatal("required markers missing")
	}
	if !(personaIdx < rulesIdx) {
		t.Error("persona block must precede rules block")
	}
}
