import { RaidBuffs } from '../core/proto/common.js';
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
import { TotemsSection } from '../core/components/totem_inputs.js';

import { TankDeathknight, TankDeathknight_Rotation as DeathKnightRotation, DeathknightTalents as DeathKnightTalents, TankDeathknight_Options as DeathKnightOptions } from '../core/proto/deathknight.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Tooltips from '../core/constants/tooltips.js';

import * as DeathKnightInputs from './inputs.js';
import * as Presets from './presets.js';

export class TankDeathknightSimUI extends IndividualSimUI<Spec.SpecTankDeathknight> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecTankDeathknight>) {
		super(parentElem, player, {
			cssClass: 'tank-deathknight-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
				"<p>Completely unfinished.</p>"
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmor,
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
				Stat.StatSpellHit,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatAttackPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatExpertise,
				Stat.StatMeleeHit,
				Stat.StatSpellHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatDefense,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			defaults: {
				// Default equipped gear.
				gear: Presets.P1_BLOOD_BIS_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 0.05,
					[Stat.StatStamina]: 1,
					[Stat.StatStrength]: 0.33,
					[Stat.StatAgility]: 0.6,
					[Stat.StatAttackPower]: 0.06,
					[Stat.StatExpertise]: 0.67,
					[Stat.StatMeleeHit]: 0.67,
					[Stat.StatMeleeCrit]: 0.28,
					[Stat.StatMeleeHaste]: 0.21,
					[Stat.StatArmorPenetration]: 0.19,
					[Stat.StatBlock]: 0.35,
					[Stat.StatBlockValue]: 0.59,
					[Stat.StatDodge]: 0.7,
					[Stat.StatParry]: 0.58,
					[Stat.StatDefense]: 0.8,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.BloodTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					retributionAura: true,
					powerWordFortitude: TristateEffect.TristateEffectImproved, 
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					swiftRetribution: true,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					icyTalons: true,
					abominationsMight: true,
					leaderOfThePack: TristateEffect.TristateEffectRegular,
					sanctifiedRetribution: true,
					bloodlust: true,
					devotionAura: TristateEffect.TristateEffectImproved,
					stoneskinTotem: TristateEffect.TristateEffectImproved,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
					blessingOfSanctuary: true,
				}),
				debuffs: Debuffs.create({
					bloodFrenzy: true,
					faerieFire: TristateEffect.TristateEffectRegular,
					sunderArmor: true,
					misery: true,
					ebonPlaguebringer: true,
					mangle: true,
					heartOfTheCrusader: true,
					demoralizingShout: TristateEffect.TristateEffectImproved,
					frostFever: TristateEffect.TristateEffectImproved,
					insectSwarm: true,
					judgementOfLight: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: DeathKnightInputs.TankDeathKnightRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.SpellDamageDebuff,
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.InFrontOfTarget,
					DeathKnightInputs.StartingRunicPower,
				],
			},
			encounterPicker: {
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.BloodTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_BLOOD_BIS_PRESET,
				],
			},
		});
	}
}
