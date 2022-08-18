package encounters

import (
	"github.com/wowsims/wotlk/sim/core"
)

func init() {
	//registerBlackTemple()
	//registerSunwellPlateau()
	registerNaxxramas10()
}

func AddSingleTargetBossEncounter(presetTarget core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
