package druid

import (
	"github.com/wowsims/wotlk/sim/core"
)

// Returns the time to wait before the next action, or 0 if innervate is on CD
// or disabled.
func (druid *Druid) registerInnervateCD() {
	innervateTargetAgent := druid.Party.Raid.GetPlayerFromRaidTarget(druid.SelfBuffs.InnervateTarget)
	if innervateTargetAgent == nil {
		return
	}
	innervateTarget := innervateTargetAgent.GetCharacter()

	actionID := core.ActionID{SpellID: 29166, Tag: druid.Index}
	var innervateSpell *core.Spell

	innervateCD := core.InnervateCD

	var innervateAura *core.Aura
	var innervateManaThreshold float64
	druid.RegisterResetEffect(func(sim *core.Simulation) {
		if innervateTarget == druid.GetCharacter() {
			if druid.StartingForm.Matches(Cat) {
				// double shift + innervate cost.
				// Prevents not having enough mana to shift back into form if more powershift are executed
				innervateManaThreshold = druid.CatForm.DefaultCast.Cost*2 + innervateSpell.DefaultCast.Cost
			} else {
				// Threshold can be lower when casting on self because its never mid-cast.
				innervateManaThreshold = 500
			}
		} else {
			innervateManaThreshold = core.InnervateManaThreshold(innervateTarget)
		}
		innervateAura = core.InnervateAura(innervateTarget, actionID.Tag)
	})

	innervateSpell = druid.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    SpellFlagOmenTrigger,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.04,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    druid.NewTimer(),
				Duration: innervateCD,
			},
		},

		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			// Technically this shouldn't be allowed in bear form either, but bear
			// doesn't have shifting implemented in its rotation.
			if druid.InForm(Cat) {
				return false
			}
			// If target already has another innervate, don't cast.
			if innervateTarget.HasActiveAuraWithTag(core.InnervateAuraTag) {
				return false
			}

			return true
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			innervateAura.Activate(sim)
		},
	})

	druid.AddMajorCooldown(core.MajorCooldown{
		Spell: innervateSpell,
		Type:  core.CooldownTypeMana,
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			// Innervate needs to be activated as late as possible to maximize DPS. The issue is that
			// innervate gives so much mana that it can cause Super Mana Potion or Dark Rune usages
			// to be delayed, if they come off CD soon after innervate. This delay is minimized by
			// activating innervate from the smallest amount of mana possible.
			return innervateTarget.CurrentMana() <= innervateManaThreshold
		},
	})
}
