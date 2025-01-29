package git

import "os/exec"

func (g *Git) signCommit() error {
	// ensure that git is available on the path
	g.log.Debug("Checking for git")

	if _, err := exec.LookPath("git"); err != nil {
		return err
	}

	g.log.Debugf("Git found in PATH, signing commit with key %s", *g.signingKey)

	// sign the commit
	return exec.Command("git", "commit", "--amend", "--no-edit", "--gpg-sign="+*g.signingKey).Run()
}
