package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (warlock *Warlock) registerImmolateSpell() {
	actionID := core.ActionID{SpellID: 27215}
	baseCost := 445.0

	effect := core.SpellEffect{
		BonusSpellCritRating: float64(warlock.Talents.Devastation) * 1 * core.SpellCritRatingPerCritChance,
		DamageMultiplier: 1 *
			(1 + (0.05 * float64(warlock.Talents.ImprovedImmolate))) *
			(1 + (0.02 * float64(warlock.Talents.Emberstorm))),
		ThreatMultiplier: 1 - 0.05*float64(warlock.Talents.DestructiveReach),
		BaseDamage:       core.BaseDamageConfigMagic(332.0, 332.0, 0.2+0.04*float64(warlock.Talents.ShadowAndFlame)),
		OutcomeApplier:   warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, core.TernaryFloat64(warlock.Talents.Ruin, 0, 1))),
		OnSpellHitDealt:  applyDotOnLanded(&warlock.ImmolateDot),
		ProcMask:         core.ProcMaskSpellDamage,
	}

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - 0.01*float64(warlock.Talents.Cataclysm)),
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond*2000 - (time.Millisecond * 100 * time.Duration(warlock.Talents.Bane)),
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	target := warlock.CurrentTarget

	// DOT: 615 dmg over 15s (123 every 3 sec, mod 0.13)
	warlock.ImmolateDot = core.NewDot(core.Dot{
		Spell: warlock.Immolate,
		Aura: target.RegisterAura(core.Aura{
			Label:    "immolate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + core.TernaryInt(ItemSetVoidheartRaiment.CharacterHasSetBonus(&warlock.Character, 4), 1, 0), // voidheart 4p gives 1 extra tick
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(615/5, 0.13),
			OutcomeApplier:   warlock.OutcomeFuncTick(),
			IsPeriodic:       true,
			ProcMask:         core.ProcMaskPeriodicDamage,
		}),
	})
}
