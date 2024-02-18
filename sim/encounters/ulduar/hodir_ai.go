package ulduar

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func addHodir10(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        32845,
			Name:      "Hodir",
			Level:     83,
			MobType:   proto.MobType_MobTypeGiant,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      8_115_990,
				stats.Armor:       10643,
				stats.AttackPower: 574,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.4,
			MinBaseDamage:    25000, // TODO: Find real value
			DamageSpread:     0.3333,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     HodirTargetInputs(),
		},
		AI: NewHodir10AI(),
	})
	core.AddPresetEncounter("Hodir", []string{
		bossPrefix + "/Hodir",
	})
}

func addHodir25(bossPrefix string) {
	core.AddPresetTarget(&core.PresetTarget{
		PathPrefix: bossPrefix,
		Config: &proto.Target{
			Id:        32845,
			Name:      "Hodir",
			Level:     83,
			MobType:   proto.MobType_MobTypeGiant,
			TankIndex: 0,

			Stats: stats.Stats{
				stats.Health:      32_477_905,
				stats.Armor:       10643,
				stats.AttackPower: 805,
				stats.BlockValue:  76,
			}.ToFloatArray(),

			SpellSchool:      proto.SpellSchool_SpellSchoolPhysical,
			SwingSpeed:       2.4,
			MinBaseDamage:    46300, // TODO: Find real value
			DamageSpread:     0.3333,
			SuppressDodge:    false,
			ParryHaste:       false,
			DualWield:        false,
			DualWieldPenalty: false,
			TargetInputs:     HodirTargetInputs(),
		},
		AI: NewHodir25AI(),
	})
	core.AddPresetEncounter("Hodir", []string{
		bossPrefix + "/Hodir",
	})
}

type HodirAI struct {
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

	raidSize int

	StormPowerPrio  bool
	StarlightUptime float64
}

func HodirTargetInputs() []*proto.TargetInput {
	return []*proto.TargetInput{
		{
			Label:     "Stormpower Prio",
			Tooltip:   "Should stormpower buff be applied when available",
			InputType: proto.InputType_Bool,
			BoolValue: true,
		},
		{
			Label:       "Starlight Uptime %",
			Tooltip:     "Uptime on Starlight haste buff (Range 0-100%)",
			InputType:   proto.InputType_Number,
			NumberValue: 80.0,
		},
	}
}

func NewHodir10AI() core.AIFactory {
	return func() core.TargetAI {
		return &HodirAI{
			raidSize: 10,
		}
	}
}

func NewHodir25AI() core.AIFactory {
	return func() core.TargetAI {
		return &HodirAI{
			raidSize: 25,
		}
	}
}

func (ai *HodirAI) Initialize(target *core.Target, config *proto.Target) {
	ai.Target = target

	ai.StormPowerPrio = config.TargetInputs[0].BoolValue
	ai.StarlightUptime = config.TargetInputs[1].NumberValue

	ai.registerBuffsDebuffs(target)
	ai.registerFlashFreeze(target)
	ai.registerFrozenBlowSpell(target)
}

func (ai *HodirAI) Reset(sim *core.Simulation) {
	ai.HasCampfire = true
	// First campfire in 15-20 seconds
	ai.ToastyFireTime = time.Duration(15+5.0*sim.RandomFloat("HodirAI Toasty Fire")) * time.Second
	// First storms in 33-38 seconds
	ai.NextStorms = time.Duration(33+5.0*sim.RandomFloat("HodirAI Next Storm")) * time.Second
}

func (ai *HodirAI) registerFlashFreeze(target *core.Target) {
	ai.FlashFreeze = target.GetOrRegisterSpell(core.SpellConfig{
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
			// Remove last fire in 0-5 seconds
			pa := &core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Duration(5.0*sim.RandomFloat("HodirAI Remove Last Fire"))*time.Second,
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

			// Activate new fires in 15-20 seconds
			pa = &core.PendingAction{
				NextActionAt: sim.CurrentTime + time.Duration(15+5.0*sim.RandomFloat("HodirAI Activate New Fires"))*time.Second,
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

			ai.NextStorms = max(ai.NextStorms, sim.CurrentTime+time.Duration(3+5.0*sim.RandomFloat("HodirAI Next Storm"))*time.Second)
		},
	})
}

