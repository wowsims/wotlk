import { RaidBuffs } from '/wotlk/core/proto/common.js';
import { PartyBuffs } from '/wotlk/core/proto/common.js';
import { IndividualBuffs } from '/wotlk/core/proto/common.js';
import { Debuffs } from '/wotlk/core/proto/common.js';
import { Class } from '/wotlk/core/proto/common.js';
import { Consumes } from '/wotlk/core/proto/common.js';
import { Encounter } from '/wotlk/core/proto/common.js';
import { ItemSlot } from '/wotlk/core/proto/common.js';
import { MobType } from '/wotlk/core/proto/common.js';
import { RaidTarget } from '/wotlk/core/proto/common.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { Stat } from '/wotlk/core/proto/common.js';
import { TristateEffect } from '/wotlk/core/proto/common.js'
import { Stats } from '/wotlk/core/proto_utils/stats.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';

import { SmitePriest, SmitePriest_Rotation as Rotation, SmitePriest_Options as Options, SmitePriest_Rotation, SmitePriest_Rotation_RotationType } from '/wotlk/core/proto/priest.js';

import * as IconInputs from '/wotlk/core/components/icon_inputs.js';
import * as OtherInputs from '/wotlk/core/components/other_inputs.js';
import * as Mechanics from '/wotlk/core/constants/mechanics.js';
import * as Tooltips from '/wotlk/core/constants/tooltips.js';

import * as SmitePriestInputs from './inputs.js';
import * as Presets from './presets.js';

export class SmitePriestSimUI extends IndividualSimUI<Spec.SpecSmitePriest> {
	constructor(parentElem: HTMLElement, player: Player<Spec.SpecSmitePriest>) {
		super(parentElem, player, {
			cssClass: 'smite-priest-sim-ui',
			// List any known bugs / issues here and they'll be shown on the site.
			knownIssues: [
			],

			// All stats for which EP should be calculated.
			epStats: [
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			// Reference stat against which to calculate EP. I think all classes use either spell power or attack power.
			epReferenceStat: Stat.StatSpellPower,
			// Which stats to display in the Character Stats section, at the bottom of the left-hand sidebar.
			displayStats: [
				Stat.StatHealth,
				Stat.StatStamina,
				Stat.StatIntellect,
				Stat.StatSpirit,
				Stat.StatSpellPower,
				Stat.StatShadowSpellPower,
				Stat.StatHolySpellPower,
				Stat.StatSpellHit,
				Stat.StatSpellCrit,
				Stat.StatSpellHaste,
				Stat.StatMP5,
			],
			modifyDisplayStats: (player: Player<Spec.SpecSmitePriest>) => {
				let stats = new Stats();
				stats = stats.addStat(Stat.StatSpellHit,
					player.getTalents().shadowFocus * 2 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE +
					player.getTalents().focusedPower * 2 * Mechanics.SPELL_HIT_RATING_PER_HIT_CHANCE);

				return {
					talents: stats,
				};
			},

			defaults: {
				// Default equipped gear.
				gear: Presets.P1_PRESET.gear,
				// Default EP weights for sorting gear in the gear picker.
				epWeights: Stats.fromMap({
					[Stat.StatIntellect]: 1.38,
					[Stat.StatSpirit]: 1.18,
					[Stat.StatSpellPower]: 1,
					[Stat.StatSpellHit]: 2.57,
					[Stat.StatShadowSpellPower]: 0.05,
					[Stat.StatHolySpellPower]: 0.95,
					[Stat.StatSpellCrit]: 0.44,
					[Stat.StatSpellHaste]: 0.28, // tricky because SP is tricky
					[Stat.StatMP5]: 2.05,
				}),
				// Default consumes settings.
				consumes: Presets.DefaultConsumes,
				// Default rotation settings.
				rotation: Presets.DefaultRotation,
				// Default talents.
				talents: Presets.StandardTalents.data,
				// Default spec-specific settings.
				specOptions: Presets.DefaultOptions,
				// Default raid/party buffs settings.
				raidBuffs: RaidBuffs.create({
					arcaneBrilliance: true,
					divineSpirit: true,
					giftOfTheWild: TristateEffect.TristateEffectImproved,
					bloodlust: true,
					manaSpringTotem: TristateEffect.TristateEffectRegular,
					wrathOfAirTotem: true,
				}),
				partyBuffs: PartyBuffs.create({
				}),
				individualBuffs: IndividualBuffs.create({
					blessingOfKings: true,
					blessingOfWisdom: 2,

				}),
				debuffs: Debuffs.create({
					judgementOfWisdom: true,
					misery: true,
					curseOfElements: true,
				}),
			},

			// IconInputs to include in the 'Self Buffs' section on the settings tab.
			selfBuffInputs: [
				SmitePriestInputs.SelfPowerInfusion,
			],
			// Inputs to include in the 'Rotation' section on the settings tab.
			rotationInputs: SmitePriestInputs.SmitePriestRotationConfig,
			// Inputs to include in the 'Other' section on the settings tab.
			otherInputs: {
				inputs: [
					OtherInputs.PrepopPotion,
					OtherInputs.TankAssignment,
				],
			},
			encounterPicker: {
				// Target stats to show for 'Simple' encounters.
				simpleTargetStats: [
				],
				// Whether to include 'Execute Duration (%)' in the 'Encounter' section of the settings tab.
				showExecuteProportion: false,
			},

			presets: {
				// Preset talents that the user can quickly select.
				talents: [
					Presets.StandardTalents,
					Presets.HolyTalents,
				],
				// Preset gear configurations that the user can quickly select.
				gear: [
					Presets.P1_PRESET,
				],
			},
		});
	}
}
