package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/gomarkdown/markdown"

	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/google/go-github/v57/github"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func createHTMLFile(filename string, data map[string]string) error {
    
    
    // Create a template context
    tmpl := template.Must(template.ParseFiles("template.html"))

// Create the file and get the file pointer
file, err := os.Create(filename)
if err != nil {
	return err
}
defer file.Close()

// Execute the template with the file pointer and data
err = tmpl.Execute(file, data)
if err != nil {
	return err
}

return nil
}


func main() {
	ctx := context.Background()

	now := time.Now()
	client := github.NewClient(nil)
	oneWeekAgo := now.AddDate(0, 0, -7)

	
	var htmlContent string

	weekfile := oneWeekAgo.Format("2006-01-02") + ".html"
	
	releases, _, err := client.Repositories.ListReleases(ctx, "coredns", "coredns", &github.ListOptions{
		PerPage: 5,
	})

	htmlContent += "<h2>Latest Releases for " + oneWeekAgo.Format("2006-01-02") + "</h2>"

	for _, release := range releases {
		if release.GetPublishedAt().Before(oneWeekAgo) {

			htmlContent += "<h2>CoreDNS/CoreDNS</h2> "
			htmlContent += "<h3>Release notes for " + release.GetName() + "</h3>"
			htmlContent += "<h4>" + release.GetPublishedAt().Format("2006-01-02") + "</h4>"

			content := []byte(release.GetBody())

			htmlContent += string(mdToHTML(content))
			htmlContent += "<br>"

		}

	}

	


	data := make(map[string]string)
	data["content"] = htmlContent
	data["title"] = "Whats new in the last week"
	data["date"] = now.Format("2006-01-02")

	err = createHTMLFile(weekfile, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}
