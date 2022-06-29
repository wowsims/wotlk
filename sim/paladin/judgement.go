package paladin

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

const JudgementManaCost = 147.0
const JudgementCDTime = time.Second * 10
const JudgementDuration = time.Second * 20

// Shared conditions required to be able to cast any Judgement.
func (paladin *Paladin) canJudgement(sim *core.Simulation) bool {
	return paladin.CurrentSeal != nil && paladin.CurrentSeal.IsActive() && paladin.JudgementOfWisdom.IsReady(sim)
}

func (paladin *Paladin) registerJudgementOfBloodSpell(cdTimer *core.Timer, sanctifiedJudgementMetrics *core.ResourceMetrics) {
	effect := core.SpellEffect{
		ProcMask: core.ProcMaskMeleeOrRangedSpecial,

		BonusCritRating:  3 * core.MeleeCritRatingPerCritChance * float64(paladin.Talents.Fanaticism),
		DamageMultiplier: paladin.WeaponSpecializationMultiplier(),
		ThreatMultiplier: 1,

		BaseDamage:     core.BaseDamageConfigMagic(295, 325, 0.429),
		OutcomeApplier: paladin.OutcomeFuncMeleeSpecialCritOnly(paladin.MeleeCritMultiplier()),

		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			paladin.sanctifiedJudgement(sim, sanctifiedJudgementMetrics, paladin.SealOfBlood.DefaultCast.Cost)
			paladin.SealOfBloodAura.Deactivate(sim)
			paladin.CurrentSeal = nil

			// Add mana from Spiritual Attunement
			// 33% of damage is self-inflicted, 10% of self-inflicted damage is returned as mana
			paladin.AddMana(sim, spellEffect.Damage*0.33*0.1, paladin.SpiritualAttunementMetrics, false)
		},
	}

	baseCost := core.TernaryFloat64(ItemSetCrystalforgeBattlegear.CharacterHasSetBonus(&paladin.Character, 2), JudgementManaCost-35, JudgementManaCost)
	paladin.JudgementOfBlood = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 31898},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagJudgement,

		ResourceType: stats.Mana,
		BaseCost:     JudgementManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost - JudgementManaCost*(0.03*float64(paladin.Talents.Benediction)),
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCDTime - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})
}

func (paladin *Paladin) CanJudgementOfBlood(sim *core.Simulation) bool {
	return paladin.canJudgement(sim) && paladin.CurrentSeal == paladin.SealOfBloodAura
}

func (paladin *Paladin) registerJudgementOfTheCrusaderSpell(cdTimer *core.Timer, sanctifiedJudgementMetrics *core.ResourceMetrics) {
	percentBonus := 1.0
	if ItemSetJusticarBattlegear.CharacterHasSetBonus(&paladin.Character, 2) {
		percentBonus = 1.15
	}
	flatBonus := 0.0
	if paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 23203 {
		flatBonus += 33.0
	} else if paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 27949 {
		flatBonus += 47.0
	}
	paladin.JudgementOfTheCrusaderAura = core.JudgementOfTheCrusaderAura(paladin.CurrentTarget, paladin.Talents.ImprovedSealOfTheCrusader, flatBonus, percentBonus)

	baseCost := core.TernaryFloat64(ItemSetCrystalforgeBattlegear.CharacterHasSetBonus(&paladin.Character, 2), JudgementManaCost-35, JudgementManaCost)
	paladin.JudgementOfTheCrusader = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27159},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagJudgement,

		ResourceType: stats.Mana,
		BaseCost:     JudgementManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost - JudgementManaCost*(0.03*float64(paladin.Talents.Benediction)),
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCDTime - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				paladin.sanctifiedJudgement(sim, sanctifiedJudgementMetrics, paladin.SealOfTheCrusader.DefaultCast.Cost)
				paladin.SealOfTheCrusaderAura.Deactivate(sim)
				paladin.CurrentSeal = nil
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskEmpty,
			OutcomeApplier: paladin.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				paladin.JudgementOfTheCrusaderAura.Activate(sim)
				paladin.CurrentJudgement = paladin.JudgementOfTheCrusaderAura
			},
		}),
	})
}

func (paladin *Paladin) CanJudgementOfTheCrusader(sim *core.Simulation) bool {
	return paladin.canJudgement(sim) && paladin.CurrentSeal == paladin.SealOfTheCrusaderAura
}

func (paladin *Paladin) registerJudgementOfWisdomSpell(cdTimer *core.Timer, sanctifiedJudgementMetrics *core.ResourceMetrics) {
	paladin.JudgementOfWisdomAura = core.JudgementOfWisdomAura(paladin.CurrentTarget)

	baseCost := core.TernaryFloat64(ItemSetCrystalforgeBattlegear.CharacterHasSetBonus(&paladin.Character, 2), JudgementManaCost-35, JudgementManaCost)
	paladin.JudgementOfWisdom = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27164},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagJudgement,

		ResourceType: stats.Mana,
		BaseCost:     JudgementManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost - JudgementManaCost*(0.03*float64(paladin.Talents.Benediction)),
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCDTime - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				paladin.sanctifiedJudgement(sim, sanctifiedJudgementMetrics, paladin.SealOfWisdom.DefaultCast.Cost)
				paladin.SealOfWisdomAura.Deactivate(sim)
				paladin.CurrentSeal = nil
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskEmpty,
			OutcomeApplier: paladin.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				paladin.JudgementOfWisdomAura.Activate(sim)
				paladin.CurrentJudgement = paladin.JudgementOfWisdomAura
			},
		}),
	})
}

func (paladin *Paladin) CanJudgementOfWisdom(sim *core.Simulation) bool {
	return paladin.canJudgement(sim) && paladin.CurrentSeal == paladin.SealOfWisdomAura
}

