package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// If two spells proc Ignite at almost exactly the same time, the latter
// overwrites the former.
const IgniteTicks = 2

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
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && result.DidCrit() {
				mage.procIgnite(sim, result)
			}
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if mage.LivingBomb != nil && result.DidCrit() {
				mage.procIgnite(sim, result)
			}
		},
	})

	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 12654},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskProc,
		Flags:       SpellFlagMage | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Ignite",
			},
			NumberOfTicks: IgniteTicks,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.SpellMetrics[target.UnitIndex].Hits++
			spell.Dot(target).ApplyOrReset(sim)
		},
	})
}

func (mage *Mage) procIgnite(sim *core.Simulation, result *core.SpellResult) {
	dot := mage.Ignite.Dot(result.Target)

	newDamage := result.Damage * 0.08 * float64(mage.Talents.Ignite)
	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	dot.SnapshotAttackerMultiplier = 1
	dot.SnapshotBaseDamage = (outstandingDamage + newDamage) / float64(IgniteTicks)
	mage.Ignite.Cast(sim, result.Target)
}

func (mage *Mage) applyEmpoweredFire() {
	if mage.Talents.EmpoweredFire == 0 {
		return
	}

	procChance := []float64{0, .33, .67, 1}[mage.Talents.EmpoweredFire]
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 67545})

	mage.RegisterAura(core.Aura{
		Label:    "Empowered Fire",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnPeriodicDamageDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell == mage.Ignite && (procChance == 1 || sim.Proc(procChance, "Empowered Fire")) {
				mage.AddMana(sim, mage.Unit.BaseMana*0.02, manaMetrics)
			}
		},
	})
}
