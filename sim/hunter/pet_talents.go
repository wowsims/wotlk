package hunter

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (hp *HunterPet) ApplyTalents() {
	talents := hp.Talents()
	// TODO:
	// Wolverine bite
	// Thunderstomp

	hp.AddStat(stats.MeleeCrit, 3*core.CritRatingPerCritChance*float64(talents.SpidersBite))
	hp.AddStat(stats.SpellCrit, 3*core.CritRatingPerCritChance*float64(talents.SpidersBite))
	hp.PseudoStats.MeleeSpeedMultiplier *= 1 + 0.15*float64(talents.CobraReflexes)
	hp.AutoAttacks.MHEffect.DamageMultiplier *= 1 - 0.075*float64(talents.CobraReflexes)
	hp.AutoAttacks.MHEffect.DamageMultiplier *= 1 + 0.03*float64(talents.SpikedCollar)
	hp.AutoAttacks.MHEffect.DamageMultiplier *= 1 + 0.03*float64(talents.SharkAttack)

	hp.PseudoStats.ArcaneDamageTakenMultiplier *= 1 + 0.05*float64(talents.GreatResistance)
	hp.PseudoStats.FireDamageTakenMultiplier *= 1 + 0.05*float64(talents.GreatResistance)
	hp.PseudoStats.FrostDamageTakenMultiplier *= 1 + 0.05*float64(talents.GreatResistance)
	hp.PseudoStats.NatureDamageTakenMultiplier *= 1 + 0.05*float64(talents.GreatResistance)
	hp.PseudoStats.ShadowDamageTakenMultiplier *= 1 + 0.05*float64(talents.GreatResistance)

	if talents.GreatStamina != 0 {
		bonus := 1 + 0.04*float64(talents.GreatStamina)
		hp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(stamina float64, _ float64) float64 {
				return stamina * bonus
			},
		})
	}

	if talents.NaturalArmor != 0 {
		bonus := 1 + 0.05*float64(talents.NaturalArmor)
		hp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Armor,
			ModifiedStat: stats.Armor,
			Modifier: func(armor float64, _ float64) float64 {
				return armor * bonus
			},
		})
	}

	if talents.BloodOfTheRhino != 0 {
		hp.PseudoStats.HealingTakenMultiplier *= 1 + 0.2*float64(talents.BloodOfTheRhino)

		bonus := 1 + 0.02*float64(talents.BloodOfTheRhino)
		hp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Stamina,
			ModifiedStat: stats.Stamina,
			Modifier: func(stamina float64, _ float64) float64 {
				return stamina * bonus
			},
		})
	}

	if talents.PetBarding != 0 {
		hp.AddStat(stats.Dodge, 1*core.DodgeRatingPerDodgeChance*float64(talents.PetBarding))

		bonus := 1 + 0.05*float64(talents.PetBarding)
		hp.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Armor,
			ModifiedStat: stats.Armor,
			Modifier: func(armor float64, _ float64) float64 {
				return armor * bonus
			},
		})
	}

	hp.applyOwlsFocus()
	hp.applyCullingTheHerd()
	hp.applyFeedingFrenzy()

	hp.registerRoarOfRecoveryCD()
	hp.registerRabidCD()
	hp.registerCallOfTheWildCD()
}

