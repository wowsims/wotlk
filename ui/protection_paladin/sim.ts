import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import {
	APLAction,
	APLListItem,
	APLRotation,
} from '../core/proto/apl.js';
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';

import { PaladinMajorGlyph, PaladinSeal } from '../core/proto/paladin.js';

import * as ProtectionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, {
			cssClass: 'protection-paladin-sim-ui',
			cssScheme: 'paladin',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatSpellHit,
				Stat.StatMeleeCrit,
				Stat.StatExpertise,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatSpellPower,
				Stat.StatArmor,
				Stat.StatBonusArmor,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
				Stat.StatNatureResistance,
				Stat.StatShadowResistance,
				Stat.StatFrostResistance,
			],
			epPseudoStats: [
				PseudoStat.PseudoStatMainHandDps,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatBonusArmor,
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatMeleeHaste,
				Stat.StatExpertise,
				Stat.StatArmorPenetration,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
				Stat.StatNatureResistance,
				Stat.StatShadowResistance,
				Stat.StatFrostResistance,
			],
			modifyDisplayStats: (player: Player<Spec.SpecProtectionPaladin>) => {
				let stats = new Stats();

				TypedEvent.freezeAllAndDo(() => {
					if (player.getMajorGlyphs().includes(PaladinMajorGlyph.GlyphOfSealOfVengeance) && (player.getSpecOptions().seal == PaladinSeal.Vengeance)) {
						stats = stats.addStat(Stat.StatExpertise, 10 * Mechanics.EXPERTISE_PER_QUARTER_PERCENT_REDUCTION);
					}
				})

				return {
					talents: stats,
				};
			},
			defaults: {
				// Default equipped gear.
				gear: Presets.P3_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatArmor]: 0.07,
					[Stat.StatBonusArmor]: 0.06,
					[Stat.StatStamina]: 1.14,
					[Stat.StatStrength]: 1.00,
					[Stat.StatAgility]: 0.62,
					[Stat.StatAttackPower]: 0.26,
					[Stat.StatExpertise]: 0.69,
					[Stat.StatMeleeHit]: 0.79,
					[Stat.StatMeleeCrit]: 0.30,
					[Stat.StatMeleeHaste]: 0.17,
					[Stat.StatArmorPenetration]: 0.04,
					[Stat.StatSpellPower]: 0.13,
					[Stat.StatBlock]: 0.52,
					[Stat.StatBlockValue]: 0.28,
					[Stat.StatDodge]: 0.46,
					[Stat.StatParry]: 0.61,
					[Stat.StatDefense]: 0.54,
				}, {
					[PseudoStat.PseudoStatMainHandDps]: 3.33,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.GenericAoeTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					powerWordFortitude: TristateEffect.TristateEffectImproved,
					strengthOfEarthTotem: TristateEffect.TristateEffectImproved,
					arcaneBrilliance: true,
					unleashedRage: true,
					leaderOfThePack: TristateEffect.TristateEffectRegular,
					icyTalons: true,
					totemOfWrath: true,
					demonicPact: 500,
					swiftRetribution: true,
					moonkinAura: TristateEffect.TristateEffectRegular,
					sanctifiedRetribution: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					bloodlust: true,
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
					misery: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					ebonPlaguebringer: true,
					totemOfWrath: true,
					shadowMastery: true,
					bloodFrenzy: true,
					mangle: true,
					exposeArmor: true,
					sunderArmor: true,
					vindication: true,
					thunderClap: TristateEffect.TristateEffectImproved,
					insectSwarm: true,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ProtectionPaladinInputs.ProtectionPaladinRotationConfig,
			// Buff and Debuff inputs to include/exclude, overriding the EP-based defaults.
			includeBuffDebuffInputs: [
				IconInputs.HealthBuff,
			],
			excludeBuffDebuffInputs: [
			],
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HealingCadenceVariation,
					OtherInputs.BurstWindow,
					OtherInputs.HpPercentForDefensives,
					OtherInputs.InspirationUptime,
					ProtectionPaladinInputs.AuraSelection,
					ProtectionPaladinInputs.UseAvengingWrath,
					ProtectionPaladinInputs.JudgementSelection,
					ProtectionPaladinInputs.StartingSealSelection,
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
					Presets.GenericAoeTalents,
				],
				// Preset rotations that the user can quickly select.
				rotations: [
					Presets.ROTATION_DEFAULT,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.PRERAID_PRESET,
					Presets.P4_PRERAID_PRESET,
					Presets.P1_PRESET,
					Presets.P2_PRESET,
					Presets.P3_PRESET,
					Presets.P4_PRESET,
				],
			},

			autoRotation: (player: Player<Spec.SpecProtectionPaladin>): APLRotation => {
				return Presets.ROTATION_DEFAULT.rotation.rotation!;
			},
		});
	}
}
