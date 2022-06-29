package rogue

import (
	"strconv"
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Returns whether any Deadly Poisons are being used.
func (rogue *Rogue) applyPoisons() {
	hasWFTotem := rogue.HasAura(core.WindfuryTotemAuraLabel)
	rogue.applyDeadlyPoison(hasWFTotem)
	rogue.applyInstantPoison(hasWFTotem)
}

func (rogue *Rogue) registerDeadlyPoisonSpell() {
	actionID := core.ActionID{SpellID: 27186}

	rogue.DeadlyPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolNature,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			BonusSpellHitRating: 5 * core.SpellHitRatingPerHitChance * float64(rogue.Talents.MasterPoisoner),
			ThreatMultiplier:    1,
			OutcomeApplier:      rogue.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					if rogue.DeadlyPoisonDot.IsActive() {
						rogue.DeadlyPoisonDot.Refresh(sim)
						rogue.DeadlyPoisonDot.AddStack(sim)
					} else {
						rogue.DeadlyPoisonDot.Apply(sim)
						rogue.DeadlyPoisonDot.SetStacks(sim, 1)
					}
				}
			},
		}),
	})

	target := rogue.CurrentTarget
	dotAura := target.RegisterAura(core.Aura{
		Label:     "DeadlyPoison-" + strconv.Itoa(int(rogue.Index)),
		ActionID:  actionID,
		MaxStacks: 5,
		Duration:  time.Second * 12,
	})
	rogue.DeadlyPoisonDot = core.NewDot(core.Dot{
		Spell:         rogue.DeadlyPoison,
		Aura:          dotAura,
		NumberOfTicks: 4,
		TickLength:    time.Second * 3,
		TickEffects: core.TickFuncApplyEffects(core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1 + 0.04*float64(rogue.Talents.VilePoisons),
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage:       core.MultiplyByStacks(core.BaseDamageConfigFlat(180/4), dotAura),
			OutcomeApplier:   rogue.OutcomeFuncTick(),
		})),
	})
}

func (rogue *Rogue) applyDeadlyPoison(hasWFTotem bool) {
	procMask := core.GetMeleeProcMaskForHands(
		!hasWFTotem && rogue.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueRogueDeadlyPoison,
		rogue.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueRogueDeadlyPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}

	procChance := 0.3 + 0.02*float64(rogue.Talents.ImprovedPoisons)

	rogue.RegisterAura(core.Aura{
		Label:    "Deadly Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Deadly Poison") > procChance {
				return
			}

			rogue.DeadlyPoison.Cast(sim, spellEffect.Target)
		},
	})
}

func (rogue *Rogue) registerInstantPoisonSpell() {
	rogue.InstantPoison = rogue.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 26891},
		SpellSchool: core.SpellSchoolNature,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskEmpty,
			DamageMultiplier:    1 + 0.04*float64(rogue.Talents.VilePoisons),
			ThreatMultiplier:    1,
			BonusSpellHitRating: 5 * core.SpellHitRatingPerHitChance * float64(rogue.Talents.MasterPoisoner),
			BaseDamage:          core.BaseDamageConfigRoll(146, 194),
			OutcomeApplier:      rogue.OutcomeFuncMagicHitAndCrit(rogue.SpellCritMultiplier()),
		}),
	})
}

func (rogue *Rogue) applyInstantPoison(hasWFTotem bool) {
	procMask := core.GetMeleeProcMaskForHands(
		!hasWFTotem && rogue.Consumes.MainHandImbue == proto.WeaponImbue_WeaponImbueRogueInstantPoison,
		rogue.Consumes.OffHandImbue == proto.WeaponImbue_WeaponImbueRogueInstantPoison)

	if procMask == core.ProcMaskUnknown {
		return
	}

	procChance := 0.2 + 0.02*float64(rogue.Talents.ImprovedPoisons)

	rogue.RegisterAura(core.Aura{
		Label:    "Instant Poison",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() || !spellEffect.ProcMask.Matches(procMask) {
				return
			}
			if sim.RandomFloat("Instant Poison") > procChance {
				return
			}

			rogue.procInstantPoison(sim, spellEffect)
		},
	})
}

func (rogue *Rogue) procInstantPoison(sim *core.Simulation, spellEffect *core.SpellEffect) {
	rogue.InstantPoison.Cast(sim, spellEffect.Target)
}
