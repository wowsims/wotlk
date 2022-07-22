import { RaidBuffs,
	PartyBuffs,
	IndividualBuffs,
	Debuffs,
	Spec,
	Stat,
	TristateEffect,
	Race,
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
				"<p>This sim is under current development for Wrath of the Lich King.</p>\
				<p>Talents and Glyphs are mostly ready but rotations are under development.</p>"
			],
			// race: Race.RaceOrc,

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
				//Pet stats for buffs in UI only
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
				specOptions: Presets.DestructionOptions,

				// Default buffs and debuffs settings.
				raidBuffs: Presets.DefaultRaidBuffs,

				partyBuffs: PartyBuffs.create({}),

				individualBuffs: Presets.DefaultIndividualBuffs,

				debuffs: Presets.DefaultDebuffs,
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				WarlockInputs.PetType,
				WarlockInputs.ArmorInput,
				WarlockInputs.WeaponImbue,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationIconInputs: [
				WarlockInputs.PrimarySpellInput,
				WarlockInputs.SecondaryDotInput,
				WarlockInputs.SpecSpellInput,
				WarlockInputs.CorruptionSpell,
			],
			rotationInputs: WarlockInputs.WarlockRotationConfig,
			
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.PrepopPotion,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
				// executeProportion20: 0.25,
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
					Presets.Naked,
					Presets.SWP_BIS,
					Presets.P1_PreBiS,
					Presets.P1_BiS,
				],
			},
		});
	}
}
