import { FeralDruid_Rotation_FinishingMove as FinishingMove } from '../core/proto/druid.js';
import { FeralDruid_Options as DruidOptions } from '../core/proto/druid.js';
import { RaidTarget } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { NO_TARGET } from '../core/proto_utils/utils.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { EventID, TypedEvent } from '../core/typed_event.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Target } from '../core/target.js';
import { getEnumValues } from '../core/utils.js';
import { ItemSlot } from '../core/proto/common.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Helper function for identifying whether 2pT6 is equipped, which impacts allowed rotation options
function numThunderheartPieces(player: Player<Spec.SpecFeralDruid>): number {
	const gear = player.getGear();
	const itemIds = [31048, 31042, 31034, 31044, 31039, 34556, 34444, 34573];
	return gear.asArray().map(equippedItem => equippedItem?.item.id).filter(id => itemIds.includes(id!)).length
}

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecFeralDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecFeralDruid>) => player.getSpecOptions().innervateTarget?.targetIndex != NO_TARGET,
	setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = RaidTarget.create({
			targetIndex: newValue ? 0 : NO_TARGET,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const LatencyMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFeralDruid>({
	fieldName: 'latencyMs',
	label: 'Latency',
	labelTooltip: 'Player latency, in milliseconds. Adds a delay to actions that cannot be spell queued.',
});

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, FinishingMove>({
			fieldName: 'finishingMove',
			label: 'Finishing Move',
			labelTooltip: 'Specify whether Rip or Ferocious Bite should be used as the primary finisher in the DPS rotation.',
			values: [
				{ name: 'Rip', value: FinishingMove.Rip },
				{ name: 'Ferocious Bite', value: FinishingMove.Bite },
				{ name: 'None', value: FinishingMove.None },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'biteweave',
			label: 'Enable Bite-weaving',
			labelTooltip: 'Spend Combo Points on Ferocious Bite when Rip is already applied on the target.',
			enableWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().finishingMove == FinishingMove.Rip,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'ripweave',
			label: 'Enable Rip-weaving',
			labelTooltip: 'Spend Combo Points on Rip when at 52 Energy or above.',
			enableWhen: (player: Player<Spec.SpecFeralDruid>) => player.getRotation().finishingMove == FinishingMove.Bite,
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, number>({
			fieldName: 'ripMinComboPoints',
			label: 'Rip CP Threshold',
			labelTooltip: 'Minimum Combo Points to accumulate before casting Rip as a finisher.',
			values: [
				{ name: '4', value: 4 },
				{ name: '5', value: 5 },
			],
			enableWhen: (player: Player<Spec.SpecFeralDruid>) => (player.getRotation().finishingMove == FinishingMove.Rip) || (player.getRotation().ripweave && (player.getRotation().finishingMove != FinishingMove.None)),
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, number>({
			fieldName: 'biteMinComboPoints',
			label: 'Bite CP Threshold',
			labelTooltip: 'Minimum Combo Points to accumulate before casting Ferocious Bite as a finisher.',
			values: [
				{ name: '4', value: 4 },
				{ name: '5', value: 5 },
			],
			enableWhen: (player: Player<Spec.SpecFeralDruid>) => (player.getRotation().finishingMove == FinishingMove.Bite) || (player.getRotation().biteweave && (player.getRotation().finishingMove != FinishingMove.None)),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'mangleTrick',
			label: 'Use Mangle trick',
			labelTooltip: 'Cast Mangle rather than Shred when between 50-56 Energy with 2pT6, or 60-61 Energy without 2pT6, and with less than 1 second remaining until the next Energy tick.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'rakeTrick',
			label: 'Use Rake/Bite tricks',
			labelTooltip: 'Cast Rake or Ferocious Bite rather than powershifting when between 35-39 Energy without 2pT6, and with more than 1 second remaining until the next Energy tick.',
			enableWhen: (player: Player<Spec.SpecFeralDruid>) => numThunderheartPieces(player) < 2,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'maintainFaerieFire',
			label: 'Maintain Faerie Fire',
			labelTooltip: 'Use Faerie Fire whenever it is not active on the target.',
		}),
	],
};
