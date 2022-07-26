import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Encounter } from '/wotlk/core/proto/common.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';
import { MobType } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Player } from '/wotlk/core/player.js';
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';

import { Deathknight, Deathknight_Rotation as DeathKnightRotation, DeathknightTalents as DeathKnightTalents, Deathknight_Options as DeathKnightOptions } from '/wotlk/core/proto/deathknight.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as DeathKnightInputs from './inputs.js';
import * as Presets from './presets.js';

export class DeathknightSimUI extends IndividualSimUI<Spec.SpecDeathknight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecDeathknight>) {
		super(parentElem, player, {
			cssClass: 'deathknight-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				"<p>Rotation logic is just hit things on CGD. It is not good don't take it as actual data.</p>\
				<p>Damage multipliers are also likely to not be properly stacking until further beta testing.</p>\
				<p>Pet scaling is likely to not be properly working until further beta testing.</p>"
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStrength,
				Stat.StatArmor,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,

				// TODO: Remove these when debuff categories support us
				Stat.StatSpellPower,
				Stat.StatSpellHit,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.P1_FROST_BIS_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatStrength]: 2.61,
					[Stat.StatAgility]: 1.14,
					[Stat.StatArmor]: 0.027,
					[Stat.StatAttackPower]: 1,
					[Stat.StatExpertise]: 1.73,
					[Stat.StatMeleeHaste]: 1.26,
					[Stat.StatMeleeHit]: 1.71,
					[Stat.StatMeleeCrit]: 1.83,
					[Stat.StatArmorPenetration]: 1.425,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.FrostTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					swiftRetribution: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					icyTalons: true,
					abominationsMight: true,
					leaderOfThePack: TristateEffect.TristateEffectRegular,
					sanctifiedRetribution: true,
					bloodlust: true,
					devotionAura: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
					heroicPresence: true,
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					bloodFrenzy: true,
					faerieFire: TristateEffect.TristateEffectRegular,
					sunderArmor: true,
					misery: true,
					ebonPlaguebringer: true,
					mangle: true,
					heartOfTheCrusader: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DeathKnightInputs.DeathKnightRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					DeathKnightInputs.StartingRunicPower,
					DeathKnightInputs.PetUptime,
					
					DeathKnightInputs.PrecastGhoulFrenzy,
					DeathKnightInputs.PrecastHornOfWinter,
					
					OtherInputs.PrepopPotion,

					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.FrostTalents,
					Presets.FrostUnholyTalents,
					Presets.UnholyDualWieldTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_FROST_PRE_BIS_PRESET,
					Presets.P1_FROST_BIS_PRESET,
					Presets.P1_UNHOLY_DW_BIS_PRESET,
				],
			},
		});
	}
}
