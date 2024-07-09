import * as InputHelpers from '../core/components/input_helpers.js';
import { Player } from '../core/player.js';
import { Spec } from '../core/proto/common.js';
import {
	Mage_Options_ArmorType as ArmorType,
	Mage_Rotation_PrimaryFireSpell as PrimaryFireSpell,
} from '../core/proto/mage.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { TypedEvent } from '../core/typed_event.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Armor = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecMage, ArmorType>({
	fieldName: 'armor',
	values: [
		{ value: ArmorType.NoArmor, tooltip: '无护甲' },
		{ actionId: ActionId.fromSpellId(43024), value: ArmorType.MageArmor, tooltip: '法师护甲' },
		{ actionId: ActionId.fromSpellId(43046), value: ArmorType.MoltenArmor, tooltip: '熔岩护甲' },
	],
});

export const WaterElementalDisobeyChance = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'waterElementalDisobeyChance',
	percent: true,
	label: '水元素不听话率 (%)',
	labelTooltip: '水元素行动失败的百分比。这代表了水元素移动或站立不动而不施法的概率。',
	changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.specOptionsChangeEmitter, player.talentsChangeEmitter]),
	showWhen: (player: Player<Spec.SpecMage>) => player.getTalents().summonWaterElemental,
});

export const FocusMagicUptime = InputHelpers.makeSpecOptionsNumberInput<Spec.SpecMage>({
	fieldName: 'focusMagicPercentUptime',
	label: '魔法专注持续时间 (%)',
	labelTooltip: '魔法专注的持续时间百分比',
	extraCssClasses: ['within-raid-sim-hide'],
});

export const MageRotationConfig = {
	inputs: [
		// ********************************************************
		//                       火焰输入
		// ********************************************************
		InputHelpers.makeRotationEnumInput<Spec.SpecMage, PrimaryFireSpell>({
			fieldName: 'primaryFireSpell',
			label: '主要法术',
			values: [
				{ name: '火球术', value: PrimaryFireSpell.Fireball },
				{ name: '霜火箭', value: PrimaryFireSpell.FrostfireBolt },
				{ name: '灼烧', value: PrimaryFireSpell.Scorch },
			],
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 1,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		// ********************************************************
		//                       冰霜输入
		// ********************************************************
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'useIceLance',
			label: '使用冰枪术',
			labelTooltip: '在寒冰指效果结束时施放冰枪术，之后使用深度冻结。',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 2,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		// ********************************************************
		//                      奥术输入
		// ********************************************************
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'only3ArcaneBlastStacksBelowManaPercent',
			percent: true,
			label: '法力值低于 % 时奥术冲击堆叠到3层',
			labelTooltip: '当法力值低于该百分比时，在3层奥术冲击时使用奥术飞弹/奥术弹幕，而不是4层。',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'blastWithoutMissileBarrageAboveManaPercent',
			percent: true,
			label: '法力值高于 % 时无导弹弹幕奥术冲击',
			labelTooltip: '当法力值高于该百分比时，连续施放奥术冲击直到出现导弹弹幕效果。',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationNumberInput<Spec.SpecMage>({
			fieldName: 'missileBarrageBelowManaPercent',
			percent: true,
			label: '法力值低于 % 时尽快使用导弹弹幕',
			labelTooltip: '当法力值低于该百分比时，尽快使用导弹弹幕效果。这有助于节约法力。',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'useArcaneBarrage',
			label: '使用奥术弹幕',
			labelTooltip: '在旋转中包括奥术弹幕。',
			enableWhen: (player: Player<Spec.SpecMage>) => player.getTalents().arcaneBarrage,
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalentTree() == 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),

		InputHelpers.makeRotationBooleanInput<Spec.SpecMage>({
			fieldName: 'maintainImprovedScorch',
			label: '保持强化灼烧',
			labelTooltip: '当灼烧层数少于5层，或减益效果剩余时间少于4秒时，始终使用灼烧。',
			showWhen: (player: Player<Spec.SpecMage>) => player.getTalents().improvedScorch > 0,
			changeEmitter: (player: Player<Spec.SpecMage>) => TypedEvent.onAny([player.rotationChangeEmitter, player.talentsChangeEmitter]),
		}),
	],
};