func (hp *HunterPet) applyOwlsFocus() {
	if hp.Talents().OwlsFocus == 0 {
		return
	}

	procChance := 0.15 * float64(hp.Talents().OwlsFocus)

	procAura := hp.RegisterAura(core.Aura{
		Label:    "Owl's Focus Proc",
		ActionID: core.ActionID{SpellID: 53515},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hp.PseudoStats.NoCost = true
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hp.PseudoStats.NoCost = false
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskSpecial) {
				aura.Deactivate(sim)
			}
		},
	})

	hp.RegisterAura(core.Aura{
		Label:    "Owls Focus",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.ProcMask.Matches(core.ProcMaskSpecial) && sim.RandomFloat("Owls Focus") < procChance {
				procAura.Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) applyCullingTheHerd() {
	if hp.Talents().CullingTheHerd == 0 {
		return
	}

	damageMult := 1 + 0.01*float64(hp.Talents().CullingTheHerd)

	makeProcAura := func(unit *core.Unit) *core.Aura {
		return unit.RegisterAura(core.Aura{
			Label:    "Culling the Herd Proc",
			ActionID: core.ActionID{SpellID: 52858},
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier *= damageMult
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.DamageDealtMultiplier /= damageMult
			},
		})
	}
	petAura := makeProcAura(&hp.Unit)
	ownerAura := makeProcAura(&hp.hunterOwner.Unit)

	hp.RegisterAura(core.Aura{
		Label:    "Culling the Herd",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			// TODO: Smack also
			if spellEffect.Outcome.Matches(core.OutcomeCrit) && (spell.IsSpellAction(BiteSpellID) || spell.IsSpellAction(ClawSpellID)) {
				petAura.Activate(sim)
				ownerAura.Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) applyFeedingFrenzy() {
	if hp.Talents().FeedingFrenzy == 0 {
		return
	}

	multiplier := 1.0 + 0.8*float64(hp.Talents().FeedingFrenzy)

	hp.RegisterResetEffect(func(sim *core.Simulation) {
		sim.RegisterExecutePhaseCallback(func(sim *core.Simulation, isExecute20 bool) {
			if !isExecute20 {
				hp.PseudoStats.DamageDealtMultiplier *= multiplier
			}
		})
	})
}

func (hp *HunterPet) registerRoarOfRecoveryCD() {
	if !hp.Talents().RoarOfRecovery {
		return
	}

	hunter := hp.hunterOwner
	actionID := core.ActionID{SpellID: 53517}
	manaMetrics := hunter.NewManaMetrics(actionID)

	rorSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 3,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period:   time.Second * 3,
				NumTicks: 3,
				OnAction: func(sim *core.Simulation) {
					hunter.AddMana(sim, hunter.MaxMana()*0.1, manaMetrics, false)
				},
			})
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: rorSpell,
		Type:  core.CooldownTypeMana,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled() && hunter.CurrentManaPercent() < 0.6
		},
	})
}

func (hp *HunterPet) registerRabidCD() {
	if !hp.Talents().Rabid {
		return
	}

	hunter := hp.hunterOwner
	actionID := core.ActionID{SpellID: 53401}
	ppmm := hp.AutoAttacks.NewPPMManager(10.0, core.ProcMaskMelee)

	var curBonusPerStack float64

	procAura := hp.RegisterAura(core.Aura{
		Label:     "Rabid Stacking",
		ActionID:  actionID.WithTag(1),
		Duration:  core.NeverExpires,
		MaxStacks: 5,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.AddStatDynamic(sim, stats.AttackPower, curBonusPerStack*float64(newStacks-oldStacks))
		},
	})

	buffAura := hp.RegisterAura(core.Aura{
		Label:    "Rabid",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			curBonusPerStack = aura.Unit.GetStat(stats.AttackPower) * 0.5
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			procAura.Deactivate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if !spellEffect.Landed() {
				return
			}
			if !spellEffect.ProcMask.Matches(core.ProcMaskMelee) {
				return
			}

			if !ppmm.Proc(sim, spellEffect.ProcMask, "Rabid") {
				return
			}

			procAura.Activate(sim)
			procAura.AddStack(sim)
		},
	})

	rabidSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 45,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			buffAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: rabidSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled()
		},
	})
}

func (hp *HunterPet) registerCallOfTheWildCD() {
	if !hp.Talents().CallOfTheWild {
		return
	}

	hunter := hp.hunterOwner
	actionID := core.ActionID{SpellID: 53434}

	makeProcAura := func(unit *core.Unit) *core.Aura {
		var curBonus stats.Stats
		return unit.RegisterAura(core.Aura{
			Label:    "Call of the Wild",
			ActionID: actionID,
			Duration: time.Second * 20,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				val := unit.GetStat(stats.AttackPower) * 0.1
				curBonus = stats.Stats{stats.AttackPower: val, stats.RangedAttackPower: val}

				unit.AddStatsDynamic(sim, curBonus)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				unit.AddStatsDynamic(sim, curBonus.Multiply(-1))
			},
		})
	}
	petAura := makeProcAura(&hp.Unit)
	ownerAura := makeProcAura(&hp.hunterOwner.Unit)

	cotwSpell := hunter.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute * 5,
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			petAura.Activate(sim)
			ownerAura.Activate(sim)
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: cotwSpell,
		Type:  core.CooldownTypeDPS,
		CanActivate: func(sim *core.Simulation, character *core.Character) bool {
			return hp.IsEnabled()
		},
	})
}