func (paladin *Paladin) registerJudgementOfLightSpell(cdTimer *core.Timer, sanctifiedJudgementMetrics *core.ResourceMetrics) {
	paladin.JudgementOfLightAura = core.JudgementOfLightAura(paladin.CurrentTarget)

	baseCost := core.TernaryFloat64(ItemSetCrystalforgeBattlegear.CharacterHasSetBonus(&paladin.Character, 2), JudgementManaCost-35, JudgementManaCost)
	paladin.JudgementOfLight = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27163},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagJudgement,

		ResourceType: stats.Mana,
		BaseCost:     JudgementManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost - JudgementManaCost*(0.03*float64(paladin.Talents.Benediction)),
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCDTime - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
			OnCastComplete: func(sim *core.Simulation, spell *core.Spell) {
				paladin.sanctifiedJudgement(sim, sanctifiedJudgementMetrics, paladin.SealOfLight.DefaultCast.Cost)
				paladin.SealOfLightAura.Deactivate(sim)
				paladin.CurrentSeal = nil
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:       core.ProcMaskEmpty,
			OutcomeApplier: paladin.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if !spellEffect.Landed() {
					return
				}
				paladin.JudgementOfLightAura.Activate(sim)
				paladin.CurrentJudgement = paladin.JudgementOfLightAura
			},
		}),
	})
}

func (paladin *Paladin) CanJudgementOfLight(sim *core.Simulation) bool {
	return paladin.canJudgement(sim) && paladin.CurrentSeal == paladin.SealOfLightAura
}

func (paladin *Paladin) registerJudgementOfRighteousnessSpell(cdTimer *core.Timer, sanctifiedJudgementMetrics *core.ResourceMetrics) {
	baseCost := core.TernaryFloat64(ItemSetCrystalforgeBattlegear.CharacterHasSetBonus(&paladin.Character, 2), JudgementManaCost-35, JudgementManaCost)
	paladin.JudgementOfRighteousness = paladin.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 27157},
		SpellSchool: core.SpellSchoolHoly,
		Flags:       SpellFlagJudgement | core.SpellFlagBinary,

		ResourceType: stats.Mana,
		BaseCost:     JudgementManaCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost - JudgementManaCost*(0.03*float64(paladin.Talents.Benediction)),
			},
			CD: core.Cooldown{
				Timer:    cdTimer,
				Duration: JudgementCDTime - (time.Second * time.Duration(paladin.Talents.ImprovedJudgement)),
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask: core.ProcMaskMeleeOrRangedSpecial,

			BonusSpellPower:  core.TernaryFloat64(paladin.Equip[proto.ItemSlot_ItemSlotRanged].ID == 33504, 94.0, 0),
			BonusCritRating:  3 * core.MeleeCritRatingPerCritChance * float64(paladin.Talents.Fanaticism),
			DamageMultiplier: 1 + 0.03*float64(paladin.Talents.ImprovedSealOfRighteousness),
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMagic(225, 246, 0.728),
			OutcomeApplier: paladin.OutcomeFuncMagicHitAndCritBinary(paladin.SpellCritMultiplier()),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				paladin.sanctifiedJudgement(sim, sanctifiedJudgementMetrics, paladin.SealOfRighteousness.DefaultCast.Cost)
				paladin.SealOfRighteousnessAura.Deactivate(sim)
				paladin.CurrentSeal = nil
			},
		}),
	})
}

// Defines judgement refresh behavior from attacks
// Returns extra mana if a different pally applied Judgement of Wisdom
func (paladin *Paladin) setupJudgementRefresh() {
	const mana = 74 / 2
	paladin.RegisterAura(core.Aura{
		Label:    "Refresh Judgement",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Landed() && spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				if paladin.CurrentJudgement != nil && paladin.CurrentJudgement.IsActive() {
					// Refresh the judgement
					paladin.CurrentJudgement.Refresh(sim)

					// Check if current judgement is not JoW and also that JoW is on the target
					if paladin.CurrentJudgement != paladin.JudgementOfWisdomAura && paladin.JudgementOfWisdomAura.IsActive() {
						// Just trigger a second JoW
						if paladin.JowManaMetrics == nil {
							paladin.JowManaMetrics = paladin.NewManaMetrics(core.ActionID{SpellID: 27164})
						}
						paladin.AddMana(sim, mana, paladin.JowManaMetrics, false)
					}
				}
			}
		},
	})
}

// Helper function to implement Sanctified Seals talent
func (paladin *Paladin) sanctifiedJudgement(sim *core.Simulation, manaMetrics *core.ResourceMetrics, mana float64) {
	if paladin.Talents.SanctifiedJudgement == 0 {
		return
	}

	var proc float64
	if paladin.Talents.SanctifiedJudgement == 3 {
		proc = 1
	} else {
		proc = 0.33 * float64(paladin.Talents.SanctifiedJudgement)
	}

	if sim.RandomFloat("Sanctified Judgement") < proc {
		paladin.AddMana(sim, 0.8*mana, manaMetrics, false)
	}
}

func (paladin *Paladin) registerJudgements() {
	cdTimer := paladin.NewTimer()
	sanctifiedJudgementMetrics := paladin.NewManaMetrics(core.ActionID{SpellID: 31930})
	paladin.registerJudgementOfBloodSpell(cdTimer, sanctifiedJudgementMetrics)
	paladin.registerJudgementOfTheCrusaderSpell(cdTimer, sanctifiedJudgementMetrics)
	paladin.registerJudgementOfWisdomSpell(cdTimer, sanctifiedJudgementMetrics)
	paladin.registerJudgementOfLightSpell(cdTimer, sanctifiedJudgementMetrics)
	paladin.registerJudgementOfRighteousnessSpell(cdTimer, sanctifiedJudgementMetrics)
}
