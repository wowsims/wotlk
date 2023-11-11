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
	cd := time.Minute*1 - dk.thassariansPlateCooldownReduction(dk.BoneShield)
	stackRemovalCd := 0 * time.Second

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
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				if sim.CurrentTime > stackRemovalCd+2*time.Second {
					stackRemovalCd = sim.CurrentTime

					aura.RemoveStack(sim)
					if aura.GetStacks() == 0 {
						aura.Deactivate(sim)
					}
				}
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyDamageModifier(0.02)
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.8
			stackRemovalCd = sim.CurrentTime
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			dk.ModifyDamageModifier(-0.02)
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.8
		},
	})

	dk.BoneShield = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			UnholyRuneCost: 1,
			RunicPowerGain: 10,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.BoneShieldAura.Activate(sim)
			dk.BoneShieldAura.SetStacks(sim, dk.BoneShieldAura.MaxStacks)
		},
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.BoneShield,
			Type:  core.CooldownTypeSurvival,
		})
	}
}
