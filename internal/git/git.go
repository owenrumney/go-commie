package git

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/owenrumney/go-commie/internal/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Git struct {
	*git.Repository
	log *logger.Log
}

func New(log *logger.Log) (*Git, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, err
	}

	return &Git{repo, log}, nil
}

func (g *Git) Commit() error {
	return g.addCommit()
}

func (g *Git) addCommit() error {
	worktree, err := g.Repository.Worktree()
	if err != nil {
		return err
	}

	title := g.getCommitMsgTitle()
	body := g.getCommitBody()

	commitMsg := title + "\n\n" + body

	g.log.Debugf("Commit message: %s", commitMsg)

	commit, err := worktree.Commit(commitMsg, &git.CommitOptions{})
	if err != nil {
		return err
	}

	obj, err := g.Repository.CommitObject(commit)
	if err != nil {
		return err
	}

	println(obj.Hash.String())

	return nil
}

func (g *Git) getCommitBody() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n Enter commit body with 2 empty lines to complete to complete: ")
	var lines []string
	consecutiveEmptyLines := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			consecutiveEmptyLines++
		} else {
			consecutiveEmptyLines = 0
		}
		if consecutiveEmptyLines == 2 {
			break
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (g *Git) getCommitMsgTitle() string {
	suggestedTitle, err := g.getBranchName()
	if err != nil {
		g.log.Debugf("Error getting branch name: %s", err)
		suggestedTitle = ""
	}

	if suggestedTitle != "" {
		suggestedTitle = strings.ReplaceAll(suggestedTitle, "-", " ")
		titleParts := strings.Split(suggestedTitle, "/")

		titleCases := cases.Title(language.English)

		if len(titleParts) > 1 {
			suggestedTitle = fmt.Sprintf("%s: %s", titleParts[0], titleCases.String(titleParts[1]))
		}

	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n Enter commit message body, enter to choose selected [%s]: ", suggestedTitle)
	title, _ := reader.ReadString('\n')

	title = strings.TrimSpace(title)
	if title == "" {
		title = suggestedTitle
	}

	return title
}

func (g *Git) getBranchName() (string, error) {
	head, err := g.Repository.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}
