package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

var TalentTreeSizes = [3]int{17, 17, 16}

type Warlock struct {
	core.Character
	Talents *proto.WarlockTalents
	Options *proto.WarlockOptions

	Pet *WarlockPet

	ChaosBolt    *core.Spell
	Conflagrate  *core.Spell
	Corruption   *core.Spell
	DarkPact     *core.Spell
	DrainSoul    *core.Spell
	Haunt        *core.Spell
	Immolate     *core.Spell
	Incinerate   *core.Spell
	LifeTap      *core.Spell
	SearingPain  *core.Spell
	ShadowBolt   *core.Spell
	ShadowCleave *core.Spell
	Shadowburn   *core.Spell
	SoulFire     *core.Spell
	DemonicGrace *core.Spell
	DrainLife    *core.Spell
	RainOfFire   *core.Spell

	CurseOfElements      *core.Spell
	CurseOfElementsAuras core.AuraArray
	CurseOfWeakness      *core.Spell
	CurseOfWeaknessAuras core.AuraArray
	CurseOfTongues       *core.Spell
	CurseOfTonguesAuras  core.AuraArray
	CurseOfAgony         *core.Spell
	CurseOfDoom          *core.Spell
	AmplifyCurse         *core.Spell

	DemonicEmpowerment      *core.Spell
	DemonicEmpowermentAura  *core.Aura
	DemonicSoulAura         *core.Aura
	HauntDebuffAuras        core.AuraArray
	ImmolationAura          *core.Spell
	IncinerateAura          *core.Aura
	Metamorphosis           *core.Spell
	MetamorphosisAura       *core.Aura
	NightfallProcAura       *core.Aura
	PyroclasmAura           *core.Aura
	DemonicGraceAura        *core.Aura
	AmplifyCurseAura        *core.Aura
	ImprovedShadowBoltAuras core.AuraArray
	LakeOfFireAuras         core.AuraArray

	// The sum total of demonic pact spell power * seconds.
	DPSPAggregate float64
	PreviousTime  time.Duration

	petStmBonusSP float64
}

func (warlock *Warlock) GetCharacter() *core.Character {
	return &warlock.Character
}

func (warlock *Warlock) GetWarlock() *Warlock {
	return warlock
}

func (warlock *Warlock) Initialize() {
	warlock.registerChaosBoltSpell()
	warlock.registerCorruptionSpell()
	warlock.registerImmolateSpell()
	warlock.registerIncinerateSpell()
	warlock.registerShadowBoltSpell()
	warlock.registerShadowCleaveSpell()
	// warlock.registerCurseOfElementsSpell()
	// warlock.registerCurseOfWeaknessSpell()
	// warlock.registerCurseOfTonguesSpell()
	warlock.registerCurseOfAgonySpell()
	warlock.registerAmplifyCurseSpell()
	// warlock.registerCurseOfDoomSpell()
	warlock.registerLifeTapSpell()
	// warlock.registerSeedSpell()
	// warlock.registerSoulFireSpell()
	// warlock.registerUnstableAfflictionSpell()
	// warlock.registerDrainSoulSpell()
	// warlock.registerConflagrateSpell()
	warlock.registerHauntSpell()
	// warlock.registerDemonicEmpowermentSpell()
	warlock.registerMetamorphosisSpell()
	// warlock.registerDarkPactSpell()
	warlock.registerShadowBurnSpell()
	warlock.registerSearingPainSpell()
	// warlock.registerInfernoSpell()
	// warlock.registerBlackBook()
	warlock.registerDemonicGraceSpell()
	warlock.registerDrainLifeSpell()
	warlock.registerRainOfFireSpell()
}

func (warlock *Warlock) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.BloodPact = max(raidBuffs.BloodPact, core.MakeTristateValue(
		warlock.Options.Summon == proto.WarlockOptions_Imp,
		warlock.Talents.ImprovedImp == 3,
	))
}

func (warlock *Warlock) Reset(sim *core.Simulation) {
	if sim.CurrentTime == 0 {
		warlock.petStmBonusSP = 0
	}

	// warlock.ItemSwap.SwapItems(sim, []proto.ItemSlot{proto.ItemSlot_ItemSlotMainHand,
	// 	proto.ItemSlot_ItemSlotOffHand, proto.ItemSlot_ItemSlotRanged}, false)
	// warlock.setupCooldowns(sim)
}

func NewWarlock(character *core.Character, options *proto.Player, warlockOptions *proto.WarlockOptions) *Warlock {
	warlock := &Warlock{
		Character: *character,
		Talents:   &proto.WarlockTalents{},
		Options:   warlockOptions,
	}
	core.FillTalentsProto(warlock.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)
	warlock.EnableManaBar()

	warlock.AddStatDependency(stats.Strength, stats.AttackPower, 1)
	warlock.AddStatDependency(stats.Intellect, stats.SpellCrit, core.CritPerIntAtLevel[warlock.Class][int(warlock.Level)]*core.SpellCritRatingPerCritChance)

	if warlock.Options.Summon != proto.WarlockOptions_NoSummon {
		warlock.Pet = warlock.NewWarlockPet()
	}

	warlock.applyWeaponImbue()

	return warlock
}

func (warlock *Warlock) HasRune(rune proto.WarlockRune) bool {
	return warlock.HasRuneById(int32(rune))
}

func (warlock *Warlock) OnGCDReady(sim *core.Simulation) {
}

// Agent is a generic way to access underlying warlock on any of the agents.
type WarlockAgent interface {
	GetWarlock() *Warlock
}
