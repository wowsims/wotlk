package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
)

func (warlock *Warlock) getCurseOfAgonyBaseConfig(rank int) core.SpellConfig {
	spellId := [7]int32{0, 980, 1014, 6217, 11711, 11712, 11713}[rank]
	spellCoeff := [7]float64{0, .046, .077, .083, .083, .083, .083}[rank]
	baseDamage := [7]float64{0, 7, 15, 27, 42, 65, 87}[rank]
	manaCost := [7]float64{0, 25, 50, 90, 130, 170, 215}[rank]
	level := [7]int{0, 8, 18, 28, 38, 48, 58}[rank]

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolShadow,
		Flags:         core.SpellFlagAPL | core.SpellFlagHauntSE | core.SpellFlagResetAttackSwing | core.SpellFlagPureDot,
		ProcMask:      core.ProcMaskSpellDamage,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		BonusHitRating:   2 * float64(warlock.Talents.Suppression) * core.CritRatingPerCritChance,
		DamageMultiplier: 1 + 0.02*float64(warlock.Talents.ImprovedCurseOfWeakness),
		ThreatMultiplier: 1,
		FlatThreatBonus:  0,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "CurseofAgony-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks: 12,
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 0.5 * (baseDamage + spellCoeff*dot.Spell.SpellPower())
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])

				if warlock.AmplifyCurseAura.IsActive() {
					dot.SnapshotAttackerMultiplier *= 1.5
					warlock.AmplifyCurseAura.Deactivate(sim)
				}
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				if dot.TickCount%4 == 0 { // CoA ramp up
					dot.SnapshotBaseDamage += 0.5 * (baseDamage + spellCoeff*dot.Spell.SpellPower())
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				//warlock.CurseOfDoom.Dot(target).Cancel(sim)
				spell.Dot(target).Apply(sim)
			}
		},
	}
}

func (warlock *Warlock) registerCurseOfAgonySpell() {
	maxRank := 6

	for i := 1; i <= maxRank; i++ {
		config := warlock.getCurseOfAgonyBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.CurseOfAgony = warlock.GetOrRegisterSpell(config)
		}
	}
}

func (warlock *Warlock) registerAmplifyCurseSpell() {
	if !warlock.Talents.AmplifyCurse {
		return
	}

	actionID := core.ActionID{SpellID: 18288}

	warlock.AmplifyCurseAura = warlock.GetOrRegisterAura(core.Aura{
		Label:    "Amplify Curse",
		ActionID: actionID,
		Duration: time.Second * 30,
	})

	warlock.AmplifyCurse = warlock.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagAPL,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    warlock.NewTimer(),
				Duration: 3 * time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			warlock.AmplifyCurseAura.Activate(sim)
		},
	})
}
