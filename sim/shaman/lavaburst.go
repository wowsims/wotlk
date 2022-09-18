package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var lavaBurstActionID = core.ActionID{SpellID: 60043}

func (shaman *Shaman) registerLavaBurstSpell() {
	baseCost := baseMana * 0.1

	spellConfig := core.SpellConfig{
		ActionID:     lavaBurstActionID,
		SpellSchool:  core.SpellSchoolFire,
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
		ThreatMultiplier: 1 - (0.1/3)*float64(shaman.Talents.ElementalPrecision),
	}

	bonusBase := core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == VentureCoLightningRod, 121, 0) +
		core.TernaryFloat64(shaman.Equip[items.ItemSlotRanged].ID == ThunderfallTotem, 215, 0)

	effect := core.SpellEffect{
		ProcMask:   core.ProcMaskSpellDamage,
		BaseDamage: core.BaseDamageConfigMagic(1192+bonusBase, 1518+bonusBase, 0.5714+(0.05*float64(shaman.Talents.Shamanism)+core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLava), 0.1, 0))),
		// TODO: does lava flows multiply or add with elemental fury? Only matters if you had <5pts which probably won't happen.
		OutcomeApplier: shaman.OutcomeFuncMagicHitAndCrit(shaman.ElementalCritMultiplier([]float64{0, 0.06, 0.12, 0.24}[shaman.Talents.LavaFlows] + core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthShatterGarb, 4), 0.1, 0))),
	}

	if shaman.HasSetBonus(ItemSetThrallsRegalia, 4) {
		lvbdotDmg := 0.0 // dynamically changing dmg
		spell := shaman.RegisterSpell(core.SpellConfig{
			Flags:    core.SpellFlagIgnoreModifiers,
			ActionID: core.ActionID{SpellID: 71824},

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
		})
		lvbdot := core.NewDot(core.Dot{
			Spell: spell,
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
				ProcMask:       core.ProcMaskEmpty,
				OutcomeApplier: shaman.OutcomeFuncTick(),
			}),
		})

		effect.OnSpellHitDealt = func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			lvbdotDmg = spellEffect.Damage * 0.1 // TODO: does this dot pool with the previous dot?
			lvbdot.Apply(sim)                    // will resnapshot dmg
		}
	}

	spellConfig.ApplyEffects = core.ApplyEffectFuncDirectDamage(effect)
	shaman.LavaBurst = shaman.RegisterSpell(spellConfig)
}
