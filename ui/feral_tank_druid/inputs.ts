import * as InputHelpers from '../core/components/input_helpers.js';
import { Spec } from '../core/proto/common.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecFeralTankDruid>({
	fieldName: 'startingRage',
	label: '初始怒气',
	labelTooltip: '每次战斗开始时的初始怒气值。',
});

export const FeralTankDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralTankDruid>({
			fieldName: 'maulRageThreshold',
			label: '重殴怒气阈值',
			labelTooltip: '当怒气高于此值时才施放重殴。',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralTankDruid>({
			fieldName: 'lacerateTime',
			label: '割伤刷新宽限时间',
			labelTooltip: '当割伤的剩余持续时间少于此值时刷新割伤（以秒为单位）。',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralTankDruid>({
			fieldName: 'maintainDemoralizingRoar',
			label: '保持挫志咆哮',
			labelTooltip: '在主要目标上保持挫志咆哮。如果已有更强的减益效果，则不会施放。',
		}),
	],
};
