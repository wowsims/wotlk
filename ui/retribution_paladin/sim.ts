import { RaidBuffs } from '../core/proto/common.js';
import { PartyBuffs } from '../core/proto/common.js';
import { IndividualBuffs } from '../core/proto/common.js';
import { Debuffs } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { Stat, PseudoStat } from '../core/proto/common.js';
import { TristateEffect } from '../core/proto/common.js'
import { Stats } from '../core/proto_utils/stats.js';
import { Player } from '../core/player.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as IconInputs from '../core/components/icon_inputs.js';
import * as OtherInputs from '../core/components/other_inputs.js';
import * as Mechanics from '../core/constants/mechanics.js';

import { PaladinMajorGlyph, PaladinSeal } from '../core/proto/paladin.js';

import * as RetributionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

export class RetributionPaladinSimUI extends IndividualSimUI<Spec.SpecRetributionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecRetributionPaladin>) {
		super(parentElem, player, {
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
			modifyDisplayStats: (player: Player<Spec.SpecRetributionPaladin>) => {
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
				gear: Presets.P1_PRESET.gear,
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
					bloodlust: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					hornOfWinter: true,
					battleShout: TristateEffect.TristateEffectImproved,
					sanctifiedRetribution: true,
					swiftRetribution: true,
					elementalOath: true,
					rampage: true,
					trueshotAura: true,
					icyTalons: true,
					totemOfWrath: true,
					wrathOfAirTotem: true,
					demonicPact: 500,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					judgementsOfTheWise: true,
					blessingOfKings: true,
					blessingOfMight: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({
					shadowMastery: true,
					totemOfWrath: true,
					judgementOfWisdom: true,
					judgementOfLight: true,
					misery: true,
					curseOfElements: true,
					bloodFrenzy: true,
					exposeArmor: true,
					sunderArmor: true,
					faerieFire: TristateEffect.TristateEffectImproved,
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
				IconInputs.ReplenishmentBuff,
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
				// Preset talents that the user can quickly select.
				talents: [
					Presets.AuraMasteryTalents,
					Presets.DivineSacTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.PRE_RAID_PRESET,
					Presets.P1_PRESET,
					Presets.P2_PRESET,
					Presets.P3_PRESET,
					Presets.P4_PRESET,
					Presets.P5_PRESET,
				],
			},
		});
	}
}
