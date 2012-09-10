package githubPagesPublish

import (
	"fmt"
	"os"
	"os/exec"
	"io/ioutil"
	"regexp"
	"strings"
)

type Publisher struct {
	Uri    string
	Branch string
	Path   string
}

func (p Publisher) gitRun(args []string) error {
  cmd := exec.Command("git",args...)
  cmd.Stdout = os.Stdout
  cmd.Stdout = os.Stdout
  cmd.Dir = p.Path
  fmt.Fprintf(os.Stderr, "Running `git %s`\n",strings.Join(args, " "))
  err := cmd.Start()
  if err != nil {
	  fmt.Fprintf(os.Stderr, "err %v\n", err)
	  return err
  }
  fmt.Fprintf(os.Stderr, "Waiting for command to finish...\n")
  err = cmd.Wait()
  fmt.Fprintf(os.Stderr, "Command finished with error: %v\n", err)
  return err
}

func New(uri, branch string) (*Publisher, error) {
	var p = &Publisher{
		uri,
		branch,
		"",
	}
	p.Path, _ = p.pathName()

	args := []string{
	  "clone",
	  "--depth", "1",
	  "-b", p.Branch,
	  p.Uri,
	  p.Path,
	}
	fmt.Fprintf(os.Stderr, "git clone: %s\n", strings.Join(args, " "))
	err := p.gitRun(args)
	if err != nil {
	  return nil, err
	}

	args = []string{
	  "checkout",
	  p.Branch,
	}
	err = p.gitRun(args)
	if err != nil {
	  // Create it!
	  	args = []string{
		  "checkout",
		  "--orphan",
		  p.Branch,
		}
		err = p.gitRun(args)
		if err != nil {
		  return nil, err
		}
	}

	args = []string{
	  "branch",
	}
	err = p.gitRun(args)
	return p, nil
}

func (p Publisher) pathName() (name string, err error) {
	reg, _ := regexp.Compile("[^\\w.]")
	s := []string{
		"Publisher",
		string(reg.ReplaceAll([]byte(p.Uri), []byte("_"))),
		p.Branch,
		fmt.Sprintf("%d", os.Getpid()),
		"",
	}
	return ioutil.TempDir("",strings.Join(s, "__")) 
}

func (p Publisher) Push() error {
  args := []string{
	"add",
	".",
  }
  err := p.gitRun(args)
  if err != nil {
	return err
  }

  args = []string{
	"commit",
	"-a",
	"-m",
	"Automatic commit from Go!",
  }
  err = p.gitRun(args)
  if err != nil {
	return err
  }

  args = []string{
	"push",
	"origin",
	p.Branch,
  }
  return p.gitRun(args)
}

func (p Publisher) Close() error {
  fmt.Fprintf(os.Stderr, "Removing path %s\n", p.Path)
  return os.RemoveAll(p.Path)
}
