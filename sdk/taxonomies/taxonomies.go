package taxonomies

import _ "embed"

//go:embed all_skills.json
var SkillsJSON []byte

//go:embed all_domains.json
var DomainsJSON []byte
