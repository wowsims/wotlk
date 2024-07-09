import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import {
	Hunter_Options_Ammo as Ammo,
	Hunter_Rotation_RotationType as RotationType,
	Hunter_Rotation_StingType as StingType,
} from '../core/proto/hunter.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { makePetTypeInputConfig } from '../core/talents/hunter_pet.js';
import { TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const WeaponAmmo = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecHunter, Ammo>({
	fieldName: 'ammo',
	numColumns: 2,
	values: [
		{ value: Ammo.AmmoNone, tooltip: 'No Ammo' },
		{ actionId: ActionId.fromItemId(52021), value: Ammo.IcebladeArrow },
		{ actionId: ActionId.fromItemId(41165), value: Ammo.SaroniteRazorheads },
		{ actionId: ActionId.fromItemId(41586), value: Ammo.TerrorshaftArrow },
		{ actionId: ActionId.fromItemId(31737), value: Ammo.TimelessArrow },
		{ actionId: ActionId.fromItemId(34581), value: Ammo.MysteriousArrow },
		{ actionId: ActionId.fromItemId(33803), value: Ammo.AdamantiteStinger },
		{ actionId: ActionId.fromItemId(28056), value: Ammo.BlackflightArrow },
	],
});

export const PetTypeInput = makePetTypeInputConfig();

export const PetUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'petUptime',
	label: '宠物存活时间(%)',
	labelTooltip: '宠物在战斗中存活的时间百分比。',
	percent: true,
});

export const UseHuntersMark = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecHunter>({
	fieldName: 'useHuntersMark',
	id: ActionId.fromSpellId(53338),
});

export const SniperTrainingUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'sniperTrainingUptime',
	label: '狙击训练持续时间(%)',
	labelTooltip: '狙击训练天赋在战斗持续时间中的持续时间百分比。',
	percent: true,
	showWhen: (player: Player<Spec.SpecHunter>) => player.getTalents().sniperTraining > 0,
	changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
});

export const TimeToTrapWeaveMs = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecHunter>({
	fieldName: 'timeToTrapWeaveMs',
	label: '陷阱舞移动时间',
	labelTooltip: '在向老板移动并重新进行远程自动攻击之间，用于爆炸陷阱的时间，以毫秒为单位。',
});


export const HunterRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter, RotationType>({
			fieldName: 'type',
			label: '类型',
			values: [
				{ name: '单体木桩', value: RotationType.SingleTarget },
				{ name: '群体AOE', value: RotationType.Aoe },
			],
		}),
		InputHelpers.makeRotationEnumInput<Spec.SpecHunter, StingType>({
			fieldName: 'sting',
			label: '钉刺',
			labelTooltip: '在主要目标上保持选定的钉刺。',
			values: [
				{ name: '无', value: StingType.NoSting },
				{ name: '毒蝎钉刺', value: StingType.ScorpidSting },
				{ name: '毒蛇钉刺', value: StingType.SerpentSting },
			],
			showWhen: (player: Player<Spec.SpecHunter>) => player.getSimpleRotation().type == RotationType.SingleTarget,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'trapWeave',
			label: '陷阱舞',
			labelTooltip: '在适当的时间使用爆炸陷阱。请注意，选择此选项将禁用黑箭，因为它们共享冷却时间。',
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'allowExplosiveShotDownrank',
			label: '允许降级爆炸射击',
			labelTooltip: '在锁定和装载触发期间陷阱舞爆炸射击等级3。这是可行的,因为等级3和等级4的点数可以叠加。',
			showWhen: (player: Player<Spec.SpecHunter>) => player.getSimpleRotation().type != RotationType.Custom && player.getTalents().explosiveShot && player.getTalents().lockAndLoad > 0,
			changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecHunter>({
			fieldName: 'multiDotSerpentSting',
			label: '多目标毒蛇钉刺',
			labelTooltip: '对多个目标施放毒蛇钉刺',
			changeEmitter: (player: Player<Spec.SpecHunter>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'viperStartManaPercent',
			label: '毒蛇开始法力百分比',
			labelTooltip: '当法力值低于此数值时切换到毒蛇姿态。',
			percent: true,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecHunter>({
			fieldName: 'viperStopManaPercent',
			label: '毒蛇停止法力百分比',
			labelTooltip: '当法力值高于此数值时切换回雄鹰姿态。',
			percent: true,
		}),
	],
};
