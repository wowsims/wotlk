package ulduar

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addHodir25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        32845,
			Name:      "Hodir 25",
			Level:     83,
			MobType:   proto.MobType_MobTypeGiant,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      32_477_905,
				stats.Armor:       10643,
				stats.AttackPower: 574,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.4,
			MinBaseDamage:    50000,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
		},
		AI: NewHodir25AI(),
	})
	core.AddPresetEncounter("Hodir 25", []string{
		bossPrefix + "/Hodir 25",
	})
}

type Hodir25AI struct {
	Target *core.Target

	// Frozen Blows Mechanics
	FrozenBlows     *core.Spell
	FrozenBlowsAura *core.Aura
	FrozenBlowsAuto *core.Spell
	FrozenBlowsCast *core.Spell

	FlashFreeze *core.Spell

	// Magic Damage Debuff
	Singed         *core.Aura
	ToastyFires    []*core.Aura
	ToastyFireTime time.Duration
	HasCampfire    bool

	// Haste Buff
	Starlight []*core.Aura

	// Crit Buff
	StormCloud []*core.Aura
	NextStorms time.Duration
}

func NewHodir25AI() core.AIFactory {
	return func() core.TargetAI {
		return &Hodir25AI{}
	}
}

func (ai *Hodir25AI) Initialize(target *core.Target) {
	ai.Target = target

	ai.registerBuffsDebuffs(target)
	ai.registerFlashFreeze(target)
	ai.registerFrozenBlowSpell(target)
}

func (ai *Hodir25AI) Reset(*core.Simulation) {
	ai.HasCampfire = true
	// First campfire in 15-20 seconds
	ai.ToastyFireTime = time.Duration(15+rand.Intn(5)) * time.Second
	ai.NextStorms = time.Duration(2+rand.Intn(5)) * time.Second
}

func (ai *Hodir25AI) registerFlashFreeze(target *core.Target) {
	ai.FlashFreeze = target.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 61968},

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 45,
			},
			DefaultCast: core.Cast{
				CastTime: time.Second * 9,
			},
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			pa := &core.PendingAction{
				// Remove last fire in 0-5 seconds
				NextActionAt: sim.CurrentTime + time.Duration(rand.Intn(5))*time.Second,
				OnAction: func(s *core.Simulation) {
					ai.HasCampfire = false
					if sim.Raid.Size() >= 10 {
						for _, toastyFire := range ai.ToastyFires {
							toastyFire.Deactivate(sim)
						}
					}
				},
			}
			sim.AddPendingAction(pa)

			pa = &core.PendingAction{
				// Activate new fires in 15-20 seconds
				NextActionAt: sim.CurrentTime + time.Duration(15+rand.Intn(5))*time.Second,
				OnAction: func(s *core.Simulation) {
					ai.HasCampfire = true
					if sim.Raid.Size() >= 10 {
						for _, toastyFire := range ai.ToastyFires {
							toastyFire.Activate(sim)
						}
					}
				},
			}
			sim.AddPendingAction(pa)
		},
	})
}

