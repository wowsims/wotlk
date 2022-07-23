import { IconEnumPicker } from '/wotlk/core/components/icon_enum_picker.js';
import { AirTotem, EarthTotem, FireTotem, WaterTotem, ShamanTotems, } from '/wotlk/core/proto/shaman.js';
import { ActionId } from '/wotlk/core/proto_utils/action_id.js';
export function TotemsSection(simUI, parentElem) {
    parentElem.innerHTML = `
		<div class="totem-dropdowns-container"></div>
		<div class="totem-inputs-container"></div>
	`;
    const totemDropdownsContainer = parentElem.getElementsByClassName('totem-dropdowns-container')[0];
    const totemInputsContainer = parentElem.getElementsByClassName('totem-inputs-container')[0];
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
        equals: (a, b) => a == b,
        zeroValue: EarthTotem.NoEarthTotem,
        changedEvent: (player) => player.rotationChangeEmitter,
        getValue: (player) => player.getRotation().totems?.earth || EarthTotem.NoEarthTotem,
        setValue: (eventID, player, newValue) => {
            const newRotation = player.getRotation();
            if (!newRotation.totems)
                newRotation.totems = ShamanTotems.create();
            newRotation.totems.earth = newValue;
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
        equals: (a, b) => a == b,
        zeroValue: AirTotem.NoAirTotem,
        changedEvent: (player) => player.rotationChangeEmitter,
        getValue: (player) => player.getRotation().totems?.air || AirTotem.NoAirTotem,
        setValue: (eventID, player, newValue) => {
            const newRotation = player.getRotation();
            if (!newRotation.totems)
                newRotation.totems = ShamanTotems.create();
            newRotation.totems.air = newValue;
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
        equals: (a, b) => a == b,
        zeroValue: FireTotem.NoFireTotem,
        changedEvent: (player) => player.rotationChangeEmitter,
        getValue: (player) => player.getRotation().totems?.fire || FireTotem.NoFireTotem,
        setValue: (eventID, player, newValue) => {
            const newRotation = player.getRotation();
            if (!newRotation.totems)
                newRotation.totems = ShamanTotems.create();
            newRotation.totems.fire = newValue;
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
        equals: (a, b) => a == b,
        zeroValue: WaterTotem.NoWaterTotem,
        changedEvent: (player) => player.rotationChangeEmitter,
        getValue: (player) => player.getRotation().totems?.water || WaterTotem.NoWaterTotem,
        setValue: (eventID, player, newValue) => {
            const newRotation = player.getRotation();
            if (!newRotation.totems)
                newRotation.totems = ShamanTotems.create();
            newRotation.totems.water = newValue;
            player.setRotation(eventID, newRotation);
        },
    });
    return 'Totems';
}
