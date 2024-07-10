import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { Spec,UnitReference, UnitReference_Type as UnitType  } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { EventID } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecBalanceDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
	extraCssClasses: [
		'within-raid-sim-hide',
	],
	getValue: (player: Player<Spec.SpecBalanceDruid>) => player.getSpecOptions().innervateTarget?.type == UnitType.Player,
	setValue: (eventID: EventID, player: Player<Spec.SpecBalanceDruid>, newValue: boolean) => {
		const newOptions = player.getSpecOptions();
		newOptions.innervateTarget = UnitReference.create({
			type: newValue ? UnitType.Player : UnitType.Unknown,
			index: 0,
		});
		player.setSpecOptions(eventID, newOptions);
	},
});

export const OkfUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecBalanceDruid>({
	fieldName: 'okfUptime',
	label: '枭兽狂暴持续时间 (%)',
	labelTooltip: '枭兽狂暴在战斗中的持续时间百分比',
	percent: true,
});
