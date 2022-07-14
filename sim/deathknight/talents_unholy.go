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
	// TODO:

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
	// TODO:

	// Blood-Caked Blade
	// TODO:

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
	// TODO:

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
