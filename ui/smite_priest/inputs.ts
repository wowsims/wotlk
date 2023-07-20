import { UnitReference, UnitReference_Type as UnitType } from '../core/proto/common.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfPowerInfusion = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'powerInfusionTarget',
	id: ActionId.fromSpellId(10060),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecSmitePriest>) => player.getSpecOptions().powerInfusionTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecSmitePriest>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.powerInfusionTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const InnerFire = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'useInnerFire',
	id: ActionId.fromSpellId(48168),
});

export const Shadowfiend = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecSmitePriest>({
	fieldName: 'useShadowfiend',
	id: ActionId.fromSpellId(34433),
});

export const SmitePriestRotationConfig = {
	inputs: [
		InputHelpers.makeRotationBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'useDevouringPlague',
			label: 'Use Devouring Plague',
			labelTooltip: 'Use Devouring Plague whenever its not active.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'useShadowWordDeath',
			label: 'Use Shadow Word: Death',
			labelTooltip: 'Use Shadow Word: Death whenever off CD.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'useMindBlast',
			label: 'Use Mind Blast',
			labelTooltip: 'Use Mind Blast whenever off CD.',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecSmitePriest>({
			fieldName: 'memeDream',
			label: 'Meme Dream',
			labelTooltip: 'Assumes 2nd Smite Priest in raid, so just spams HF + Smite with permanent HF uptime.',
			extraCssClasses: ['within-raid-sim-hide'],
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecSmitePriest>({
			fieldName: 'allowedHolyFireDelayMs',
			label: 'Allowed Delay for HF',
			labelTooltip: 'Time, in milliseconds, the player is allowed to wait for Holy Fire if it is about to come off CD.',
		}),
	],
};
