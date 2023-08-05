package collect

// 构建动态规则模型 TaskModel
type (
	TaskModel struct {
		Property
		Root  string      `json:"root_script"`
		Rules []RuleModel `json:"rule"`
	}
	RuleModel struct {
		Name      string `json:"name"`
		ParseFunc string `json:"parse_script"`
	}
)
