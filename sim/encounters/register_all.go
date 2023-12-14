package encounters

import (
	"github.com/wowsims/sod/sim/core"
)

func init() {
	// TODO: Classic encounters?
	// naxxramas.Register()
	addLevel25("SoD")
	addLevel40("SoD")
	addLevel50("SoD")
	addLevel60("SoD")
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
