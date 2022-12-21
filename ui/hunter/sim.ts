import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Class } from '../core/proto/common.js';
import { Consumes } from '../core/proto/common.js';
import { Encounter } from '../core/proto/common.js';
import { ItemSlot } from '../core/proto/common.js';
import { MobType } from '../core/proto/common.js';
import { Race } from '../core/proto/common.js';
import { RangedWeaponType } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getTalentPoints } from '../core/proto_utils/utils.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { getPetTalentsConfig } from '../core/talents/hunter_pet.js';
import { protoToTalentString } from '../core/talents/factory.js';

import {
	Hunter,
	Hunter_Rotation as HunterRotation,
	Hunter_Options as HunterOptions,
	Hunter_Options_PetType as PetType,
	HunterPetTalents,
} from '../core/proto/hunter.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as HunterInputs from './inputs.js';
import * as Presets from './presets.js';

export class HunterSimUI extends IndividualSimUI<Spec.SpecHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parentElem, player, {
			cssClass: 'hunter-sim-ui',
			cssScheme: 'hunter',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
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
				// Warning when too many Pet talent points are used without BM talented.
				(simUI: IndividualSimUI<Spec.SpecHunter>) => {
					return {
						updateOn: TypedEvent.onAny([simUI.player.talentsChangeEmitter, simUI.player.specOptionsChangeEmitter]),
						getContent: () => {
							const specOptions = simUI.player.getSpecOptions();
							const petTalents = specOptions.petTalents || HunterPetTalents.create();
							const petTalentString = protoToTalentString(petTalents, getPetTalentsConfig(specOptions.petType));
							const talentPoints = getTalentPoints(petTalentString);

							const isBM = simUI.player.getTalents().beastMastery;
							const maxPoints = isBM ? 20 : 16;

							if (talentPoints == 0) {
								// Just return here, so we don't show a warning during page load.
								return '';
							} else if (talentPoints < maxPoints) {
								return 'Unspent pet talent points.';
							} else if (talentPoints > maxPoints) {
								return 'More than 16 points spent in pet talents, but Beast Mastery is not talented.';
							} else {
								return '';
							}
						},
					};
				},
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatAgility,
				Stat.StatRangedAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatMP5,
			],
			epPseudoStats: [
				PseudoStat.PseudoStatRangedDps,
			],
			// Reference stat against which to calculate EP.
			epReferenceStat: Stat.StatRangedAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatAgility,
				Stat.StatIntellect,
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

				const rangedWeapon = player.getEquippedItem(ItemSlot.ItemSlotRanged);
				if (rangedWeapon?.enchant?.effectId == 3608) {
					stats = stats.addStat(Stat.StatMeleeCrit, 40);
				}
				if (player.getRace() == Race.RaceDwarf && rangedWeapon?.item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeGun) {
					stats = stats.addStat(Stat.StatMeleeCrit, 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
				}
				if (player.getRace() == Race.RaceTroll && rangedWeapon?.item.rangedWeaponType == RangedWeaponType.RangedWeaponTypeBow) {
					stats = stats.addStat(Stat.StatMeleeCrit, 1 * Mechanics.MELEE_CRIT_RATING_PER_CRIT_CHANCE);
				}

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.SV_P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatStamina]: 0.5,
					[Stat.StatAgility]: 2.65,
					[Stat.StatIntellect]: 1.1,
					[Stat.StatRangedAttackPower]: 1.0,
					[Stat.StatMeleeHit]: 2,
					[Stat.StatMeleeCrit]: 1.5,
					[Stat.StatMeleeHaste]: 1.39,
					[Stat.StatArmorPenetration]: 1.32,
				}, {
					[PseudoStat.PseudoStatRangedDps]: 6.32,
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
				HunterInputs.UseHuntersMark,
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
					HunterInputs.SniperTrainingUptime,
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
					Presets.MM_PRERAID_PRESET,
					Presets.MM_P1_PRESET,
					Presets.SV_PRERAID_PRESET,
					Presets.SV_P1_PRESET,
				],
			},
		});
	}
}
