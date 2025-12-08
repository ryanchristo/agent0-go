package core

// SkillsData is data based on the OASF skills taxonomy.
type SkillsData struct {
	Skills map[string]any `json:"skills"`
}

// DomainData is data based on the OASF domains taxonomy.
type DomainData struct {
	Domains map[string]any `json:"domains"`
}

// ValidateSkill validates if a skill slug exists in the OASF taxonomy.
func ValidateSkill(slug string) bool {

	// TODO: implementation

	return true
}

// ValidateDomain validates if a domain slug exists in the OASF taxonomy.
func ValidateDomain(slug string) bool {

	// TODO: implementation

	return true
}
