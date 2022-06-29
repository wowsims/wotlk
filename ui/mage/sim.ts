import { IndividualBuffs } from '/tbc/core/proto/common.js';
import { Class } from '/tbc/core/proto/common.js';
import { Consumes } from '/tbc/core/proto/common.js';
import { Debuffs } from '/tbc/core/proto/common.js';
import { Encounter } from '/tbc/core/proto/common.js';
import { ItemSlot } from '/tbc/core/proto/common.js';
import { MobType } from '/tbc/core/proto/common.js';
import { PartyBuffs } from '/tbc/core/proto/common.js';
import { RaidBuffs } from '/tbc/core/proto/common.js';
import { RaidTarget } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { TristateEffect } from '/tbc/core/proto/common.js'
import { Stats } from '/tbc/core/proto_utils/stats.js';
import { Player } from '/tbc/core/player.js';
import { Sim } from '/tbc/core/sim.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';

import { Alchohol } from '/tbc/core/proto/common.js';
import { BattleElixir } from '/tbc/core/proto/common.js';
import { Flask } from '/tbc/core/proto/common.js';
import { Food } from '/tbc/core/proto/common.js';
import { GuardianElixir } from '/tbc/core/proto/common.js';
import { Conjured } from '/tbc/core/proto/common.js';
import { Drums } from '/tbc/core/proto/common.js';
import { PetFood } from '/tbc/core/proto/common.js';
import { Potions } from '/tbc/core/proto/common.js';
import { WeaponImbue } from '/tbc/core/proto/common.js';

import { Mage, Mage_Rotation as MageRotation, MageTalents as MageTalents, Mage_Options as MageOptions } from '/tbc/core/proto/mage.js';

import * as IconInputs from '/tbc/core/components/icon_inputs.js';
import * as OtherInputs from '/tbc/core/components/other_inputs.js';
import * as Tooltips from '/tbc/core/constants/tooltips.js';

import * as MageInputs from './inputs.js';
import * as Presets from './presets.js';

export class MageSimUI extends IndividualSimUI<Spec.SpecMage> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecMage>) {
		super(parentElem, player, {
			cssClass: 'mage-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatArcaneSpellPower,
				Stat.StatFireSpellPower,
				Stat.StatFrostSpellPower,
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
				Stat.StatArcaneSpellPower,
				Stat.StatFireSpellPower,
				Stat.StatFrostSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_ARCANE_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 1.29,
					[Stat.StatSpirit]: 0.89,
					[Stat.StatSpellPower]: 1,
					[Stat.StatArcaneSpellPower]: 0.78,
					[Stat.StatFireSpellPower]: 0,
					[Stat.StatFrostSpellPower]: 0.21,
					[Stat.StatSpellCrit]: 0.77,
					[Stat.StatSpellHaste]: 0.84,
					[Stat.StatMP5]: 0.61,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultArcaneConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultArcaneRotation,
				// Default talents.
				talents: Presets.ArcaneTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultArcaneOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					drums: Drums.DrumsOfBattle,
					bloodlust: 1,
					manaSpringTotem: TristateEffect.TristateEffectImproved,
					manaTideTotems: 1,
					wrathOfAirTotem: TristateEffect.TristateEffectRegular,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					blessingOfSalvation: true,
					innervates: 1,
				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					misery: true,
					curseOfElements: TristateEffect.TristateEffectRegular,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				MageInputs.MageArmor,
				MageInputs.MoltenArmor,
			],
			// IconInputs to include in the 'Other Buffs' section on the settings tab.
			raidBuffInputs: [
				IconInputs.DivineSpirit,
				IconInputs.GiftOfTheWild,
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
				IconInputs.CurseOfElements,
				IconInputs.Misery,
				IconInputs.ImprovedScorch,
				IconInputs.WintersChill,
			],
			// Which options are selectable in the 'Consumes' section.
			consumeOptions: {
				potions: [
					Potions.SuperManaPotion,
					Potions.DestructionPotion,
				],
				conjured: [
					Conjured.ConjuredMageManaEmerald,
					Conjured.ConjuredDarkRune,
					Conjured.ConjuredFlameCap,
				],
				flasks: [
					Flask.FlaskOfBlindingLight,
					Flask.FlaskOfPureDeath,
					Flask.FlaskOfSupremePower,
				],
				battleElixirs: [
					BattleElixir.AdeptsElixir,
					BattleElixir.ElixirOfMajorFirePower,
					BattleElixir.ElixirOfMajorFrostPower,
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
			rotationInputs: MageInputs.MageRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					MageInputs.EvocationTicks,
					OtherInputs.ShadowPriestDPS,
					OtherInputs.StartingPotion,
					OtherInputs.NumStartingPotions,
					OtherInputs.StartingConjured,
					OtherInputs.NumStartingConjured,
					OtherInputs.SnapshotImprovedWrathOfAirTotem,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Target stats to show for 'Simple' encounters.
				simpleTargetStats: [
					Stat.StatArcaneResistance,
					Stat.StatFireResistance,
					Stat.StatFrostResistance,
				],
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: true,
			},

			// If true, the talents on the talents tab will not be individually modifiable by the user.
			// Note that the use can still pick between preset talents, if there is more than 1.
			freezeTalents: false,

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.ArcaneTalents,
					Presets.FireTalents,
					Presets.FrostTalents,
					Presets.DeepFrostTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_ARCANE_PRESET,
					Presets.P2_ARCANE_PRESET,
					Presets.P3_ARCANE_PRESET,
					Presets.P4_ARCANE_PRESET,
					Presets.P5_ARCANE_PRESET,
					Presets.P1_FIRE_PRESET,
					Presets.P2_FIRE_PRESET,
					Presets.P3_FIRE_PRESET,
					Presets.P4_FIRE_PRESET,
					Presets.P5_FIRE_PRESET,
					Presets.P1_FROST_PRESET,
					Presets.P2_FROST_PRESET,
					Presets.P3_FROST_PRESET,
					Presets.P4_FROST_PRESET,
					Presets.P5_FROST_PRESET,
				],
			},
		});
	}
}
