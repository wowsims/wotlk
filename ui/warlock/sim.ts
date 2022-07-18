import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { TypedEvent } from '/wotlk/core/typed_event.js';

import { Warlock, Warlock_Rotation as WarlockRotation, WarlockTalents as WarlockTalents, Warlock_Options as WarlockOptions, Warlock_Options_Armor, Warlock_Options_Summon } from '/wotlk/core/proto/warlock.js';

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
				gear: Presets.P5_DESTRO.gear,

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
				rotation: Presets.AfflictionRotation,
				// Default talents.
				talents: Presets.AfflictionTalents.data,
				// Default spec-specific settings.
				specOptions: WarlockOptions.create({
					armor: Warlock_Options_Armor.FelArmor,
					summon: Warlock_Options_Summon.Felhunter,
				}),
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					divineSpirit: true,
					bloodlust: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					totemOfWrath: true,
					wrathOfAirTotem: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,

				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					misery: true,
					curseOfElements: true,
					faerieFire: TristateEffect.TristateEffectImproved,
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
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: WarlockInputs.WarlockRotationConfig,

			spellInputs: [
				WarlockInputs.PrimarySpellShadowbolt,
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
				// Preset rotations that the user can quickly select.
				// rotation: [
				// 	Presets.AfflictionRotation,
				// 	Presets.DemonologyRotation,
				// 	Presets.DestructionRotation,
				// ],
				//Preset gear configurations that the user can quickly select.
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
