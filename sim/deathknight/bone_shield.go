package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (deathKnight *DeathKnight) registerBoneShieldSpell() {
	if !deathKnight.Talents.BoneShield {
		return
	}

	actionID := core.ActionID{SpellID: 49222}
	cdTimer := deathKnight.NewTimer()
	cd := time.Minute * 1

	deathKnight.BoneShieldAura = deathKnight.RegisterAura(core.Aura{
		Label:     "Bone Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 5,
		MaxStacks: 3 + core.TernaryInt32(deathKnight.HasMajorGlyph(proto.DeathKnightMajorGlyph_GlyphOfBoneShield), 1, 0),
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.BoneShieldAura.Activate(sim)
			deathKnight.BoneShieldAura.UpdateExpires(sim.CurrentTime + time.Minute*4)
			deathKnight.BoneShieldAura.SetStacks(sim, 3)
		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			aura.RemoveStack(sim)
			if aura.GetStacks() == 0 {
				aura.Deactivate(sim)
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.ModifyAdditiveDamageModifier(sim, 0.02)

			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.8
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.ModifyAdditiveDamageModifier(sim, -0.02)

			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	deathKnight.BoneShield = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			deathKnight.BoneShieldAura.Activate(sim)
			deathKnight.BoneShieldAura.Prioritize()

			dkSpellCost := deathKnight.DetermineOptimalCost(sim, 0, 0, 1)
			deathKnight.Spend(sim, spell, dkSpellCost)

			amountOfRunicPower := 10.0
			deathKnight.AddRunicPower(sim, amountOfRunicPower, deathKnight.BoneShield.RunicPowerMetrics())
		},
	})
}

func (deathKnight *DeathKnight) CanBoneShield(sim *core.Simulation) bool {
	return deathKnight.CastCostPossible(sim, 0.0, 0, 0, 1) && deathKnight.BoneShield.IsReady(sim)
}

func (deathKnight *DeathKnight) CastBoneShield(sim *core.Simulation, target *core.Target) bool {
	if deathKnight.CanBoneShield(sim) {
		deathKnight.CastBoneShield(sim, target)
		return true
	}
	return false
}
