package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerBoneShieldSpell() {
	if !dk.Talents.BoneShield {
		return
	}

	actionID := core.ActionID{SpellID: 49222}
	cdTimer := dk.NewTimer()
	cd := time.Minute*1 - dk.thassariansPlateCooldownReduction(dk.BoneShield)

	dk.BoneShieldAura = dk.RegisterAura(core.Aura{
		Label:     "Bone Shield",
		ActionID:  actionID,
		Duration:  time.Minute * 5,
		MaxStacks: 3 + core.TernaryInt32(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfBoneShield), 1, 0),
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			dk.BoneShieldAura.Activate(sim)
			dk.BoneShieldAura.UpdateExpires(sim.CurrentTime + time.Minute*5)
			dk.BoneShieldAura.SetStacks(sim, dk.BoneShieldAura.MaxStacks)
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

	baseCost := float64(core.NewRuneCost(10, 0, 0, 1, 0))
	dk.BoneShield = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID:     actionID,
		Flags:        core.SpellFlagNoOnCastComplete,
		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:  core.GCDDefault,
				Cost: baseCost,
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
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 0.0, 0, 0, 1) && dk.BoneShield.IsReady(sim)
	}, nil)
}
