import * as InputHelpers from '../core/components/input_helpers.js';
import { Spec } from '../core/proto/common.js';
import {
	ShamanImbue,
	ShamanShield,
	ShamanSyncType,
} from '../core/proto/shaman.js';
import { ActionId } from '../core/proto_utils/action_id.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const ShamanShieldInput = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanShield>({
	fieldName: 'shield',
	values: [
		{ value: ShamanShield.NoShield, tooltip: 'No Shield' },
		{ actionId: ActionId.fromSpellId(57960), value: ShamanShield.WaterShield },
		{ actionId: ActionId.fromSpellId(49281), value: ShamanShield.LightningShield },
	],
});

export const ShamanImbueMH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueMh',
	values: [
		{ value: ShamanImbue.NoImbue, tooltip: 'No Main Hand Enchant' },
		{ actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon, text: 'R10' },
		{ actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank, text: 'R9' },
		{ actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
	],
});

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueOh',
	values: [
		{ value: ShamanImbue.NoImbue, tooltip: 'No Off Hand Enchant' },
		{ actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon, text: 'R10' },
		{ actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank, text: 'R9' },
		{ actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
	],
});

export const SyncTypeInput = InputHelpers.makeSpecOptionsEnumInput<Spec.SpecEnhancementShaman, ShamanSyncType>({
	fieldName: 'syncType',
	label: '攻速同步设置',
	labelTooltip:
		`选择你的同步或错开选项：
		<ul>
			<li><div>自动：将根据你的武器攻击速度自动选择同步选项</div></li>
			<li><div>无：不进行同步或错开，适用于武器速度不匹配的情况</div></li>
			<li><div>完美同步：使主手和副手武器总是同时攻击，适用于匹配的武器速度</div></li>
			<li><div>副手延迟：在保持0.5秒怒火急速触发窗口内的同时，为副手攻击增加一个轻微的延迟</div></li>
		</ul>`,
	values: [
		{ name: "自动", value: ShamanSyncType.Auto },
		{ name: '无', value: ShamanSyncType.NoSync },
		{ name: '完美同步', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: '副手延迟', value: ShamanSyncType.DelayOffhandSwings },
	],
});
