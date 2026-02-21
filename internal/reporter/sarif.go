package reporter

import (
	"encoding/json"
	"fmt"
	"gopherflush/pkg/types"
	"os"
	"path/filepath"
)

type SARIFReporter struct {
	outputPath string
}

func NewSARIFReporter(outputPath string) *SARIFReporter {
	return &SARIFReporter{
		outputPath: outputPath,
	}
}

type SARIFReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SARIFRule `json:"rules"`
}

type SARIFRule struct {
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	ShortDescription     SARIFText   `json:"shortDescription"`
	FullDescription      SARIFText   `json:"fullDescription"`
	DefaultConfiguration SARIFConfig `json:"defaultConfiguration"`
}

type SARIFText struct {
	Text string `json:"text"`
}

type SARIFConfig struct {
	Level string `json:"level"`
}

type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   SARIFText       `json:"message"`
	Locations []SARIFLocation `json:"locations"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           SARIFRegion           `json:"region"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

func (r *SARIFReporter) Generate(report *types.Report) error {
	sarif := r.buildSARIF(report)

	data, err := json.MarshalIndent(sarif, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化SARIF报告失败: %w", err)
	}

	dir := filepath.Dir(r.outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建输出目录失败: %w", err)
		}
	}

	if err := os.WriteFile(r.outputPath, data, 0644); err != nil {
		return fmt.Errorf("写入SARIF报告失败: %w", err)
	}

	return nil
}

func (r *SARIFReporter) buildSARIF(report *types.Report) *SARIFReport {
	rules := r.extractRules(report)

	return &SARIFReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "GopherFlush",
						Version:        "1.0.0",
						InformationURI: "https://github.com/gopherflush/gopherflush",
						Rules:          rules,
					},
				},
				Results: r.buildResults(report),
			},
		},
	}
}

func (r *SARIFReporter) extractRules(report *types.Report) []SARIFRule {
	ruleMap := make(map[string]*SARIFRule)

	for _, v := range report.Violations {
		if _, exists := ruleMap[v.RuleName]; !exists {
			level := r.severityToLevel(v.Severity)
			ruleMap[v.RuleName] = &SARIFRule{
				ID:   v.RuleName,
				Name: v.RuleName,
				ShortDescription: SARIFText{
					Text: v.Message,
				},
				FullDescription: SARIFText{
					Text: v.Suggestion,
				},
				DefaultConfiguration: SARIFConfig{
					Level: level,
				},
			}
		}
	}

	rules := make([]SARIFRule, 0, len(ruleMap))
	for _, rule := range ruleMap {
		rules = append(rules, *rule)
	}

	return rules
}

func (r *SARIFReporter) buildResults(report *types.Report) []SARIFResult {
	results := make([]SARIFResult, 0, len(report.Violations))

	for _, v := range report.Violations {
		results = append(results, SARIFResult{
			RuleID: v.RuleName,
			Level:  r.severityToLevel(v.Severity),
			Message: SARIFText{
				Text: v.Message,
			},
			Locations: []SARIFLocation{
				{
					PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{
							URI: v.FilePath,
						},
						Region: SARIFRegion{
							StartLine:   v.Line,
							StartColumn: v.Column,
						},
					},
				},
			},
		})
	}

	return results
}

func (r *SARIFReporter) severityToLevel(severity types.Severity) string {
	switch severity {
	case types.SeverityCritical:
		return "error"
	case types.SeverityHigh:
		return "error"
	case types.SeverityMedium:
		return "warning"
	case types.SeverityLow:
		return "note"
	default:
		return "none"
	}
}
