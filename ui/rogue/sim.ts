import { Race, RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { Encounter } from '../core/proto/common.js';
import { ItemSlot } from '../core/proto/common.js';
import { MobType } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import { Rogue, Rogue_Rotation as RogueRotation, Rogue_Options as RogueOptions } from '../core/proto/rogue.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as RogueInputs from './inputs.js';
import * as Presets from './presets.js';

export class RogueSimUI extends IndividualSimUI<Spec.SpecRogue> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRogue>) {
		super(parentElem, player, {
			cssClass: 'rogue-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				'Rotations are not fully optimized, especially for non-standard setups.',
			],
			warnings: [
				(simUI: IndividualSimUI<Spec.SpecRogue>) => {
					return {
						updateOn: simUI.player.changeEmitter,
						getContent: () => {
							if (simUI.player.getRotation().maintainExposeArmor && simUI.player.getTalents().improvedExposeArmor < 2) {
								return '\'Maintain Expose Armor\' selected, but missing points in Improved Expose Armor!';
							} else {
								return '';
							}
						},
					};
				},
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatExpertise,
			],
			// Reference stat against which to calculate EP.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatExpertise,
			],

			defaults: {
				// Default equipped gear.
				gear: Presets.PRERAID_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatAgility]: 1.89,
					[Stat.StatStrength]: 1.17,
					[Stat.StatAttackPower]: 1,
					[Stat.StatMeleeHit]: 1.65,
					[Stat.StatMeleeCrit]: 1.11,
					[Stat.StatMeleeHaste]: 1.27,
					[Stat.StatArmorPenetration]: 0.3,
					[Stat.StatExpertise]: 1.69,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.CombatTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					icyTalons: true,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
					abominationsMight: true,
					swiftRetribution: true,
					elementalOath: true,
					sanctifiedRetribution: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					heartOfTheCrusader: true,
					mangle: true,
					sunderArmor: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					shadowMastery: true,
					earthAndMoon: true,
					bloodFrenzy: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				RogueInputs.MainHandImbue,
				RogueInputs.OffHandImbue,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: RogueInputs.RogueRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.SpellCritBuff,
				IconInputs.SpellCritDebuff,
				IconInputs.SpellHitDebuff,
				IconInputs.SpellDamageDebuff
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.StartingConjured,
					OtherInputs.PrepopPotion,
					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
				],
			},
			additionalIconSections: {
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.CombatTalents,
					Presets.AssassinationTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.PRERAID_PRESET,
					Presets.P1_PRESET,
				],
			},
		});
	}
}