func (ai *Hodir25AI) registerBuffsDebuffs(target *core.Target) {
	// Create aura for stacking singed in raid sim
	if ai.Target.Env.Raid.Size() >= 10 {
		ai.ToastyFires = make([]*core.Aura, 0)
		for _, party := range ai.Target.Env.Raid.Parties {
			for _, player := range party.Players {
				character := player.GetCharacter()
				aura := character.RegisterAura(core.Aura{
					Label:    "Toasty Fire" + strconv.Itoa(int(character.Index)),
					ActionID: core.ActionID{SpellID: 62821},
					Duration: core.NeverExpires,
					OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if !spell.ProcMask.Matches(core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) {
							return
						}

						if sim.Proc(0.3, "Singed") {
							ai.Singed.Activate(sim)
							ai.Singed.AddStack(sim)
						}
					},
				})
				ai.ToastyFires = append(ai.ToastyFires, aura)
			}
		}
	}

	ai.Singed = target.RegisterAura(core.Aura{
		Label:     "Singed",
		ActionID:  core.ActionID{SpellID: 65280},
		MaxStacks: 25,
		Duration:  time.Second * 25,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks, newStacks int32) {
			oldValue := 1.0 + float64(oldStacks)*0.02
			newValue := 1.0 + float64(newStacks)*0.02

			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= oldValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= oldValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= oldValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= oldValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= oldValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= oldValue

			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= newValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= newValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= newValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= newValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= newValue
			aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= newValue
		},
	})

	ai.Starlight = make([]*core.Aura, 0)
	for _, party := range ai.Target.Env.Raid.Parties {
		for _, player := range party.PlayersAndPets {
			character := player.GetCharacter()
			aura := character.RegisterAura(core.Aura{
				Label:    "Starlight" + strconv.Itoa(int(character.Index)),
				ActionID: core.ActionID{SpellID: 62807},
				Duration: time.Second * 30,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					character.MultiplyAttackSpeed(sim, 1.5)
					character.MultiplyCastSpeed(1.5)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					character.MultiplyAttackSpeed(sim, 1/1.5)
					character.MultiplyCastSpeed(1 / 1.5)
				},
			})

			core.ApplyFixedUptimeAura(aura, 0.8, time.Second*5)
			ai.Starlight = append(ai.Starlight, aura)
		}
	}

	ai.StormCloud = make([]*core.Aura, 0)
	for _, party := range ai.Target.Env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			aura := character.RegisterAura(core.Aura{
				Label:    "Stormcloud" + strconv.Itoa(int(character.Index)),
				ActionID: core.ActionID{SpellID: 63711},
				Duration: 30 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexPhysical] += 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexArcane] += 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexFire] += 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexFrost] += 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexHoly] += 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexNature] += 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexShadow] += 1.35
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexPhysical] -= 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexArcane] -= 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexFire] -= 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexFrost] -= 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexHoly] -= 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexNature] -= 1.35
					aura.Unit.PseudoStats.SchoolCritMultiplier[stats.SchoolIndexShadow] -= 1.35
				},
			})

			ai.StormCloud = append(ai.StormCloud, aura)

			if ai.Target.Env.Raid.Size() < 10 {
				core.ApplyFixedUptimeAura(aura, 0.6, time.Second*30)
			}
		}
	}
}

func (ai *Hodir25AI) registerFrozenBlowSpell(target *core.Target) {
	ai.FrozenBlowsAura = target.RegisterAura(core.Aura{
		Label:    "Hodir Frozen Blows",
		ActionID: core.ActionID{SpellID: 63512},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] *= 0.3
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] /= 0.3
		},
	})

	ai.FrozenBlows = target.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63512},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHAuto,
		Flags:       core.SpellFlagMeleeMetrics,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Minute,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			ai.FrozenBlowsAura.Activate(sim)
		},
	})

	// Replace MH Hit when under Frozen Blows buff
	ai.Target.Unit.AutoAttacks.ReplaceMHSwing = func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
		if ai.FrozenBlowsAura.IsActive() {
			return ai.FrozenBlowsAuto
		} else {
			return nil
		}
	}

	ai.FrozenBlowsAuto = target.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63511},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHAuto,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1,
		CritMultiplier:   ai.Target.AutoAttacks.MH.CritMultiplier,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeEnemyMeleeWhite)

			// Deal extra frost damage if hit landed
			if result.Landed() {
				ai.FrozenBlowsCast.Cast(sim, target)
			}
		},
	})

	ai.FrozenBlowsCast = target.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 63511},
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, 40000, spell.OutcomeAlwaysHit)
		},
	})
}

