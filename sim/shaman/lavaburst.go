package shaman

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
)

func (shaman *Shaman) registerLavaBurstSpell() {
	actionID := core.ActionID{SpellID: 60043}
	dmgBonus := core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == VentureCoLightningRod, 121, 0) +
		core.TernaryFloat64(shaman.Equip[core.ItemSlotRanged].ID == ThunderfallTotem, 215, 0)
	spellCoeff := 0.5714 +
		0.05*float64(shaman.Talents.Shamanism) +
		core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfLava), 0.1, 0)

	var lvbDot *core.Dot
	if shaman.HasSetBonus(ItemSetThrallsRegalia, 4) {
		dotSpell := shaman.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 71824},
			SpellSchool: core.SpellSchoolFire,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagIgnoreModifiers,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				lvbDot.Apply(sim)
				spell.CalcAndDealOutcome(sim, target, spell.OutcomeAlwaysHit)
			},
		})
		lvbDot = core.NewDot(core.Dot{
			Spell: dotSpell,
			Aura: shaman.CurrentTarget.RegisterAura(core.Aura{
				Label:    "LavaBursted-" + strconv.Itoa(int(shaman.Index)),
				ActionID: dotSpell.ActionID,
			}),
			TickLength:    time.Second * 2,
			NumberOfTicks: 3,

			SnapshotAttackerMultiplier: 1,

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		})
	}

	shaman.LavaBurst = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagFocusable,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.1,
			Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
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
		DamageMultiplier: 1 + 0.01*float64(shaman.Talents.Concussion) + 0.02*float64(shaman.Talents.CallOfFlame),
		CritMultiplier:   shaman.ElementalCritMultiplier([]float64{0, 0.06, 0.12, 0.24}[shaman.Talents.LavaFlows] + core.TernaryFloat64(shaman.HasSetBonus(ItemSetEarthShatterGarb, 4), 0.1, 0)),
		ThreatMultiplier: shaman.spellThreatMultiplier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := dmgBonus + sim.Roll(1192, 1518) + spellCoeff*spell.SpellPower()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if lvbDot != nil && result.Landed() {
				lvbDot.SnapshotBaseDamage = result.Damage * 0.1 / float64(lvbDot.NumberOfTicks)
				lvbDot.Spell.Cast(sim, target)
			}
			spell.DealDamage(sim, result)
		},
	})
}
