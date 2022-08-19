package deathknight

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (dk *Deathknight) registerAntiMagicShellSpell() {
	actionID := core.ActionID{SpellID: 48707}
	cdTimer := dk.NewTimer()
	cd := time.Second*45 - time.Second*time.Duration(core.TernaryInt32(dk.HasMajorGlyph(proto.DeathknightMajorGlyph_GlyphOfAntiMagicShell), 2, 0))

	baseCost := float64(core.NewRuneCost(20.0, 0, 0, 0, 0))
	dk.AntiMagicShell = dk.RegisterSpell(nil, core.SpellConfig{
		ActionID: actionID,

		ResourceType: stats.RunicPower,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: cd,
			},
			IgnoreHaste: true,
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			dk.AntiMagicShellAura.Activate(sim)
		},
	}, func(sim *core.Simulation) bool {
		return dk.CastCostPossible(sim, 20.0, 0, 0, 0) && dk.AntiMagicShell.IsReady(sim)
	}, nil)

	rpMetrics := dk.AntiMagicShell.RunicPowerMetrics()

	var targetDummySpell *core.Spell = nil
	dk.AntiMagicShellAura = dk.RegisterAura(core.Aura{
		Label:    "Anti-Magic Shell",
		ActionID: actionID,
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			target := aura.Unit.CurrentTarget
			if targetDummySpell == nil && target != nil {
				targetDummySpell = aura.Unit.CurrentTarget.RegisterSpell(core.SpellConfig{
					ActionID:    core.ActionID{SpellID: 49375},
					SpellSchool: core.SpellSchoolMagic,

					Cast: core.CastConfig{},

					ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
						ProcMask: core.ProcMaskSpellDamage,

						DamageMultiplier: 1,

						BaseDamage:     core.BaseDamageConfigRoll(dk.Inputs.AvgAMSHit*0.9, dk.Inputs.AvgAMSHit*1.1),
						OutcomeApplier: target.OutcomeFuncAlwaysHit(),
					}),
				})
			}

			pa := &core.PendingAction{}
			pa.NextActionAt = sim.CurrentTime + time.Duration(sim.RandomFloat("ams induced damage")*5.0*float64(time.Second))
			pa.Priority = core.ActionPriorityAuto
			pa.OnAction = func(sim *core.Simulation) {
				if sim.RandomFloat("AMS trigger chance") < core.MinFloat(dk.Inputs.AvgAMSSuccessRate, 1.0) {
					targetDummySpell.Cast(sim, aura.Unit)
				}
			}
			sim.AddPendingAction(pa)
		},

		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
		},

		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage > 0 {
				absorvedDamage := core.MinFloat(0.75*spellEffect.Damage, 0.5*dk.MaxHealth())
				dk.RemoveHealth(sim, spellEffect.Damage-absorvedDamage)
				dk.AddRunicPower(sim, spellEffect.Damage/69.0, rpMetrics)
			}
		},
	})
}
