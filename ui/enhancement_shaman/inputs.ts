import { BooleanPicker } from '../core/components/boolean_picker.js';
import { EnumPicker } from '../core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../core/components/icon_enum_picker.js';
import { IconPickerConfig } from '../core/components/icon_picker.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	EnhancementShaman_Options as ShamanOptions,
	ShamanTotems,
	ShamanShield,
	ShamanImbue,
	ShamanSyncType,
	EnhancementShaman_Rotation_PrimaryShock as PrimaryShock,
	EnhancementShaman_Rotation_RotationType as RotationType,
	EnhancementShaman_Rotation_CustomRotationSpell as CustomRotationSpell
} from '../core/proto/shaman.js';
import { CustomSpell, Spec, ItemSwap, ItemSlot } from '../core/proto/common.js';
import { ActionId } from '../core/proto_utils/action_id.js';
import { Player } from '../core/player.js';
import { Sim } from '../core/sim.js';
import { IndividualSimUI } from '../core/individual_sim_ui.js';
import { Target } from '../core/target.js';
import { EventID, TypedEvent } from '../core/typed_event.js';

import * as InputHelpers from '../core/components/input_helpers.js';

// Configuration for spec-specific UI elements on the settings tab.
// These don't need to be in a separate file but it keeps things cleaner.

export const Bloodlust = InputHelpers.makeSpecOptionsBooleanIconInput<Spec.SpecEnhancementShaman>({
	fieldName: 'bloodlust',
	id: ActionId.fromSpellId(2825),
});
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
		{ actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon, text: 'R10'},
		{ actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank, text: 'R9'},
		{ actionId: ActionId.fromSpellId(58796), value: ShamanImbue.FrostbrandWeapon },
	],
});

export const ShamanImbueOH = InputHelpers.makeSpecOptionsEnumIconInput<Spec.SpecEnhancementShaman, ShamanImbue>({
	fieldName: 'imbueOh',
	values: [
		{ value: ShamanImbue.NoImbue, tooltip: 'No Off Hand Enchant' },
		{ actionId: ActionId.fromSpellId(58804), value: ShamanImbue.WindfuryWeapon },
		{ actionId: ActionId.fromSpellId(58790), value: ShamanImbue.FlametongueWeapon, text: 'R10'},
		{ actionId: ActionId.fromSpellId(58789), value: ShamanImbue.FlametongueWeaponDownrank, text: 'R9'},
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
		{ name: "Automatic", value: ShamanSyncType.Auto},
		{ name: 'None', value: ShamanSyncType.NoSync },
		{ name: 'Perfect Sync', value: ShamanSyncType.SyncMainhandOffhandSwings },
		{ name: 'Delayed Offhand', value: ShamanSyncType.DelayOffhandSwings },
	],
});

export const EnhancmentItemSwapInputs = InputHelpers.MakeItemSwapInput<Spec.SpecEnhancementShaman>({
	fieldName: 'itemSwap',
	values: [
		ItemSlot.ItemSlotMainHand,
		ItemSlot.ItemSlotOffHand,
		//ItemSlot.ItemSlotRanged, Not support yet
	],
	labelTooltip: 'Start with the swapped items until Fire Elemntal has been summoned, swap back to normal gear set. Weapons come pre enchanted with FT9 and FT10. If a slot is empty it will not be used in the swap',
	showWhen: (player: Player<Spec.SpecEnhancementShaman>) => (player.getRotation().totems?.useFireElemental && player.getRotation().enableItemSwap) || false
})

