package enhancement

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/shaman"
)

func RegisterEnhancementShaman() {
	core.RegisterAgentFactory(
		proto.Player_EnhancementShaman{},
		proto.Spec_SpecEnhancementShaman,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewEnhancementShaman(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_EnhancementShaman)
			if !ok {
				panic("Invalid spec value for Enhancement Shaman!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewEnhancementShaman(character core.Character, options *proto.Player) *EnhancementShaman {
	enhOptions := options.GetEnhancementShaman()

	selfBuffs := shaman.SelfBuffs{
		Bloodlust: enhOptions.Options.Bloodlust,
		Shield:    enhOptions.Options.Shield,
		ImbueMH:   enhOptions.Options.ImbueMh,
		ImbueOH:   enhOptions.Options.ImbueOh,
	}

	totems := &proto.ShamanTotems{}
	if enhOptions.Rotation.Totems != nil {
		totems = enhOptions.Rotation.Totems
	}

	enh := &EnhancementShaman{
		Shaman: shaman.NewShaman(character, enhOptions.Talents, totems, selfBuffs, true),
	}

	enh.EnableResumeAfterManaWait(enh.OnGCDReady)
	enh.rotation = NewPriorityRotation(enh, enhOptions.Rotation)

	// Enable Auto Attacks for this spec
	enh.EnableAutoAttacks(enh, core.AutoAttackOptions{
		MainHand:       enh.WeaponFromMainHand(enh.DefaultMeleeCritMultiplier()),
		OffHand:        enh.WeaponFromOffHand(enh.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
		SyncType:       int32(enhOptions.Options.SyncType),
	})

	if enhOptions.Options.ItemSwap != nil {
		if enhOptions.Options.ItemSwap.MhItem != nil {
			itemSpec := core.ItemSpec{
				ID:      enhOptions.Options.ItemSwap.MhItem.Id,
				Gems:    enhOptions.Options.ItemSwap.MhItem.Gems,
				Enchant: enhOptions.Options.ItemSwap.MhItem.Enchant,
			}
			item := core.NewItem(itemSpec)
			enh.mh = &item
		}

		if enhOptions.Options.ItemSwap.OhItem != nil {
			itemSpec := core.ItemSpec{
				ID:      enhOptions.Options.ItemSwap.OhItem.Id,
				Gems:    enhOptions.Options.ItemSwap.OhItem.Gems,
				Enchant: enhOptions.Options.ItemSwap.OhItem.Enchant,
			}
			item := core.NewItem(itemSpec)
			enh.oh = &item
		}

	}

	if !enh.HasMHWeapon() {
		enh.SelfBuffs.ImbueMH = proto.ShamanImbue_NoImbue
	}
	if !enh.HasOHWeapon() {
		enh.SelfBuffs.ImbueOH = proto.ShamanImbue_NoImbue
	}
	enh.ApplyWindfuryImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_WindfuryWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_WindfuryWeapon)
	enh.ApplyFlametongueImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeapon)
	enh.ApplyFlametongueDownrankImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeaponDownrank,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FlametongueWeaponDownrank)
	enh.ApplyFrostbrandImbue(
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FrostbrandWeapon,
		enh.SelfBuffs.ImbueOH == proto.ShamanImbue_FrostbrandWeapon)

	if enh.SelfBuffs.ImbueMH == proto.ShamanImbue_WindfuryWeapon ||
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeapon ||
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FlametongueWeaponDownrank ||
		enh.SelfBuffs.ImbueMH == proto.ShamanImbue_FrostbrandWeapon {
	}

	enh.SpiritWolves = &shaman.SpiritWolves{
		SpiritWolf1: enh.NewSpiritWolf(1),
		SpiritWolf2: enh.NewSpiritWolf(2),
	}

	enh.ShamanisticRageManaThreshold = enhOptions.Rotation.ShamanisticRageManaThreshold

	return enh
}

type EnhancementShaman struct {
	*shaman.Shaman

	rotation Rotation

	mh *core.Item
	oh *core.Item

	scheduler common.GCDScheduler
}

func (enh *EnhancementShaman) GetShaman() *shaman.Shaman {
	return enh.Shaman
}

func (enh *EnhancementShaman) Initialize() {
	enh.Shaman.Initialize()
	enh.DelayDPSCooldowns(3 * time.Second)
}

func (enh *EnhancementShaman) Reset(sim *core.Simulation) {
	enh.Shaman.Reset(sim)

	mcd := enh.GetMajorCooldown(enh.FireElementalTotem.ActionID)
	oldShouldActive := mcd.ShouldActivate
	mcd.ShouldActivate = func(s *core.Simulation, c *core.Character) bool {
		success := oldShouldActive(s, c)

		if success {
			swapped := false
			if enh.mh != nil {
				swappMh := enh.mh
				currentMh := enh.GetMHWeapon()
				newStats := swappMh.Stats.Add(currentMh.Stats.Multiply(-1))

				spBonus := 211.0
				spMod := 1.0 + 0.1*float64(enh.Talents.ElementalWeapons)
				newStats = newStats.Add(stats.Stats{stats.SpellPower: spBonus * spMod})

				enh.AddStatsDynamic(s, newStats)

				if sim.Log != nil {
					sim.Log("Swapping Main Hand: %v", newStats)
				}
				swapped = true
			}

			if enh.oh != nil {
				swappMh := enh.oh
				currentWep := enh.GetOHWeapon()
				newStats := swappMh.Stats.Add(currentWep.Stats.Multiply(-1))

				spBonus := 211.0
				spMod := 1.0 + 0.1*float64(enh.Talents.ElementalWeapons)
				newStats = newStats.Add(stats.Stats{stats.SpellPower: spBonus * spMod})

				enh.AddStatsDynamic(s, newStats)

				if sim.Log != nil {
					sim.Log("Swapping Off Hand: %v", newStats)
				}
				swapped = true
			}

			if swapped {
				enh.AutoAttacks.StopMeleeUntil(s, s.CurrentTime)
				core.StartDelayedAction(s, core.DelayedActionOptions{
					DoAt: s.CurrentTime + 1500*time.Millisecond,
					OnAction: func(s *core.Simulation) {
						newStats := stats.Stats{}
						if enh.mh != nil {
							newStats = enh.GetMHWeapon().Stats.Add(enh.mh.Stats.Multiply(-1))
						} else if enh.oh != nil {
							newStats = newStats.Add(enh.GetOHWeapon().Stats.Add(enh.oh.Stats.Multiply(-1)))
						}

						enh.AddStatsDynamic(s, newStats)
						enh.AutoAttacks.StopMeleeUntil(s, s.CurrentTime)
					},
				})
			}
		}

		return success
	}
}

func (enh *EnhancementShaman) CastLightningBoltWeave(sim *core.Simulation, reactionTime time.Duration) bool {
	previousAttack := sim.CurrentTime - enh.AutoAttacks.PreviousSwingAt
	reactionTime = core.TernaryDuration(previousAttack < reactionTime, reactionTime-previousAttack, 0)

	//calculate cast times for weaving
	lbCastTime := enh.ApplyCastSpeed(enh.LightningBolt.DefaultCast.CastTime-(time.Millisecond*time.Duration(500*enh.MaelstromWeaponAura.GetStacks()))) + reactionTime
	//calculate swing times for weaving
	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime

	if lbCastTime < timeUntilSwing {
		if reactionTime > 0 {
			reactionTime += sim.CurrentTime

			enh.HardcastWaitUntil(sim, reactionTime, func(_ *core.Simulation, _ *core.Unit) {
				enh.GCD.Reset()
				enh.LightningBolt.Cast(sim, enh.CurrentTarget)
			})

			enh.WaitUntil(sim, reactionTime)
			return true
		}
		return enh.LightningBolt.Cast(sim, enh.CurrentTarget)
	}

	return false
}

func (enh *EnhancementShaman) CastLavaBurstWeave(sim *core.Simulation, reactionTime time.Duration) bool {
	previousAttack := sim.CurrentTime - enh.AutoAttacks.PreviousSwingAt
	reactionTime = core.TernaryDuration(previousAttack < reactionTime, reactionTime-previousAttack, 0)

	//calculate cast times for weaving
	lvbCastTime := enh.ApplyCastSpeed(enh.LavaBurst.DefaultCast.CastTime) + reactionTime
	//calculate swing times for weaving
	timeUntilSwing := enh.AutoAttacks.NextAttackAt() - sim.CurrentTime

	if lvbCastTime < timeUntilSwing {
		if reactionTime > 0 {
			reactionTime += sim.CurrentTime

			enh.HardcastWaitUntil(sim, reactionTime, func(_ *core.Simulation, _ *core.Unit) {
				enh.GCD.Reset()
				enh.LavaBurst.Cast(sim, enh.CurrentTarget)
			})

			enh.WaitUntil(sim, reactionTime)
			return true
		}

		return enh.LavaBurst.Cast(sim, enh.CurrentTarget)
	}

	return false
}
