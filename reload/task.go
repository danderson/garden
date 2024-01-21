package reload

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Task struct {
	Name    string
	Match   []string
	Ignore  []string
	Command []string
	Notify  func()

	needsRun bool
	matchRe  *regexp.Regexp
	ignoreRe *regexp.Regexp
}

func (t *Task) prepare() {
	t.needsRun = true
	t.matchRe = globsToRe(t.Match)
	t.ignoreRe = globsToRe(t.Ignore)
}

func (t *Task) match(files []string) bool {
	if t.needsRun {
		return true
	}
	for _, f := range files {
		if len(t.Ignore) > 0 && t.ignoreRe.MatchString(f) {
			continue
		}
		if len(t.Match) > 0 && t.matchRe.MatchString(f) {
			t.needsRun = true
			return true
		}
	}
	return false
}

func (t *Task) run() error {
	t.needsRun = false
	log.Printf("Running: %s", strings.Join(t.Command, " "))
	cmd := exec.Command(t.Command[0], t.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func globsToRe(globs []string) *regexp.Regexp {
	var reParts []string
	for _, m := range globs {
		part := globToRe(m)
		reParts = append(reParts, "(?:"+part+")")
	}
	reStr := fmt.Sprintf("(?:%s)$", strings.Join(reParts, "|"))
	return regexp.MustCompile(reStr)
}

func globToRe(glob string) string {
	prefix, suffix := "", ""

	if strings.HasPrefix(glob, "**/") {
		prefix = "(|.*/)"
		glob = glob[3:]
	}
	if strings.HasSuffix(glob, "/**") {
		suffix = "/.*"
		glob = glob[:len(glob)-3]
	}

	aroundDoubleStar := strings.Split(glob, "/**/")
	for i, p := range aroundDoubleStar {
		aroundStar := strings.Split(p, "*")
		for j := range aroundStar {
			aroundStar[j] = regexp.QuoteMeta(aroundStar[j])
		}
		aroundDoubleStar[i] = strings.Join(aroundStar, "[^/]*")
	}

	middle := strings.Join(aroundDoubleStar, "(?:/|/.*/)")

	return prefix + middle + suffix
}
