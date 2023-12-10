package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getDrainLifeBaseConfig(rank int) core.SpellConfig {
	hasRune := warlock.HasRune(proto.WarlockRune_RuneChestMasterChanneler)

	spellId := [7]int32{0, 689, 699, 709, 7651, 11699, 11700}[rank]
	spellCoeff := [7]float64{0, .078, .1, .1, .1, .1, .1}[rank]
	baseDamage := [7]float64{0, 10, 17, 29, 41, 55, 71}[rank]
	manaCost := [7]float64{0, 55, 85, 135, 185, 240, 300}[rank]
	level := [7]int{0, 14, 22, 30, 38, 46, 54}[rank]

	ticks := core.TernaryInt32(hasRune, 15, 5)

	if hasRune {
		manaCost *= 2
	}

	actionID := core.ActionID{SpellID: spellId}
	healthMetrics := warlock.NewHealthMetrics(actionID)

	spellConfig := core.SpellConfig{
		ActionID:      actionID,
		SpellSchool:   core.SpellSchoolShadow,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagHauntSE | core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
				// ChannelTime: channelTime,
			},
		},

		BonusHitRating:           float64(warlock.Talents.Suppression) * 2 * core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1,
		DamageMultiplier:         1 + 0.02*float64(warlock.Talents.ImprovedDrainLife),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Drain Life",
			},
			NumberOfTicks:       ticks,
			TickLength:          1 * time.Second,
			AffectedByCastSpeed: false,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				baseDmg := baseDamage + spellCoeff*dot.Spell.SpellPower()
				dot.SnapshotBaseDamage = baseDmg
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				result := dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTickCounted)
				health := result.Damage
				if hasRune {
					health *= 1.5
				}
				warlock.GainHealth(sim, health, healthMetrics)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.SpellMetrics[target.UnitIndex].Hits--
				dot := spell.Dot(target)
				dot.Apply(sim)
				dot.UpdateExpires(dot.ExpiresAt())

				warlock.EverlastingAfflictionRefresh(sim, target)
			}
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, useSnapshot bool) *core.SpellResult {
			if useSnapshot {
				dot := spell.Dot(target)
				return dot.CalcSnapshotDamage(sim, target, spell.OutcomeExpectedMagicAlwaysHit)
			} else {
				baseDmg := baseDamage + spellCoeff*spell.SpellPower()
				return spell.CalcPeriodicDamage(sim, target, baseDmg, spell.OutcomeExpectedMagicAlwaysHit)
			}
		},
	}

	if hasRune {
		spellConfig.Cast.CD = core.Cooldown{
			Timer:    warlock.NewTimer(),
			Duration: 15 * time.Second,
		}
	} else {
		spellConfig.Flags |= core.SpellFlagChanneled
		spellConfig.Cast.DefaultCast.ChannelTime = time.Second * 5
	}

	return spellConfig
}

func (warlock *Warlock) registerDrainLifeSpell() {
	maxRank := 6

	for i := 1; i <= maxRank; i++ {
		config := warlock.getDrainLifeBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.DrainLife = warlock.GetOrRegisterSpell(config)
		}
	}
}
