package warrior

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (warrior *Warrior) registerShieldSlamSpell() {
	hasGlyph := warrior.HasMajorGlyph(proto.WarriorMajorGlyph_GlyphOfBlocking)
	var glyphOfBlockingAura *core.Aura = nil
	if hasGlyph {
		glyphOfBlockingAura = warrior.GetOrRegisterAura(core.Aura{
			Label:    "Glyph of Blocking",
			ActionID: core.ActionID{SpellID: 58397},
			Duration: 10 * time.Second,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				warrior.PseudoStats.BlockValueMultiplier += 0.1
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				warrior.PseudoStats.BlockValueMultiplier -= 0.1
			},
		})
	}

	warrior.ShieldSlam = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 47488},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial, // TODO: Is this right?
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL,

		RageCost: core.RageCostOptions{
			Cost:   20 - float64(warrior.Talents.FocusedRage),
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock
		},

		BonusCritRating: 5 * core.CritRatingPerCritChance * float64(warrior.Talents.CriticalBlock),
		DamageMultiplier: 1 +
			.05*float64(warrior.Talents.GagOrder) +
			core.TernaryFloat64(warrior.HasSetBonus(ItemSetOnslaughtArmor, 4), .10, 0) +
			core.TernaryFloat64(warrior.HasSetBonus(ItemSetDreadnaughtPlate, 2), .10, 0) +
			core.TernaryFloat64(warrior.HasSetBonus(ItemSetYmirjarLordsPlate, 2), .20, 0), // TODO: All additive multipliers?
		CritMultiplier:   warrior.critMultiplier(mh),
		ThreatMultiplier: 1.3,
		FlatThreatBonus:  770,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {

			// Apply SBV cap with special bypass rules for Shield Block and Glyph of Blocking
			// TODO: Verify that this bypass behavior and DR curve are correct

			sbvMod := warrior.PseudoStats.BlockValueMultiplier
			sbvMod /= (sbvMod - core.TernaryFloat64(warrior.ShieldBlockAura.IsActive(), 1, 0) - core.TernaryFloat64(glyphOfBlockingAura.IsActive(), 0.1, 0))

			sbv := warrior.BlockValue() / sbvMod

			sbv = sbvMod * (core.TernaryFloat64(sbv <= 1960.0, sbv, 0.0) + core.TernaryFloat64(sbv > 1960.0 && sbv <= 3160.0, 0.09333333333*sbv+1777.06666667, 0.0) + core.TernaryFloat64(sbv > 3160.0, 2072.0, 0.0))

			baseDamage := sim.Roll(990, 1040) + sbv
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				if glyphOfBlockingAura != nil {
					glyphOfBlockingAura.Activate(sim)
				}
			} else {
				spell.IssueRefund(sim)
			}
		},
	})
}
