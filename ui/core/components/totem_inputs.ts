import { BooleanPicker } from '/wotlk/core/components/boolean_picker.js';
import { EnumPicker } from '/wotlk/core/components/enum_picker.js';
import { IconEnumPicker, IconEnumPickerConfig } from '/wotlk/core/components/icon_enum_picker.js';
import { IconPickerConfig } from '/wotlk/core/components/icon_picker.js';
import {
	AirTotem,
	EarthTotem,
	FireTotem,
	WaterTotem,
	ShamanTotems,
} from '/wotlk/core/proto/shaman.js';
import { Spec } from '/wotlk/core/proto/common.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
import { Player } from '/wotlk/core/player.js';
import { Sim } from '/wotlk/core/sim.js';
import { IndividualSimUI } from '/wotlk/core/individual_sim_ui.js';
import { EventID, TypedEvent } from '/wotlk/core/typed_event.js';

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
			{ actionId: ActionId.fromSpellId(25528), value: EarthTotem.StrengthOfEarthTotem },
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
			{ actionId: ActionId.fromSpellId(25908), value: AirTotem.TranquilAirTotem },
			{ actionId: ActionId.fromSpellId(25587), value: AirTotem.WindfuryTotem },
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
			{ actionId: ActionId.fromSpellId(25552), value: FireTotem.MagmaTotem },
			{ actionId: ActionId.fromSpellId(25533), value: FireTotem.SearingTotem },
			{ actionId: ActionId.fromSpellId(30706), value: FireTotem.TotemOfWrath },
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
			{ actionId: ActionId.fromSpellId(25570), value: WaterTotem.ManaSpringTotem },
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
