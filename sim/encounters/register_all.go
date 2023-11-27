package encounters

import (
	"github.com/wowsims/classic/sod/sim/core"
)

func init() {
	// TODO: Classic encounters?
	// naxxramas.Register()
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
