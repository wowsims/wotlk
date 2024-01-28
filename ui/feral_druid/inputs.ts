import { Player } from '../core/player.js';
import { Spec, UnitReference, UnitReference_Type as UnitType } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecFeralDruid>({
	fieldName: 'innervateTarget',
	actionId: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecFeralDruid>) => player.getSpecOptions().innervateTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecFeralDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const LatencyMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFeralDruid>({
	fieldName: 'latencyMs',
	label: 'Latency',
	labelTooltip: 'Player latency, in milliseconds. Adds a delay to actions that cannot be spell queued.',
});

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'assumeBleedActive',
	label: 'Assume Bleed Always Active',
	labelTooltip: 'Assume bleed always exists for \'Rend and Tear\' calculations. Otherwise will only calculate based on own rip/rake/lacerate.',
	extraCssClasses: ['within-raid-sim-hide'],
})

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'minCombosForRip',
			label: 'Min Rip CP',
			labelTooltip: 'Combo Point threshold for allowing a Rip cast',
			float: false,
			positive: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'maxWaitTime',
			label: 'Max Wait Time',
			labelTooltip: 'Max seconds to wait for an Energy tick to cast rather than powershifting',
			float: true,
			positive: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'preroarDuration',
			label: 'Pre-Roar Duration',
			labelTooltip: 'Seconds remaining on a prior Savage Roar buff at the start of the pull',
			float: true,
			positive: true,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'maintainFaerieFire',
			label: 'Maintain Faerie Fire',
			labelTooltip: 'If checked, bundle Faerie Fire refreshes with powershifts. Ignored if an external Faerie Fire debuff is selected in settings.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'precastTigersFury',
			label: `Pre-cast Tiger's Fury`,
			labelTooltip: `If checked, model a pre-pull Tiger's Fury cast 3 seconds before initiating combat.`,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useShredTrick',
			label: `Use Shred Trick`,
			labelTooltip: `If checked, enable the "Shred trick" micro-optimization. This should only be used on short fight lengths with full powershifting uptime.`,
		}),
	],
};
