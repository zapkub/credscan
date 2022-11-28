package reposcan

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/zapkub/credscan/internal"
)

type LookForRule struct {
	RuleID   string
	Keywords []string
}

func Lookfor(rules []LookForRule) RuleRunner {
	return func(ctx context.Context, filepath string, reader io.Reader) ([]*internal.Finding, error) {
		var findings []*internal.Finding
		// TODO may do more optimization for
		// a really long line
		iter := bufio.NewScanner(reader)
		line := 0
		for iter.Scan() {

			for _, rule := range rules {

				for _, keyword := range rule.Keywords {
					if strings.Contains(string(iter.Bytes()), keyword) {
						var f = internal.Finding{
							Type:   "sast",
							RuleID: rule.RuleID,
							Location: &internal.FindingLocation{
								Path: filepath,
								Position: &internal.FindingLocationPosition{
									Begin: &internal.FindingLocationPositionBegin{
										Line: line,
									},
								},
							},
						}
						findings = append(findings, &f)
					}
				}

			}

			line++
		}
		return findings, nil
	}
}
