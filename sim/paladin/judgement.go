package paladin

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
)

// const JudgementDuration = time.Second * 20

// Shared conditions required to be able to cast any Judgement.
//
//nolint:unused
func (paladin *Paladin) canJudgement(sim *core.Simulation) bool {
	return paladin.CurrentSeal != nil && paladin.CurrentSeal.IsActive() && paladin.JudgementOfLight.IsReady(sim)
}

func (paladin *Paladin) registerJudgementOfWisdomSpell(cdTimer *core.Timer) {
	jowAuras := paladin.NewEnemyAuraArray(core.JudgementOfWisdomAura)

	paladin.JudgementOfWisdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 53408},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc, // can proc TaJ itself and from seal
		Flags:       SpellFlagPrimaryJudgement | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: cdTimer,
				Duration: (time.Second * 10) -
					(time.Second * time.Duration(paladin.Talents.ImprovedJudgements)) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetRedemptionBattlegear, 4), 1*time.Second, 0) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetGladiatorsVindication, 4), 1*time.Second, 0),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Primary Judgements cannot crit or be dodged, parried, or blocked-- only miss. (Unless target is a hunter.)
			jow := jowAuras.Get(target)
			if jow.IsActive() {
				jow.Refresh(sim)
			} else {
				jow.Activate(sim)
			}
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)
		},

		RelatedAuras: []core.AuraArray{jowAuras},
	})
}

func (paladin *Paladin) registerJudgementOfLightSpell(cdTimer *core.Timer) {
	jolAuras := paladin.NewEnemyAuraArray(core.JudgementOfLightAura)

	paladin.JudgementOfLight = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 20271},
		SpellSchool: core.SpellSchoolHoly,
		ProcMask:    core.ProcMaskProc,
		Flags:       SpellFlagPrimaryJudgement | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.05,
			Multiplier: 1 - 0.02*float64(paladin.Talents.Benediction),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer: cdTimer,
				Duration: (time.Second * 10) -
					(time.Second * time.Duration(paladin.Talents.ImprovedJudgements)) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetRedemptionBattlegear, 4), 1*time.Second, 0) -
					core.TernaryDuration(paladin.HasSetBonus(ItemSetGladiatorsVindication, 4), 1*time.Second, 0),
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			// Primary Judgements cannot crit or be dodged, parried, or blocked-- only miss. (Unless target is a hunter.)
			jol := jolAuras.Get(target)
			if jol.IsActive() {
				jol.Refresh(sim)
			} else {
				jol.Activate(sim)
			}
			spell.CalcAndDealOutcome(sim, target, spell.OutcomeRangedHit)
		},

		RelatedAuras: []core.AuraArray{jolAuras},
	})
}

// Defines judgement refresh behavior from attacks
// Returns extra mana if a different pally applied Judgement of Wisdom
// func (paladin *Paladin) setupJudgementRefresh() {
// 	const mana = 74 / 2
// 	paladin.RegisterAura(core.Aura{
// 		Label:    "Refresh Judgement",
// 		Duration: core.NeverExpires,
// 		OnReset: func(aura *core.Aura, sim *core.Simulation) {
// 			aura.Activate(sim)
// 		},
// 		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
// 			if result.Landed() && spell.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
// 				if paladin.CurrentJudgement != nil && paladin.CurrentJudgement.IsActive() {
// 					// Refresh the judgement
// 					paladin.CurrentJudgement.Refresh(sim)

// 					// Check if current judgement is not JoW and also that JoW is on the target
// 					if paladin.CurrentJudgement != paladin.JudgementOfWisdomAura && paladin.JudgementOfWisdomAura.IsActive() {
// 						// Just trigger a second JoW
// 						if paladin.JowManaMetrics == nil {
// 							paladin.JowManaMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 27164})
// 						}
// 						paladin.AddMana(sim, mana, paladin.JowManaMetrics, false)
// 					}
// 				}
// 			}
// 		},
// 	})
// }

func (paladin *Paladin) registerJudgements() {
	// Shared CD for all judgements.
	cdTimer := paladin.NewTimer()
	paladin.registerJudgementOfWisdomSpell(cdTimer)
	paladin.registerJudgementOfLightSpell(cdTimer)
}
