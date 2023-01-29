import {
	ArmorType,
	ItemSlot,
} from '../proto/common.js';
import {
	armorTypeNames,
	weaponTypeNames,
} from '../proto_utils/names.js';
import {
	classToEligibleWeaponTypes,
	classToMaxArmorType,
	isDualWieldSpec,
} from '../proto_utils/utils.js';
import { Player } from '../player.js';
import { Sim } from '../sim.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { getEnumValues } from '../utils.js';

import { BooleanPicker } from './boolean_picker.js';
import { NumberPicker } from './number_picker.js';
import { Popup } from './popup.js';
import { Input } from './input.js';
import { BaseModal } from './base_modal.js';

declare var tippy: any;

export class FiltersMenu extends BaseModal {
	constructor(rootElem: HTMLElement, player: Player<any>, slot: ItemSlot) {
		super(rootElem, 'filters-menu', {size: 'md', title: 'Filters'});

		if (Player.ARMOR_SLOTS.includes(slot)) {
			const maxArmorType = classToMaxArmorType[player.getClass()];
			if (maxArmorType >= ArmorType.ArmorTypeLeather) {
				const section = this.newSection('Armor Type');
				section.classList.add('filters-menu-section-bool-list');
				const armorTypes = (getEnumValues(ArmorType) as Array<ArmorType>)
						.filter(at => at != ArmorType.ArmorTypeUnknown)
						.filter(at => at <= maxArmorType);

				armorTypes.forEach(armorType => {
					new BooleanPicker<Sim>(section, player.sim, {
						label: armorTypeNames[armorType],
						inline: true,
						changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
						getValue: (sim: Sim) => sim.getFilters().armorTypes.includes(armorType),
						setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
							const filters = sim.getFilters();
							if (newValue) {
								filters.armorTypes.push(armorType);
							} else {
								filters.armorTypes = filters.armorTypes.filter(at => at != armorType);
							}
							sim.setFilters(eventID, filters);
						},
					});
				});
			}
		} else if (Player.WEAPON_SLOTS.includes(slot)) {
			const weaponTypeSection = this.newSection('Weapon Type');
			weaponTypeSection.classList.add('filters-menu-section-bool-list');
			const weaponTypes = classToEligibleWeaponTypes[player.getClass()].map(ewt => ewt.weaponType);

			weaponTypes.forEach(weaponType => {
				new BooleanPicker<Sim>(weaponTypeSection, player.sim, {
					label: weaponTypeNames[weaponType],
					inline: true,
					changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
					getValue: (sim: Sim) => sim.getFilters().weaponTypes.includes(weaponType),
					setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
						const filters = sim.getFilters();
						if (newValue) {
							filters.weaponTypes.push(weaponType);
						} else {
							filters.weaponTypes = filters.weaponTypes.filter(at => at != weaponType);
						}
						sim.setFilters(eventID, filters);
					},
				});
			});

			const weaponSpeedSection = this.newSection('Weapon Speed');
			weaponSpeedSection.classList.add('filters-menu-section-number-list');
			new NumberPicker<Sim>(weaponSpeedSection, player.sim, {
				label: 'Min MH Speed',
				//labelTooltip: 'Maximum speed for the mainhand weapon. If 0, no maximum value is applied.',
				float: true,
				positive: true,
				changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
				getValue: (sim: Sim) => sim.getFilters().minMhWeaponSpeed,
				setValue: (eventID: EventID, sim: Sim, newValue: number) => {
					const filters = sim.getFilters();
					filters.minMhWeaponSpeed = newValue;
					sim.setFilters(eventID, filters);
				},
			});
			new NumberPicker<Sim>(weaponSpeedSection, player.sim, {
				label: 'Max MH Speed',
				//labelTooltip: 'Maximum speed for the mainhand weapon. If 0, no maximum value is applied.',
				float: true,
				positive: true,
				changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
				getValue: (sim: Sim) => sim.getFilters().maxMhWeaponSpeed,
				setValue: (eventID: EventID, sim: Sim, newValue: number) => {
					const filters = sim.getFilters();
					filters.maxMhWeaponSpeed = newValue;
					sim.setFilters(eventID, filters);
				},
			});
			if (isDualWieldSpec(player.spec)) {
				new NumberPicker<Sim>(weaponSpeedSection, player.sim, {
					label: 'Min OH Speed',
					//labelTooltip: 'Minimum speed for the offhand weapon. If 0, no minimum value is applied.',
					float: true,
					positive: true,
					changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
					getValue: (sim: Sim) => sim.getFilters().minOhWeaponSpeed,
					setValue: (eventID: EventID, sim: Sim, newValue: number) => {
						const filters = sim.getFilters();
						filters.minOhWeaponSpeed = newValue;
						sim.setFilters(eventID, filters);
					},
				});
				new NumberPicker<Sim>(weaponSpeedSection, player.sim, {
					label: 'Max OH Speed',
					//labelTooltip: 'Maximum speed for the offhand weapon. If 0, no maximum value is applied.',
					float: true,
					positive: true,
					changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
					getValue: (sim: Sim) => sim.getFilters().maxOhWeaponSpeed,
					setValue: (eventID: EventID, sim: Sim, newValue: number) => {
						const filters = sim.getFilters();
						filters.maxOhWeaponSpeed = newValue;
						sim.setFilters(eventID, filters);
					},
				});
			}
		}
	}

	private newSection(name: string): HTMLElement {
		const section = document.createElement('div');
		section.classList.add('menu-section');
		this.body.appendChild(section);
		section.innerHTML = `
			<div class="menu-section-header">
				<h6 class="menu-section-title">${name}</h6>
			</div>
			<div class="menu-section-content"></div>
		`;
		return section.getElementsByClassName('menu-section-content')[0] as HTMLElement;
	}

	static anyFiltersForSlot(slot: ItemSlot) {
		return Player.ARMOR_SLOTS.includes(slot) || Player.WEAPON_SLOTS.includes(slot);
	}
}
