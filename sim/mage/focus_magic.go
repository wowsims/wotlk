package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) applyFocusMagic() {
	if !mage.Talents.FocusMagic {
		return
	}

	// This is used only for the individual sim.
	if mage.Party.Raid.Size() == 1 {
		if mage.Options.FocusMagicPercentUptime > 0 {
			selfAura, _ := core.FocusMagicAura(mage.GetCharacter(), nil)
			core.ApplyFixedUptimeAura(selfAura, float64(mage.Options.FocusMagicPercentUptime)/100, time.Second*10, 1)
		}
		return
	}

	focusMagicTargetAgent := mage.Party.Raid.GetPlayerFromUnitReference(mage.Options.FocusMagicTarget)
	if focusMagicTargetAgent == nil {
		return
	} else if focusMagicTargetAgent.GetCharacter() == mage.GetCharacter() {
		// When self is selected, give permanent self buff.
		selfAura, _ := core.FocusMagicAura(mage.GetCharacter(), nil)
		core.MakePermanent(selfAura)
		return
	}

	core.FocusMagicAura(mage.GetCharacter(), focusMagicTargetAgent.GetCharacter())
}
