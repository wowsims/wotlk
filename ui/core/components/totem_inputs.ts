import { BooleanPicker } from '../components/boolean_picker.js';
import { EnumPicker } from '../components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '../components/icon_enum_picker.js';
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

export type DpsShaman = Spec.SpecEnhancementShaman | Spec.SpecElementalShaman;

export function TotemsSection(simUI: IndividualSimUI<DpsShaman>, parentElem: HTMLElement): string {
	parentElem.innerHTML = `
		<div class="totem-dropdowns-container"></div>
		<div class="totem-inputs-container"></div>
	`;
	const totemDropdownsContainer = parentElem.getElementsByClassName('totem-dropdowns-container')[0] as HTMLElement;
	const totemInputsContainer = parentElem.getElementsByClassName('totem-inputs-container')[0] as HTMLElement;

	const earthTotemPicker = new IconEnumPicker(totemDropdownsContainer, simUI.player, {
		extraCssClasses: [
			'earth-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#ffdfba', value: EarthTotem.NoEarthTotem },
			{ actionId: ActionId.fromSpellId(58643), value: EarthTotem.StrengthOfEarthTotem },
			{ actionId: ActionId.fromSpellId(8143), value: EarthTotem.TremorTotem },
		],
		equals: (a: EarthTotem, b: EarthTotem) => a == b,
		zeroValue: EarthTotem.NoEarthTotem,
		changedEvent: (player: Player<DpsShaman>) => player.rotationChangeEmitter,
		getValue: (player: Player<DpsShaman>) => player.getRotation().totems?.earth || EarthTotem.NoEarthTotem,
		setValue: (eventID: EventID, player: Player<DpsShaman>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.earth = newValue;
			player.setRotation(eventID, newRotation);
		},
	});

	const airTotemPicker = new IconEnumPicker(totemDropdownsContainer, simUI.player, {
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
		changedEvent: (player: Player<DpsShaman>) => player.rotationChangeEmitter,
		getValue: (player: Player<DpsShaman>) => player.getRotation().totems?.air || AirTotem.NoAirTotem,
		setValue: (eventID: EventID, player: Player<DpsShaman>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.air = newValue;
			player.setRotation(eventID, newRotation);
		},
	});

	const fireTotemPicker = new IconEnumPicker(totemDropdownsContainer, simUI.player, {
		extraCssClasses: [
			'fire-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#ffb3ba', value: FireTotem.NoFireTotem },
			{ actionId: ActionId.fromSpellId(58734), value: FireTotem.MagmaTotem },
			{ actionId: ActionId.fromSpellId(58704), value: FireTotem.SearingTotem },
			{ actionId: ActionId.fromSpellId(57722), value: FireTotem.TotemOfWrath },
		],
		equals: (a: FireTotem, b: FireTotem) => a == b,
		zeroValue: FireTotem.NoFireTotem,
		changedEvent: (player: Player<DpsShaman>) => player.rotationChangeEmitter,
		getValue: (player: Player<DpsShaman>) => player.getRotation().totems?.fire || FireTotem.NoFireTotem,
		setValue: (eventID: EventID, player: Player<DpsShaman>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.fire = newValue;
			player.setRotation(eventID, newRotation);
		},
	});

	const waterTotemPicker = new IconEnumPicker(totemDropdownsContainer, simUI.player, {
		extraCssClasses: [
			'water-totem-picker',
		],
		numColumns: 1,
		values: [
			{ color: '#bae1ff', value: WaterTotem.NoWaterTotem },
			{ actionId: ActionId.fromSpellId(58774), value: WaterTotem.ManaSpringTotem },
		],
		equals: (a: WaterTotem, b: WaterTotem) => a == b,
		zeroValue: WaterTotem.NoWaterTotem,
		changedEvent: (player: Player<DpsShaman>) => player.rotationChangeEmitter,
		getValue: (player: Player<DpsShaman>) => player.getRotation().totems?.water || WaterTotem.NoWaterTotem,
		setValue: (eventID: EventID, player: Player<DpsShaman>, newValue: number) => {
			const newRotation = player.getRotation();
			if (!newRotation.totems)
				newRotation.totems = ShamanTotems.create();
			newRotation.totems!.water = newValue;
			player.setRotation(eventID, newRotation);
		},
	});

	return 'Totems';
}
