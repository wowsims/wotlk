package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

// TODO Mind Flay (48156) now "periodically triggers" Mind Flay (58381), probably to allow haste to work.
//  The first never deals damage, so the latter should probably be used as ActionID here.
func (priest *Priest) MindFlayActionID(numTicks int) core.ActionID {
	return core.ActionID{SpellID: 48156, Tag: int32(numTicks)}
}

func (priest *Priest) newMindFlaySpell(numTicks int) *core.Spell {
	baseCost := priest.BaseMana * 0.09

	channelTime := time.Second * time.Duration(numTicks)
	for _, gem := range priest.Equip[proto.ItemSlot_ItemSlotHead].Gems {
		if gem.ID == 25895 || gem.ID == 41335 {
			channelTime -= time.Duration(numTicks) * (time.Millisecond * 100)
		}
	}
	// ADDED TROLL MF BUG 15% REDUCED CHANNEL TIME DUE TO DA VOODOO SHUFFLE
	if priest.GetCharacter().Race == proto.Race_RaceTroll {
		channelTime -= time.Duration(numTicks) * (time.Millisecond * 150)
	}
	if priest.HasSetBonus(ItemSetCrimsonAcolyte, 4) {
		channelTime -= time.Duration(numTicks) * (time.Millisecond * 170)
	}

	effect := core.SpellEffect{
		OutcomeApplier: priest.OutcomeFuncMagicHit(),
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if priest.ShadowWordPainDot.IsActive() {
				if priest.Talents.PainAndSuffering == 3 {
					priest.ShadowWordPainDot.Rollover(sim)
				} else if sim.RandomFloat("Pain and Suffering") < (float64(priest.Talents.PainAndSuffering) * 0.33) {
					priest.ShadowWordPainDot.Rollover(sim)
				}
			}
			priest.MindFlayDot[numTicks].Apply(sim)
		},
	}

	return priest.RegisterSpell(core.SpellConfig{
		ActionID:     priest.MindFlayActionID(numTicks),
		SpellSchool:  core.SpellSchoolShadow,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagChanneled,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:        baseCost * (1 - 0.05*float64(priest.Talents.FocusedMind)),
				GCD:         core.GCDDefault,
				ChannelTime: channelTime,
			},
			ModifyCast: func(sim *core.Simulation, spell *core.Spell, cast *core.Cast) {
				// if our channel is longer than GCD it will have human latency to end it beause you can't queue the next spell.
				wait := priest.ApplyCastSpeed(channelTime)
				gcd := core.MaxDuration(core.GCDMin, priest.ApplyCastSpeed(core.GCDDefault))
				if wait > gcd && priest.Latency > 0 {
					base := priest.Latency * 0.25
					variation := base + sim.RandomFloat("spriest latency")*base // should vary from 0.66 - 1.33 of given latency
					variation = core.MaxFloat(variation, 10)                    // no player can go under XXXms response time
					cast.AfterCastDelay += time.Duration(variation) * time.Millisecond
				}
			},
		},

		BonusHitRating:   float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
		BonusCritRating:  float64(priest.Talents.MindMelt)*2*core.CritRatingPerCritChance + core.TernaryFloat64(priest.HasSetBonus(ItemSetZabras, 4), 5, 0)*core.CritRatingPerCritChance,
		DamageMultiplier: 1,
		CritMultiplier:   priest.SpellCritMultiplier(1, float64(priest.Talents.ShadowPower)/5),
		ThreatMultiplier: 1 - 0.08*float64(priest.Talents.ShadowAffinity),

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (priest *Priest) newMindFlayDot(numTicks int) *core.Dot {
	target := priest.CurrentTarget

	normalCoeff := 0.257
	miseryCoeff := normalCoeff * (1 + 0.05*float64(priest.Talents.Misery))

	normMod := 1 + float64(priest.Talents.Darkness)*0.02 + float64(priest.Talents.TwinDisciplines)*0.01 // initialize modifier
	swpMod := normMod * (1 +
		core.TernaryFloat64(priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfMindFlay)), 0.1, 0) +
		0.02*float64(priest.Talents.TwistedFaith))

	hasGlyphOfShadow := priest.HasGlyph(int32(proto.PriestMajorGlyph_GlyphOfShadow))

	var mfReducTime time.Duration
	mfReducTime = time.Millisecond * 0
	// ADDED Bug where root resist gem reduces MF by 10%
	// ADDED TROLL MF BUG 15% REDUCED CHANNEL TIME DUE TO DA VOODOO SHUFFLE
	if priest.GetCharacter().Race == proto.Race_RaceTroll {
		mfReducTime = time.Millisecond * 150
	}

	for _, gem := range priest.Equip[proto.ItemSlot_ItemSlotHead].Gems {
		if gem.ID == 25895 || gem.ID == 41335 {
			mfReducTime = time.Millisecond * 100
			if priest.GetCharacter().Race == proto.Race_RaceTroll {
				mfReducTime = time.Millisecond*150 + time.Millisecond*100
			}
		}
	}
	if priest.HasSetBonus(ItemSetCrimsonAcolyte, 4) {
		mfReducTime = mfReducTime + time.Millisecond*170
	}

	return core.NewDot(core.Dot{
		Spell: priest.MindFlay[numTicks],
		Aura: target.RegisterAura(core.Aura{
			Label:    "MindFlay-" + strconv.Itoa(numTicks) + "-" + strconv.Itoa(int(priest.Index)),
			ActionID: priest.MindFlayActionID(numTicks),
		}),

		NumberOfTicks:       numTicks,
		TickLength:          time.Second - mfReducTime,
		AffectedByCastSpeed: true,

		OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
			dmg := 588.0 / 3
			if priest.MiseryAura.IsActive() {
				dmg += miseryCoeff * dot.Spell.SpellPower()
			} else {
				dmg += normalCoeff * dot.Spell.SpellPower()
			}
			if priest.ShadowWordPainDot.IsActive() {
				dmg *= swpMod
			} else {
				dmg *= normMod
			}
			dot.SnapshotBaseDamage = dmg * (1 + 0.02*float64(priest.ShadowWeavingAura.GetStacks()))

			dot.SnapshotCritChance = dot.Spell.SpellCritChance(target)
			dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
		},
		OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
			result := dot.CalcSnapshotDamage(sim, target, dot.OutcomeMagicHitAndSnapshotCrit)
			dot.Spell.DealPeriodicDamage(sim, &result)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim)
			}
			if result.DidCrit() && hasGlyphOfShadow {
				priest.ShadowyInsightAura.Activate(sim)
			}
			if result.DidCrit() && priest.ImprovedSpiritTap != nil && sim.RandomFloat("Improved Spirit Tap") > 0.5 {
				priest.ImprovedSpiritTap.Activate(sim)
			}
		},
	})
}

func (priest *Priest) MindFlayTickDuration() time.Duration {
	mfReducTime := time.Millisecond * 0
	for _, gem := range priest.Equip[proto.ItemSlot_ItemSlotHead].Gems {
		if gem.ID == 25895 || gem.ID == 41335 {
			mfReducTime = time.Millisecond * 100
		}
	}
	return priest.ApplyCastSpeed(time.Second - core.TernaryDuration(priest.T10FourSetBonus, time.Millisecond*170, 0) - core.TernaryDuration(priest.GetCharacter().Race == proto.Race_RaceTroll, time.Millisecond*150, 0) - mfReducTime)
}
