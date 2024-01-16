package wotlk

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

type ProcStatBonusEffect struct {
	Name       string
	ID         int32
	AuraID     int32
	Bonus      stats.Stats
	Duration   time.Duration
	Callback   core.AuraCallback
	ProcMask   core.ProcMask
	Outcome    core.HitOutcome
	Harmful    bool
	ProcChance float64
	PPM        float64
	ICD        time.Duration

	// For ignoring a hardcoded spell.
	IgnoreSpellID int32
}

func newProcStatBonusEffect(config ProcStatBonusEffect) {
	core.NewItemEffect(config.ID, func(agent core.Agent) {
		character := agent.GetCharacter()

		procID := core.ActionID{SpellID: config.AuraID}
		if procID.IsEmptyAction() {
			procID = core.ActionID{ItemID: config.ID}
		}
		procAura := character.NewTemporaryStatsAura(config.Name+" Proc", procID, config.Bonus, config.Duration)

		handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
			procAura.Activate(sim)
		}
		if config.IgnoreSpellID != 0 {
			ignoreSpellID := config.IgnoreSpellID
			handler = func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				if !spell.IsSpellAction(ignoreSpellID) {
					procAura.Activate(sim)
				}
			}
		}

		triggerAura := core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
			ActionID:   core.ActionID{ItemID: config.ID},
			Name:       config.Name,
			Callback:   config.Callback,
			ProcMask:   config.ProcMask,
			Outcome:    config.Outcome,
			Harmful:    config.Harmful,
			ProcChance: config.ProcChance,
			PPM:        config.PPM,
			ICD:        config.ICD,
			Handler:    handler,
		})
		procAura.Icd = triggerAura.Icd
	})
}

