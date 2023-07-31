package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (paladin *Paladin) RegisterAvengingWrathCD() {
	actionID := core.ActionID{SpellID: 31884}

	paladin.AvengingWrathAura = paladin.RegisterAura(core.Aura{
		Label:    "Avenging Wrath",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfAvengingWrath) {
				paladin.HammerOfWrath.CD.Duration /= 2
			}
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.2
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			if paladin.HasMajorGlyph(proto.PaladinMajorGlyph_GlyphOfAvengingWrath) {
				paladin.HammerOfWrath.CD.Duration *= 2
			}
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.2
		},
	})
	core.RegisterPercentDamageModifierEffect(paladin.AvengingWrathAura, 1.2)

	paladin.AvengingWrath = paladin.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.08,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    paladin.NewTimer(),
				Duration: time.Minute*3 - (time.Second * time.Duration(30*paladin.Talents.SanctifiedWrath)),
			},
			SharedCD: core.Cooldown{
				Timer:    paladin.GetMutualLockoutDPAW(),
				Duration: 30 * time.Second,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			paladin.AvengingWrathAura.Activate(sim)
		},
	})

	paladin.AddMajorCooldown(core.MajorCooldown{
		Spell: paladin.AvengingWrath,
		Type:  core.CooldownTypeDPS,
		// modify this logic if it should ever not be spammed on CD / maybe should synced with other CDs
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			if paladin.HoldLastAvengingWrathUntilExecution && float64(sim.CurrentTime+paladin.AvengingWrath.CD.Duration) >= float64(sim.Duration) {
				if float64(sim.CurrentTime+paladin.AvengingWrathAura.Duration) >= float64(sim.Duration) {
					// If we're cutting it close on iteration end time, pop AW anyways to try and get full duration.
					return true
				}
				return false
			}

			if paladin.CurrentSeal == paladin.SealOfVengeanceAura {
				if paladin.SovDotSpell.Dot(paladin.CurrentTarget).GetStacks() < 5 {
					return false
				}
			}

			return true
		},
	})
}
