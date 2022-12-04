package mage

import (
	"time"

	"github.com/wowsims/wotlk/sim/common"
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

const (
	SpellFlagMage   = core.SpellFlagAgentReserved1
	BarrageSpells   = core.SpellFlagAgentReserved2
	HotStreakSpells = core.SpellFlagAgentReserved3
)

func RegisterMage() {
	core.RegisterAgentFactory(
		proto.Player_Mage{},
		proto.Spec_SpecMage,
		func(character core.Character, options *proto.Player) core.Agent {
			return NewMage(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Mage)
			if !ok {
				panic("Invalid spec value for Mage!")
			}
			player.Spec = playerSpec
		},
	)
}

type Mage struct {
	core.Character
	Talents *proto.MageTalents

	Options  *proto.Mage_Options
	Rotation *proto.Mage_Rotation

	isMissilesBarrage        bool
	isMissilesBarrageVisible bool
	numCastsDone             int32
	num4CostAB               int32
	extraABsAP               int32
	disabledMCDs             []*core.MajorCooldown

	waterElemental *WaterElemental
	mirrorImage    *MirrorImage

	// Cached values for a few mechanics.
	spellDamageMultiplier float64

	// Current bonus crit from AM+CC interaction.
	bonusAMCCCrit   float64
	bonusCritDamage float64

	ArcaneBlast     *core.Spell
	ArcaneExplosion *core.Spell
	ArcaneMissiles  *core.Spell
	Blizzard        *core.Spell
	DeepFreeze      *core.Spell
	Ignite          *core.Spell
	LivingBomb      *core.Spell
	Fireball        *core.Spell
	FireBlast       *core.Spell
	Flamestrike     *core.Spell
	Frostbolt       *core.Spell
	FrostfireBolt   *core.Spell
	Pyroblast       *core.Spell
	Scorch          *core.Spell
	WintersChill    *core.Spell
	MirrorImage     *core.Spell

	IcyVeins             *core.Spell
	SummonWaterElemental *core.Spell

	ArcaneMissilesDot *core.Dot
	LivingBombDot     *core.Dot // living bomb is used for single-target only, currently
	FireballDot       *core.Dot
	FlamestrikeDot    *core.Dot
	FrostfireDot      *core.Dot
	PyroblastDot      *core.Dot

	ArcaneBlastAura    *core.Aura
	MissileBarrageAura *core.Aura
	ClearcastingAura   *core.Aura
	ScorchAura         *core.Aura
	HotStreakAura      *core.Aura
	CombustionAura     *core.Aura
	FingersOfFrostAura *core.Aura
	BrainFreezeAura    *core.Aura

	// Used to prevent utilising Brain Freeze immediately after proccing it.
	BrainFreezeActivatedAt time.Duration

	IgniteDots []*core.Dot

	manaTracker common.ManaSpendingRateTracker
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) HasMajorGlyph(glyph proto.MageMajorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}
func (mage *Mage) HasMinorGlyph(glyph proto.MageMinorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true

	if mage.Talents.ArcaneEmpowerment == 3 {
		raidBuffs.ArcaneEmpowerment = true
	}
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) Initialize() {
	mage.registerArcaneBlastSpell()
	mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerBlizzardSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlamestrikeSpell()
	mage.registerFrostboltSpell()
	mage.registerPyroblastSpell()
	mage.registerScorchSpell()
	mage.registerWintersChillSpell()
	mage.registerLivingBombSpell()
	mage.registerFrostfireBoltSpell()

	mage.registerEvocationCD()
	mage.registerManaGemsCD()
	mage.registerMirrorImageCD()

	mage.num4CostAB = 0
	mage.extraABsAP = mage.Rotation.ExtraBlastsDuringFirstAp
}

func (mage *Mage) launchExecuteCDOptimizer(sim *core.Simulation) {

	pa := &core.PendingAction{
		Priority: core.ActionPriorityRegen,
	}
	pa.OnAction = func(sim *core.Simulation) {
		if sim.IsExecutePhase35() {
			for _, mcd := range mage.disabledMCDs {
				mcd.Enable()
			}
			// TODO looks fishy, since disabledMCDs isn't emptied; also, this could use an executePhaseCallback instead
		} else {
			for _, mcd := range mage.GetMajorCooldowns() {
				isBloodLust := mcd.Spell.ActionID == core.ActionID{SpellID: 2825, Tag: -1} //ignore blood lust as it shouldn't be saved
				isFlameCap := mcd.Spell.ActionID == core.ActionID{ItemID: 22788}           //ignore flame cap because it's so long
				isPotionOfSpeed := mcd.Spell.ActionID == core.ActionID{ItemID: 40211}
				if mcd.Spell.CD.Duration > (sim.Duration-sim.CurrentTime) && mcd.Type.Matches(core.CooldownTypeDPS) &&
					!isBloodLust && !isFlameCap || isPotionOfSpeed {
					mcd.Disable()
					mage.disabledMCDs = append(mage.disabledMCDs, mcd)
				}
			}

			pa.NextActionAt = sim.CurrentTime + core.MinDuration(40*time.Second, time.Duration(.35*float64(sim.Duration)))

			executeTime := time.Duration(.7 * float64(sim.Duration))

			if pa.NextActionAt > executeTime {
				pa.NextActionAt = executeTime
			}
			if pa.NextActionAt < sim.Duration {
				sim.AddPendingAction(pa)
			}

		}

	}

	pa.OnAction(sim) // immediately activate first pending action
}

func (mage *Mage) Reset(sim *core.Simulation) {
	mage.numCastsDone = 0
	mage.num4CostAB = 0
	mage.extraABsAP = mage.Rotation.ExtraBlastsDuringFirstAp
	mage.manaTracker.Reset()
	mage.bonusAMCCCrit = 0

	if mage.Rotation.Type == proto.Mage_Rotation_Fire && mage.Rotation.OptimizeCdsForExecute { // make this an option
		mage.disabledMCDs = make([]*core.MajorCooldown, 0, 10)
		mage.launchExecuteCDOptimizer(sim)
	}
}

func NewMage(character core.Character, options *proto.Player) *Mage {
	mageOptions := options.GetMage()

	mage := &Mage{
		Character: character,
		Talents:   mageOptions.Talents,
		Options:   mageOptions.Options,
		Rotation:  mageOptions.Rotation,

		spellDamageMultiplier: 1.0,
		manaTracker:           common.NewManaSpendingRateTracker(),
	}
	mage.bonusCritDamage = .25*float64(mage.Talents.SpellPower) + .1*float64(mage.Talents.Burnout)
	mage.EnableManaBar()
	mage.EnableResumeAfterManaWait(mage.tryUseGCD)

	if mage.Options.Armor == proto.Mage_Options_MageArmor {
		mage.PseudoStats.SpiritRegenRateCasting += .5
		if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMageArmor) {
			mage.PseudoStats.SpiritRegenRateCasting += .2
		}
		if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
			mage.PseudoStats.SpiritRegenRateCasting += .1
		}
	} else if mage.Options.Armor == proto.Mage_Options_MoltenArmor {
		//Need to switch to spirit crit calc
		multi := 0.35
		if mage.HasMajorGlyph(proto.MageMajorGlyph_GlyphOfMoltenArmor) {
			multi += .2
		}
		if mage.HasSetBonus(ItemSetKhadgarsRegalia, 2) {
			multi += .15
		}
		mage.Character.AddStatDependency(stats.Spirit, stats.SpellCrit, multi)
	}

	mage.mirrorImage = mage.NewMirrorImage()

	if mage.Talents.SummonWaterElemental {
		mage.waterElemental = mage.NewWaterElemental(mage.Rotation.WaterElementalDisobeyChance)
	}

	return mage
}

