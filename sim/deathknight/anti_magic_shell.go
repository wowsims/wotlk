package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core/proto"

	"github.com/wowsims/wotlk/sim/core"
)

func (dk *Deathknight) registerAntiMagicShellSpell() {
	actionID := core.ActionID{SpellID: 48707}

	dk.AntiMagicShell = dk.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagAPL,

		RuneCost: core.RuneCostOptions{
			RunicPowerCost: 20,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    dk.NewTimer(),
				Duration: time.Second * 45,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.AntiMagicShellAura.Activate(sim)
		},
	})

	rpMetrics := dk.AntiMagicShell.RunicPowerMetrics()

	physDmgTakenMult := dk.darkrunedPlateAMSBonus()
	spellDmgTakenMult := 0.25

	var targetDummySpell *core.Spell = nil
	var totalDamageAbsorbed float64
	dk.AntiMagicShellAura = dk.RegisterAura(core.Aura{
		Label:    "Anti-Magic Shell",
		ActionID: actionID,
		Duration: time.Second*5 + core.TernaryDuration(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfAntiMagicShell), 2*time.Second, 0),
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			if dk.Inputs.IsDps {
				target := aura.Unit.CurrentTarget
				if targetDummySpell == nil && target != nil {
					targetDummySpell = aura.Unit.CurrentTarget.RegisterSpell(core.SpellConfig{
						ActionID:    core.ActionID{SpellID: 49375},
						SpellSchool: core.SpellSchoolMagic,
						ProcMask:    core.ProcMaskSpellDamage,
						Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagNoMetrics,

						Cast: core.CastConfig{},

						DamageMultiplier: 1,

						ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
							baseDamage := dk.Inputs.AvgAMSHit * sim.Roll(0.9, 1.1)
							spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeAlwaysHit)
						},
					})
				}

				pa := &core.PendingAction{}
				pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomFloat("ams induced damage")*5.0*float64(time.Second))
				pa.Priority = core.ActionPriorityAuto
				pa.OnAction = func(sim *core.Simulation) {
					if sim.RandomFloat("AMS trigger chance") < min(dk.Inputs.AvgAMSSuccessRate, 1.0) {
						targetDummySpell.Cast(sim, aura.Unit)
					}
				}
				sim.AddPendingAction(pa)
			}

			totalDamageAbsorbed = 0.0
		},
	})

	dk.AddDynamicDamageTakenModifier(func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
		if dk.AntiMagicShellAura.IsActive() && (result.Damage > 0) {
			absorbFrac := 1.0 - core.TernaryFloat64(spell.SpellSchool == core.SpellSchoolPhysical, physDmgTakenMult, spellDmgTakenMult)
			absorbedDmg := absorbFrac * result.Damage

			if absorbedDmg > 0 {
				result.Damage -= absorbedDmg
				dk.AddRunicPower(sim, absorbedDmg/69.0, rpMetrics)
				totalDamageAbsorbed += absorbedDmg

				if totalDamageAbsorbed >= 0.5*dk.MaxHealth() {
					dk.AntiMagicShellAura.Deactivate(sim)
				}
			}
		}
	})

	if !dk.Inputs.IsDps {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell: dk.AntiMagicShell,
			Type:  core.CooldownTypeSurvival,
		})
	} else if dk.Inputs.UseAMS {
		dk.AddMajorCooldown(core.MajorCooldown{
			Spell:    dk.AntiMagicShell,
			Type:     core.CooldownTypeDPS,
			Priority: core.CooldownPriorityLow,
		})
	}
}
