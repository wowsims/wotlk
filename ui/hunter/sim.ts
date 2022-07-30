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
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import {
	Hunter,
	Hunter_Rotation as HunterRotation,
	Hunter_Options as HunterOptions,
	Hunter_Options_PetType as PetType,
} from '/wotlk/core/proto/hunter.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Mechanics from '/wotlk/core/constants/mechanics.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as HunterInputs from './inputs.js';
import * as Presets from './presets.js';

export class HunterSimUI extends IndividualSimUI<Spec.SpecHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parentElem, player, {
			cssClass: 'hunter-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				'Sim uses a simple priority-based rotation.',
				'Melee weaving and trap weaving are not included in the rotation.',
			],
			warnings: [
				// Warning when using exotic pet without BM talented.
				(simUI: IndividualSimUI<Spec.SpecHunter>) => {
					return {
						updateOn: TypedEvent.onAny([simUI.player.talentsChangeEmitter, simUI.player.specOptionsChangeEmitter]),
						getContent: () => {
							const petIsExotic = [
								PetType.Chimaera,
								PetType.CoreHound,
								PetType.Devilsaur,
								PetType.Silithid,
								PetType.SpiritBeast,
								PetType.Worm,
							].includes(simUI.player.getSpecOptions().petType);

							const isBM = simUI.player.getTalents().beastMastery;

							if (petIsExotic && !isBM) {
								return 'Cannot use exotic pets without the Beast Mastery talent.';
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
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatAttackPower,
				Stat.StatRangedAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatMP5,
			],
			// Reference stat against which to calculate EP.
			epReferenceStat: Stat.StatRangedAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatAgility,
				Stat.StatStrength,
				Stat.StatIntellect,
				Stat.StatAttackPower,
				Stat.StatRangedAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatMP5,
			],
			modifyDisplayStats: (player: Player<Spec.SpecHunter>) => {
				let stats = new Stats();
				stats = stats.addStat(Stat.StatMeleeCrit, player.getTalents().lethalShots * 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.PRERAID_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 0.7,
					[Stat.StatAgility]: 3.2,
					[Stat.StatStrength]: 0.05,
					[Stat.StatAttackPower]: 0.05,
					[Stat.StatRangedAttackPower]: 1.0,
					[Stat.StatMeleeHit]: 3,
					[Stat.StatMeleeCrit]: 1.2,
					[Stat.StatMeleeHaste]: 2.4,
					[Stat.StatArmorPenetration]: 0.4,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.SurvivalTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					powerWordFortitude: TristateEffect.TristateEffectImproved,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					windfuryTotem: TristateEffect.TristateEffectImproved,
					battleShout: TristateEffect.TristateEffectImproved,
					leaderOfThePack: TristateEffect.TristateEffectImproved,
					sanctifiedRetribution: true,
					unleashedRage: true,
					moonkinAura: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: 2,
					blessingOfMight: 2,
					vampiricTouch: true,
				}),
				debuffs: Debuffs.create({
					sunderArmor: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					judgementOfWisdom: true,
					curseOfElements: true,
					heartOfTheCrusader: true,
					savageCombat: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
				HunterInputs.PetTypeInput,
				HunterInputs.WeaponAmmo,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: HunterInputs.HunterRotationConfig,
			petConsumeInputs: [
				IconInputs.SpicedMammothTreats,
			],
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.StaminaBuff,
				IconInputs.SpellDamageDebuff,
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					HunterInputs.PetUptime,
					//HunterInputs.PetSingleAbility,
					HunterInputs.SniperTrainingUptime,
					//HunterInputs.LatencyMs,
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
					Presets.BeastMasteryTalents,
					Presets.MarksmanTalents,
					Presets.SurvivalTalents,
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
