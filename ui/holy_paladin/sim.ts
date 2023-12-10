import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Player } from '../core/player.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Debuffs, IndividualBuffs, PartyBuffs, RaidBuffs, Spec, Stat, TristateEffect } from '../core/proto/common.js';
import { Stats } from '../core/proto_utils/stats.js';

import * as OtherInputs from '../core/components/other_inputs.js';

import * as HolyPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

export class HolyPaladinSimUI extends IndividualSimUI<Spec.SpecHolyPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHolyPaladin>) {
		super(parentElem, player, {
			cssClass: 'holy-paladin-sim-ui',
			cssScheme: 'paladin',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
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
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.DefaultGear.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.38,
					[Stat.StatSpirit]: 0.34,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellCrit]: 0.69,
					[Stat.StatSpellHaste]: 0.77,
					[Stat.StatMP5]: 0.00,
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
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					powerWordFortitude: TristateEffect.TristateEffectImproved,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					arcaneBrilliance: true,
					moonkinAura: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					thorns: TristateEffect.TristateEffectImproved,
					devotionAura: TristateEffect.TristateEffectImproved,
					shadowProtection: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfSanctuary: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					judgementOfLight: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					exposeArmor: true,
					sunderArmor: true,
					thunderClap: TristateEffect.TristateEffectImproved,
					insectSwarm: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: HolyPaladinInputs.HolyPaladinRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.InspirationUptime,
					HolyPaladinInputs.AuraSelection,
					HolyPaladinInputs.JudgementSelection,
				],
			},
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
					Presets.DefaultGear,
				],
			},

			autoRotation: (_player: Player<Spec.SpecHolyPaladin>): APLRotation => {
				return APLRotation.create();
			},
		});
	}
}
