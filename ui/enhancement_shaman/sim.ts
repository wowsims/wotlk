import {
	Class,
	Faction,
	IndividualBuffs,
	ItemSlot,
	PartyBuffs,
	PseudoStat,
	Race,
	Spec,
	Stat,
	TristateEffect,
} from '../core/proto/common.js';
import {
	APLRotation,
} from '../core/proto/apl.js';
import { ShamanImbue } from '../core/proto/shaman.js';
import { Player } from '../core/player.js';
import { Stats } from '../core/proto_utils/stats.js';
import { getSpecIcon, specNames } from '../core/proto_utils/utils.js';
import { IndividualSimUI, registerSpecConfig } from '../core/individual_sim_ui.js';
import { TotemsSection } from '../core/components/totem_inputs.js';

import * as BuffDebuffInputs from '../core/components/inputs/buffs_debuffs.js';
import * as OtherInputs from '../core/components/other_inputs.js';

import * as ShamanInputs from './inputs.js';
import * as Presets from './presets.js';
import { FireElementalSection } from '../core/components/fire_elemental_inputs.js';

const SPEC_CONFIG = registerSpecConfig(Spec.SpecEnhancementShaman, {
	cssClass: 'enhancement-shaman-sim-ui',
	cssScheme: 'shaman',
	// List any known bugs / issues here and they'll be shown on the site.
	knownIssues: [
	],

	// All stats for which EP should be calculated.
	epStats: [
		Stat.StatIntellect,
		Stat.StatAgility,
		Stat.StatStrength,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatArmorPenetration,
		Stat.StatExpertise,
		Stat.StatSpellPower,
		Stat.StatSpellCrit,
		Stat.StatSpellHit,
		Stat.StatSpellHaste,
	],
	epPseudoStats: [
		PseudoStat.PseudoStatMainHandDps,
		PseudoStat.PseudoStatOffHandDps,
	],
	// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
	epReferenceStat: Stat.StatAttackPower,
	// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
	displayStats: [
		Stat.StatHealth,
		Stat.StatStamina,
		Stat.StatStrength,
		Stat.StatAgility,
		Stat.StatIntellect,
		Stat.StatAttackPower,
		Stat.StatMeleeHit,
		Stat.StatMeleeCrit,
		Stat.StatMeleeHaste,
		Stat.StatExpertise,
		Stat.StatArmorPenetration,
		Stat.StatSpellPower,
		Stat.StatSpellHit,
		Stat.StatSpellCrit,
		Stat.StatSpellHaste,
	],

	defaults: {
		// Default equipped gear.
		gear: Presets.P4_PRESET_WF.gear,
		// Default EP weights for sorting gear in the gear picker.
		epWeights: Stats.fromMap({
			[Stat.StatIntellect]: 1.48,
			[Stat.StatAgility]: 1.59,
			[Stat.StatStrength]: 1.1,
			[Stat.StatSpellPower]: 1.13,
			[Stat.StatSpellHit]: 0, //default EP assumes cap
			[Stat.StatSpellCrit]: 0.91,
			[Stat.StatSpellHaste]: 0.37,
			[Stat.StatAttackPower]: 1.0,
			[Stat.StatMeleeHit]: 1.38,
			[Stat.StatMeleeCrit]: 0.81,
			[Stat.StatMeleeHaste]: 1.61, //haste is complicated
			[Stat.StatArmorPenetration]: 0.48,
			[Stat.StatExpertise]: 0, //default EP assumes cap
		}, {
			[PseudoStat.PseudoStatMainHandDps]: 5.21,
			[PseudoStat.PseudoStatOffHandDps]: 2.21,
		}),
		// Default consumes settings.
		consumes: Presets.DefaultConsumes,
		// Default talents.
		talents: Presets.StandardTalents.data,
		// Default spec-specific settings.
		specOptions: Presets.DefaultOptions,
		// Default raid/party buffs settings.
		raidBuffs: Presets.DefaultRaidBuffs,
		partyBuffs: PartyBuffs.create({
		}),
		individualBuffs: IndividualBuffs.create({
			blessingOfKings: true,
			blessingOfWisdom: TristateEffect.TristateEffectImproved,
			blessingOfMight: TristateEffect.TristateEffectImproved,
			judgementsOfTheWise: true,
		}),
		debuffs: Presets.DefaultDebuffs,
	},

	// IconInputs to include in the 'Player' section on the settings tab.
	playerIconInputs: [
		ShamanInputs.ShamanShieldInput,
		ShamanInputs.ShamanImbueMH,
		ShamanInputs.ShamanImbueOH,
	],
	// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
	includeBuffDebuffInputs: [
		BuffDebuffInputs.ReplenishmentBuff,
		BuffDebuffInputs.MP5Buff,
		BuffDebuffInputs.SpellHasteBuff,
		BuffDebuffInputs.SpiritBuff,
	],
	excludeBuffDebuffInputs: [
		BuffDebuffInputs.BleedDebuff,
	],
	// Inputs to include in the 'Other' section on the settings tab.
	otherInputs: {
		inputs: [
			ShamanInputs.SyncTypeInput,
			OtherInputs.TankAssignment,
			OtherInputs.InFrontOfTarget,
		],
	},
	itemSwapSlots: [ItemSlot.ItemSlotMainHand, ItemSlot.ItemSlotOffHand],
	customSections: [
		TotemsSection,
		FireElementalSection
	],
	encounterPicker: {
		// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
		showExecuteProportion: false,
	},

	presets: {
		// Preset talents that the user can quickly select.
		talents: [
			Presets.StandardTalents,
			Presets.Phase3Talents,
		],
		// Preset rotations that the user can quickly select.
		rotations: [
			Presets.ROTATION_FT_DEFAULT,
			Presets.ROTATION_WF_DEFAULT,
			Presets.ROTATION_PHASE_3,
		],
		// Preset gear configurations that the user can quickly select.
		gear: [
			Presets.PRERAID_PRESET,
			Presets.P1_PRESET,
			Presets.P2_PRESET_FT,
			Presets.P2_PRESET_WF,
			Presets.P3_PRESET_ALLIANCE,
			Presets.P3_PRESET_HORDE,
			Presets.P4_PRESET_FT,
			Presets.P4_PRESET_WF,
		],
	},

	autoRotation: (player: Player<Spec.SpecEnhancementShaman>): APLRotation => {
		const hasT94P = player.getCurrentStats().sets.includes('Triumphant Nobundo\'s Battlegear (4pc)')
			|| player.getCurrentStats().sets.includes('Nobundo\'s Battlegear (4pc)')
			|| player.getCurrentStats().sets.includes('Triumphant Thrall\'s Battlegear (4pc)')
			|| player.getCurrentStats().sets.includes('Thrall\'s Battlegear (4pc)');
		const options = player.getSpecOptions();

		if (hasT94P) {
			console.log("has set");
			return Presets.ROTATION_PHASE_3.rotation.rotation!;
		} else if (options.imbueMh == ShamanImbue.FlametongueWeapon) {
			return Presets.ROTATION_FT_DEFAULT.rotation.rotation!;
		} else {
			return Presets.ROTATION_WF_DEFAULT.rotation.rotation!;
		}
	},

	raidSimPresets: [
		{
			spec: Spec.SpecEnhancementShaman,
			tooltip: specNames[Spec.SpecEnhancementShaman],
			defaultName: 'Enhancement',
			iconUrl: getSpecIcon(Class.ClassShaman, 1),

			talents: Presets.StandardTalents.data,
			specOptions: Presets.DefaultOptions,
			consumes: Presets.DefaultConsumes,
			defaultFactionRaces: {
				[Faction.Unknown]: Race.RaceUnknown,
				[Faction.Alliance]: Race.RaceDraenei,
				[Faction.Horde]: Race.RaceOrc,
			},
			defaultGear: {
				[Faction.Unknown]: {},
				[Faction.Alliance]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET_FT.gear,
					3: Presets.P3_PRESET_ALLIANCE.gear,
					4: Presets.P4_PRESET_FT.gear,
				},
				[Faction.Horde]: {
					1: Presets.P1_PRESET.gear,
					2: Presets.P2_PRESET_FT.gear,
					3: Presets.P3_PRESET_HORDE.gear,
					4: Presets.P4_PRESET_FT.gear,
				},
			},
		},
	],
});

export class EnhancementShamanSimUI extends IndividualSimUI<Spec.SpecEnhancementShaman> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecEnhancementShaman>) {
		super(parentElem, player, SPEC_CONFIG);
	}
}
