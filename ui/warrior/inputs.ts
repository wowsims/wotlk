import * as InputHelpers from '../core/components/input_helpers.js';
import { Spec } from '../core/proto/common.js';
import {
	WarriorShout,
} from '../core/proto/warrior.js';
import { ActionId } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Recklessness = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useRecklessness',
	id: ActionId.fromSpellId(1719),
});

export const ShatteringThrow = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecWarrior>({
	fieldName: 'useShatteringThrow',
	id: ActionId.fromSpellId(64382),
});

export const StartingRage = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecWarrior>({
	fieldName: 'startingRage',
	label: '初始怒气',
	labelTooltip: '每次战斗开始时的初始怒气值。',
});

export const StanceSnapshot = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'stanceSnapshot',
	label: '姿态快照',
	labelTooltip: '在切换姿态时同时施放的技能将受益于切换前的姿态加成。',
});

// 允许在自动镶嵌宝石时忽略精准上限
// （对武器战有用）
export const DisableExpertiseGemming = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecWarrior>({
	fieldName: 'disableExpertiseGemming',
	label: '禁用精准宝石',
	labelTooltip: '禁用精准属性的自动镶嵌。',
});

export const ShoutPicker = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecWarrior, WarriorShout>({
	fieldName: 'shout',
	values: [
		{ color: 'c79c6e', value: WarriorShout.WarriorShoutNone },
		{ actionId: ActionId.fromSpellId(2048), value: WarriorShout.WarriorShoutBattle },
		{ actionId: ActionId.fromSpellId(469), value: WarriorShout.WarriorShoutCommanding },
	],
});
