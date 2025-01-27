package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/owenrumney/go-commie/internal/logger"
	"github.com/owenrumney/go-commie/ui"
)

var statusMap = map[git.StatusCode]string{
	git.Unmodified: " ",
	git.Modified:   "modified:",
	git.Added:      "added:   ",
	git.Deleted:    "deleted: ",
	git.Renamed:    "renamed: ",
	git.Copied:     "copied:  ",
	git.Untracked:  "new:     ",
}

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

	wt, err := g.Repository.Worktree()
	if err != nil {
		g.log.Fatal(err)
	}

	status, err := wt.Status()
	if err != nil {
		g.log.Fatal(err)
	}

	stagedFiles, unstagedFiles, untrackedFiles := g.listFiles(status)

	if len(stagedFiles) > 0 {
		fmt.Printf("\nStaged files:\n  %s", strings.Join(stagedFiles, "\n  "))
	}
	if len(unstagedFiles) > 0 {
		fmt.Printf("\nUnstaged files:\n  %s", strings.Join(unstagedFiles, "\n  "))
	}
	if len(untrackedFiles) > 0 {
		fmt.Printf("\nUntracked files:\n  %s", strings.Join(untrackedFiles, "\n  "))
	}

	fmt.Println("\n\n")

	if len(stagedFiles) > 0 {
		return g.addCommit()
	}

	fmt.Println("No files staged for commit")
	return nil
}

func (g *Git) listFiles(status git.Status) (stagedFiles, unstagedFiles, untrackedFiles []string) {

	for path, st := range status {
		if st.Staging != git.Unmodified {
			file := "\033[32m" + statusMap[st.Staging] + " " + path
			if st.Staging == git.Renamed {
				file = file + " -> " + st.Extra
			}
			file = file + "\033[0m"
			stagedFiles = append(stagedFiles, file)
		}

		if st.Worktree == git.Untracked {
			file := "\033[33m" + statusMap[st.Worktree] + " " + path + "\033[0m"
			untrackedFiles = append(untrackedFiles, file)
		} else if st.Worktree != git.Unmodified {
			file := "\033[31m" + statusMap[st.Worktree] + " " + path
			if st.Worktree == git.Renamed {
				file = file + " -> " + st.Extra
			}
			file = file + "\033[0m"
			unstagedFiles = append(unstagedFiles, file)
		}
	}

	return stagedFiles, unstagedFiles, untrackedFiles
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
	return ui.GetMultilineInput("Enter commit body")
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

		if len(titleParts) > 1 {
			title := titleParts[1]
			suggestedTitle = fmt.Sprintf("%s: %s", titleParts[0],
				strings.ToUpper(string(title[0]))+strings.ToLower(title[1:]))
		}

	}

	useSuggested, err := ui.YesNoQuestion(fmt.Sprintf(`Use suggested title "%s"`, suggestedTitle), true)
	if err != nil {
		g.log.Fatal(err)
	}

	if useSuggested {
		return suggestedTitle
	}

	_, prefix, err := ui.ChooseFromList("Choose the appropriate prefix", []string{
		"feat",
		"fix",
		"docs",
		"style",
		"refactor",
		"test",
	})
	if err != nil {
		g.log.Fatal(err)
	}

	title := ui.GetInput("Enter commit title")
	return fmt.Sprintf("%s: %s", prefix, title)
}

func (g *Git) getBranchName() (string, error) {
	head, err := g.Repository.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}
