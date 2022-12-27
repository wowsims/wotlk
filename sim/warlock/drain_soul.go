package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerDrainSoulSpell() {
	actionID := core.ActionID{SpellID: 47855}
	soulSiphonMultiplier := 0.03 * float64(warlock.Talents.SoulSiphon)
	baseCost := warlock.BaseMana * 0.14

	warlock.DrainSoul = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - 0.02*float64(warlock.Talents.Suppression)),
				GCD:  core.GCDDefault,
				// ChannelTime: channelTime,
			},
		},

		DamageMultiplierAdditive: 1 +
			warlock.GrandSpellstoneBonus() +
			0.03*float64(warlock.Talents.ShadowMastery),
		// For performance optimization, the execute modifier is basekit since we never use it before execute
		DamageMultiplier: (4.0 + 0.04*float64(warlock.Talents.DeathsEmbrace)) / (1 + 0.04*float64(warlock.Talents.DeathsEmbrace)),
		ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.ImprovedDrainSoul),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				warlock.DrainSoulDot.Apply(sim)
				warlock.DrainSoulDot.Aura.UpdateExpires(warlock.DrainSoulDot.Aura.ExpiresAt())
			}
		},
	})

	warlock.DrainSoulDot = core.NewDot(core.Dot{
		Spell: warlock.DrainSoul,
		Aura: warlock.CurrentTarget.RegisterAura(core.Aura{
			Label:    "Drain Soul-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       5,
		TickLength:          3 * time.Second,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
			baseDmg := 142 + 0.429*dot.Spell.SpellPower()

			auras := []*core.Aura{
				warlock.HauntDebuffAura,
				warlock.UnstableAfflictionDot.Aura,
				warlock.CorruptionDot.Aura,
				warlock.SeedDots[target.Index].Aura,
				warlock.CurseOfAgonyDot.Aura,
				warlock.CurseOfDoomDot.Aura,
				warlock.CurseOfElementsAura,
				warlock.CurseOfWeaknessAura,
				warlock.CurseOfTonguesAura,
				warlock.ShadowEmbraceDebuffAura(target),
				// missing: death coil
			}
			numActive := 0
			for _, aura := range auras {
				if aura.IsActive() {
					numActive++
				}
			}
			dot.SnapshotBaseDamage = baseDmg * (1.0 + float64(core.MinInt(3, numActive))*soulSiphonMultiplier)

			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
		},
	})
}
