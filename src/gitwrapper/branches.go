package gitwrapper

import (
	"errors"
	git "github.com/libgit2/git2go"
)

// Create a new branch from the current head with the specified name.
// Returns the branch
func CreateBranch(name string, repo *git.Repository) (*git.Branch, error) {
	commit, err := getCommit("HEAD", repo)
	if err != nil {
		return nil, err
	}

	branch, err := repo.CreateBranch(name, commit, false)
	if err != nil {
		return nil, err
	}

	return branch, nil
}

// Checks out the given branch by name
// name - Plain Text branch name (e.g. 'master')
// repo - Repo to checkout from
func CheckoutBranch(name string, repo *git.Repository) error {
	err := checkout(name, repo)
	if err != nil {
		return err
	}

	err = moveHead(name, repo)
	return err
}

// Moves the head to the given branch
// Files are not changed
func moveHead(name string, repo *git.Repository) error {
	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {
		return err
	}

	err = repo.SetHead(branch.Reference.Name())
	return err
}

// Checks out the given branch without moving head
// Doesn't change current branch tag
func checkout(name string, repo *git.Repository) error {
	commit, err := getCommit(name, repo)
	if err != nil {
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	checkoutOps := git.CheckoutOpts{}
	checkoutOps.Strategy = git.CheckoutSafe
	err = repo.CheckoutTree(tree, &checkoutOps)

	return err
}

func BranchExists(name string, repo *git.Repository) bool {
	_, err := getCommit(name, repo)
	return err == nil
}

func DeleteBranch(name string, repo *git.Repository) error {
	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {
		return err
	}
	err = branch.Delete()
	if err != nil {
		return err
	}

	return nil
}

func CurrentBranchName(repo *git.Repository) (string, error) {
	iterator, err := repo.NewBranchIterator(git.BranchLocal)
	if err != nil {
		return "", err
	}
	for branch, _, err := iterator.Next(); err == nil; branch, _, err = iterator.Next() {
		head, err := branch.IsHead()
		if err != nil {
			return "", err
		}
		if head {
			return branch.Name()
		}
	}
	return "", errors.New("Could not find current line.")
}

// Merge the head of the named branch into the current branch head.
func Merge(name string, repo *git.Repository) (bool, error) {
	otherHead, err := getCommit(name, repo)
	if err != nil {
		return false, err
	}
	annOther, err := repo.LookupAnnotatedCommit(otherHead.Id())
	if err != nil {
		return false, err
	}
	sources := []*git.AnnotatedCommit{annOther}

	analysis, _, err := repo.MergeAnalysis(sources)
	if err != nil {
		return false, err
	}
	if analysis&git.MergeAnalysisNone != 0 || analysis&git.MergeAnalysisUpToDate != 0 {
		return false, errors.New("Nothing to absorb")
	}
	// TODO: This probably doesn't support fast-forward
	if analysis&git.MergeAnalysisNormal == 0 {
		return false, errors.New("Non-normal absorb")
	}

	mergeOptions, err := git.DefaultMergeOptions()
	if err != nil {
		return false, err
	}
	checkoutOptions := git.CheckoutOpts{
		Strategy: git.CheckoutForce | git.CheckoutAllowConflicts,
	}

	err = repo.Merge(sources, &mergeOptions, &checkoutOptions)
	if err != nil {
		return false, err
	}

	index, err := repo.Index()
	if err != nil {
		return false, err
	}

	if index.HasConflicts() {
		// Make Git let us commit.
		//index.CleanupConflicts()
		//err = Commit(repo, "Began absorbing "+name+" with conflicts\n\n#CONFLICTS", "HEAD^{commit}", name+"^{commit}")
		//if err != nil {
		//	return false, err
		//}
		return true, nil
	} else {
		err = Commit(repo, "Absorbed "+name, "HEAD^{commit}", name+"^{commit}")
		if err != nil {
			return false, err
		}
		return false, nil
	}
}
