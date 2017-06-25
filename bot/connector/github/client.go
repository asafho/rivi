package github

import (
	"context"
	"github.com/bivas/rivi/bot"
	"github.com/bivas/rivi/util"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ghClient struct {
	client *github.Client
	secret []byte
	owner  string
	repo   string
}

func (c *ghClient) GetState(issue int) string {
	response, _, err := c.client.Issues.Get(context.Background(), c.owner, c.repo, issue)
	if err != nil {
		util.Logger.Error("Unable to get issue %d state. %s", issue, err)
		return ""
	}
	return *response.State
}

func handleLabelsResult(labels []*github.Label, err error, logError func(error)) []string {
	result := make([]string, 0)
	if err != nil {
		logError(err)
	} else {
		for _, p := range labels {
			result = append(result, *p.Name)
		}
	}
	return result
}

func (c *ghClient) GetAvailableLabels() []string {
	util.Logger.Debug("Getting available labels")
	labels, _, e := c.client.Issues.ListLabels(context.Background(), c.owner, c.repo, nil)
	return handleLabelsResult(labels, e, func(err error) {
		util.Logger.Error("Unable to get available labels. %s", err)
	})
}

func (c *ghClient) GetLabels(issue int) []string {
	util.Logger.Debug("Getting labels for issue %d", issue)
	labels, _, err := c.client.Issues.ListLabelsByIssue(context.Background(), c.owner, c.repo, issue, nil)
	return handleLabelsResult(labels, err, func(err error) {
		util.Logger.Error("Unable to get labels for issue %d. %s", issue, err)
	})
}

func (c *ghClient) AddLabel(issue int, label string) []string {
	util.Logger.Debug("Adding label '%s' to issue %d", label, issue)
	labels, _, err := c.client.Issues.AddLabelsToIssue(context.Background(), c.owner, c.repo, issue, []string{label})
	return handleLabelsResult(labels, err, func(err error) {
		util.Logger.Error("Unable to add label %s for issue %d. %s", label, issue, err)
	})
}

func (c *ghClient) RemoveLabel(issue int, label string) []string {
	util.Logger.Debug("Removing label '%s' from issue %d", label, issue)
	c.client.Issues.RemoveLabelForIssue(context.Background(), c.owner, c.repo, issue, label)
	return c.GetLabels(issue)
}

func (c *ghClient) GetAssignees(issue int) []string {
	util.Logger.Debug("Getting assignees for issue %d", issue)
	result := make([]string, 0)
	users, _, err := c.client.Issues.ListAssignees(context.Background(), c.owner, c.repo, nil)
	if err != nil {
		util.Logger.Error("Unable to get assignees for issue %d. %s", issue, err)
	} else {
		for _, p := range users {
			result = append(result, *p.Login)
		}
	}
	return result
}

func (c *ghClient) AddAssignees(issue int, assignees ...string) []string {
	util.Logger.Debug("Adding assignees '%s' for issue %d", assignees, issue)
	response, _, err := c.client.Issues.AddAssignees(context.Background(), c.owner, c.repo, issue, assignees)
	result := make([]string, 0)
	if err != nil {
		util.Logger.Error("Unable to add assignees %s for issue %d. %s", assignees, issue, err)
	} else {
		for _, p := range response.Assignees {
			result = append(result, *p.Login)
		}
	}
	return result

}

func (c *ghClient) RemoveAssignees(issue int, assignees ...string) []string {
	util.Logger.Debug("Removing assignees '%s' from issue %d", assignees, issue)
	response, _, err := c.client.Issues.RemoveAssignees(context.Background(), c.owner, c.repo, issue, assignees)
	result := make([]string, 0)
	if err != nil {
		util.Logger.Error("Unable to remove assignees %s for issue %d. %s", assignees, issue, err)
	} else {
		for _, p := range response.Assignees {
			result = append(result, *p.Login)
		}
	}
	return result
}

func (c *ghClient) GetFileNames(issue int) []string {
	files, _, err := c.client.PullRequests.ListFiles(context.Background(), c.owner, c.repo, issue, nil)
	result := make([]string, 0)
	if err != nil {
		util.Logger.Error("Unable to get file names for issue %d. %s", issue, err)
	} else {
		for _, p := range files {
			result = append(result, *p.Filename)
		}
	}
	return result
}

func (c *ghClient) GetComments(issue int) []string {
	comments, _, err := c.client.Issues.ListComments(context.Background(), c.owner, c.repo, issue, nil)
	result := make([]string, 0)
	if err != nil {
		util.Logger.Error("Unable to get comments for issue %d. %s", issue, err)
	} else {
		for _, p := range comments {
			result = append(result, *p.Body)
		}
	}
	return result
}

func (c *ghClient) AddComment(issue int, comment string) {
	commentObject := &github.IssueComment{Body: github.String(comment)}
	_, _, err := c.client.Issues.CreateComment(context.Background(), c.owner, c.repo, issue, commentObject)
	if err != nil {
		util.Logger.Error("Unable to add comment for issue %d. %s", issue, err)
	}
}

func newClient(config bot.ClientConfig, owner, repo string) *ghClient {
	source := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GetOAuthToken()},
	)
	c := oauth2.NewClient(context.Background(), source)

	return &ghClient{
		client: github.NewClient(c),
		secret: []byte(config.GetSecret()),
		owner:  owner,
		repo:   repo,
	}
}