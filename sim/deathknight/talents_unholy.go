package deathknight

import (
	//"github.com/wowsims/wotlk/sim/core/proto"

	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	//"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (deathKnight *DeathKnight) ApplyUnholyTalents() {
	// Vicious Strikes
	// Implemented outside

	// Virulence
	deathKnight.AddStat(stats.SpellHit, core.SpellHitRatingPerHitChance*float64(deathKnight.Talents.Virulence))

	// Epidemic
	// Implemented outside

	// Morbidity
	// Implemented outside

	// Ravenous Dead
	// TODO: Ghoul part
	if deathKnight.Talents.RavenousDead > 0 {
		strengthCoeff := 0.01 * float64(deathKnight.Talents.RavenousDead)
		deathKnight.AddStatDependency(stats.StatDependency{
			SourceStat:   stats.Strength,
			ModifiedStat: stats.Strength,
			Modifier: func(strength float64, _ float64) float64 {
				return strength * (1.0 + strengthCoeff)
			},
		})
	}

	// Outbreak
	// Implemented outside

	// Necrosis
	deathKnight.applyNecrosis()

	// Blood-Caked Blade
	deathKnight.applyBloodCakedBlade()

	// Night of the Dead
	// TODO:

	// Unholy Blight
	deathKnight.applyUnholyBlight()

	// Impurity
	// TODO:

	// Dirge
	// Implemented outside

	// Reaping
	// TODO:

	// Master of Ghouls
	// TODO:

	// Desolation
	deathKnight.applyDesolation()

	// Ghoul Frenzy
	// TODO:

	// Bone Shield
	// TODO:

	// Wandering Plague
	deathKnight.applyWanderingPlague()

	// Crypt Fever
	// Ebon Plaguebringer
	// TODO: Diseases damage increase still missing
	deathKnight.applyEbonPlaguebringer()

	// Scourge Strike
	// Implemented outside

	// Rage of Rivendare
	// TODO: % bonus damage to spells/abilities (not white hits)
	deathKnight.AddStat(stats.Expertise, float64(deathKnight.Talents.RageOfRivendare)*core.ExpertisePerQuarterPercentReduction)

	// Summon Gargoyle
	// TODO:
}

func (deathKnight *DeathKnight) applyWanderingPlague() {
	if deathKnight.Talents.WanderingPlague == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 49655}

	deathKnight.WanderingPlague = deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNone,

		ApplyEffects: core.ApplyEffectFuncAOEDamage(deathKnight.Env, core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return deathKnight.LastDiseaseDamage
				},
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})
}

func (deathKnight *DeathKnight) applyNecrosis() {
	if deathKnight.Talents.Necrosis == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 51465}
	target := deathKnight.CurrentTarget

	var curDmg float64
	necrosisHit := deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		Flags:       core.SpellFlagNone,

		ApplyEffects: core.ApplyEffectFuncDirectDamageTargetModifiersOnly(core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(_ *core.Simulation, _ *core.SpellEffect, _ *core.Spell) float64 {
					return curDmg * 0.04 * float64(deathKnight.Talents.Necrosis)
				},
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})

	deathKnight.NecrosisAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Necrosis",
		ActionID: actionID,
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.NecrosisAura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			curDmg = spellEffect.Damage
			necrosisHit.Cast(sim, target)
		},
	})
}

func (deathKnight *DeathKnight) applyBloodCakedBlade() {
	if deathKnight.Talents.BloodCakedBlade == 0 {
		return
	}

	target := deathKnight.CurrentTarget

	mhBaseDamage := core.BaseDamageFuncMeleeWeapon(core.MainHand, false, 0, 1.0, true)
	ohBaseDamage := core.BaseDamageFuncMeleeWeapon(core.OffHand, false, 0, 1.0, true)

	var isMH = false
	bloodCakedBladeHit := deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 50463},
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagMeleeMetrics,

		ApplyEffects: core.ApplyEffectFuncDirectDamageTargetModifiersOnly(core.SpellEffect{
			// No proc mask, so it won't proc itself.
			ProcMask: core.ProcMaskEmpty,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, spellEffect *core.SpellEffect, spell *core.Spell) float64 {
					if isMH {
						return mhBaseDamage(sim, spellEffect, spell) * (0.25 + float64(deathKnight.countActiveDiseases())*0.125)
					} else {
						return ohBaseDamage(sim, spellEffect, spell) * (0.25 + float64(deathKnight.countActiveDiseases())*0.125)
					}
				},
			},
			OutcomeApplier: deathKnight.OutcomeFuncMeleeWeaponSpecialNoHitNoCrit(),
		}),
	})

	deathKnight.BloodCakedBladeAura = deathKnight.RegisterAura(core.Aura{
		Label:    "Blood-Caked Blade",
		ActionID: core.ActionID{SpellID: 49628},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			deathKnight.BloodCakedBladeAura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spellEffect.Damage == 0 || !spellEffect.ProcMask.Matches(core.ProcMaskMeleeWhiteHit) {
				return
			}

			if sim.RandomFloat("Blood-Caked Blade Roll") < 0.30 {
				isMH = spellEffect.ProcMask.Matches(core.ProcMaskMeleeMHAuto)
				bloodCakedBladeHit.Cast(sim, target)
			}
		},
	})
}

func (deathKnight *DeathKnight) applyEbonPlaguebringer() {
	if deathKnight.Talents.EbonPlaguebringer == 0 {
		return
	}

	ebonPlaguebringerBonusCrit := core.CritRatingPerCritChance * float64(deathKnight.Talents.EbonPlaguebringer)
	deathKnight.AddStat(stats.MeleeCrit, ebonPlaguebringerBonusCrit)
	deathKnight.AddStat(stats.SpellCrit, ebonPlaguebringerBonusCrit)

	target := deathKnight.CurrentTarget

	epAura := core.EbonPlaguebringerAura(target)
	epAura.Duration = time.Second * (15 + 3*time.Duration(deathKnight.Talents.Epidemic))

	deathKnight.EbonPlagueAura = epAura
}

func (deathKnight *DeathKnight) applyDesolation() {
	if deathKnight.Talents.Desolation == 0 {
		return
	}

	actionID := core.ActionID{SpellID: 66803}

	deathKnight.DesolationAura = deathKnight.RegisterAura(core.Aura{
		ActionID: actionID,
		Label:    "Desolation",
		Duration: time.Second * 20.0,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.0 + 0.01*float64(deathKnight.Talents.Desolation)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.0 + 0.01*float64(deathKnight.Talents.Desolation)
		},
	})
}

func (deathKnight *DeathKnight) applyUnholyBlight() {
	actionID := core.ActionID{SpellID: 50536}
	target := deathKnight.CurrentTarget

	unholyBlightSpell := deathKnight.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
	})

	deathKnight.LastDeathCoilDamage = 1500
	deathKnight.UnholyBlight = core.NewDot(core.Dot{
		Spell: unholyBlightSpell,
		Aura: target.RegisterAura(core.Aura{
			Label:    "UnholyBlight-" + strconv.Itoa(int(deathKnight.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 10,
		TickLength:    time.Second * 1,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			ProcMask:         core.ProcMaskPeriodicDamage,
			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			IsPeriodic:       true,
			BaseDamage: core.BaseDamageConfig{
				Calculator: func(sim *core.Simulation, hitEffect *core.SpellEffect, spell *core.Spell) float64 {
					return (0.10 * deathKnight.LastDeathCoilDamage) / 10
				},
				TargetSpellCoefficient: 1,
			},
			OutcomeApplier: deathKnight.OutcomeFuncAlwaysHit(),
		}),
	})
}
