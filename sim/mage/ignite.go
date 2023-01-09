package mage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// If two spells proc Ignite at almost exactly the same time, the latter
// overwrites the former.
const IgniteMunchWindow = time.Millisecond * 10
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
			if spell == mage.LivingBombDot.Spell && result.DidCrit() {
				mage.procIgnite(sim, result)
			}
		},
	})

	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 12654},
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       SpellFlagMage | core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			mage.IgniteDots[target.Index].ApplyOrReset(sim)
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
		},
	})

	mage.IgniteDots = make([]*core.Dot, mage.Env.GetNumTargets())
	for i := range mage.IgniteDots {
		mage.IgniteDots[i] = core.NewDot(core.Dot{
			Spell: mage.Ignite,
			Aura: mage.Env.GetTargetUnit(int32(i)).RegisterAura(core.Aura{
				Label:    "Ignite-" + strconv.Itoa(int(mage.Index)),
				ActionID: mage.Ignite.ActionID,
				Tag:      "Ignite",
			}),
			NumberOfTicks: IgniteTicks,
			TickLength:    time.Second * 2,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
			SnapshotAttackerMultiplier: 1,
		})
	}
}

func (mage *Mage) procIgnite(sim *core.Simulation, result *core.SpellResult) {
	dot := mage.IgniteDots[result.Target.Index]

	newDamage := result.Damage * 0.08 * float64(mage.Talents.Ignite)
	outstandingDamage := core.TernaryFloat64(dot.IsActive(), dot.SnapshotBaseDamage*float64(dot.NumberOfTicks-dot.TickCount), 0)

	var munchPenalty float64
	if mage.Options.IgniteMunching && sim.CurrentTime <= mage.igniteMunchTime+IgniteMunchWindow {
		munchPenalty = mage.igniteMunchDmg
		if sim.Log != nil {
			mage.Log(sim, fmt.Sprintf("Ignite munched: %0.01f", mage.igniteMunchDmg))
		}
	}

	dot.SnapshotBaseDamage = (outstandingDamage + newDamage - munchPenalty) / float64(IgniteTicks)
	mage.igniteMunchDmg = newDamage
	mage.igniteMunchTime = sim.CurrentTime
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
				mage.AddMana(sim, mage.Unit.BaseMana*0.02, manaMetrics, false)
			}
		},
	})
}
