import {
	ShamanShield,
	ShamanImbue,
	ShamanSyncType,
} from '../core/proto/shaman.js';
import { Spec } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';

import * as InputHelpers from '../core/components/input_helpers.js';

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
	label: 'Sync/Stagger Setting',
	labelTooltip:
		`Choose your sync or stagger option Perfect
		<ul>
			<li><div>Auto: Will auto pick sync options based on your weapons attack speeds</div></li>
			<li><div>None: No Sync or Staggering, used for mismatched weapon speeds</div></li>
			<li><div>Perfect Sync: Makes your weapons always attack at the same time, for match weapon speeds</div></li>
			<li><div>Delayed Offhand: Adds a slight delay to the offhand attacks while staying within the 0.5s flurry ICD window</div></li>
		</ul>`,
	values: [
		{ name: "Automatic", value: ShamanSyncType.Auto },
		{ name: 'None', value: ShamanSyncType.NoSync },
		{ name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
	],
});
