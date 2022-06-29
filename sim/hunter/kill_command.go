package hunter

import (
	"time"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/stats"
)

func (hunter *Hunter) applyKillCommand() {
	if hunter.pet == nil {
		return
	}

	hunter.RegisterAura(core.Aura{
		Label:    "Kill Command Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Outcome.Matches(core.OutcomeCrit) {
				hunter.killCommandEnabledUntil = sim.CurrentTime + time.Second*5
				hunter.TryKillCommand(sim, hunter.CurrentTarget)
			}
		},
	})
}

func (hunter *Hunter) registerKillCommandSpell() {
	baseCost := 75.0

	hunter.KillCommand = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34026},
		SpellSchool: core.SpellSchoolPhysical,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 5,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskEmpty,
			ThreatMultiplier: 1,
			OutcomeApplier:   hunter.OutcomeFuncAlwaysHit(),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				hunter.killCommandEnabledUntil = 0
				hunter.pet.KillCommand.Cast(sim, hunter.CurrentTarget)
			},
		}),
	})
}

func (hp *HunterPet) registerKillCommandSpell() {
	var beastLordProcAura *core.Aura
	if ItemSetBeastLord.CharacterHasSetBonus(&hp.hunterOwner.Character, 4) {
		beastLordProcAura = hp.hunterOwner.NewTemporaryStatsAura("Beast Lord Proc", core.ActionID{SpellID: 37483}, stats.Stats{stats.ArmorPenetration: 600}, time.Second*15)
	}

	hp.KillCommand = hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 34027},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:         core.ProcMaskMeleeMHSpecial,
			BonusCritRating:  float64(hp.hunterOwner.Talents.FocusedFire) * 10 * core.MeleeCritRatingPerCritChance,
			DamageMultiplier: hp.config.DamageMultiplier,
			ThreatMultiplier: 1,

			BaseDamage:     core.BaseDamageConfigMeleeWeapon(core.MainHand, false, 127, 1, true),
			OutcomeApplier: hp.OutcomeFuncMeleeSpecialHitAndCrit(2),

			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if beastLordProcAura != nil {
					beastLordProcAura.Activate(sim)
				}
			},
		}),
	})
}

func (hunter *Hunter) TryKillCommand(sim *core.Simulation, target *core.Unit) {
	if hunter.pet == nil || !hunter.pet.IsEnabled() {
		return
	}

	if hunter.killCommandEnabledUntil < sim.CurrentTime || hunter.killCommandBlocked {
		return
	}

	if hunter.CurrentMana() < 75 {
		return
	}

	if hunter.KillCommand.IsReady(sim) {
		hunter.KillCommand.Cast(sim, target)
	}
}
