import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import {
	Rogue_Options_PoisonImbue as Poison,
} from '../core/proto/rogue.js';
import { ActionId } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const MainHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'mhImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Main Hand Poison' },
		{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
		{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
		{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
	],
});

export const OffHandImbue = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecRogue, Poison>({
	fieldName: 'ohImbue',
	numColumns: 1,
	values: [
		{ value: Poison.NoPoison, tooltip: 'No Off Hand Poison' },
		{ actionId: ActionId.fromItemId(43233), value: Poison.DeadlyPoison },
		{ actionId: ActionId.fromItemId(43231), value: Poison.InstantPoison },
		{ actionId: ActionId.fromItemId(43235), value: Poison.WoundPoison },
	],
});

export const StartingOverkillDuration = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'startingOverkillDuration',
	label: '起始灭绝持续时间',
	labelTooltip: '每次战斗开始时的初始灭绝增益持续时间。',
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().overkill || player.getTalents().masterOfSubtlety > 0
});

export const VanishBreakTime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'vanishBreakTime',
	label: '消失打断时间',
	labelTooltip: '施放消失后开始攻击所需的时间。',
	extraCssClasses: ['experimental'],
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().overkill || player.getTalents().masterOfSubtlety > 0
});

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'assumeBleedActive',
	label: '假设流血始终存在',
	labelTooltip: '假设\'血之饥渴\'激活时流血效果始终存在。否则仅根据自身的锁喉/割裂计算。',
	extraCssClasses: ['within-raid-sim-hide'],
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().hungerForBlood
});

export const HonorOfThievesCritRate = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecRogue>({
	fieldName: 'honorOfThievesCritRate',
	label: '盗贼的尊严暴击率',
	labelTooltip: '其他组员在100秒内产生的暴击次数',
	showWhen: (player: Player<Spec.SpecRogue>) => player.getTalents().honorAmongThieves > 0
});

export const ApplyPoisonsManually = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecRogue>({
	fieldName: 'applyPoisonsManually',
	label: '手动配置毒药',
	labelTooltip: '防止基于已装备武器的自动毒药配置。',
});