func init() {
	// Keep these separated by stat, ordered by item ID within each group.
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Meteorite Whetstone",
		ID:         37390,
		AuraID:     60302,
		Bonus:      stats.Stats{stats.MeleeHaste: 444, stats.SpellHaste: 444},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	//newProcStatBonusEffect(ProcStatBonusEffect{
	//	Name:       "Serrah's Star",
	//	ID:         37559,
	//	Bonus:      stats.Stats{stats.MeleeCrit: 167, stats.SpellCrit: 167},
	//	Duration:   time.Second * 10,
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask:   core.ProcMaskSpellDamage,
	//	Outcome:    core.OutcomeCrit,
	//	ProcChance: 0.45,
	//	ICD:        time.Second * 45,
	//})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Spark of Life",
		ID:         37657,
		AuraID:     60520,
		Bonus:      stats.Stats{stats.MP5: 220},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Forge Ember",
		ID:         37660,
		AuraID:     60479,
		Bonus:      stats.Stats{stats.SpellPower: 512},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})

	core.AddEffectsToTest = false

	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Je'Tze's Bell",
		ID:         37835,
		AuraID:     49623,
		Bonus:      stats.Stats{stats.MP5: 125},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	//newProcStatBonusEffect(ProcStatBonusEffect{
	//	Name:       "Valonforth's Remembrance",
	//	ID:         38071,
	//	Bonus:      stats.Stats{stats.Spirit: 222},
	//	Duration:   time.Second * 10,
	//	Callback: core.CallbackOnSpellHitDealt,
	//	ProcMask:   core.ProcMaskSpellDamage,
	//	Outcome:    core.OutcomeLanded,
	//	ProcChance: 0.15,
	//	ICD:        time.Second * 45,
	//})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Embrace of the Spider",
		ID:         39229,
		AuraID:     60492,
		Bonus:      stats.Stats{stats.MeleeHaste: 505, stats.SpellHaste: 505},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Dying Curse",
		ID:         40255,
		AuraID:     60494,
		Bonus:      stats.Stats{stats.SpellPower: 765},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpellOrProc,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Grim Toll",
		ID:         40256,
		AuraID:     60437,
		Bonus:      stats.Stats{stats.ArmorPenetration: 612},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sundial of the Exiled",
		ID:         40682,
		AuraID:     60064,
		Bonus:      stats.Stats{stats.SpellPower: 590},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpellOrProc,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Mirror of Truth",
		ID:         40684,
		AuraID:     60065,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "The Egg of Mortal Essence",
		ID:         40685,
		AuraID:     60062,
		Bonus:      stats.Stats{stats.MeleeHaste: 505, stats.SpellHaste: 505},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sonic Booster",
		ID:         40767,
		AuraID:     55018,
		Bonus:      stats.Stats{stats.AttackPower: 430, stats.RangedAttackPower: 430},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.35,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Tears of Bitter Anguish",
		ID:         43573,
		AuraID:     58904,
		Bonus:      stats.Stats{stats.MeleeHaste: 410, stats.SpellHaste: 410},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Crusader's Locket",
		ID:         43829,
		AuraID:     61671,
		Bonus:      stats.Stats{stats.Expertise: 258},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Chuchu's Tiny Box of Horrors",
		ID:         43838,
		AuraID:     61619,
		Bonus:      stats.Stats{stats.MeleeCrit: 258, stats.SpellCrit: 258},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Signet of Edward the Odd",
		ID:         44308,
		AuraID:     60318,
		Bonus:      stats.Stats{stats.MeleeHaste: 125, stats.SpellHaste: 125},
		Duration:   time.Second * 13,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Flow of Knowledge",
		ID:         44912,
		AuraID:     60064,
		Bonus:      stats.Stats{stats.SpellPower: 590},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Anvil of Titans",
		ID:         44914,
		AuraID:     60065,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Jouster's Fury Alliance",
		ID:         45131,
		AuraID:     63250,
		Bonus:      stats.Stats{stats.MeleeCrit: 328, stats.SpellCrit: 328},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Jouster's Fury Horde",
		ID:         45219,
		AuraID:     63250,
		Bonus:      stats.Stats{stats.MeleeCrit: 328, stats.SpellCrit: 328},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Pyrite Infuser",
		ID:         45286,
		AuraID:     65014,
		Bonus:      stats.Stats{stats.AttackPower: 1305, stats.RangedAttackPower: 1305},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.1,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Pandora's Plea",
		ID:         45490,
		AuraID:     64741,
		Bonus:      stats.Stats{stats.SpellPower: 794},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Blood of the Old God",
		ID:         45522,
		AuraID:     64790,
		Bonus:      stats.Stats{stats.AttackPower: 1358, stats.RangedAttackPower: 1358},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.1, // wowhead shows proc chance: 10% but a comment says 1.12PPM? TODO: validate.
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Flare of the Heavens",
		ID:         45518,
		AuraID:     64713,
		Bonus:      stats.Stats{stats.SpellPower: 959},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnCastComplete,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.1,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Show of Faith",
		ID:         45535,
		AuraID:     64739,
		Bonus:      stats.Stats{stats.MP5: 272},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Comet's Trail",
		ID:         45609,
		AuraID:     64772,
		Bonus:      stats.Stats{stats.SpellHaste: 819, stats.MeleeHaste: 819},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Elemental Focus Stone",
		ID:         45866,
		AuraID:     65004,
		Bonus:      stats.Stats{stats.MeleeHaste: 552, stats.SpellHaste: 552},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sif's Remembrance",
		ID:         45929,
		AuraID:     65003,
		Bonus:      stats.Stats{stats.MP5: 220},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Mjolnir Runestone",
		ID:         45931,
		AuraID:     65019,
		Bonus:      stats.Stats{stats.ArmorPenetration: 751},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Dark Matter",
		ID:         46038,
		AuraID:     65024,
		Bonus:      stats.Stats{stats.MeleeCrit: 692, stats.SpellCrit: 692},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.15,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Abyssal Rune",
		ID:         47213,
		AuraID:     67669,
		Bonus:      stats.Stats{stats.SpellPower: 590},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellOrProc,
		Harmful:    true,
		ProcChance: 0.25,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Banner of Victory",
		ID:         47214,
		AuraID:     67671,
		Bonus:      stats.Stats{stats.AttackPower: 1008, stats.RangedAttackPower: 1008},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.20,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "The Black Heart",
		ID:         47216,
		AuraID:     67631,
		Bonus:      stats.Stats{stats.Armor: 7056},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.25,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Coren's Chromium Coaster",
		ID:         49074,
		AuraID:     60065,
		Bonus:      stats.Stats{stats.AttackPower: 1000, stats.RangedAttackPower: 1000},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Mithril Pocketwatch",
		ID:         49076,
		AuraID:     60064,
		Bonus:      stats.Stats{stats.SpellPower: 590},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellOrProc,
		Harmful:    true,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Ancient Pickled Egg",
		ID:         49078,
		AuraID:     60062,
		Bonus:      stats.Stats{stats.MeleeHaste: 505, stats.SpellHaste: 505},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		ProcChance: 0.10,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Needle-Encrusted Scorpion",
		ID:         50198,
		AuraID:     71403,
		Bonus:      stats.Stats{stats.ArmorPenetration: 678},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Outcome:    core.OutcomeCrit,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Whispering Fanged Skull",
		ID:         50342,
		AuraID:     71401,
		Bonus:      stats.Stats{stats.AttackPower: 1100, stats.RangedAttackPower: 1100},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Harmful:    true,
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Whispering Fanged Skull H",
		ID:         50343,
		AuraID:     71541,
		Bonus:      stats.Stats{stats.AttackPower: 1250, stats.RangedAttackPower: 1250},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Harmful:    true,
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Purified Lunar Dust",
		ID:         50358,
		AuraID:     71584,
		Bonus:      stats.Stats{stats.MP5: 304},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnCastComplete,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Phylactery of the Nameless Lich",
		ID:         50360,
		AuraID:     71605,
		Bonus:      stats.Stats{stats.SpellPower: 1074},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.30,
		ICD:        time.Second * 100,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Phylactery of the Nameless Lich H",
		ID:         50365,
		AuraID:     71636,
		Bonus:      stats.Stats{stats.SpellPower: 1207},
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		ProcChance: 0.30,
		ICD:        time.Second * 100,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Unmatched Destruction
		Name:       "Frostforged Sage",
		ID:         50397,
		AuraID:     72416,
		Bonus:      stats.Stats{stats.SpellPower: 285},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Endless Destruction
		Name:       "Frostforged Sage",
		ID:         50398,
		AuraID:     72416,
		Bonus:      stats.Stats{stats.SpellPower: 285},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskSpellDamage,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.10,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Unmatched Vengeance
		Name:     "Frostforged Champion",
		ID:       50401,
		AuraID:   72412,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Endless Vengeance
		Name:     "Frostforged Champion",
		ID:       50402,
		AuraID:   72412,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Unmatched Courage
		Name:       "Frostforged Defender",
		ID:         50403,
		AuraID:     72414,
		Bonus:      stats.Stats{stats.Armor: 2400},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.03,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Endless Courage
		Name:       "Frostforged Defender",
		ID:         50404,
		AuraID:     72414,
		Bonus:      stats.Stats{stats.Armor: 2400},
		Duration:   time.Second * 10,
		Callback:   core.CallbackOnSpellHitTaken,
		ProcMask:   core.ProcMaskMelee,
		Outcome:    core.OutcomeLanded,
		ProcChance: 0.03,
		ICD:        time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Unmatched Might
		Name:     "Frostforged Champion",
		ID:       52571,
		AuraID:   72412,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		// Ashen Band of Endless Might
		Name:     "Frostforged Champion",
		ID:       52572,
		AuraID:   72412,
		Bonus:    stats.Stats{stats.AttackPower: 480, stats.RangedAttackPower: 480},
		Duration: time.Second * 10,
		Callback: core.CallbackOnSpellHitDealt,
		ProcMask: core.ProcMaskMeleeOrRanged,
		Outcome:  core.OutcomeLanded,
		PPM:      1,
		ICD:      time.Second * 60,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sharpened Twilight Scale",
		ID:         54569,
		AuraID:     75458,
		Bonus:      stats.Stats{stats.AttackPower: 1304, stats.RangedAttackPower: 1304},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Harmful:    true, // doesn't matter what, just that 'when you deal damage'
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Sharpened Twilight Scale H",
		ID:         54590,
		AuraID:     75456,
		Bonus:      stats.Stats{stats.AttackPower: 1472, stats.RangedAttackPower: 1472},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		ProcMask:   core.ProcMaskDirect | core.ProcMaskProc,
		Harmful:    true, // doesn't matter what, just that 'when you deal damage'
		ProcChance: 0.35,
		ICD:        time.Second * 45,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Charred Twilight Scale",
		ID:         54572,
		AuraID:     75466,
		Bonus:      stats.Stats{stats.SpellPower: 763},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellOrProc,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})
	newProcStatBonusEffect(ProcStatBonusEffect{
		Name:       "Charred Twilight Scale H",
		ID:         54588,
		AuraID:     75473,
		Bonus:      stats.Stats{stats.SpellPower: 861},
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt | core.CallbackOnPeriodicDamageDealt,
		ProcMask:   core.ProcMaskSpellOrProc,
		ProcChance: 0.10,
		ICD:        time.Second * 50,
	})

	core.AddEffectsToTest = true
}
