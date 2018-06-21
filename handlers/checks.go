package handlers

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gojp/goreportcard/check"
	"github.com/gojp/goreportcard/download"
	"os"
)

type notFoundError struct {
	repo string
}

func (n notFoundError) Error() string {
	return fmt.Sprintf("%q not found in cache", n.repo)
}

func dirName(repo string) string {
	return fmt.Sprintf("_repos/src/%s", repo)
}

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Weight        float64             `json:"weight"`
	Percentage    float64             `json:"percentage"`
	Error         string              `json:"error"`
}

type checksResp struct {
	Checks               []score   `json:"checks"`
	Average              float64   `json:"average"`
	Grade                Grade     `json:"grade"`
	Files                int       `json:"files"`
	Issues               int       `json:"issues"`
	Repo                 string    `json:"repo"`
	ResolvedRepo         string    `json:"resolvedRepo"`
	LastRefresh          time.Time `json:"last_refresh"`
	LastRefreshFormatted string    `json:"formatted_last_refresh"`
	LastRefreshHumanized string    `json:"humanized_last_refresh"`
}

func newChecksRespPrivate(repo string) (checksResp, error) {
	// fetch the repo and grade it
	repoRoot, err := download.Download(repo, "_repos/src")
	if err != nil {
		return checksResp{}, fmt.Errorf("could not clone repo: %v", err)
	}

	repo = repoRoot.Root

	dir := dirName(repo)
	defer os.RemoveAll("_repos")

	resp, err := checkResp(dir)
	resp.Repo = repo //allows abstracting checkResp for local and remote checks, yet leaving code largely same
	resp.ResolvedRepo = repoRoot.Root

	return resp, nil
}

func checkResp(dir string) (checksResp, error) {
	filenames, skipped, err := check.GoFiles(dir)
	if err != nil {
		return checksResp{}, fmt.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return checksResp{}, fmt.Errorf("no .go files found")
	}

	err = check.RenameFiles(skipped)
	if err != nil {
		log.Println("Could not remove files:", err)
	}
	defer check.RevertFiles(skipped)

	//for _, fn := range filenames {
	//	fmt.Println(fn)
	//}

	checks := []check.Check{
		check.GoFmt{Dir: dir, Filenames: filenames},
		check.GoVet{Dir: dir, Filenames: filenames},
		check.GoLint{Dir: dir, Filenames: filenames},
		check.GoCyclo{Dir: dir, Filenames: filenames},
		check.License{Dir: dir, Filenames: []string{}},
		check.Misspell{Dir: dir, Filenames: filenames},
		check.IneffAssign{Dir: dir, Filenames: filenames},
		// check.ErrCheck{Dir: dir, Filenames: filenames}, // disable errcheck for now, too slow and not finalized
	}

	ch := make(chan score)
	for _, c := range checks {
		go func(c check.Check) {
			p, summaries, err := c.Percentage()
			errMsg := ""
			if err != nil {
				log.Printf("ERROR: (%s) %v", c.Name(), err)
				errMsg = err.Error()
			}
			s := score{
				Name:          c.Name(),
				Description:   c.Description(),
				FileSummaries: summaries,
				Weight:        c.Weight(),
				Percentage:    p,
				Error:         errMsg,
			}
			ch <- s
		}(c)
	}

	t := time.Now().UTC()
	resp := checksResp{
		Repo:                 dir,
		Files:                len(filenames),
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	var total, totalWeight float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		resp.Checks = append(resp.Checks, s)
		total += s.Percentage * s.Weight
		totalWeight += s.Weight
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}

	total /= totalWeight

	sort.Sort(ByWeight(resp.Checks))
	resp.Average = total
	resp.Issues = len(issues)
	resp.Grade = grade(total * 100)

	return resp, err
}

// ByWeight implements sorting for checks by weight descending
type ByWeight []score

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
