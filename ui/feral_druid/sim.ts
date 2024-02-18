import {
	Class,
	Cooldowns,
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	PseudoStat,
	Race,
	RaidBuffs,
	Spec,
	Stat,
	TristateEffect,
} from '../core/proto/common.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { Gear } from '../core/proto_utils/gear.js';
import { PhysicalDPSGemOptimizer } from '../core/components/suggest_gems_action.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { Player } from '../core/player.js';

import {
	FeralDruid_Rotation as DruidRotation,
} from '../core/proto/druid.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as AplUtils from '../core/proto_utils/apl_utils.js';

import * as DruidInputs from './inputs.js';
import * as Presets from './presets.js';
import {
	APLAction,
	APLPrepullAction,
	APLListItem,
	APLRotation,
} from '../core/proto/apl.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecFeralDruid, {
	cssClass: 'feral-druid-sim-ui',
	cssScheme: 'druid',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],
	warnings: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
		Stat.StatMana,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStrength]: 2.40,
			[Stat.StatAgility]: 2.39,
			[Stat.StatAttackPower]: 1,
			[Stat.StatMeleeHit]: 2.51,
			[Stat.StatMeleeCrit]: 2.23,
			[Stat.StatMeleeHaste]: 1.83,
			[Stat.StatArmorPenetration]: 2.08,
			[Stat.StatExpertise]: 2.44,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 16.5,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			arcaneBrilliance: true,
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			bloodlust: true,
			manaSpringTotem: TristateEffect.TristateEffectRegular,
			strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
			battleShout: TristateEffect.TristateEffectImproved,
			unleashedRage: true,
			icyTalons: true,
			swiftRetribution: true,
			sanctifiedRetribution: true,
		}),
		partyBuffs: PartyBuffs.create({
			heroicPresence: true,
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfKings: true,
			blessingOfMight: TristateEffect.TristateEffectImproved,
		}),
		debuffs: Debuffs.create({
			judgementOfWisdom: true,
			bloodFrenzy: true,
			giftOfArthas: true,
			exposeArmor: true,
			faerieFire: TristateEffect.TristateEffectImproved,
			sunderArmor: true,
			curseOfWeakness: TristateEffect.TristateEffectRegular,
			heartOfTheCrusader: true,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: DruidInputs.FeralDruidRotationConfig,
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.IntellectBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.JudgementOfWisdom,
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			DruidInputs.LatencyMs,
			DruidInputs.AssumeBleedActive,
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
			Presets.StandardTalents,
		],
		rotations: [
			Presets.SIMPLE_ROTATION_DEFAULT,
			Presets.APL_ROTATION_DEFAULT,
			Presets.APL_ROTATION_CUSTOM_EXAMPLE,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.P1_PRESET,
			Presets.P2_PRESET,
			Presets.P3_PRESET,
			Presets.P4_PRESET,
		],
	},
	
	autoRotation: (_player: Player<Spec.SpecFeralDruid>): APLRotation => {
		return Presets.APL_ROTATION_DEFAULT.rotation.rotation!;
	},

	simpleRotation: (player: Player<Spec.SpecFeralDruid>, simple: DruidRotation, cooldowns: Cooldowns): APLRotation => {
		let [prepullActions, actions] = AplUtils.standardCooldownDefaults(cooldowns);

		const preOmen = APLPrepullAction.fromJsonString(`{"action":{"activateAura":{"auraId":{"spellId":16870}}},"doAtValue":{"const":{"val":"-1s"}}}`);
		const preZerk = APLPrepullAction.fromJsonString(`{"action":{"castSpell":{"spellId":{"spellId":50334}}},"doAtValue":{"const":{"val":"-1s"}}}`);
		const blockZerk = APLAction.fromJsonString(`{"condition":{"const":{"val":"false"}},"castSpell":{"spellId":{"spellId":50334}}}`);
		const doRotation = APLAction.fromJsonString(`{"catOptimalRotationAction":{"rotationType":${simple.rotationType},"manualParams":${simple.manualParams},"maxFfDelay":${simple.maxFfDelay.toFixed(2)},"minRoarOffset":${simple.minRoarOffset.toFixed(2)},"ripLeeway":${simple.ripLeeway.toFixed(0)},"useRake":${simple.useRake},"useBite":${simple.useBite},"biteTime":${simple.biteTime.toFixed(2)},"flowerWeave":${simple.flowerWeave}}}`);

		prepullActions.push(...[
			simple.prePopOoc ? preOmen: null,
			simple.prePopBerserk ? preZerk: null,
		].filter(a => a) as Array<APLPrepullAction>)

		actions.push(...[
			blockZerk,
			doRotation,
		].filter(a => a) as Array<APLAction>)

		return APLRotation.create({
			prepullActions: prepullActions,
			priorityList: actions.map(action => APLListItem.create({
				action: action,
			}))
		});
	},

	raidSimPresets: [
		{
			spec: Spec.SpecFeralDruid,
			tooltip: specNames[Spec.SpecFeralDruid],
			defaultName: 'Cat',
			iconUrl: getSpecIcon(Class.ClassDruid, 3),

			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceNightElf,
				[Faction.Horde]: Race.RaceTauren,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET.gear,
					3: Presets.P3_PRESET.gear,
					4: Presets.P4_PRESET.gear,
				},
			},
		},
	],
});

export class FeralDruidSimUI extends IndividualSimUI<Spec.SpecFeralDruid> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecFeralDruid>) {
		super(parentElem, player, SPEC_CONFIG);

		const _gemOptimizer = new FeralGemOptimizer(this);
	}
}

class FeralGemOptimizer extends PhysicalDPSGemOptimizer {
	constructor(simUI: IndividualSimUI<Spec.SpecFeralDruid>) {
		super(simUI, true, true, true, true);
	}

	calcCritCap(gear: Gear): Stats {
		const baseCritCapPercentage = 77.8; // includes 3% Crit debuff
		let agiProcs = 0;

		if (gear.hasRelic(47668)) {
			agiProcs += 200;
		}

		if (gear.hasRelic(50456)) {
			agiProcs += 44*5;
		}

		if (gear.hasTrinket(47131) || gear.hasTrinket(47464)) {
			agiProcs += 510;
		}

		if (gear.hasTrinket(47115) || gear.hasTrinket(47303)) {
			agiProcs += 450;
		}

		if (gear.hasTrinket(44253) || gear.hasTrinket(42987)) {
			agiProcs += 300;
		}

		return new Stats().withStat(Stat.StatMeleeCrit, (baseCritCapPercentage - agiProcs*1.1*1.06*1.02/83.33) * 45.91);
	}
}
