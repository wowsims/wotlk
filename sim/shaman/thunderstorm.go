package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var thunderstormActionID = core.ActionID{SpellID: 51490}

// newThunderstormSpell returns a precomputed instance of lightning bolt to use for casting.
func (shaman *Shaman) newThunderstormSpell(doDamage bool) *core.Spell {
	manaRestore := 0.08
	if shaman.HasMinorGlyph(proto.ShamanMinorGlyph_GlyphOfThunderstorm) {
		manaRestore = 0.1
	}

	manaMetrics := shaman.NewManaMetrics(thunderstormActionID)

	spellConfig := core.SpellConfig{
		ActionID:    thunderstormActionID,
		SpellSchool: core.SpellSchoolNature,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 45,
			},
		},
		ApplyEffects: func(sim *core.Simulation, u *core.Unit, s2 *core.Spell) {
			shaman.AddMana(sim, shaman.MaxMana()*manaRestore, manaMetrics, true)
		},
	}

	if doDamage {
		effect := core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: float64(shaman.Talents.ElementalPrecision) * 2 * core.SpellHitRatingPerHitChance,
			DamageMultiplier:    1 * (1 + 0.01*float64(shaman.Talents.Concussion)),
			ThreatMultiplier:    1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),
			BaseDamage:          core.BaseDamageConfigMagic(566, 644, 0.172),
			OutcomeApplier:      shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier(0)),
		}
		aoeApply := core.ApplyEffectFuncAOEDamageCapped(shaman.Env, (605+0.172*shaman.GetStat(stats.NatureSpellPower))*10, effect)
		spellConfig.ApplyEffects = func(sim *core.Simulation, unit *core.Unit, spell *core.Spell) {
			aoeApply(sim, unit, spell)                                           // Calculates hits/crits/dmg on each target
			shaman.AddMana(sim, shaman.MaxMana()*manaRestore, manaMetrics, true) // adds mana no matter what
		}
	}

	return shaman.RegisterSpell(spellConfig)
}
