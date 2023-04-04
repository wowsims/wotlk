import { BooleanPicker } from '../components/boolean_picker.js';
import { EnumPicker } from '../components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../components/icon_enum_picker.js';
import { IconPicker } from '../components/icon_picker.js';
import { IconPickerConfig } from '../components/icon_picker.js';
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
import { Sim } from '../sim.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import * as InputHelpers from '../components/input_helpers.js';
import { SpecRotation, ShamanSpecs } from '../proto_utils/utils.js';
import { ContentBlock } from './content_block.js';
import { Input } from './input.js';

export function TotemsSection(parentElem: HTMLElement, simUI: IndividualSimUI<ShamanSpecs>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'totems-settings', {
		header: {title: 'Totems'}
	});

	let totemDropdownGroup = Input.newGroupContainer();
	totemDropdownGroup.classList.add('totem-dropdowns-container', 'icon-group');

	let fireElementalContainer = document.createElement('div');
	fireElementalContainer.classList.add('fire-elemental-inputs-container');

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
		changedEvent: (player: Player<ShamanSpecs>) => player.rotationChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems?.earth || EarthTotem.NoEarthTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.earth = newValue;
			player.setRotation(eventID, newRotation);
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
		changedEvent: (player: Player<ShamanSpecs>) => player.rotationChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems?.water || WaterTotem.NoWaterTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.water = newValue;
			player.setRotation(eventID, newRotation);
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
		changedEvent: (player: Player<ShamanSpecs>) => player.rotationChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems?.fire || FireTotem.NoFireTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.fire = newValue;
			player.setRotation(eventID, newRotation);
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
		changedEvent: (player: Player<ShamanSpecs>) => player.rotationChangeEmitter,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems?.air || AirTotem.NoAirTotem,
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.air = newValue;
			player.setRotation(eventID, newRotation);
		},
	});

	const fireElementalBooleanIconInput = InputHelpers.makeBooleanIconInput<ShamanSpecs, ShamanTotems, Player<ShamanSpecs>>({
		getModObject: (player: Player<ShamanSpecs>) => player,
		getValue: (player: Player<ShamanSpecs>) => player.getRotation().totems || ShamanTotems.create(),
		setValue: (eventID: EventID, player: Player<ShamanSpecs>, newVal: ShamanTotems) => {
			const newRotation = player.getRotation();
			newRotation.totems = newVal;
			player.setRotation(eventID, newRotation);
		},
		changeEmitter: (player: Player<Spec.SpecEnhancementShaman>) => player.rotationChangeEmitter,
	}, ActionId.fromSpellId(2894), "useFireElemental");

	new IconPicker(fireElementalContainer, simUI.player, fireElementalBooleanIconInput);

	return contentBlock;
}
