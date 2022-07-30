package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type ProcStatBonusEffect struct {
	Name       string
	ID         int32
	Bonus      stats.Stats
	Duration   time.Duration
	Callback   Callback
	ProcMask   core.ProcMask
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	PPM        float64
	ICD        time.Duration
}

func newProcStatBonusEffect(config ProcStatBonusEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()
		procAura := character.NewTemporaryStatsAura(config.Name+" Proc", core.ActionID{ItemID: config.ID}, config.Bonus, config.Duration)

		MakeProcTriggerAura(&character.Unit, ProcTrigger{
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			PPM:        config.PPM,
			ICD:        config.ICD,
			Handler: func(sim *core.Simulation, _ *core.Spell, _ *core.SpellEffect) {
				procAura.Activate(sim)
			},
		})
	})
}

func init() {
	// TODO: <https://www.wowhead.com/wotlk/item=45507/the-generals-heart>
	//  Equip: When struck in combat has a chance of shielding you in a protective barrier which will reduce damage from each attack by 205.  Lasts 10 secs. (Proc chance: 5%, 50s cooldown)

	// Keep these separated by stat, ordered by item ID within each group.
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Meteorite Whetstone",
		ID:         37390,
		Bonus:      stats.Stats{stats.MeleeHaste: 444, stats.SpellHaste: 444},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Serrah's Star",
		ID:         37559,
		Bonus:      stats.Stats{stats.MeleeCrit: 167, stats.SpellCrit: 167},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.45,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Spark of Life",
		ID:         37657,
		Bonus:      stats.Stats{stats.MP5: 220},
		Duration:   time.Second * 15,
		Callback:   OnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Forge Ember",
		ID:         37660,
		Bonus:      stats.Stats{stats.SpellPower: 512, stats.HealingPower: 512},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})

	core.AddEffectsToTest = false

	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Je'Tze's Bell",
		ID:         37835,
		Bonus:      stats.Stats{stats.MP5: 125},
		Duration:   time.Second * 15,
		Callback:   OnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Valonforth's Remembrance",
		ID:         38071,
		Bonus:      stats.Stats{stats.Spirit: 222},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Embrace of the Spider",
		ID:         39229,
		Bonus:      stats.Stats{stats.MeleeHaste: 505, stats.SpellHaste: 505},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Dying Curse",
		ID:         40255,
		Bonus:      stats.Stats{stats.SpellPower: 765, stats.HealingPower: 765},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Grim Toll",
		ID:         40256,
		Bonus:      stats.Stats{stats.ArmorPenetration: 612},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sundial of the Exiled",
		ID:         40682,
		Bonus:      stats.Stats{stats.SpellPower: 590, stats.HealingPower: 590},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Mirror of Truth",
		ID:         40684,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sonic Booster",
		ID:         40767,
		Bonus:      stats.Stats{stats.AttackPower: 430, stats.RangedAttackPower: 430},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt | OnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.35,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Tears of Bitter Anguish",
		ID:         43573,
		Bonus:      stats.Stats{stats.MeleeHaste: 410, stats.SpellHaste: 410},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Crusader's Locket",
		ID:         43829,
		Bonus:      stats.Stats{stats.Expertise: 258},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Chuchu's Tiny Box of Horrors",
		ID:         43838,
		Bonus:      stats.Stats{stats.MeleeCrit: 258, stats.SpellCrit: 258},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Signet of Edward the Odd",
		ID:         44308,
		Bonus:      stats.Stats{stats.MeleeHaste: 125, stats.SpellHaste: 125},
		Duration:   time.Second * 13,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Flow of Knowledge",
		ID:         44912,
		Bonus:      stats.Stats{stats.SpellPower: 590, stats.HealingPower: 590},
		Duration:   time.Second * 10,
		Callback:   OnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Anvil of Titans",
		ID:         44914,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Jouster's Fury Alliance",
		ID:         45131,
		Bonus:      stats.Stats{stats.MeleeCrit: 328, stats.SpellCrit: 328},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Jouster's Fury Horde",
		ID:         45219,
		Bonus:      stats.Stats{stats.MeleeCrit: 328, stats.SpellCrit: 328},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Pyrite Infuser",
		ID:         45286,
		Bonus:      stats.Stats{stats.AttackPower: 1234},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Pandora's Plea",
		ID:         45490,
		Bonus:      stats.Stats{stats.SpellPower: 751, stats.HealingPower: 751},
		Duration:   time.Second * 10,
		Callback:   OnCastComplete,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Blood of the Old God",
		ID:         45522,
		Bonus:      stats.Stats{stats.AttackPower: 1284},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.1, // wowhead shows proc chance: 10% but a comment says 1.12PPM? TODO: validate.
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Flare of the Heavens",
		ID:         45518,
		Bonus:      stats.Stats{stats.SpellPower: 850, stats.HealingPower: 850},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Show of Faith",
		ID:         45535,
		Bonus:      stats.Stats{stats.MP5: 241},
		Duration:   time.Second * 15,
		Callback:   OnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Comet's Trail",
		ID:         45609,
		Bonus:      stats.Stats{stats.SpellHaste: 726, stats.MeleeHaste: 726},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Elemental Focus Stone",
		ID:         45866,
		Bonus:      stats.Stats{stats.MeleeHaste: 522, stats.SpellHaste: 522},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sif's Remembrance",
		ID:         45929,
		Bonus:      stats.Stats{stats.MP5: 195},
		Duration:   time.Second * 15,
		Callback:   OnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Mjolnir Runestone",
		ID:         45931,
		Bonus:      stats.Stats{stats.ArmorPenetration: 665},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Dark Matter",
		ID:         46038,
		Bonus:      stats.Stats{stats.MeleeCrit: 612, stats.SpellCrit: 612},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Abyssal Rune",
		ID:         47213,
		Bonus:      stats.Stats{stats.SpellPower: 590, stats.HealingPower: 590},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Harmful:    true,
		ProcChance: 0.25,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Banner of Victory",
		ID:         47214,
		Bonus:      stats.Stats{stats.AttackPower: 1008, stats.RangedAttackPower: 1008},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "The Black Heart",
		ID:         47216,
		Bonus:      stats.Stats{stats.Armor: 7056},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Coren's Chromium Coaster",
		ID:         49074,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Mithril Pocketwatch",
		ID:         49076,
		Bonus:      stats.Stats{stats.SpellPower: 590, stats.HealingPower: 590},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Needle-Encrusted Scorpion",
		ID:         50198,
		Bonus:      stats.Stats{stats.ArmorPenetration: 678},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Whispering Fanged Skull",
		ID:         50342,
		Bonus:      stats.Stats{stats.AttackPower: 1100, stats.RangedAttackPower: 1100},
		Duration:   time.Second * 15,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskPeriodicDamage,
		Harmful:    true,
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Whispering Fanged Skull H",
		ID:         50343,
		Bonus:      stats.Stats{stats.AttackPower: 1250, stats.RangedAttackPower: 1250},
		Duration:   time.Second * 15,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskPeriodicDamage,
		Harmful:    true,
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Purified Lunar Dust",
		ID:         50358,
		Bonus:      stats.Stats{stats.MP5: 304},
		Duration:   time.Second * 15,
		Callback:   OnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Phylactery of the Nameless Lich",
		ID:         50360,
		Bonus:      stats.Stats{stats.SpellPower: 1074, stats.HealingPower: 1074},
		Duration:   time.Second * 20,
		Callback:   OnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskPeriodicDamage,
		ProcChance: 0.30,
		ICD:        time.Second * 100,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Phylactery of the Nameless Lich H",
		ID:         50365,
		Bonus:      stats.Stats{stats.SpellPower: 1207, stats.HealingPower: 1207},
		Duration:   time.Second * 20,
		Callback:   OnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskPeriodicDamage,
		ProcChance: 0.30,
		ICD:        time.Second * 100,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Ashen Band of Unmatched Destruction",
		ID:         50397,
		Bonus:      stats.Stats{stats.SpellPower: 285, stats.HealingPower: 285},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Ashen Band of Endless Destruction",
		ID:         50398,
		Bonus:      stats.Stats{stats.SpellPower: 285, stats.HealingPower: 285},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:     "Ashen Band of Unmatched Vengeance",
		ID:       50401,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: OnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:     "Ashen Band of Endless Vengeance",
		ID:       50402,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: OnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Ashen Band of Unmatched Courage",
		ID:         50403,
		Bonus:      stats.Stats{stats.Armor: 2400},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.03,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Ashen Band of Endless Courage",
		ID:         50404,
		Bonus:      stats.Stats{stats.Armor: 2400},
		Duration:   time.Second * 10,
		Callback:   OnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.03,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:     "Ashen Band of Unmatched Might",
		ID:       52571,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: OnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:     "Ashen Band of Endless Might",
		ID:       52572,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: OnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sharpened Twilight Scale",
		ID:         54569,
		Bonus:      stats.Stats{stats.AttackPower: 1304, stats.RangedAttackPower: 1304},
		Duration:   time.Second * 15,
		Callback:   OnSpellHitDealt,
		Harmful:    true, // doesn't matter what, just that 'when you deal damage'
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sharpened Twilight Scale H",
		ID:         54590,
		Bonus:      stats.Stats{stats.AttackPower: 1472, stats.RangedAttackPower: 1472},
		Duration:   time.Second * 15,
		Callback:   OnSpellHitDealt,
		Harmful:    true, // doesn't matter what, just that 'when you deal damage'
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Charred Twilight Scale",
		ID:         54572,
		Bonus:      stats.Stats{stats.SpellPower: 763, stats.HealingPower: 763},
		Duration:   time.Second * 15,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Charred Twilight Scale H",
		ID:         54588,
		Bonus:      stats.Stats{stats.SpellPower: 861, stats.HealingPower: 861},
		Duration:   time.Second * 15,
		Callback:   OnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	// TODO: Provides a hot on each heal cast. Add later when we do healing sim
	// newProcStatBonusEffect(ProcStatBonusEffect{
	// 	Name:       "Glowing Twilight Scale",
	// 	ID:         54573,
	// 	Duration:   time.Second * 15,
	// 	Callback:   OnSpellHitDealt,
	// 	ProcMask:   core.ProcMaskHealingDone,
	// 	ProcChance: 0.10,
	// 	ICD:        time.Second * 50,
	// })
	// TODO: adds dodge when hit taken which drops health under 35%
	// newProcStatBonusEffect(ProcStatBonusEffect{
	// 	Name:       "Petrified Twilight Scale",
	// 	ID:         54571,
	// 	Bonus:      stats.Stats{stats.Dodge: 733},
	// 	Duration:   time.Second * 10,
	// 	Callback:   OnSpellHitTaken,
	// 	// TODO: when health < 35%
	// 	ICD:        time.Second * 45,
	// })

	core.AddEffectsToTest = true
}
