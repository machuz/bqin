package bqin_test

import (
	"os"
	"testing"

	"github.com/kayac/bqin"
)

var ExpectedDefault = []string{}

func TestLoadConfig(t *testing.T) {
	os.Setenv("AWS_REGION", "ap-northeast-1")

	type pattern struct {
		path                 string
		expectedRulesStrings []string
	}

	t.Run("success", func(t *testing.T) {
		patterns := []pattern{
			{
				"testdata/config/default.yaml",
				[]string{
					"s3://bqin.bucket.test/data/user => bqin-test-gcp.test.user",
					"s3://bqin.bucket.test/data/(.+)/part-([0-9]+).csv => bqin-test-gcp.test.$1_$2",
				},
			},
			{
				"testdata/config/standard.yaml",
				[]string{
					"s3://bqin.bucket.test/data/user => bqin-test-gcp.test.user",
				},
			},
			{
				"testdata/config/hive_format.yaml",
				[]string{
					"s3://bqin.bucket.test/data/(.+)/snapshot_at=([0-9]{8})/.+ => bqin-test-gcp.test.$1_$2",
				},
			},
		}
		for _, p := range patterns {
			t.Run(p.path, func(t *testing.T) {
				conf, err := bqin.LoadConfig(p.path)
				if err != nil {
					t.Errorf("unexpected error :%s", err)
					return
				}
				for i, rule := range conf.Rules {
					if rule.String() != p.expectedRulesStrings[i] {
						t.Logf("rule[%d]      got: %s", i, rule.String())
						t.Logf("rule[%d] expected: %s", i, p.expectedRulesStrings[i])
						t.Errorf("rule[%d] is unexpected", i)
					}
				}
			})
		}
	})

	t.Run("fail", func(t *testing.T) {
		patterns := []pattern{
			{path: "testdata/config/not_found.yaml"},
			{path: "testdata/config/broken_invalid_key_regexp.yaml"},
			{path: "testdata/config/broken_no_queue_name.yaml"},
			{path: "testdata/config/broken_no_key_matcher.yaml"},
			{path: "testdata/config/broken_no_tempbucket_option.yaml"},
		}
		for _, p := range patterns {
			t.Run(p.path, func(t *testing.T) {
				_, err := bqin.LoadConfig(p.path)
				if err == nil {
					t.Errorf("LoadConfig(%s) must be failed", p.path)
				}
			})
		}
	})
}
