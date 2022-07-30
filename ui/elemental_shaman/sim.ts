import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Player } from '/wotlk/core/player.js';
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';
import { TotemsSection } from '/wotlk/core/components/totem_inputs.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Mechanics from '/wotlk/core/constants/mechanics.js';

import * as ShamanInputs from './inputs.js';
import * as Presets from './presets.js';
import { shamanGlyphsConfig } from '../core/talents/shaman.js';

export class ElementalShamanSimUI extends IndividualSimUI<Spec.SpecElementalShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecElementalShaman>) {
		super(parentElem, player, {
			cssClass: 'elemental-shaman-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],
			warnings: [
				// Warning to use all 4 totems if T6 2pc bonus is active.
				(simUI: IndividualSimUI<Spec.SpecElementalShaman>) => {
					return {
						updateOn: TypedEvent.onAny([simUI.player.rotationChangeEmitter, simUI.player.currentStatsEmitter]),
						getContent: () => {
							const hasT62P = simUI.player.getCurrentStats().sets.includes('Skyshatter Regalia (2pc)');
							const totems = simUI.player.getRotation().totems!;
							const hasAll4Totems = totems && totems.earth && totems.air && totems.fire && totems.water;
							if (hasT62P && !hasAll4Totems) {
								return 'T6 2pc bonus is equipped, but inactive because not all 4 totem types are being used.';
							} else {
								return '';
							}
						},
					};
				},
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpellPower,
				Stat.StatNatureSpellPower,
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
				Stat.StatMana,
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatSpellPower,
				Stat.StatNatureSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			modifyDisplayStats: (player: Player<Spec.SpecElementalShaman>) => {
				let stats = new Stats();
				stats = stats.addStat(Stat.StatSpellHit, player.getTalents().elementalPrecision * 2 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);
				stats = stats.addStat(Stat.StatSpellCrit,
					player.getTalents().tidalMastery * 1 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE +
					player.getTalents().thunderingStrikes * 1 * Mechanics.SPELL_CRIT_RATING_PER_CRIT_CHANCE);

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.17,
					[Stat.StatSpellPower]: 1,
					[Stat.StatNatureSpellPower]: 1,
					[Stat.StatSpellCrit]: 1,
					[Stat.StatSpellHaste]: 1,
					[Stat.StatMP5]: 0.08,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.StandardTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					divineSpirit: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					moonkinAura: TristateEffect.TristateEffectImproved,
					sanctifiedRetribution: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: 2,
					vampiricTouch: true,
				}),
				debuffs: Debuffs.create({
					faerieFire: TristateEffect.TristateEffectImproved,
					judgementOfWisdom: true,
					misery: true,
					curseOfElements: true,
					shadowMastery: true,
				}),
			},
			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				ShamanInputs.ShamanShieldInput,
				ShamanInputs.Bloodlust,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ShamanInputs.ElementalShamanRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.PrepopPotion,
					OtherInputs.TankAssignment,
				],
			},
			customSections: [
				TotemsSection,
			],
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.StandardTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.SWP_PRESET,
					Presets.PRE_RAID_PRESET,
					Presets.P1_PRESET,
				],
			},
		});
	}
}
