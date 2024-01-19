import {
	Class,
	Debuffs,
	Faction,
	IndividualBuffs,
	PartyBuffs,
	Race,
	RaidBuffs,
	Spec,
	Stat, PseudoStat,
	TristateEffect,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { getSpecIcon } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';

import * as OtherInputs from '../core/components/other_inputs.js';
import * as RetributionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecRetributionPaladin, {
	cssClass: 'retribution-paladin-sim-ui',
	cssScheme: 'paladin',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatMP5,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmorPenetration,
		Stat.StatSpellPower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatSpellHaste,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatMP5,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmorPenetration,
		Stat.StatSpellHaste,
		Stat.StatSpellPower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatMana,
		Stat.StatHealth,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.DefaultGear.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatStrength]: 2.53,
			[Stat.StatAgility]: 1.13,
			[Stat.StatIntellect]: 0.15,
			[Stat.StatSpellPower]: 0.32,
			[Stat.StatSpellHit]: 0.41,
			[Stat.StatSpellCrit]: 0.01,
			[Stat.StatSpellHaste]: 0.12,
			[Stat.StatMP5]: 0.05,
			[Stat.StatAttackPower]: 1,
			[Stat.StatMeleeHit]: 1.96,
			[Stat.StatMeleeCrit]: 1.16,
			[Stat.StatMeleeHaste]: 1.44,
			[Stat.StatArmorPenetration]: 0.76,
			[Stat.StatExpertise]: 1.80,
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 7.33,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default rotation settings.
		rotation: Presets.DefaultRotation,
		// Default talents.
		talents: Presets.AuraMasteryTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: RaidBuffs.create({
			arcaneBrilliance: true,
			divineSpirit: true,
			giftOfTheWild: TristateEffect.TristateEffectImproved,
			manaSpringTotem: TristateEffect.TristateEffectRegular,
			battleShout: TristateEffect.TristateEffectImproved,
			trueshotAura: true,
		}),
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfKings: true,
			blessingOfMight: TristateEffect.TristateEffectImproved,
		}),
		debuffs: Debuffs.create({
			judgementOfWisdom: true,
			judgementOfLight: true,
			curseOfElements: true,
			exposeArmor: TristateEffect.TristateEffectImproved,
			sunderArmor: true,
			faerieFire: true,
			curseOfWeakness: TristateEffect.TristateEffectRegular,
		}),
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		RetributionPaladinInputs.AuraSelection,
		RetributionPaladinInputs.JudgementSelection,
		RetributionPaladinInputs.StartingSealSelection,
	],
	// Inputs to include in the 'Rotation' section on the settings tab.
	rotationInputs: {
		inputs: [
			RetributionPaladinInputs.RotationSelector,
			RetributionPaladinInputs.RetributionPaladinRotationDivinePleaSelection,
			RetributionPaladinInputs.RetributionPaladinRotationAvoidClippingConsecration,
			RetributionPaladinInputs.RetributionPaladinRotationHoldLastAvengingWrathUntilExecution,
			RetributionPaladinInputs.RetributionPaladinRotationCancelChaosBane,
			RetributionPaladinInputs.RetributionPaladinRotationDivinePleaSelectionAlternate,
			RetributionPaladinInputs.RetributionPaladinRotationDivinePleaPercentageConfig,
			RetributionPaladinInputs.RetributionPaladinRotationConsSlackConfig,
			RetributionPaladinInputs.RetributionPaladinRotationExoSlackConfig,
			RetributionPaladinInputs.RetributionPaladinRotationHolyWrathConfig,
			RetributionPaladinInputs.RetributionPaladinSoVTargets,
			RetributionPaladinInputs.RetributionPaladinRotationPriorityConfig,
			RetributionPaladinInputs.RetributionPaladinCastSequenceConfig
		]
	},
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
	],
	excludeBuffDebuffInputs: [
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		rotations: [
			Presets.ROTATION_PRESET_LEGACY_DEFAULT,
			Presets.ROTATION_PRESET_DEFAULT,
		],
		// Preset talents that the user can quickly select.
		talents: [
			Presets.AuraMasteryTalents,
			Presets.DivineSacTalents,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.DefaultGear,
		],
	},

	autoRotation: (): APLRotation => {
		return Presets.ROTATION_PRESET_DEFAULT.rotation.rotation!;
	},

	raidSimPresets: [
		{
			spec: Spec.SpecRetributionPaladin,
			tooltip: 'Retribution Paladin',
			defaultName: 'Retribution',
			iconUrl: getSpecIcon(Class.ClassPaladin, 2),

			talents: Presets.AuraMasteryTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceHuman,
				[Faction.Horde]: Race.RaceUnknown,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.DefaultGear.gear,
				},
				[Faction.Horde]: {
					1: Presets.DefaultGear.gear,
				},
			},
		},
	],
});

export class RetributionPaladinSimUI extends IndividualSimUI<Spec.SpecRetributionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
