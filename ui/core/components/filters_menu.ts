import {
	ArmorType,
	ItemSlot,
} from '../proto/common.js';
import {
	armorTypeNames,
	rangedWeaponTypeNames,
	weaponTypeNames,
} from '../proto_utils/names.js';
import {
	classToEligibleRangedWeaponTypes,
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

declare var tippy: any;

export class FiltersMenu extends Popup {
	private readonly contentElem: HTMLElement;

	constructor(rootElem: HTMLElement, player: Player<any>, slot: ItemSlot) {
		super(rootElem);

		this.rootElem.classList.add('filters-menu');
		this.rootElem.innerHTML = `
			<div class="menu-title">
				<span>FILTERS</span>
			</div>
			<div class="menu-content">
			</div>
		`;
		this.addCloseButton();

		this.contentElem = this.rootElem.getElementsByClassName('menu-content')[0] as HTMLElement;

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
			const weaponTypesGroup = Input.newGroupContainer();
			weaponTypeSection.appendChild(weaponTypesGroup);

			weaponTypes.forEach(weaponType => {
				new BooleanPicker<Sim>(weaponTypesGroup, player.sim, {
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
		} else if (slot == ItemSlot.ItemSlotRanged) {
			const rangedWeaponTypes = classToEligibleRangedWeaponTypes[player.getClass()];
			if (rangedWeaponTypes.length <= 1) {
				return;
			}
			const rangedWeaponTypeSection = this.newSection('Ranged Weapon Type');
			rangedWeaponTypeSection.classList.add('filters-menu-section-bool-list');
			const rangedWeaponTypesGroup = Input.newGroupContainer();
			rangedWeaponTypeSection.appendChild(rangedWeaponTypesGroup);

			rangedWeaponTypes.forEach(rangedWeaponType => {
				new BooleanPicker<Sim>(rangedWeaponTypesGroup, player.sim, {
					label: rangedWeaponTypeNames[rangedWeaponType],
					inline: true,
					changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
					getValue: (sim: Sim) => sim.getFilters().rangedWeaponTypes.includes(rangedWeaponType),
					setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
						const filters = sim.getFilters();
						if (newValue) {
							filters.rangedWeaponTypes.push(rangedWeaponType);
						} else {
							filters.rangedWeaponTypes = filters.rangedWeaponTypes.filter(at => at != rangedWeaponType);
						}
						sim.setFilters(eventID, filters);
					},
				});
			});

			const rangedWeaponSpeedSection = this.newSection('Ranged Weapon Speed');
			rangedWeaponSpeedSection.classList.add('filters-menu-section-number-list');
			new NumberPicker<Sim>(rangedWeaponSpeedSection, player.sim, {
				label: 'Min Ranged Speed',
				//labelTooltip: 'Maximum speed for the ranged weapon. If 0, no maximum value is applied.',
				float: true,
				positive: true,
				changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
				getValue: (sim: Sim) => sim.getFilters().minRangedWeaponSpeed,
				setValue: (eventID: EventID, sim: Sim, newValue: number) => {
					const filters = sim.getFilters();
					filters.minRangedWeaponSpeed = newValue;
					sim.setFilters(eventID, filters);
				},
			});
			new NumberPicker<Sim>(rangedWeaponSpeedSection, player.sim, {
				label: 'Max Ranged Speed',
				//labelTooltip: 'Maximum speed for the ranged weapon. If 0, no maximum value is applied.',
				float: true,
				positive: true,
				changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
				getValue: (sim: Sim) => sim.getFilters().maxRangedWeaponSpeed,
				setValue: (eventID: EventID, sim: Sim, newValue: number) => {
					const filters = sim.getFilters();
					filters.maxRangedWeaponSpeed = newValue;
					sim.setFilters(eventID, filters);
				},
			});
		}
	}

	private newSection(name: string): HTMLElement {
		const section = document.createElement('div');
		section.classList.add('menu-section');
		this.contentElem.appendChild(section);
		section.innerHTML = `
			<div class="menu-section-header">
				<span class="menu-section-title">${name}</span>
			</div>
			<div class="menu-section-content">
			</div>
		`;
		return section.getElementsByClassName('menu-section-content')[0] as HTMLElement;
	}

	static anyFiltersForSlot(slot: ItemSlot) {
		return [
			Player.ARMOR_SLOTS,
			Player.WEAPON_SLOTS,
			ItemSlot.ItemSlotRanged,
		].flat().includes(slot);
	}
}