func (ai *HodirAI) registerBuffsDebuffs(target *core.Target) {
	// Create aura for stacking singed in raid sim
	if ai.Target.Env.Raid.Size() > 1 {
		ai.ToastyFires = make([]*core.Aura, 0)
		for _, party := range ai.Target.Env.Raid.Parties {
			for _, player := range party.Players {
				character := player.GetCharacter()
				aura := character.GetOrRegisterAura(core.Aura{
					Label:    "Toasty Fire" + strconv.Itoa(int(character.Index)),
					ActionID: core.ActionID{SpellID: 62821},
					Duration: core.NeverExpires,
					OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
						if !spell.ProcMask.Matches(core.ProcMaskRangedSpecial | core.ProcMaskSpellDamage) {
							return
						}

						if sim.Proc(0.33, "Singed") {
							ai.Singed.Activate(sim)
							ai.Singed.AddStack(sim)
						}
					},
				})
				ai.ToastyFires = append(ai.ToastyFires, aura)
			}
		}
	}

	ai.Singed = target.GetOrRegisterAura(core.Aura{
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
			aura := character.GetOrRegisterAura(core.Aura{
				Label:    "Starlight" + strconv.Itoa(int(character.UnitIndex)),
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

			core.ApplyFixedUptimeAura(aura, min(max(ai.StarlightUptime, 0.0), 100.0)/100.0, time.Second*15, time.Second*10)
			ai.Starlight = append(ai.Starlight, aura)
		}
	}

	ai.StormCloud = make([]*core.Aura, 0)
	for _, party := range ai.Target.Env.Raid.Parties {
		for _, player := range party.Players {
			character := player.GetCharacter()
			aura := character.GetOrRegisterAura(core.Aura{
				Label:    "Stormcloud" + strconv.Itoa(int(character.Index)),
				ActionID: core.ActionID{SpellID: 63711},
				Duration: 30 * time.Second,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range aura.Unit.Spellbook {
						spell.CritMultiplier *= 2.35
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					for _, spell := range aura.Unit.Spellbook {
						spell.CritMultiplier /= 2.35
					}
				},
			})

			ai.StormCloud = append(ai.StormCloud, aura)
		}
	}
}

func (ai *HodirAI) registerFrozenBlowSpell(target *core.Target) {
	ai.FrozenBlowsAura = target.GetOrRegisterAura(core.Aura{
		Label:    "Hodir Frozen Blows",
		ActionID: core.ActionID{SpellID: core.TernaryInt32(ai.raidSize == 25, 63512, 62478)},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] *= 0.3
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolDamageDealtMultiplier[core.SpellSchoolPhysical] /= 0.3
		},
	})

	ai.FrozenBlows = target.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: core.TernaryInt32(ai.raidSize == 25, 63512, 62478)},
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
	ai.Target.Unit.AutoAttacks.SetReplaceMHSwing(func(sim *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
		if ai.FrozenBlowsAura.IsActive() {
			return ai.FrozenBlowsAuto
		}
		return mhSwingSpell
	})

	ai.FrozenBlowsAuto = target.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: core.TernaryInt32(ai.raidSize == 25, 63511, 62867)}.WithTag(1),
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHAuto,
		Flags:       core.SpellFlagMeleeMetrics,

		DamageMultiplier: 1,
		CritMultiplier:   ai.Target.AutoAttacks.MH().CritMultiplier,

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

	ai.FrozenBlowsCast = target.GetOrRegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: core.TernaryInt32(ai.raidSize == 25, 63511, 62867)}.WithTag(2),
		SpellSchool: core.SpellSchoolFrost,
		ProcMask:    core.ProcMaskSpellDamage,

		DamageMultiplier: 1,
		CritMultiplier:   1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.CalcAndDealDamage(sim, target, core.TernaryFloat64(ai.raidSize == 25, 40000, 31062), spell.OutcomeAlwaysHit)
		},
	})
}

