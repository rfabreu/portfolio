package chat

import (
	"strings"
	"testing"
)

func TestSystemPromptIncludesResumeContent(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "Rafael Abreu") {
		t.Error("SystemPrompt should include resume content (e.g., the canonical name)")
	}
	if !strings.Contains(p, "## RESUME") {
		t.Error("SystemPrompt should include the RESUME section header")
	}
}

func TestSystemPromptIncludesAllProjects(t *testing.T) {
	p := SystemPrompt()
	if !strings.Contains(p, "## PROJECT CASE STUDIES") {
		t.Error("SystemPrompt should include the PROJECT CASE STUDIES section header")
	}
	// Each project markdown's frontmatter title should appear.
	// Read embedded directory to know the count.
	entries, err := projectsFS.ReadDir("embedded/projects")
	if err != nil {
		t.Fatalf("failed to read embedded projects: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("no embedded projects found — sync step must run before tests")
	}
	// Spot-check: expect at least one project's content to be in the prompt.
	for _, entry := range entries {
		content, err := projectsFS.ReadFile("embedded/projects/" + entry.Name())
		if err != nil {
			t.Errorf("failed to read %s: %v", entry.Name(), err)
			continue
		}
		// First 100 chars should appear if the file was concatenated.
		snippet := string(content)
		if len(snippet) > 100 {
			snippet = snippet[:100]
		}
		if !strings.Contains(p, snippet) {
			t.Errorf("project file %s content not found in SystemPrompt", entry.Name())
		}
	}
}

func TestSystemPromptOrderingResumeBeforeProjects(t *testing.T) {
	p := SystemPrompt()
	resumeIdx := strings.Index(p, "## RESUME")
	projectsIdx := strings.Index(p, "## PROJECT CASE STUDIES")
	if resumeIdx == -1 || projectsIdx == -1 {
		t.Fatal("required section headers missing")
	}
	if !(resumeIdx < projectsIdx) {
		t.Error("RESUME section must precede PROJECT CASE STUDIES")
	}
}

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
