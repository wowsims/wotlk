import { RaidBuffs } from '/tbc/core/proto/common.js';
import { PartyBuffs } from '/tbc/core/proto/common.js';
import { IndividualBuffs } from '/tbc/core/proto/common.js';
import { Debuffs } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js'
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { Player } from '/tbc/core/player.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { TypedEvent } from '/tbc/core/typed_event.js';

import { Alchohol } from '/tbc/core/proto/common.js';
import { BattleElixir } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { GuardianElixir } from '/tbc/core/proto/common.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { Drums } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';

import { Warlock, Warlock_Rotation as WarlockRotation, WarlockTalents as WarlockTalents, Warlock_Options as WarlockOptions, Warlock_Options_Armor, Warlock_Options_Summon } from '/tbc/core/proto/warlock.js';

import * as IconInputs from '/tbc/core/components/icon_inputs.js';
import * as OtherInputs from '/tbc/core/components/other_inputs.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

import * as WarlockInputs from './inputs.js';
import * as Presets from './presets.js';

export class WarlockSimUI extends IndividualSimUI<Spec.SpecWarlock> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecWarlock>) {
		super(parentElem, player, {
			cssClass: 'warlock-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [

			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpellPower,
				Stat.StatShadowSpellPower,
				Stat.StatFireSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatShadowSpellPower,
				Stat.StatFireSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P4_DESTRO.gear,

				// TODO: FIND EPS FOR WARLOCKS
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.4,
					[Stat.StatSpirit]: 0.1,
					[Stat.StatSpellPower]: 1,
					[Stat.StatShadowSpellPower]: 1,
					[Stat.StatFireSpellPower]: 1,
					[Stat.StatSpellCrit]: 0.8,
					[Stat.StatSpellHaste]: 1.2,
					[Stat.StatMP5]: 0.00,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.DestructionTalents.data,
				// Default spec-specific settings.
				specOptions: WarlockOptions.create({
					armor: Warlock_Options_Armor.FelArmor,
					sacrificeSummon: true,
					summon: Warlock_Options_Summon.Succubus,
				}),
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					divineSpirit: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					drums: Drums.DrumsOfBattle,
					bloodlust: 1,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					totemOfWrath: 1,
					wrathOfAirTotem: TristateEffect.TristateEffectRegular,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					blessingOfSalvation: true,
				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					misery: true,
					shadowWeaving: true,
					curseOfElements: TristateEffect.TristateEffectRegular,
					faerieFire: TristateEffect.TristateEffectImproved,
					sunderArmor: true,
					isbUptime: 0.65,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				WarlockInputs.FelArmor,
				WarlockInputs.DemonArmor,
				WarlockInputs.DemonSummon,
				WarlockInputs.Sacrifice,
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.ArcaneBrilliance,
				IconInputs.DivineSpirit,
			],
			partyBuffInputs: [
				IconInputs.MoonkinAura,
				IconInputs.DrumsOfBattleBuff,
				IconInputs.DrumsOfRestorationBuff,
				IconInputs.Bloodlust,
				IconInputs.WrathOfAirTotem,
				IconInputs.TotemOfWrath,
				IconInputs.ManaSpringTotem,
				IconInputs.ManaTideTotem,
				IconInputs.SanctityAura,
				IconInputs.DraeneiRacialCaster,
				IconInputs.EyeOfTheNight,
				IconInputs.ChainOfTheTwilightOwl,
				IconInputs.JadePendantOfBlasting,
				IconInputs.AtieshWarlock,
				IconInputs.AtieshMage,
			],
			playerBuffInputs: [
				IconInputs.BlessingOfKings,
				IconInputs.BlessingOfWisdom,
				IconInputs.BlessingOfSalvation,
				IconInputs.Innervate,
				IconInputs.PowerInfusion,
			],
			// IconInputs to include in the 'Debuffs' section on the settings tab.
			debuffInputs: [
				IconInputs.JudgementOfWisdom,
				IconInputs.ImprovedSealOfTheCrusader,
				IconInputs.ShadowWeaving,
				IconInputs.Misery,
				IconInputs.ImprovedScorch,
				IconInputs.CurseOfElements,
				IconInputs.ExposeArmor,
				IconInputs.SunderArmor,
				IconInputs.BloodFrenzy,
				IconInputs.HuntersMark,
				IconInputs.FaerieFire,
				IconInputs.CurseOfRecklessness,
			],
			// Which options are selectable in the 'Consumes' section.
			consumeOptions: {
				potions: [
					Potions.SuperManaPotion,
					Potions.DestructionPotion,
				],
				conjured: [
					Conjured.ConjuredDarkRune,
					Conjured.ConjuredFlameCap,
				],
				flasks: [
					Flask.FlaskOfPureDeath,
					Flask.FlaskOfSupremePower,
				],
				battleElixirs: [
					BattleElixir.AdeptsElixir,
					BattleElixir.ElixirOfMajorShadowPower,
					BattleElixir.ElixirOfMajorFirePower,
				],
				guardianElixirs: [
					GuardianElixir.ElixirOfDraenicWisdom,
					GuardianElixir.ElixirOfMajorMageblood,
				],
				food: [
					Food.FoodBlackenedBasilisk,
					Food.FoodSkullfishSoup,
				],
				alcohol: [
					Alchohol.AlchoholKreegsStoutBeatdown,
				],
				weaponImbues: [
					WeaponImbue.WeaponImbueBrilliantWizardOil,
					WeaponImbue.WeaponImbueSuperiorWizardOil,
				],
				other: [
					IconInputs.ScrollOfSpiritV,
				],
			},
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: WarlockInputs.WarlockRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.ISBUptime,
					OtherInputs.ShadowPriestDPS,
					OtherInputs.StartingPotion,
					OtherInputs.NumStartingPotions,
					OtherInputs.SnapshotImprovedWrathOfAirTotem,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Target stats to show for 'Simple' encounters.
				simpleTargetStats: [
					Stat.StatShadowResistance,
					Stat.StatFireResistance,
					Stat.StatArmor,
				],
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			// If true, the talents on the talents tab will not be individually modifiable by the user.
			// Note that the use can still pick between preset talents, if there is more than 1.
			freezeTalents: false,

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.AfflicationTalents,
					Presets.DemonologistTalents,
					Presets.DestructionTalents,
					Presets.T6DestroTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_DESTRO,
					Presets.P2_DESTRO,
					Presets.P3_DESTRO,
					Presets.P4_DESTRO,
					Presets.P5_DESTRO,
				],
			},
		});
	}
}
