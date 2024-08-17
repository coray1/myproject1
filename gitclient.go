package githubsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/google/go-github/v63/github"
	"github.com/madneal/gshark/model"
	"net/http"

	"github.com/madneal/gshark/service"
)

var (
	GithubClients map[string]*Client
	GithubClient  *Client
)

type Client struct {
	Client *github.Client
	Token  string
}

func InitGithubClients(tokens []model.Token) map[string]*Client {
	githubClients := make(map[string]*Client)
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: httpTransport}
	for _, token := range tokens {
		githubToken := token.Content
		if githubToken != "" {
			gitClient := github.NewClient(httpClient).WithAuthToken(githubToken)
			githubClients[token.Content] = NewGitClient(gitClient, githubToken)
		}
	}
	return githubClients
}

func GetGithubClient() (*Client, error) {
	var c *Client
	err, tokens := service.ListTokenByType("github")
	if err != nil {
		return c, err
	}
	clients := InitGithubClients(tokens)

	for _, client := range clients {
		c = client
		break
	}
	if c == nil {
		err = errors.New("github Client initial failed, please add token")
	}
	return c, err
}

func NewGitClient(GithubClient *github.Client, token string) *Client {
	return &Client{Client: GithubClient, Token: token}
}

func (c *Client) GetUserInfo(username string) (*github.User, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Users.Get(ctx, username)
}

func (c *Client) GetOrgsMembers(org string) ([]*github.User, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Organizations.ListMembers(ctx, org, nil)
}

func (c *Client) GetOrgsRepos(org string) ([]*github.Repository, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Repositories.ListByOrg(ctx, org, nil)
}

func (c *Client) GetUserRepos(username string) ([]*github.Repository, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Repositories.List(ctx, username, nil)
}

//func (c *Client) GetUsersRepos(users []*github.User) map[string][]*github.Repository {
//	result := make(map[string][]*github.Repository)
//	for _, u := range users {
//		repos, resp, _ := c.GetUserRepos(*u.Login)
//		model.UpdateRate(c.Token, resp)
//		result[*u.Login] = repos
//	}
//	return result
//}
//
//func (c *Client) GetStrUsersRepos(users []string) map[string][]*github.Repository {
//	result := make(map[string][]*github.Repository)
//	for _, u := range users {
//		repos, resp, _ := c.GetUserRepos(u)
//		model.UpdateRate(c.Token, resp)
//		result[u] = repos
//	}
//	return result
//}

func (c *Client) GetUserOrgs(username string) ([]*github.Organization, *github.Response, error) {
	ctx := context.Background()
	return c.Client.Organizations.List(ctx, username, nil)
}