func (ai *HodirAI) ExecuteCustomRotation(sim *core.Simulation) {
	singedStacks := ai.Singed.GetStacks()

	if sim.CurrentTime >= ai.ToastyFireTime && ai.HasCampfire {
		// Refresh Singed approximately in individual sims
		if ai.Target.Env.Raid.Size() == 1 {
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

	// Stormclouds are cast every 30-35 seconds
	// Affects 2 people - each spread storm power to 6 others
	if sim.CurrentTime >= ai.NextStorms {
		ai.NextStorms = sim.CurrentTime + 30*time.Second + time.Duration(5.0*sim.RandomFloat("HodirAI Cast Storm Cloud"))*time.Second

		if ai.Target.Env.Raid.Size() > 1 {
			// Raid sim we simulate storm clouds and storm power spreading
			// Assign random storm spreader
			storm1 := int(float64(ai.Target.Env.Raid.Size()) * sim.RandomFloat("HodirAI Random Storm Spreader"))
			storm2 := storm1

			// Set max possible spreads
			maxBuffs := min(6, sim.Raid.Size()-1)

			// 2 storms on 25m
			if ai.raidSize == 25 {
				for storm1 == storm2 {
					// Assign 2nd random spreader
					storm2 = int(float64(ai.Target.Env.Raid.Size()) * sim.RandomFloat("HodirAI Random Storm Spreader"))
				}

				// Set max possible spreads
				maxBuffs = min(12, sim.Raid.Size()-2)
			}

			// Prio order for storm power
			mages := make([]int, 0)
			boomies := make([]int, 0)
			warlocks := make([]int, 0)
			shamans := make([]int, 0)
			spriests := make([]int, 0)
			dks := make([]int, 0)

			for _, party := range sim.Raid.Parties {
				for _, player := range party.Players {
					character := player.GetCharacter()
					raidIndex := int(character.Index)

					// Can't prio if its the storm spreader
					if raidIndex == storm1 || raidIndex == storm2 {
						continue
					}

					switch character.Class {
					case proto.Class_ClassMage:
						mages = append(mages, raidIndex)
					case proto.Class_ClassDruid:
						if character.PrimaryTalentTree == 0 {
							boomies = append(boomies, raidIndex)
						}
					case proto.Class_ClassWarlock:
						warlocks = append(warlocks, raidIndex)
					case proto.Class_ClassShaman:
						if character.PrimaryTalentTree != 2 {
							shamans = append(shamans, raidIndex)
						}
					case proto.Class_ClassPriest:
						if character.PrimaryTalentTree == 2 {
							spriests = append(spriests, raidIndex)
						}
					case proto.Class_ClassDeathknight:
						dks = append(dks, raidIndex)
					}
				}
			}

			maxBuffs = ai.stormCloudPrioApply(sim, maxBuffs, mages)
			maxBuffs = ai.stormCloudPrioApply(sim, maxBuffs, boomies)
			maxBuffs = ai.stormCloudPrioApply(sim, maxBuffs, warlocks)
			maxBuffs = ai.stormCloudPrioApply(sim, maxBuffs, shamans)
			maxBuffs = ai.stormCloudPrioApply(sim, maxBuffs, spriests)
			maxBuffs = ai.stormCloudPrioApply(sim, maxBuffs, dks)

			// Spread randomly whats left
			for maxBuffs > 0 {
				target := -1
				for target == -1 || target == storm1 || target == storm2 || ai.StormCloud[target].IsActive() {
					target = int(float64(ai.Target.Env.Raid.Size()) * sim.RandomFloat("HodirAI Random Storm Receiver"))
				}
				ai.StormCloud[target].Activate(sim)
				maxBuffs = maxBuffs - 1
			}
		} else {
			// Individual sim we assume actor is prioritized for every storm power
			// so just activate them
			if ai.StormPowerPrio {
				for _, stormCloud := range ai.StormCloud {
					stormCloud.Activate(sim)
				}
			}
		}
	}

	if ai.Target.CurrentTarget != nil {
		if ai.FrozenBlows.IsReady(sim) && sim.CurrentTime >= ai.FrozenBlows.CD.Duration {
			ai.FrozenBlows.Cast(sim, nil)
		}

		if ai.FlashFreeze.IsReady(sim) && sim.CurrentTime >= ai.FlashFreeze.CD.Duration {
			ai.FlashFreeze.Cast(sim, nil)
		}
	}

	if ai.Target.GCD.IsReady(sim) {
		nextEventAt := sim.CurrentTime + time.Minute

		// All possible next events
		events := []time.Duration{
			max(ai.FrozenBlows.ReadyAt(), ai.FrozenBlows.CD.Duration),
			max(ai.FlashFreeze.ReadyAt(), ai.FlashFreeze.CD.Duration),
			ai.NextStorms,
		}

		if ai.Target.Env.Raid.Size() == 1 {
			// Individual Sim approximation - taken from some random logs
			timeBetweenStacks := 400 * time.Millisecond // TODO: Expose this
			events = append(events, max(ai.ToastyFireTime, sim.CurrentTime+timeBetweenStacks))
		} else {
			timeBetweenNewCampfires := 3 * time.Second // TODO: Improve on Fires Approximation by actually simulating active campfires
			events = append(events, max(ai.ToastyFireTime, sim.CurrentTime+timeBetweenNewCampfires))
		}

		// if ai.Target.CurrentTarget != nil {
		// 	events = append(events, max(ai.PhasePunch.ReadyAt(), ai.PhasePunch.CD.Duration))
		// 	events = append(events, max(ai.QuantumStrike.ReadyAt(), ai.QuantumStrike.CD.Duration))
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

func (ai *HodirAI) stormCloudPrioApply(sim *core.Simulation, maxBuffs int, targets []int) int {
	// Loop over prio targets
	for _, target := range targets {
		if maxBuffs > 0 {
			if !ai.StormCloud[target].IsActive() {
				ai.StormCloud[target].Activate(sim)
				maxBuffs = maxBuffs - 1
			}
		}
	}
	return maxBuffs
}
