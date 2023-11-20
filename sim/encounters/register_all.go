package encounters

import (
	"github.com/wowsims/classic/sim/core"
	"github.com/wowsims/classic/sim/encounters/icc"
	"github.com/wowsims/classic/sim/encounters/naxxramas"
	"github.com/wowsims/classic/sim/encounters/toc"
	"github.com/wowsims/classic/sim/encounters/ulduar"
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
