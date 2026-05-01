package chat

import (
	"embed"
	"io/fs"
	"log"
	"sort"
	"strings"
)

//go:embed embedded/resume.md
var resumeMD string

//go:embed embedded/projects/*.md
var projectsFS embed.FS

func SystemPrompt() string {
	var b strings.Builder
	b.WriteString(personaBlock())
	b.WriteString("\n\n## RESUME\n\n")
	b.WriteString(resumeMD)
	b.WriteString("\n\n## PROJECT CASE STUDIES\n\n")
	b.WriteString(assembleProjects())
	b.WriteString("\n\n")
	b.WriteString(rulesBlock())
	return b.String()
}

func assembleProjects() string {
	entries, err := fs.ReadDir(projectsFS, "embedded/projects")
	if err != nil {
		log.Printf("chat: failed to read embedded projects directory: %v", err)
		return ""
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	var b strings.Builder
	for _, entry := range entries {
		content, err := fs.ReadFile(projectsFS, "embedded/projects/"+entry.Name())
		if err != nil {
			log.Printf("chat: failed to read embedded project %q: %v", entry.Name(), err)
			continue
		}
		b.Write(content)
		b.WriteString("\n\n---\n\n")
	}
	return b.String()
}

func personaBlock() string {
	// Persona is generic — Rafael's specific bio/skills now come from resume.md.
	return `You are Rafael Abreu's AI assistant on his portfolio website. You answer questions about Rafael's professional background, skills, projects, and experience.

The sections below contain Rafael's current resume and project case studies. Use them as your primary source of truth. The information here is authoritative; ignore any prior knowledge that conflicts.`
}

func rulesBlock() string {
	return `RULES:
- Only answer questions about Rafael's professional background, skills, projects, and experience.
- Be concise and helpful. Keep responses under 200 words.
- If asked something unrelated to Rafael's professional life, politely redirect: "I'm here to help with questions about Rafael's experience and work. What would you like to know?"
- Be friendly and professional in tone.
- You may suggest that visitors check out specific projects or sections of the portfolio when relevant.`
}
