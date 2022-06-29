package warrior

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warrior *Warrior) registerMortalStrikeSpell(cdTimer *core.Timer) {
	cost := 30.0
	if ItemSetDestroyerBattlegear.CharacterHasSetBonus(&warrior.Character, 4) {
		cost -= 5
	}
	refundAmount := cost * 0.8

	warrior.MortalStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 30330},
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
				core.TernaryFloat64(ItemSetOnslaughtBattlegear.CharacterHasSetBonus(&warrior.Character, 4), 1.05, 1),
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, true, 210, 1, true),
			OutcomeApplier: warrior.OutcomeFuncMeleeWeaponSpecialHitAndCrit(warrior.critMultiplier(true)),

			OnInit: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if warrior.StanceMatches(DefensiveStance) {
					spellEffect.ThreatMultiplier *= 1 + 0.21*float64(warrior.Talents.TacticalMastery)
				}
			},
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					warrior.AddRage(sim, refundAmount, warrior.RageRefundMetrics)
				}
			},
		}),
	})
}

func (warrior *Warrior) CanMortalStrike(sim *core.Simulation) bool {
	return warrior.Talents.MortalStrike && warrior.CurrentRage() >= warrior.MortalStrike.DefaultCast.Cost && warrior.MortalStrike.IsReady(sim)
}