func init() {
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceBloodElf, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  30,
		stats.Agility:   41,
		stats.Stamina:   50,
		stats.Intellect: 185,
		stats.Spirit:    173,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.926,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceDraenei, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  34,
		stats.Agility:   36,
		stats.Stamina:   50,
		stats.Intellect: 182,
		stats.Spirit:    176,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.933,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceGnome, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  28,
		stats.Agility:   42,
		stats.Stamina:   50,
		stats.Intellect: 193, // Gnomes start with 162 int, we assume this include racial so / 1.05
		stats.Spirit:    174,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.93,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceHuman, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  33,
		stats.Agility:   39,
		stats.Stamina:   51,
		stats.Intellect: 181,
		stats.Spirit:    179,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.926,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceTroll, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  34,
		stats.Agility:   41,
		stats.Stamina:   52,
		stats.Intellect: 177,
		stats.Spirit:    175,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.935,
	}
	core.BaseStats[core.BaseStatsKey{Race: proto.Race_RaceUndead, Class: proto.Class_ClassMage}] = stats.Stats{
		stats.Health:    3213,
		stats.Strength:  32,
		stats.Agility:   37,
		stats.Stamina:   52,
		stats.Intellect: 179,
		stats.Spirit:    179,
		stats.Mana:      3268,
		stats.SpellCrit: core.CritRatingPerCritChance * 0.930,
	}
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}
