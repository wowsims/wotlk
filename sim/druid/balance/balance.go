package balance

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/druid"
	"sort"
)

func RegisterBalanceDruid() {
	core.RegisterAgentFactory(
		proto.Player_BalanceDruid{},
		proto.Spec_SpecBalanceDruid,
		func(character core.Character, options proto.Player) core.Agent {
			return NewBalanceDruid(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_BalanceDruid)
			if !ok {
				panic("Invalid spec value for Balance Druid!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewBalanceDruid(character core.Character, options proto.Player) *BalanceDruid {
	balanceOptions := options.GetBalanceDruid()
	selfBuffs := druid.SelfBuffs{}

	if balanceOptions.Options.InnervateTarget != nil {
		selfBuffs.InnervateTarget = *balanceOptions.Options.InnervateTarget
	} else {
		selfBuffs.InnervateTarget.TargetIndex = -1
	}

	moonkin := &BalanceDruid{
		Druid:    druid.New(character, druid.Moonkin, selfBuffs, *balanceOptions.Talents),
		Rotation: *balanceOptions.Rotation,
	}

	moonkin.ResetTalentsBonuses()
	moonkin.RegisterTalentsBonuses()
	moonkin.RegisterCooldownRankings()
	moonkin.EnableResumeAfterManaWait(moonkin.tryUseGCD)
	return moonkin
}

type BalanceCooldown struct {
	Name     string
	Stat     stats.Stat
	Priority float64
	ID       core.ActionID
}

type BalanceDruid struct {
	*druid.Druid

	Rotation           proto.BalanceDruid_Rotation
	CooldownsRanking   []BalanceCooldown
	CooldownsAvailable []*core.MajorCooldown

	// CDS
	hyperSpeedMCD *core.MajorCooldown
	potionMCD     *core.MajorCooldown
	onUseTrinket1 *core.MajorCooldown
	onUseTrinket2 *core.MajorCooldown
	potionUsed    bool
}

func (moonkin *BalanceDruid) GetDruid() *druid.Druid {
	return moonkin.Druid
}

func (moonkin *BalanceDruid) Initialize() {
	moonkin.Druid.Initialize()
	moonkin.RegisterBalanceSpells()
}

func (moonkin *BalanceDruid) Reset(sim *core.Simulation) {
	moonkin.Druid.Reset(sim)
	moonkin.RebirthTiming = moonkin.Env.BaseDuration.Seconds() * sim.RandomFloat("Rebirth Timing")

	if moonkin.Rotation.UseSmartCooldowns == true {
		moonkin.potionUsed = false
		consumes := &moonkin.Consumes

		if consumes.DefaultPotion == proto.Potions_PotionOfSpeed {
			moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40211})
		}
		if consumes.DefaultPotion == proto.Potions_PotionOfWildMagic {
			moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: 40212})
		}
		if moonkin.HasProfession(proto.Profession_Engineering) {
			moonkin.getBalanceMajorCooldown(core.ActionID{SpellID: 54758})
		}
		moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[items.ItemSlotTrinket1].ID})
		moonkin.getBalanceMajorCooldown(core.ActionID{ItemID: moonkin.Equip[items.ItemSlotTrinket2].ID})

		// Sort this array depending on moonkin.CooldownRankings
		sort.Slice(moonkin.CooldownsAvailable, func(i, j int) bool {
			return moonkin.CooldownsAvailable[i].Priority < moonkin.CooldownsAvailable[j].Priority
		})
	}
}

// Takes out a Cooldown from the generic MajorCooldownManager and adds it to a custom Slice of Cooldowns
// Sets a stat and a priority for the cooldown, depending on moonkin.CooldownsRanking
func (moonkin *BalanceDruid) getBalanceMajorCooldown(actionID core.ActionID) {
	if moonkin.Character.HasMajorCooldown(actionID) {
		majorCd := moonkin.Character.GetMajorCooldown(actionID)
		majorCd.ShouldActivate = func(sim *core.Simulation, character *core.Character) bool {
			return false
		}
		for _, v := range moonkin.CooldownsRanking {
			if v.ID.SameAction(majorCd.Spell.ActionID) {
				majorCd.Priority = v.Priority
				majorCd.Spell.ResourceType = v.Stat
				moonkin.CooldownsAvailable = append(moonkin.CooldownsAvailable, majorCd)
			}
		}
	}
}

func (moonkin *BalanceDruid) RegisterCooldownRankings() {

	moonkin.CooldownsRanking = []BalanceCooldown{
		{
			Name:     "Potion of Speed",
			Stat:     stats.SpellHaste,
			Priority: 500,
			ID:       core.ActionID{ItemID: 40211},
		},
		{
			Name:     "Hyperspeed Acceleration",
			Stat:     stats.SpellHaste,
			Priority: 340,
			ID:       core.ActionID{SpellID: 54758},
		},
		{
			Name:     "Potion of Wild Magic",
			Stat:     stats.SpellCrit,
			Priority: 200,
			ID:       core.ActionID{ItemID: 40212},
		},
		{
			Name:     "Mark of the War Prisoner",
			Stat:     stats.SpellPower,
			Priority: 346,
			ID:       core.ActionID{ItemID: 37873},
		},
	}
}
