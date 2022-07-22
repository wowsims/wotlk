import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Mechanics from '/wotlk/core/constants/mechanics.js';

import { PaladinMajorGlyph, PaladinSeal } from '/wotlk/core/proto/paladin.js';

import * as ProtectionPaladinInputs from './inputs.js';
import * as Presets from './presets.js';

export class ProtectionPaladinSimUI extends IndividualSimUI<Spec.SpecProtectionPaladin> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecProtectionPaladin>) {
		super(parentElem, player, {
			cssClass: 'protection-paladin-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatStamina,
				Stat.StatStrength,
				Stat.StatAgility,
				Stat.StatIntellect,
				Stat.StatAttackPower,
				Stat.StatMeleeHit,
				Stat.StatMeleeCrit,
				Stat.StatExpertise,
				Stat.StatMeleeHaste,
				Stat.StatArmorPenetration,
				Stat.StatSpellPower,
				Stat.StatSpellCrit,
				Stat.StatSpellHit,
				Stat.StatArmor,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatArmor,
				Stat.StatStamina,
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
				Stat.StatHolySpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatDefense,
				Stat.StatBlock,
				Stat.StatBlockValue,
				Stat.StatDodge,
				Stat.StatParry,
				Stat.StatResilience,
			],
			modifyDisplayStats: (player: Player<Spec.SpecProtectionPaladin>) => {
				let stats = new Stats();

				TypedEvent.freezeAllAndDo(() => {
					if (player.getMajorGlyphs().includes(PaladinMajorGlyph.GlyphOfSealOfVengeance) && (player.getSpecOptions().seal == PaladinSeal.Vengeance)) {
						stats = stats.addStat(Stat.StatExpertise, 10 * Mechanics.EXPERTISE_RATING_PER_EXPERTISE);
					}
				})

				return {
					talents: stats,
				};
			},
			defaults: {
				// Default equipped gear.
				gear: Presets.P4_PRESET.gear,
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
				talents: Presets.GenericAoeTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					powerWordFortitude: TristateEffect.TristateEffectRegular,
					shadowProtection: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					thorns: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					unleashedRage: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfSanctuary: true,
					blessingOfWisdom: TristateEffect.TristateEffectImproved,
				}),
				debuffs: Debuffs.create({

					misery: true,
					bloodFrenzy: true,
					exposeArmor: true,
					sunderArmor: true,
					faerieFire: TristateEffect.TristateEffectImproved,
					curseOfWeakness: TristateEffect.TristateEffectRegular,
				}),
			},

			// IconInputs to include in the 'Player' section on the settings tab.
			playerIconInputs: [
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: ProtectionPaladinInputs.ProtectionPaladinRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.TankAssignment,
					OtherInputs.IncomingHps,
					OtherInputs.HealingCadence,
					OtherInputs.HpPercentForDefensives,
					ProtectionPaladinInputs.AuraSelection,
					ProtectionPaladinInputs.UseAvengingWrath,
					ProtectionPaladinInputs.JudgementSelection,
					ProtectionPaladinInputs.StartingSealSelection,
					ProtectionPaladinInputs.DamageTakenPerSecond,
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
				// Preset gear configurations that the user can quickly select.
				gear: [
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
