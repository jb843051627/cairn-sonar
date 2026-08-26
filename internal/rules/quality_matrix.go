package rules

import (
	"sort"

	"cairn-sonar/internal/model"
)

type QualityRule struct {
	Code        string
	Description string
	Weight      float64
	Threshold   float64
	Critical    bool
}

func DefaultQualityRules() []QualityRule {
	return []QualityRule{
		{Code: "contrast", Description: "回波对比度足够", Weight: 0.25, Threshold: 6},
		{Code: "confidence", Description: "回波置信度足够", Weight: 0.25, Threshold: 0.7},
		{Code: "coverage", Description: "巡检采样覆盖完整", Weight: 0.2, Threshold: 0.95},
		{Code: "clipping", Description: "采样未大面积削波", Weight: 0.15, Threshold: 0.99},
		{Code: "calibration", Description: "仪器校准有效", Weight: 0.15, Threshold: 1},
	}
}

type RuleResult struct {
	Code     string
	Passed   bool
	Score    float64
	Weight   float64
	Critical bool
	Message  string
	Evidence []string
}

func (r RuleResult) Clone() RuleResult {
	r.Evidence = append([]string(nil), r.Evidence...)
	return r
}

func EvaluateFeatureRule(rule QualityRule, feature model.SignalFeature) RuleResult {
	result := RuleResult{Code: rule.Code, Weight: rule.Weight, Critical: rule.Critical}
	switch rule.Code {
	case "contrast":
		result.Score = bound(feature.ContrastDB / rule.Threshold)
		result.Passed = feature.ContrastDB >= rule.Threshold
		result.Message = "对比度"
	case "confidence":
		result.Score = bound(feature.Confidence / rule.Threshold)
		result.Passed = feature.Confidence >= rule.Threshold
		result.Message = "置信度"
	case "clipping":
		result.Score = 1
		result.Passed = !feature.HasFlag("clipped")
		if !result.Passed {
			result.Score = 0
		}
		result.Message = "削波"
	default:
		result.Score = 0
		result.Message = "未识别规则"
	}
	if result.Passed {
		result.Evidence = []string{feature.WindowID, result.Message + "通过"}
	} else {
		result.Evidence = []string{feature.WindowID, result.Message + "需要复核"}
	}
	return result
}

func EvaluateSurveyQuality(batch model.FeatureBatch, rules []QualityRule) (model.SurveyQuality, []RuleResult) {
	quality := batch.Quality.Clone()
	results := make([]RuleResult, 0)
	if len(batch.Items) == 0 {
		quality.AddWarning("没有可评估的回波")
		quality.Band = "poor"
		return quality, results
	}
	for _, feature := range batch.Items {
		for _, rule := range rules {
			results = append(results, EvaluateFeatureRule(rule, feature))
		}
	}
	score := 0.0
	weight := 0.0
	failures := 0
	for _, result := range results {
		score += result.Score * result.Weight
		weight += result.Weight
		if !result.Passed {
			failures++
		}
	}
	if weight > 0 {
		quality.Score = score / float64(len(batch.Items))
	}
	quality.Band = qualityBand(quality.Score)
	if failures > 0 {
		quality.AddWarning("存在未通过的声学质量规则")
	}
	quality.Coverage = coverage(batch.Items)
	return quality, results
}

func coverage(items []model.SignalFeature) float64 {
	if len(items) == 0 {
		return 0
	}
	usable := 0
	for _, item := range items {
		if item.ContrastDB > 0 && item.Confidence > 0 {
			usable++
		}
	}
	return float64(usable) / float64(len(items))
}

func SummarizeRules(results []RuleResult) map[string]float64 {
	out := make(map[string]float64)
	counts := make(map[string]int)
	for _, result := range results {
		out[result.Code] += result.Score
		counts[result.Code]++
	}
	for code, total := range out {
		if counts[code] > 0 {
			out[code] = total / float64(counts[code])
		}
	}
	return out
}

func FailedRules(results []RuleResult) []RuleResult {
	out := make([]RuleResult, 0)
	for _, result := range results {
		if !result.Passed {
			out = append(out, result.Clone())
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Critical == out[j].Critical {
			return out[i].Score < out[j].Score
		}
		return out[i].Critical
	})
	return out
}

func bound(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

type SeverityRule struct {
	Kind        string
	MinPeakDB   float64
	MinContrast float64
	Severity    int
	Tag         string
}

func DefaultSeverityRules() []SeverityRule {
	return []SeverityRule{
		{Kind: "delamination", MinPeakDB: 70, MinContrast: 14, Severity: 3, Tag: "structural"},
		{Kind: "void", MinPeakDB: 55, MinContrast: 10, Severity: 2, Tag: "subsurface"},
		{Kind: "moisture", MinPeakDB: 40, MinContrast: 7, Severity: 1, Tag: "environmental"},
	}
}

func ClassifyFeature(feature model.SignalFeature, rules []SeverityRule) (string, int, string) {
	for _, rule := range rules {
		if feature.PeakDB >= rule.MinPeakDB && feature.ContrastDB >= rule.MinContrast {
			return rule.Kind, rule.Severity, rule.Tag
		}
	}
	return "background", 0, "noise"
}
