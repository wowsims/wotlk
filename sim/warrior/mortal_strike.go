package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warrior *Warrior) registerMortalStrikeSpell(cdTimer *core.Timer) {
	cost := 30.0
	refundAmount := cost * 0.8

	warrior.MortalStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 71552},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ResourceType: stats.Rage,
		BaseCost:     cost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: cost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: time.Second*6 - time.Millisecond*200*time.Duration(warrior.Talents.ImprovedMortalStrike),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeMHSpecial,

			DamageMultiplier: 1 *
				(1 + 0.01*float64(warrior.Talents.ImprovedMortalStrike)) *
				core.TernaryFloat64(warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfMortalStrike), 1.1, 1) *
				core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtBattlegear, 4), 1.05, 1),
			BonusCritRating:  core.TernaryFloat64(warrior.HasSetBonus(ItemSetSiegebreakerBattlegear, 4), 10, 0) * core.CritRatingPerCritChance,
			ThreatMultiplier: 1,
			DynamicThreatMultiplier: func(spellEffect *core.SpellEffect, spell *core.Spell) float64 {
				if warrior.StanceMatches(DefensiveStance) {
					return 1 + 0.21*float64(warrior.Talents.TacticalMastery)
				}
				return 1.0
			},

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 80, 1, 1, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanMortalStrike(sim *core.Simulation) bool {
	if warrior.Talents.MortalStrike {
		return warrior.CurrentRage() >= warrior.MortalStrike.DefaultCast.Cost && warrior.MortalStrike.IsReady(sim)
	} else {
		return false
	}
}
