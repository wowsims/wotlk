package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

var IgniteActionID = core.ActionID{SpellID: 12848}
var empoweredFireActionId = core.ActionID{SpellID: 31658}

// TODO: Global variables very bad. This will break the raid sim, where there can be multiple mages.
var manaMetrics *core.ResourceMetrics

func (mage *Mage) registerIgniteSpell() {
	manaMetrics = mage.NewManaMetrics(empoweredFireActionId)
	mage.Ignite = mage.RegisterSpell(core.SpellConfig{
		ActionID:    IgniteActionID,
		SpellSchool: core.SpellSchoolFire,
		Flags:       SpellFlagMage | core.SpellFlagIgnoreModifiers,
	})
}

func (mage *Mage) newIgniteDot(target *core.Unit) *core.Dot {
	return core.NewDot(core.Dot{
		Spell: mage.Ignite,
		Aura: target.RegisterAura(core.Aura{
			Label:    "Ignite-" + strconv.Itoa(int(mage.Index)),
			ActionID: IgniteActionID,
		}),
		NumberOfTicks: 2,
		TickLength:    time.Second * 2,
	})
}

func (mage *Mage) procIgnite(sim *core.Simulation, target *core.Unit, damageFromProccingSpell float64) {
	igniteDot := mage.IgniteDots[target.Index]

	newIgniteDamage := damageFromProccingSpell * float64(mage.Talents.Ignite) * 0.08
	if igniteDot.IsActive() {
		newIgniteDamage += mage.IgniteTickDamage[target.Index] * float64(2-igniteDot.TickCount)
	}

	newTickDamage := newIgniteDamage / 2
	mage.IgniteTickDamage[target.Index] = newTickDamage

	// Hacky: mimic the logs in sim/core/cast.go, to get Ignite to show up in the timeline as a cast.
	// TODO: Just make this a spell.
	if sim.Log != nil {
		mage.Log(sim, "Casting %s (Cost = %0.03f, Cast Time = %s)", IgniteActionID, 0.0, time.Duration(0))
		mage.Log(sim, "Completed cast %s", IgniteActionID)
	}
	mage.Ignite.SpellMetrics[target.TableIndex].Casts++
	mage.Ignite.SpellMetrics[target.TableIndex].Hits++

	// Reassign the effect to apply the new damage value.
	igniteDot.TickEffects = core.TickFuncSnapshot(target, core.SpellEffect{
		ProcMask:         core.ProcMaskPeriodicDamage,
		DamageMultiplier: 1,
		ThreatMultiplier: 1 - 0.1*float64(mage.Talents.BurningSoul),
		IsPeriodic:       true,
		BaseDamage:       core.BaseDamageConfigFlat(newTickDamage),
		OutcomeApplier: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect, attackTable *core.AttackTable) {
			if float64(mage.Talents.EmpoweredFire)/3.0 > sim.RandomFloat("EmpoweredFireIgniteMana") {
				mage.AddMana(sim, mage.Unit.BaseMana*.02, manaMetrics, false)
			}
			mage.OutcomeFuncTick()
		},
	})
	igniteDot.Apply(sim)
}

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
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spell.SpellSchool == core.SpellSchoolFire && spellEffect.Outcome.Matches(core.OutcomeCrit) {
				mage.procIgnite(sim, spellEffect.Target, spellEffect.Damage)
			}
		},
	})
}
