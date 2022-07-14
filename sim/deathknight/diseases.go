package deathknight

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (deathKnight *DeathKnight) registerDiseaseDots() {
	deathKnight.registerFrostFever()
	deathKnight.registerBloodPlague()
	deathKnight.registerEbonPlague()
}

func (deathKnight *DeathKnight) registerEbonPlague() {
	target := deathKnight.CurrentTarget

	epAura := core.EbonPlaguebringerAura(target)
	epAura.Duration = core.NeverExpires
	deathKnight.EbonPlagueAura = epAura
}

func (deathKnight *DeathKnight) checkForEbonPlague(sim *core.Simulation) {
	if deathKnight.Talents.EbonPlaguebringer == 0 {
		return
	}
	if deathKnight.DiseasesAreActive() {
		deathKnight.EbonPlagueAura.Activate(sim)
	} else {
		deathKnight.EbonPlagueAura.Deactivate(sim)
	}
}

func (deathKnight *DeathKnight) registerFrostFever() {
	actionID := core.ActionID{SpellID: 55095}
	target := deathKnight.CurrentTarget

	deathKnight.FrostFever = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFrost,
	})

	deathKnight.FrostFeverDisease = core.NewDot(core.Dot{
		Spell: deathKnight.FrostFever,
		Aura: target.RegisterAura(core.Aura{
			Label:    "FrostFever-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				deathKnight.checkForEbonPlague(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				deathKnight.checkForEbonPlague(sim)
			},
		}),
		NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return ((127.0 + 80.0*0.32) + hitEffect.MeleeAttackPower(spell.Unit)*0.055) * (1.0 +
						core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})
}

func (deathKnight *DeathKnight) registerBloodPlague() {
	actionID := core.ActionID{SpellID: 55078}
	target := deathKnight.CurrentTarget

	deathKnight.BloodPlague = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
	})

	deathKnight.BloodPlagueDisease = core.NewDot(core.Dot{
		Spell: deathKnight.BloodPlague,
		Aura: target.RegisterAura(core.Aura{
			Label:    "BloodPlague-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				deathKnight.checkForEbonPlague(sim)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				deathKnight.checkForEbonPlague(sim)
			},
		}),
		NumberOfTicks: 5 + int(deathKnight.Talents.Epidemic),
		TickLength:    time.Second * 3,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return ((127.0 + 80.0*0.32) + hitEffect.MeleeAttackPower(spell.Unit)*0.055) * (1.0 +
						core.TernaryFloat64(deathKnight.DiseasesAreActive(), 0.05*float64(deathKnight.Talents.TundraStalker), 0.0))
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})
}
