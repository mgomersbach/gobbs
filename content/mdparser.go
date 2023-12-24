package content

import (
	"strings"

	"github.com/russross/blackfriday/v2"
)

// MenuItem represents a menu item which links to a chapter.
type MenuItem struct {
	Title string
	Link  string
}

// MDParser holds the parsed menu items and chapters.
type MDParser struct {
	MenuItems []MenuItem
	Chapters  map[string]string
}

// NewMDParser creates a new MDParser instance.
func NewMDParser(mdContent string) *MDParser {
	parser := &MDParser{
		Chapters: make(map[string]string),
	}
	parser.parse(mdContent)
	return parser
}

// parse processes the Markdown content and extracts menu items and chapters.
func (p *MDParser) parse(mdContent string) {
	currentChapter := ""
	chapterContent := strings.Builder{}

	node := blackfriday.New().Parse([]byte(mdContent))
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		switch node.Type {
		case blackfriday.Heading:
			if !entering {
				// Save the previous chapter content before starting a new one.
				if currentChapter != "" {
					p.Chapters[currentChapter] = chapterContent.String()
					chapterContent.Reset()
				}
				currentChapter = extractTextFromNode(node)
			}
		case blackfriday.Link:
			if entering {
				p.MenuItems = append(p.MenuItems, MenuItem{
					Title: extractTextFromNode(node),
					Link:  string(node.LinkData.Destination),
				})
			}
		case blackfriday.Text:
			chapterContent.WriteString(string(node.Literal))
		}
		return blackfriday.GoToNext
	})

	// Save the last chapter.
	if currentChapter != "" {
		p.Chapters[currentChapter] = chapterContent.String()
	}
}

// extractTextFromNode traverses a node and its children to extract text.
func extractTextFromNode(node *blackfriday.Node) string {
	var text strings.Builder
	node.Walk(func(subNode *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if subNode.Type == blackfriday.Text {
			text.WriteString(string(subNode.Literal))
		}
		return blackfriday.GoToNext
	})
	return strings.TrimSpace(text.String())
}

// Example usage:
// mdContent := `# Main Menu
// - [Chapter 1](#chapter-1)
// - [Chapter 2](#chapter-2)
// ...
// `

// parser := content.NewMDParser(mdContent)
// Now you can use parser.MenuItems and parser.Chapters
