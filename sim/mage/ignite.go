package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

func (mage *Mage) applyIgnite() {
	if mage.Talents.Ignite == 0 {
		return
	}

	mage.RegisterAura(core.Aura{
		Label:    "Ignite Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && result.Outcome.Matches(core.OutcomeCrit) {
				mage.procIgnite(sim, result.Target, result.Damage)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if spell == mage.LivingBombDot.Spell && result.Outcome.Matches(core.OutcomeCrit) {
				mage.procIgnite(sim, result.Target, result.Damage)
			}
		},
	})

	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: 12654},
		SpellSchool:      core.SpellSchoolFire,
		ProcMask:         core.ProcMaskEmpty,
		Flags:            SpellFlagMage | core.SpellFlagIgnoreModifiers,
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),
	})

	mage.IgniteDots = make([]*core.Dot, mage.Env.GetNumTargets())
	mage.IgniteDamageBuffers = make([]float64, mage.Env.GetNumTargets())
	for i := range mage.IgniteDots {
		mage.IgniteDots[i] = core.NewDot(core.Dot{
			Spell: mage.Ignite,
			Aura: mage.Env.GetTargetUnit(int32(i)).RegisterAura(core.Aura{
				Label:    "Ignite-" + strconv.Itoa(int(mage.Index)),
				ActionID: mage.Ignite.ActionID,
				Tag:      "Ignite",
			}),
			NumberOfTicks: 2,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				mage.IgniteDamageBuffers[target.Index] -= dot.SnapshotBaseDamage
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
			SnapshotAttackerMultiplier: 1,
		})
	}
}

func (mage *Mage) procIgnite(sim *core.Simulation, target *core.Unit, spellDamage float64) {
	dot := mage.IgniteDots[target.Index]

	if !dot.IsActive() {
		mage.IgniteDamageBuffers[target.Index] = 0
	}

	mage.IgniteDamageBuffers[target.Index] += spellDamage * float64(mage.Talents.Ignite) * 0.08
	dot.SnapshotBaseDamage = mage.IgniteDamageBuffers[target.Index] / 2
	dot.Apply(sim)
}

func (mage *Mage) applyEmpoweredFire() {
	if mage.Talents.EmpoweredFire == 0 {
		return
	}

	var actionId = core.ActionID{SpellID: 67545}

	procChance := []float64{0, .33, .67, 1}[mage.Talents.EmpoweredFire]
	manaMetrics := mage.NewManaMetrics(actionId)

	mage.RegisterAura(core.Aura{
		Label:    "Empowered Fire",
		ActionID: actionId,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == mage.Ignite && sim.Proc(procChance, "Empowered Fire") {
				mage.AddMana(sim, mage.Unit.BaseMana*0.02, manaMetrics, false)
			}
		},
	})
}
