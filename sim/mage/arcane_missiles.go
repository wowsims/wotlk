package mage

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// Note: AM doesn't charge its mana up-front, instead it charges 1/5 of the mana on each tick.
// This is probably not worth simming since no other spell in the game does this and AM isn't
// even a popular choice for arcane mages.

// I don't think the above note is true
func (mage *Mage) registerArcaneMissilesSpell() {
	actionID := core.ActionID{SpellID: 42846}
	baseCost := .31 * mage.BaseMana

	bonusCrit := 0.0
	if mage.MageTier.t9_4 {
		bonusCrit += 5 * core.CritRatingPerCritChance
	}

	// bonusCrit := float64(mage.Talents.ArcanePotency) * 10 * core.CritRatingPerCritChance

	mage.ArcaneMissiles = mage.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolArcane,
		Flags:       SpellFlagMage | core.SpellFlagChanneled,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost * (1 - .01*float64(mage.Talents.ArcaneFocus)),

				GCD:         core.GCDDefault,
				ChannelTime: time.Second * 5,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				if mage.MissileBarrageAura.IsActive() {
					if mage.MageTier.t10_2 {
						bloodmageHasteAura.Activate(sim)
					}
					mage.PseudoStats.NoCost = true
					cast.ChannelTime = cast.ChannelTime / 2
				}
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: float64(mage.Talents.ArcaneFocus+FrostTalents.Precision) * core.SpellHitRatingPerHitChance,

			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			OutcomeApplier: mage.OutcomeFuncMagicHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					mage.PseudoStats.NoCost = false
					if mage.MissileBarrageAura.IsActive() {
						mage.isMissilesBarrage = true
						mage.MultiplyCastSpeed(2)
					}

					mage.ArcaneMissilesDot.Apply(sim)
				}
			},
		}),
	})

	bonusCritDamage := mage.bonusCritDamage
	if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfArcaneMissiles) {
		bonusCritDamage += .25
	}
	target := mage.CurrentTarget
	mage.ArcaneMissilesDot = core.NewDot(core.Dot{
		Spell: mage.ArcaneMissiles,
		Aura: target.RegisterAura(core.Aura{
			Label:    "ArcaneMissiles-" + strconv.Itoa(int(mage.Index)),
			ActionID: actionID,
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if mage.isMissilesBarrage {
					mage.MultiplyCastSpeed(.5)
					mage.isMissilesBarrage = false
					if !mage.MageTier.t8_4 || sim.RandomFloat("MageT84PC") > .1 {
						mage.MissileBarrageAura.Deactivate(sim)
					}
				}
				mage.ArcaneBlastAura.Deactivate(sim)
			},
		}),

		NumberOfTicks:       5,
		TickLength:          time.Second,
		AffectedByCastSpeed: true,

		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:             core.ProcMaskSpellDamage,
			BonusSpellHitRating:  float64(mage.Talents.ArcaneFocus) * core.SpellHitRatingPerHitChance,
			BonusSpellCritRating: bonusCrit,

			DamageMultiplier: mage.spellDamageMultiplier * (1 + .04*float64(mage.Talents.TormentTheWeak)),
			ThreatMultiplier: 1 - 0.2*float64(mage.Talents.ArcaneSubtlety),

			BaseDamage:     core.BaseDamageConfigMagicNoRoll(362, 1/3.5+0.03*float64(mage.Talents.ArcaneEmpowerment)),
			OutcomeApplier: mage.OutcomeFuncMagicHitAndCrit(mage.SpellCritMultiplier(1, bonusCritDamage)),
		})),
	})
}
