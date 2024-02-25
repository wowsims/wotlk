import {
	Class,
	Cooldowns,
	Debuffs,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	Race,
	RaidBuffs,
	RangedWeaponType,
	Spec,
	Stat, PseudoStat,
	TristateEffect,
} from '../core/proto/common.js';
import {
	APLAction,
	APLListItem,
	APLRotation,
} from '../core/proto/apl.js';
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getTalentPoints, getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { TypedEvent } from '../core/typed_event.js';
import { getPetTalentsConfig } from '../core/talents/hunter_pet.js';
import { protoToTalentString } from '../core/talents/factory.js';
import { Gear } from '../core/proto_utils/gear.js';
import { PhysicalDPSGemOptimizer } from '../core/components/suggest_gems_action.js';

import {
	Hunter_Rotation as HunterRotation,
	Hunter_Rotation_StingType as StingType,
	Hunter_Options_PetType as PetType,
	HunterPetTalents,
	Hunter_Rotation_RotationType,
} from '../core/proto/hunter.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as ConsumablesInputs from '../core/components/inputs/consumables.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';
import * as AplUtils from '../core/proto_utils/apl_utils.js';

import * as HunterInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecHunter, {
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
		ConsumablesInputs.SpicedMammothTreats,
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.StaminaBuff,
		BuffDebuffInputs.SpellDamageDebuff,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			HunterInputs.PetUptime,
			HunterInputs.TimeToTrapWeaveMs,
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
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_PRESET_SIMPLE_DEFAULT,
			Presets.ROTATION_PRESET_BM,
			Presets.ROTATION_PRESET_MM,
			Presets.ROTATION_PRESET_MM_ADVANCED,
			Presets.ROTATION_PRESET_SV,
			Presets.ROTATION_PRESET_SV_ADVANCED,
			Presets.ROTATION_PRESET_AOE,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.MM_PRERAID_PRESET,
			Presets.MM_P1_PRESET,
			Presets.MM_P2_PRESET,
			Presets.MM_P3_PRESET,
			Presets.MM_P4_PRESET,
			Presets.MM_P5_PRESET,
			Presets.SV_PRERAID_PRESET,
			Presets.SV_P1_PRESET,
			Presets.SV_P2_PRESET,
			Presets.SV_P3_PRESET,
			Presets.SV_P4_PRESET,
			Presets.SV_P5_PRESET,
		],
	},

	autoRotation: (player: Player<Spec.SpecHunter>): APLRotation => {
		const talentTree = player.getTalentTree();
		const numTargets = player.sim.encounter.targets.length;
		if (numTargets >= 4) {
			return Presets.ROTATION_PRESET_AOE.rotation.rotation!;
		} else if (talentTree == 0) {
			return Presets.ROTATION_PRESET_BM.rotation.rotation!;
		} else if (talentTree == 1) {
			return Presets.ROTATION_PRESET_MM.rotation.rotation!;
		} else {
			return Presets.ROTATION_PRESET_SV.rotation.rotation!;
		}
	},

	simpleRotation: (player: Player<Spec.SpecHunter>, simple: HunterRotation, cooldowns: Cooldowns): APLRotation => {
		let [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const serpentSting = APLAction.fromJsonString(`{"condition":{"cmp":{"op":"OpGt","lhs":{"remainingTime":{}},"rhs":{"const":{"val":"6s"}}}},"multidot":{"spellId":{"spellId":49001},"maxDots":${simple.multiDotSerpentSting ? 3 : 1},"maxOverlap":{"const":{"val":"0ms"}}}}`);
		const scorpidSting = APLAction.fromJsonString(`{"condition":{"auraShouldRefresh":{"auraId":{"spellId":3043},"maxOverlap":{"const":{"val":"0ms"}}}},"castSpell":{"spellId":{"spellId":3043}}}`);
		const trapWeave = APLAction.fromJsonString(`{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":49067}}}}},"castSpell":{"spellId":{"tag":1,"spellId":49067}}}`);
		const volley = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":58434}}}`);
		const killShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":61006}}}`);
		const aimedShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49050}}}`);
		const multiShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49048}}}`);
		const steadyShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49052}}}`);
		const silencingShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":34490}}}`);
		const chimeraShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":53209}}}`);
		const blackArrow = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":63672}}}`);
		const explosiveShot4 = APLAction.fromJsonString(`{"condition":{"not":{"val":{"dotIsActive":{"spellId":{"spellId":60053}}}}},"castSpell":{"spellId":{"spellId":60053}}}`);
		const explosiveShot3 = APLAction.fromJsonString(`{"condition":{"dotIsActive":{"spellId":{"spellId":60053}}},"castSpell":{"spellId":{"spellId":60052}}}`);
		//const arcaneShot = APLAction.fromJsonString(`{"castSpell":{"spellId":{"spellId":49045}}}`);

		if (simple.viperStartManaPercent != 0) {
			actions.push(APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":34074}}}}},{"cmp":{"op":"OpLt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(simple.viperStartManaPercent * 100).toFixed(0)}%"}}}}]}},"castSpell":{"spellId":{"spellId":34074}}}`));
		}
		if (simple.viperStopManaPercent != 0) {
			actions.push(APLAction.fromJsonString(`{"condition":{"and":{"vals":[{"not":{"val":{"auraIsActive":{"auraId":{"spellId":61847}}}}},{"cmp":{"op":"OpGt","lhs":{"currentManaPercent":{}},"rhs":{"const":{"val":"${(simple.viperStopManaPercent * 100).toFixed(0)}%"}}}}]}},"castSpell":{"spellId":{"spellId":61847}}}`));
		}

		const talentTree = player.getTalentTree();
		if (simple.type == Hunter_Rotation_RotationType.Aoe) {
			actions.push(...[
				simple.sting == StingType.ScorpidSting ? scorpidSting : null,
				simple.sting == StingType.SerpentSting ? serpentSting : null,
				simple.trapWeave ? trapWeave : null,
				volley,
			].filter(a => a) as Array<APLAction>)
		} else if (talentTree == 0) { // BM
			actions.push(...[
				killShot,
				simple.trapWeave ? trapWeave : null,
				simple.sting == StingType.ScorpidSting ? scorpidSting : null,
				simple.sting == StingType.SerpentSting ? serpentSting : null,
				aimedShot,
				multiShot,
				steadyShot,
			].filter(a => a) as Array<APLAction>)
		} else if (talentTree == 1) { // MM
			actions.push(...[
				silencingShot,
				killShot,
				simple.sting == StingType.ScorpidSting ? scorpidSting : null,
				simple.sting == StingType.SerpentSting ? serpentSting : null,
				simple.trapWeave ? trapWeave : null,
				chimeraShot,
				aimedShot,
				multiShot,
				steadyShot,
			].filter(a => a) as Array<APLAction>)
		} else if (talentTree == 2) { // SV
			actions.push(...[
				killShot,
				explosiveShot4,
				simple.allowExplosiveShotDownrank ? explosiveShot3 : null,
				simple.trapWeave ? trapWeave : null,
				simple.sting == StingType.ScorpidSting ? scorpidSting : null,
				simple.sting == StingType.SerpentSting ? serpentSting : null,
				blackArrow,
				aimedShot,
				multiShot,
				steadyShot,
			].filter(a => a) as Array<APLAction>)
		}

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action => APLListItem.create({
				action: action,
			}))
		});
	},

	raidSimPresets: [
		{
			spec: Spec.SpecHunter,
			tooltip: 'Beast Mastery Hunter',
			defaultName: 'Beast Mastery',
			iconUrl: getSpecIcon(Class.ClassHunter, 0),

			talents: Presets.BeastMasteryTalents.data,
			specOptions: Presets.BMDefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.MM_P1_PRESET.gear,
					2: Presets.MM_P2_PRESET.gear,
					3: Presets.MM_P3_PRESET.gear,
					4: Presets.MM_P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.MM_P1_PRESET.gear,
					2: Presets.MM_P2_PRESET.gear,
					3: Presets.MM_P3_PRESET.gear,
					4: Presets.MM_P4_PRESET.gear,
				},
			},
		},
		{
			spec: Spec.SpecHunter,
			tooltip: 'Marksmanship Hunter',
			defaultName: 'Marksmanship',
			iconUrl: getSpecIcon(Class.ClassHunter, 1),

			talents: Presets.MarksmanTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.MM_P1_PRESET.gear,
					2: Presets.MM_P2_PRESET.gear,
					3: Presets.MM_P3_PRESET.gear,
					4: Presets.MM_P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.MM_P1_PRESET.gear,
					2: Presets.MM_P2_PRESET.gear,
					3: Presets.MM_P3_PRESET.gear,
					4: Presets.MM_P4_PRESET.gear,
				},
			},
		},
		{
			spec: Spec.SpecHunter,
			tooltip: 'Survival Hunter',
			defaultName: 'Survival',
			iconUrl: getSpecIcon(Class.ClassHunter, 2),

			talents: Presets.SurvivalTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.SV_P1_PRESET.gear,
					2: Presets.SV_P2_PRESET.gear,
					3: Presets.SV_P3_PRESET.gear,
					4: Presets.SV_P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.SV_P1_PRESET.gear,
					2: Presets.SV_P2_PRESET.gear,
					3: Presets.SV_P3_PRESET.gear,
					4: Presets.SV_P4_PRESET.gear,
				},
			},
		},
	],
});

export class HunterSimUI extends IndividualSimUI<Spec.SpecHunter> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecHunter>) {
		super(parentElem, player, SPEC_CONFIG);
		const _gemOptimizer = new HunterGemOptimizer(this);
	}
}

class HunterGemOptimizer extends PhysicalDPSGemOptimizer {
	readonly player: Player<Spec.SpecHunter>;
	arpSlop: number = 4;
	hitSlop: number = 11;

	constructor(simUI: IndividualSimUI<Spec.SpecHunter>) {
		super(simUI, true, false, true, false);
		this.player = simUI.player;
	}

	detectArpStackConfiguration(ungemmedGear: Gear): boolean {
		/*
		 * Allow ArP gems only for Marksmanship specialization. Additionally,
		 * unlike the Warrior and Feral sims, Marksmanship gemming algorithm has
		 * an additional restriction of only gemming ArP in hard cap setups where
		 * the passive ArP on the ungemmed gear set is already very high.
		 */
		this.useArpGems = (this.player.getTalentTree() === 1) && (this.arpTarget > 1000) && (this.passiveArp > 648);
		return super.detectArpStackConfiguration(ungemmedGear);
	}
}
