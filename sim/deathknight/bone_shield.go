package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (dk *Deathknight) registerBoneShieldSpell() {
	if !dk.Talents.BoneShield {
		return
	}

	actionID := core.ActionID{SpellID: 49222}
	cdTimer := dk.NewTimer()
	cd := time.Minute * 1

	dk.BoneShieldAura = dk.RegisterAura(core.Aura{
		Label:     "Bone Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 5,
		MaxStacks: 3 + core.TernaryInt32(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfBoneShield), 1, 0),
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			dk.BoneShieldAura.Activate(sim)
			dk.BoneShieldAura.UpdateExpires(sim.CurrentTime + time.Minute*4)
			dk.BoneShieldAura.SetStacks(sim, 3)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			aura.RemoveStack(sim)
			if aura.GetStacks() == 0 {
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyAdditiveDamageModifier(sim, 0.02)

			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyAdditiveDamageModifier(sim, -0.02)

			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	dk.BoneShield = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				cast.GCD = dk.getModifiedGCD()
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.BoneShieldAura.Activate(sim)
			dk.BoneShieldAura.Prioritize()

			dkSpellCost := dk.DetermineOptimalCost(sim, 0, 0, 1)
			dk.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 10.0
			dk.AddRunicPower(sim, amountOfRunicPower, dk.BoneShield.RunicPowerMetrics())
		},
	})
}

func (dk *Deathknight) CanBoneShield(sim *core.Simulation) bool {
	return dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.BoneShield.IsReady(sim)
}

func (dk *Deathknight) CastBoneShield(sim *core.Simulation, target *core.Unit) bool {
	if dk.CanBoneShield(sim) {
		dk.BoneShield.Cast(sim, target)
		return true
	}
	return false
}