export const EnhancementShamanRotationConfig = {
	inputs:
		[	
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
				fieldName: 'enableItemSwap',
				label: 'Enable Item Swapping',
				labelTooltip: 'Toggle on/off item swapping',
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) =>  player.getRotation().totems?.useFireElemental || false
			}),
			EnhancmentItemSwapInputs,
			InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, RotationType>({
				fieldName: 'rotationType',
				label: 'Type',
				labelTooltip:
					`<ul>
					<li>
						<div>Standard: Priority Rotation</div>
					</li>
					<li>
						<div>Custom: Highest spell that is ready will be cast.</div>
					</li>
				</ul>`,
				values: [
					{ name: 'Standard', value: RotationType.Priority },
					{ name: 'Custom', value: RotationType.Custom },
				],
			}),
			InputHelpers.makeCustomRotationInput<Spec.SpecEnhancementShaman, CustomRotationSpell>({
				fieldName: 'customRotation',
				numColumns: 2,
				values: [
					{ actionId: ActionId.fromSpellId(49238), value: CustomRotationSpell.LightningBolt},
					{ actionId: ActionId.fromSpellId(49238), value: CustomRotationSpell.LightningBoltWeave, text: "Weave" },
					{ actionId: ActionId.fromSpellId(49238), value: CustomRotationSpell.LightningBoltDelayedWeave, text: "Delay" },
					{ actionId: ActionId.fromSpellId(17364), value: CustomRotationSpell.StormstrikeDebuffMissing, text: "Debuff"  },
					{ actionId: ActionId.fromSpellId(17364), value: CustomRotationSpell.Stormstrike },
					{ actionId: ActionId.fromSpellId(49233), value: CustomRotationSpell.FlameShock },
					{ actionId: ActionId.fromSpellId(49231), value: CustomRotationSpell.EarthShock },
					{ actionId: ActionId.fromSpellId(58734), value: CustomRotationSpell.MagmaTotem },
					{ actionId: ActionId.fromSpellId(61657), value: CustomRotationSpell.FireNova },
					{ actionId: ActionId.fromSpellId(60103), value: CustomRotationSpell.LavaLash },
					{ actionId: ActionId.fromSpellId(49281), value: CustomRotationSpell.LightningShield },
					{ actionId: ActionId.fromSpellId(60043), value: CustomRotationSpell.LavaBurst, text: "Weave" },
					{ actionId: ActionId.fromSpellId(49236), value: CustomRotationSpell.FrostShock},
				],
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType == RotationType.Custom,
			}),
			InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, PrimaryShock>({
				fieldName: 'primaryShock',
				label: 'Primary Shock',
				values: [
					{ name: 'None', value: PrimaryShock.None },
					{ name: 'Earth Shock', value: PrimaryShock.Earth },
					{ name: 'Frost Shock', value: PrimaryShock.Frost },
				],
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
				fieldName: 'weaveFlameShock',
				label: 'Weave Flame Shock',
				labelTooltip: 'Use Flame Shock whenever the target does not already have the DoT.',
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
            InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
                fieldName:  'flameShockClipTicks',
                label:  'Refresh Flame Shock at ticks remaining',
                labelTooltip: 'Set to 0 to require the debuff be missing. A tick is 3s, affected by spell haste',
                enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().weaveFlameShock,
                showWhen: (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.FlameShock) != undefined
					}

					return player.getRotation().weaveFlameShock
				}
            }),
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({
				fieldName: 'lightningboltWeave',
				label: 'Enable Weaving Lightning Bolt',
				labelTooltip: 'Will provide a DPS increase, but is harder to execute',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().maelstromWeapon > 0,
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
			InputHelpers.makeRotationEnumInput<Spec.SpecEnhancementShaman, number>({
				fieldName: 'maelstromweaponMinStack',
				label: 'Minimum Maelstrom Stacks to Weave',
				labelTooltip: '3 stacks is the most realistic, however there are cases where lower might be possible, just much harder to do in practice',
				values: [
					{ name: '1', value: 1 },
					{ name: '2', value: 2 },
					{ name: '3', value: 3 },
					{ name: '4', value: 4 },
				],
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.LightningBoltWeave) != undefined
					}

					return player.getRotation().lightningboltWeave
				},
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.LightningBoltWeave) != undefined
					}

					return true
				}
			}),		
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'autoWeaveDelay',
				label: 'Weaving Delay After Auto Attack',
				labelTooltip: 'The amount of time to wait after an auto attack before weaveing, in milliseconds',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.LightningBoltWeave) != undefined
					}

					return player.getRotation().lightningboltWeave
				},
				showWhen:  (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.LightningBoltWeave) != undefined
					}

					return true
				},
			}),InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'delayGcdWeave',
				label: 'Delay LL to Weave',
				labelTooltip: 'The amount of time to hold Lava Lash to weave in milliseconds. Setting to 0 will disable delaying',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return false
					}

					return player.getRotation().lightningboltWeave
				},
				showWhen:  (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return false
					}

					return true
				},
			}),
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'delayGcdWeave',
				label: 'Delay Weave Time',
				labelTooltip: 'The amount of time to hold a GCD to weave in milliseconds. Setting to 0 will disable delaying',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.LightningBoltDelayedWeave) != undefined
					}

					return false
				},
				showWhen:  (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.LightningBoltDelayedWeave) != undefined
					}

					return false
				},
			}),
			InputHelpers.makeRotationBooleanInput<Spec.SpecEnhancementShaman>({ 
				fieldName: 'lavaburstWeave',
				label: 'Enable Weaving Lava Burst',
				labelTooltip: 'Not particularily useful for dual wield, mostly a 2h option',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().lightningboltWeave,
				showWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getRotation().rotationType != RotationType.Custom
			}),
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'firenovaManaThreshold',
				label: 'Minimum mana to cast Fire Nova',
				labelTooltip: 'Fire Nova will not be cast when mana is below this value. Set this medium-low, it has a bad mana-to-damage ratio',
				showWhen:  (player: Player<Spec.SpecEnhancementShaman>) => {
					if (player.getRotation().rotationType == RotationType.Custom){
						return player.getRotation().customRotation?.spells.find(customSpell => customSpell.spell == CustomRotationSpell.FireNova) != undefined
					}

					return true
				},
			}),
			InputHelpers.makeRotationNumberInput<Spec.SpecEnhancementShaman>({
				fieldName: 'shamanisticRageManaThreshold',
				label: 'Mana % to use Shamanistic Rage',
				enableWhen: (player: Player<Spec.SpecEnhancementShaman>) => player.getTalents().shamanisticRage,
			}),
		],
};


