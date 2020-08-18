package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"regexp"
	"time"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	releases, _, err := client.Repositories.ListReleases(ctx, "istio", "istio", nil)
	if err != nil {
		log.Fatalf("Failed to list Istio releases: %v", err)
	}

	for _, release := range releases {
		if time.Since(release.GetCreatedAt().Time) > 24*time.Hour {
			log.Printf("Skipping release %q because it is too old\n", release.GetName())
			continue
		}
		match, err := regexp.MatchString("\\d+\\.\\d+\\.\\d+$", release.GetName())
		if err != nil {
			log.Fatalf("Failed to apply regex match: %v", err)
		}
		if !match {
			log.Printf("Skipping release %q because it doesn't match the official release pattern\n", release.GetName())
			continue
		}

		title := fmt.Sprintf("Test Istio release %q", release.GetName())
		body := fmt.Sprintf("Istio recently (%s) released [%q](%s). Let's test it :rocket:.", release.GetCreatedAt(), release.GetName(), release.GetHTMLURL())
		issueReq := &github.IssueRequest{
			Title: &title,
			Body:  &body,
		}
		issue, _, err := client.Issues.Create(ctx, "knative-sandbox", "net-istio", issueReq)
		if err != nil {
			log.Fatalf("Failed to create issue for release %q: %v", release.GetName(), err)
		}
		log.Printf("Created issue %s for release %q\n", issue.GetHTMLURL(), release.GetName())
	}
}