func (ai *Hodir25AI) DoAction(sim *core.Simulation) {

	singedStacks := ai.Singed.GetStacks()

	if sim.CurrentTime >= ai.ToastyFireTime && ai.HasCampfire {
		// Refresh Singed approximately in individual sims
		if ai.Target.Env.Raid.Size() < 10 {
			if singedStacks < 25 {
				ai.Singed.Activate(sim)
				ai.Singed.AddStack(sim)
			} else {
				ai.Singed.Refresh(sim)
			}
		} else {
			// Activate Toasty Fire slowly for whole raid
			// TODO: Improve this with actual campfires
			firesActivated := 0
			for _, toastyFire := range ai.ToastyFires {
				if !toastyFire.IsActive() {
					toastyFire.Activate(sim)

					firesActivated = firesActivated + 1
					if firesActivated >= 4 {
						break
					}
				}
			}
		}
	}

	if ai.Target.Env.Raid.Size() >= 10 {
		if sim.CurrentTime >= ai.NextStorms {
			ai.NextStorms = sim.CurrentTime + 30*time.Second
			storm1 := rand.Intn(ai.Target.Env.Raid.Size())
			storm2 := storm1
			maxBuffs := core.MinInt(6, sim.Raid.Size()-1)

			// 2 storms on 25m
			if sim.Raid.Size() > 10 {
				for storm1 == storm2 {
					storm2 = rand.Intn(sim.Raid.Size())
				}
				maxBuffs = core.MinInt(12, sim.Raid.Size()-2)
			}

			for maxBuffs > 0 {
				target := -1
				for target == -1 || target == storm1 || target == storm2 || ai.StormCloud[target].IsActive() {
					target = rand.Intn(sim.Raid.Size())
				}
				ai.StormCloud[target].Activate(sim)
				maxBuffs = maxBuffs - 1
			}
		}
	}

	if ai.FrozenBlows.IsReady(sim) && sim.CurrentTime >= ai.FrozenBlows.CD.Duration {
		ai.FrozenBlows.Cast(sim, nil)
	}

	if ai.FlashFreeze.IsReady(sim) && sim.CurrentTime >= ai.FlashFreeze.CD.Duration {
		ai.FlashFreeze.Cast(sim, nil)
	}

	// Stormcloud CD - 30, starts casting 1-5 seconds in
	// Affects 2 people - each spread storm power to 6 others

	if ai.Target.GCD.IsReady(sim) {
		nextEventAt := sim.CurrentTime + time.Minute

		// All possible next events
		events := []time.Duration{
			core.MaxDuration(ai.FrozenBlows.ReadyAt(), ai.FrozenBlows.CD.Duration),
			core.MaxDuration(ai.FlashFreeze.ReadyAt(), ai.FlashFreeze.CD.Duration),
			ai.NextStorms,
		}

		if ai.Target.Env.Raid.Size() < 10 {
			// Individual Sim approximation - taken from some random logs
			timeBetweenStacks := 400 * time.Millisecond // TODO: Expose this
			events = append(events, core.MaxDuration(ai.ToastyFireTime, sim.CurrentTime+timeBetweenStacks))
		} else {
			timeBetweenFireStacks := 3 * time.Second // TODO: Improve on Fires Approximation by actually simulating active campfires
			events = append(events, core.MaxDuration(ai.ToastyFireTime, sim.CurrentTime+timeBetweenFireStacks))
		}

		// if ai.Target.CurrentTarget != nil {
		// 	events = append(events, core.MaxDuration(ai.PhasePunch.ReadyAt(), ai.PhasePunch.CD.Duration))
		// 	events = append(events, core.MaxDuration(ai.QuantumStrike.ReadyAt(), ai.QuantumStrike.CD.Duration))
		// }

		for _, elem := range events {
			if elem > sim.CurrentTime && elem < nextEventAt {
				nextEventAt = elem
			}
		}

		if nextEventAt == 0 {
			nextEventAt = time.Millisecond * 100
		}

		ai.Target.WaitUntil(sim, nextEventAt)
	}
}
