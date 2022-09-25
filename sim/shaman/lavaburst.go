package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (shaman *Shaman) registerLavaBurstSpell() {
	actionID := core.ActionID{SpellID: 60043}
	baseCost := baseMana * 0.1
	dmgBonus := core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == VentureCoLightningRod, 121, 0) +
		core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == ThunderfallTotem, 215, 0)
	spellCoeff := 0.5714 +
		0.05*float64(shaman.Talents.Shamanism) +
		core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLava), 0.1, 0)

	applyDot := shaman.HasSetBonus(ItemSetThrallsRegalia, 4)
	lvbdotDmg := 0.0 // dynamically changing dmg
	var lvbDot *core.Dot
	if applyDot {
		dotSpell := shaman.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{SpellID: 71824},
			// TODO: No spell school?
			ProcMask: core.ProcMaskEmpty,
			Flags:    core.SpellFlagIgnoreModifiers,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
		})
		lvbDot = core.NewDot(core.Dot{
			Spell: dotSpell,
			Aura: shaman.CurrentTarget.RegisterAura(core.Aura{
				Label:    "LavaBursted-" + strconv.Itoa(int(shaman.Index)),
				ActionID: core.ActionID{SpellID: 71824},
			}),
			TickLength:    time.Second * 2,
			NumberOfTicks: 3,
			TickEffects: core.TickFuncSnapshot(shaman.CurrentTarget, core.SpellEffect{
				BaseDamage: core.BaseDamageConfig{
					Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
						return lvbdotDmg / 3 //spread dot over 3 ticks
					},
				},
				IsPeriodic:     true,
				OutcomeApplier: shaman.OutcomeFuncTick(),
			}),
		})
	}

	shaman.LavaBurst = shaman.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagFocusable,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * (1 - float64(shaman.Talents.Convection)*0.02),
				CastTime: time.Second*2 - time.Millisecond*100*time.Duration(shaman.Talents.LightningMastery),
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * 8,
			},
			ModifyCast: func(_ *core.Simulation, spell *core.Spell, cast *core.Cast) {
				shaman.modifyCastClearcasting(spell, cast)
				if shaman.ElementalMasteryAura.IsActive() {
					cast.CastTime = 0
				} else if shaman.NaturesSwiftnessAura.IsActive() {
					cast.CastTime = 0
				}
			},
		},

		BonusHitRating:   float64(shaman.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance,
		DamageMultiplier: 1 * (1 + 0.01*float64(shaman.Talents.Concussion)) * (1.0 + 0.02*float64(shaman.Talents.CallOfFlame)),
		// TODO: does lava flows multiply or add with elemental fury? Only matters if you had <5pts which probably won't happen.
		CritMultiplier:   shaman.ElementalCritMultiplier([]float64{0, 0.06, 0.12, 0.24}[shaman.Talents.LavaFlows] + core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthShatterGarb, 4), 0.1, 0)),
		ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dmgBonus + sim.Roll(1192, 1518) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamageMagicHitAndCrit(sim, target, baseDamage)
			if applyDot && result.Landed() {
				lvbdotDmg = result.Damage * 0.1 // TODO: does this dot pool with the previous dot?
				lvbDot.Apply(sim)               // will resnapshot dmg
			}
			spell.DealDamage(sim, &result)
		},
	})
}
