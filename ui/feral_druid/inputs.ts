import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { APLRotation_Type } from '../core/proto/apl.js';
import { Spec,UnitReference, UnitReference_Type as UnitType  } from '../core/proto/common.js';
import {
	FeralDruid_Rotation_AplType as AplType,
	FeralDruid_Rotation_BiteModeType as BiteModeType,
} from '../core/proto/druid.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const SelfInnervate = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecFeralDruid>({
	fieldName: 'innervateTarget',
	id: ActionId.fromSpellId(29166),
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
	label: '网络延迟',
	labelTooltip: '玩家延迟，以毫秒为单位。会给无法法术排队的动作增加延迟。',
});

export const AssumeBleedActive = InputHelpers.makeSpecOptionsBooleanInput<Spec.SpecFeralDruid>({
	fieldName: 'assumeBleedActive',
	label: '假设流血始终存在',
	labelTooltip: '假设流血始终存在于“狂乱撕扯”天赋的计算中。否则，将仅基于自己的割伤/斜掠/割裂进行计算。',
	extraCssClasses: ['within-raid-sim-hide'],
});

function ShouldShowAdvParamST(player: Player<Spec.SpecFeralDruid>): boolean {
	const rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.SingleTarget;
}

function ShouldShowAdvParamAoe(player: Player<Spec.SpecFeralDruid>): boolean {
	const rot = player.getSimpleRotation();
	return rot.manualParams && rot.rotationType == AplType.Aoe;
}

export const FeralDruidRotationConfig = {
	inputs: [
		InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, AplType>({
			fieldName: 'rotationType',
			label: '类型',
			values: [
				{ name: '单体木桩', value: AplType.SingleTarget },
				{ name: '群体AOE', value: AplType.Aoe },
			],
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'prePopOoc',
			label: '开怪前用爪子舞触发清晰预兆',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().omenOfClarity,
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'prePopBerserk',
			label: '开怪前用狂暴',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getTalents().berserk,
			changeEmitter: (player: Player<Spec.SpecFeralDruid>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'manualParams',
			label: '进阶参数调整',
			labelTooltip: '修改割裂/咆哮/撕咬等参数,结果可以应用到WA和Hekili插件里以优化输出循环',
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'maxFfDelay',
			label: '最大精灵活延迟',
			labelTooltip: '精灵火CD到了后能允许最多的间歇时间,一般情况下我们希望卡CD施放',
			float: true,
			positive: true,
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().manualParams,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'minRoarOffset',
			label: '咆哮割裂差值(Offset)',
			labelTooltip: '指的是一个通过大样本模拟计算出来的最佳常数.当这个常数计算在条件里时,他能有效的最大化DPS且保证覆盖的有效性.过晚的覆盖咆哮会造成同步,过早的覆盖会造成DPS的损失.因此一个正确的Offset值会达到最平衡的效果,而这个值会根据阶段和装备不同有可能不同.',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'ripLeeway',
			label: '割裂延迟(Leeway)',
			labelTooltip: '指的是,比如说当你补了一个割裂,但你的能量来不及回复到足够你多打一个星然后补野蛮咆哮.因此Leeway作为一个常数是需要计算在里面以保证前后两个技能之间是有缓冲空间进行下一步决策.',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useRake',
			label: '加入斜掠',
			labelTooltip: '在某些配装下,斜掠不一定会带来正面收益',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'useBite',
			label: '加入凶猛撕咬',
			showWhen: ShouldShowAdvParamST,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'biteTime',
			label: '撕咬常数(Bite Rule)',
			labelTooltip: '是通过大样本计算出来在副本环境里最佳的一个咆哮/割裂剩余时间条件值,一般是4或者5.这个数字意味着"当我当前割裂/咆哮还剩大于X秒的时候我就可以打一个撕咬,这样能保证我后面的循环不被(大幅度的)影响"',
			showWhen: (player: Player<Spec.SpecFeralDruid>) =>
				ShouldShowAdvParamST(player) && player.getSimpleRotation().useBite == true && player.getSimpleRotation().biteModeType == BiteModeType.Emperical,
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecFeralDruid>({
			fieldName: 'flowerWeave',
			label: '使用爪子舞',
			labelTooltip: '在空能且精灵火CD时进行爪子舞来获取额外的AOE资源',
			showWhen: ShouldShowAdvParamAoe,
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			extraCssClasses: ['used-in-apl'],
			fieldName: 'raidTargets',
			label: '爪子舞时团队友方目标数量',
			labelTooltip: '团队人数,一个25人团队可以约等于30个目标(包括召唤和宠物)',
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.aplRotation.type != APLRotation_Type.TypeSimple || (ShouldShowAdvParamAoe(player) && player.getSimpleRotation().flowerWeave == true),
		}),
		// Can be uncommented if/when analytical bite mode is added
		//InputHelpers.makeRotationEnumInput<Spec.SpecFeralDruid, BiteModeType>({
		//	fieldName: 'biteModeType',
		//	label: 'Bite Mode',
		//	labelTooltip: 'Underlying "Bite logic" to use',
		//	values: [
		//		{ name: 'Emperical', value: BiteModeType.Emperical },
		//	],
		//	showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().useBite == true
		//}),
		InputHelpers.makeRotationNumberInput<Spec.SpecFeralDruid>({
			fieldName: 'hotUptime',
			label: 'Revitalize Hot Uptime',
			labelTooltip: 'Hot uptime percentage to assume when theorizing energy gains',
			percent: true,
			showWhen: (player: Player<Spec.SpecFeralDruid>) => player.getSimpleRotation().useBite == true && player.getSimpleRotation().biteModeType == BiteModeType.Analytical,
		}),
	],
};
