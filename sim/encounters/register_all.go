package encounters

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/encounters/icc"
	"github.com/wowsims/wotlk/sim/encounters/naxxramas"
	"github.com/wowsims/wotlk/sim/encounters/toc"
	"github.com/wowsims/wotlk/sim/encounters/ulduar"
)

func init() {
	naxxramas.Register()
	ulduar.Register()
	toc.Register()
	icc.Register()
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
