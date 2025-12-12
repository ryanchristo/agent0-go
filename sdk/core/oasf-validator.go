package core

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/ryanchristo/agent0-go/sdk/taxonomies"
)

// SkillsData is data based on the OASF skills taxonomy.
type SkillsData struct {
	Skills map[string]any
}

// DomainsData is data based on the OASF domains taxonomy.
type DomainsData struct {
	Domains map[string]any
}

// ValidateSkill validates if a skill slug exists in the OASF taxonomy.
func ValidateSkill(slug string) bool {
	data := loadSkillsData()
	return data.Skills[slug] != nil
}

// ValidateDomain validates if a domain slug exists in the OASF taxonomy.
func ValidateDomain(slug string) bool {
	data := loadDomainsData()
	return data.Domains[slug] != nil
}

var (
	skillsDataCache  *SkillsData
	domainsDataCache *DomainsData
	skillsOnce       sync.Once
	domainsOnce      sync.Once
)

// loadSkillsData loads and caches the skills taxonomy data.
func loadSkillsData() *SkillsData {
	skillsOnce.Do(func() {
		var skillsData SkillsData
		err := json.Unmarshal(taxonomies.SkillsJSON, &skillsData)
		if err != nil {
			log.Fatalf("failed to unmarshal skills: %v", err)
		}
		skillsDataCache = &skillsData
	})
	return skillsDataCache
}

// loadDomainsData loads and caches the domains taxonomy data.
func loadDomainsData() *DomainsData {
	domainsOnce.Do(func() {
		var domainsData DomainsData
		err := json.Unmarshal(taxonomies.DomainsJSON, &domainsData)
		if err != nil {
			log.Fatalf("failed to unmarshal domains: %v", err)
		}
		domainsDataCache = &domainsData
	})
	return domainsDataCache
}
