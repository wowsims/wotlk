import { RaidBuffs,
	PartyBuffs,
	IndividualBuffs,
	Debuffs,
	Spec,
	Stat,
	TristateEffect,
} from '/wotlk/core/proto/common.js';

import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { TypedEvent } from '/wotlk/core/typed_event.js';

import {
	Warlock,
	Warlock_Rotation as WarlockRotation,
	WarlockTalents as WarlockTalents,
	Warlock_Options as WarlockOptions,
	Warlock_Options_Armor as Armor,
	Warlock_Options_Summon as Summon,
	Warlock_Options_WeaponImbue as WeaponImbue,
} from '/wotlk/core/proto/warlock.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

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
				Stat.StatStamina,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatShadowSpellPower,
				Stat.StatFireSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
				//Pet stats
				Stat.StatStrength,
				Stat.StatAttackPower,
				Stat.StatAgility,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
			],
			// Reference stat against which to calculate EP. DPS classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
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
				gear: Presets.SWP_BIS.gear,

				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.15,
					[Stat.StatSpirit]: 0.2,
					[Stat.StatSpellPower]: 1,
					[Stat.StatShadowSpellPower]: 1,
					[Stat.StatFireSpellPower]: 0,
					[Stat.StatSpellHit]: 0.6,
					[Stat.StatSpellCrit]: 0.4,
					[Stat.StatSpellHaste]: 0.6,
					[Stat.StatMP5]: 0.00,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.AfflictionRotation,
				// Default talents.
				talents: Presets.AfflictionTalents.data,
				// Default spec-specific settings.
				specOptions: WarlockOptions.create({
					armor: Armor.FelArmor,
					summon: Summon.Felhunter,
					weaponImbue: WeaponImbue.GrandSpellstone,
				}),
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					powerWordFortitude: TristateEffect.TristateEffectImproved,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					arcaneBrilliance: true,
					divineSpirit: true,
					trueshotAura: true,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
					icyTalons: true,
					totemOfWrath: true,
					moonkinAura: TristateEffect.TristateEffectImproved,
					wrathOfAirTotem: true,
					swiftRetribution: true,
					sanctifiedRetribution: true,
					bloodlust: true,
				}),

				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					ebonPlaguebringer: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					judgementOfWisdom: true,
					misery: true,
					heartOfTheCrusader: true,
					sunderArmor: true,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				WarlockInputs.FelArmor,
				WarlockInputs.DemonArmor,
			],
			petInputs: [
				WarlockInputs.SummonImp,
				WarlockInputs.SummonSuccubus,
				WarlockInputs.SummonFelhunter,
				WarlockInputs.SummonFelguard,
			],
			weaponImbueInputs: [
				WarlockInputs.GrandSpellstone,
				WarlockInputs.GrandFirestone,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: WarlockInputs.WarlockRotationConfig,

			spellInputs: [
				WarlockInputs.PrimarySpellShadowBolt,
				WarlockInputs.PrimarySpellIncinerate,
				WarlockInputs.PrimarySpellSeed,
				WarlockInputs.SecondaryDotImmolate,
				WarlockInputs.SecondaryDotUnstableAffliction,
				WarlockInputs.SpecSpellChaosBolt,
				WarlockInputs.SpecSpellHaunt,
				WarlockInputs.CorruptionSpell,
			],
			
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.PrepopPotion,
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

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.AfflictionTalents,
					Presets.DemonologyTalents,
					Presets.DestructionTalents,
				],
				//Preset gear configurations that the user can quickly select.
				gear: [
					Presets.SWP_BIS,
				],
			},
		});
	}
}
