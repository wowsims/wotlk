import { IconEnumPicker } from '../components/icon_enum_picker.js';
import { IconPicker } from '../components/icon_picker.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
} from '../proto/shaman.js';
import { Spec } from '../proto/common.js';
import { ActionId } from '../proto_utils/action_id.js';
import { Player } from '../player.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { EventID } from '../typed_event.js';
import * as InputHelpers from '../components/input_helpers.js';
import { ShamanSpecs } from '../proto_utils/utils.js';
import { ContentBlock } from './content_block.js';
import { Input } from './input.js';

export function TotemsSection(parentElem: HTMLElement, simUI: IndividualSimUI<ShamanSpecs>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'totems-settings', {
		header: { title: 'Totems' }
	});

	let totemDropdownGroup = Input.newGroupContainer();
	totemDropdownGroup.classList.add('totem-dropdowns-container', 'icon-group');

	let fireElementalContainer = document.createElement('div');
	fireElementalContainer.classList.add('fire-elemental-input-container');

	contentBlock.bodyElement.appendChild(totemDropdownGroup);
	contentBlock.bodyElement.appendChild(fireElementalContainer);

	const earthTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'earth-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#ffdfba', value: EarthTotem.NoEarthTotem },
			{ actionId: ActionId.fromSpellId(58643), value: EarthTotem.StrengthOfEarthTotem },
			{ actionId: ActionId.fromSpellId(58753), value: EarthTotem.StoneskinTotem },
			{ actionId: ActionId.fromSpellId(8143), value: EarthTotem.TremorTotem },
		],
		equals: (a: EarthTotem, b: EarthTotem) => a == b,
		zeroValue: EarthTotem.NoEarthTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.earth || EarthTotem.NoEarthTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.earth = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const waterTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'water-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#bae1ff', value: WaterTotem.NoWaterTotem },
			{ actionId: ActionId.fromSpellId(58774), value: WaterTotem.ManaSpringTotem },
			{ actionId: ActionId.fromSpellId(58757), value: WaterTotem.HealingStreamTotem },
		],
		equals: (a: WaterTotem, b: WaterTotem) => a == b,
		zeroValue: WaterTotem.NoWaterTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.water || WaterTotem.NoWaterTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.water = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const fireTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'fire-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#ffb3ba', value: FireTotem.NoFireTotem },
			{ actionId: ActionId.fromSpellId(58734), value: FireTotem.MagmaTotem },
			{ actionId: ActionId.fromSpellId(58704), value: FireTotem.SearingTotem },
			{ actionId: ActionId.fromSpellId(57722), value: FireTotem.TotemOfWrath, showWhen: (player: Player<ShamanSpecs>) => player.getTalents().totemOfWrath },
			{ actionId: ActionId.fromSpellId(58656), value: FireTotem.FlametongueTotem },
		],
		equals: (a: FireTotem, b: FireTotem) => a == b,
		zeroValue: FireTotem.NoFireTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.fire || FireTotem.NoFireTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.fire = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	const airTotemPicker = new IconEnumPicker(totemDropdownGroup, simUI.player, {
		extraCssClasses: [
			'air-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#baffc9', value: AirTotem.NoAirTotem },
			{ actionId: ActionId.fromSpellId(8512), value: AirTotem.WindfuryTotem },
			{ actionId: ActionId.fromSpellId(3738), value: AirTotem.WrathOfAirTotem },
		],
		equals: (a: AirTotem, b: AirTotem) => a == b,
		zeroValue: AirTotem.NoAirTotem,
		changedEvent: (player: Player<ShamanSpecs>) => player.specOptionsChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems?.air || AirTotem.NoAirTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newOptions = player.getSpecOptions();
			if (!newOptions.totems)
				newOptions.totems = ShamanTotems.create();
			newOptions.totems!.air = newValue;
			player.setSpecOptions(eventID, newOptions);
		},
	});

	// Enchancement Shaman uses the Fire Elemental Inputs with custom inputs.
	if (simUI.player.spec != Spec.SpecEnhancementShaman) {
		const fireElementalBooleanIconInput = InputHelpers.makeBooleanIconInput<ShamanSpecs, ShamanTotems, Player<ShamanSpecs>>({
			getModObject: (player: Player<ShamanSpecs>) => player,
			getValue: (player: Player<ShamanSpecs>) => player.getSpecOptions().totems || ShamanTotems.create(),
			setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
				const newOptions = player.getSpecOptions();
				newOptions.totems = newVal;
				player.setSpecOptions(eventID, newOptions);
			},
			changeEmitter: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,
		}, ActionId.fromSpellId(2894), "useFireElemental");

		new IconPicker(fireElementalContainer, simUI.player, fireElementalBooleanIconInput);
	}

	return contentBlock;
}
