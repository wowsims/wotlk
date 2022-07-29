package retribution

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/paladin"
)

func RegisterRetributionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_RetributionPaladin{},
		proto.Spec_SpecRetributionPaladin,
		func(character core.Character, options proto.Player) core.Agent {
			return NewRetributionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RetributionPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Retribution Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRetributionPaladin(character core.Character, options proto.Player) *RetributionPaladin {
	retOptions := options.GetRetributionPaladin()

	ret := &RetributionPaladin{
		Paladin:       paladin.NewPaladin(character, *retOptions.Talents),
		Rotation:      *retOptions.Rotation,
		Judgement:     retOptions.Options.Judgement,
		Seal:          retOptions.Options.Seal,
		UseDivinePlea: retOptions.Options.UseDivinePlea,
		ExoSlack:      retOptions.Rotation.ExoSlack,
		ConsSlack:     retOptions.Rotation.ConsSlack,
	}
	ret.PaladinAura = retOptions.Options.Aura

	// Convert DTPS option to bonus MP5
	spAtt := retOptions.Options.DamageTakenPerSecond * 5.0 / 10.0
	ret.AddStat(stats.MP5, spAtt)

	ret.EnableAutoAttacks(ret, core.AutoAttackOptions{
		MainHand:       ret.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		AutoSwingMelee: true,
	})

	ret.EnableResumeAfterManaWait(ret.OnGCDReady)

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	Judgement     proto.PaladinJudgement
	Seal          proto.PaladinSeal
	UseDivinePlea bool
	ExoSlack      int32
	ConsSlack     int32

	SealInitComplete       bool
	DivinePleaInitComplete bool

	Rotation proto.RetributionPaladin_Rotation
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
	ret.RegisterAvengingWrathCD()

	ret.DelayDPSCooldownsForArmorDebuffs()
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
	ret.AutoAttacks.CancelAutoSwing(sim)
	ret.SealInitComplete = false
	ret.DivinePleaInitComplete = false
}
