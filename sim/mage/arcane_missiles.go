package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (mage *Mage) getArcaneMissilesTickSpell(rank int, numTicks int32, baseDotDamage float64) *core.Spell {
	spellCoeff := [9]float64{0, .132, .204, .24, .24, .24, .24, .24, .24}[rank]
	spellId := [9]int32{0, 7268, 7269, 7270, 8419, 8418, 10273, 10274, 25346}[rank]

	return mage.GetOrRegisterSpell(core.SpellConfig{
		ActionID:         core.ActionID{SpellID: spellId},
		Rank:             rank,
		SpellSchool:      core.SpellSchoolArcane,
		ProcMask:         core.ProcMaskProc | core.ProcMaskSpellDamage | core.ProcMaskNotInSpellbook,
		Flags:            SpellFlagMage,
		MissileSpeed:     20,
		BonusHitRating:   100 * core.SpellHitRatingPerHitChance, // Not an independent hit once initial lands
		BonusCritRating:  0,
		DamageMultiplier: 1,
		CritMultiplier:   1, // No crit on channels
		ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			damage := baseDotDamage + (spellCoeff * spell.SpellPower())

			// TODO: Classic review Arcane missiles snap shot mechanics for Arcane Blast rune
			if mage.ArcaneBlastAura != nil && mage.ArcaneBlastAura.IsActive() {
				damage *= 0.15 * float64(mage.ArcaneBlastAura.GetStacks())
			}

			result := spell.CalcPeriodicDamage(sim, target, damage, spell.OutcomeExpectedTick)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

func (mage *Mage) getArcaneMissilesSpellConfig(rank int) core.SpellConfig {
	numTicks := [9]int32{0, 3, 4, 5, 5, 5, 5, 5, 5}[rank]
	spellCoeff := [9]float64{0, .132, .204, .24, .24, .24, .24, .24, .24}[rank]
	baseDotDamage := [9]float64{0, 26, 38, 58, 86, 118, 155, 196, 230}[rank]
	spellId := [9]int32{0, 5143, 5144, 5145, 8416, 8417, 10211, 10212, 25345}[rank]
	manaCost := [9]float64{0, 85, 140, 235, 320, 410, 500, 595, 655}[rank]
	level := [9]int{0, 8, 16, 24, 32, 40, 48, 56, 56}[rank]

	tickLength := time.Second
	tickSpell := mage.getArcaneMissilesTickSpell(rank, numTicks, baseDotDamage)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolArcane,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         SpellFlagMage | core.SpellFlagAPL | core.SpellFlagNoMetrics | core.SpellFlagChanneled,
		RequiredLevel: level,
		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   1,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ArcaneMissles-" + strconv.Itoa(int(rank)) + "-" + strconv.Itoa(int(numTicks)),
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					if mage.ArcaneBlastAura != nil {
						mage.ArcaneBlastAura.Deactivate(sim)
					}
				},
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: true,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				tickSpell.Cast(sim, target)
				tickSpell.SpellMetrics[target.UnitIndex].Casts -= 1
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			tickSpell.SpellMetrics[target.UnitIndex].Casts += 1

			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
			spell.DealOutcome(sim, result)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := (baseDotDamage + (spellCoeff * spell.SpellPower()))

			return spell.CalcPeriodicDamage(sim, target, baseDamage, spell.OutcomeExpectedTick)
		},
	}
}

func (mage *Mage) ArcaneMisslesTickDuration() time.Duration {
	return mage.ApplyCastSpeed(time.Second)
}

func (mage *Mage) registerArcaneMissilesSpell() {
	maxRank := 8

	for i := 1; i < maxRank; i++ {
		config := mage.getArcaneMissilesSpellConfig(i)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ArcaneMissiles = mage.GetOrRegisterSpell(config)
		}
	}
}
