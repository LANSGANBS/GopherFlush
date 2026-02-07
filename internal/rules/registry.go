package rules

// Registry 规则注册表
type Registry struct {
	rules map[string]Rule
}

// NewRegistry 创建新的规则注册表
func NewRegistry() *Registry {
	return &Registry{
		rules: make(map[string]Rule),
	}
}

// Register 注册规则
func (r *Registry) Register(rule Rule) {
	r.rules[rule.Name()] = rule
}

// GetAll 获取所有规则
func (r *Registry) GetAll() []Rule {
	result := make([]Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		if rule.Enabled() {
			result = append(result, rule)
		}
	}
	return result
}

// Get 根据名称获取规则
func (r *Registry) Get(name string) (Rule, bool) {
	rule, ok := r.rules[name]
	return rule, ok
}
