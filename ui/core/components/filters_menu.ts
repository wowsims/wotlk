import { Player } from '../player.js';
import {
	ArmorType,
	ItemSlot,
} from '../proto/common.js';
import {
	RaidFilterOption,
	SourceFilterOption,
	UIItem_FactionRestriction,
} from '../proto/ui.js';
import {
	armorTypeNames,
	raidNames,
	rangedWeaponTypeNames,
	sourceNames,
	weaponTypeNames,
} from '../proto_utils/names.js';
import {
	classToEligibleRangedWeaponTypes,
	classToEligibleWeaponTypes,
	classToMaxArmorType,
	isDualWieldSpec,
} from '../proto_utils/utils.js';
import { Sim } from '../sim.js';
import { EventID } from '../typed_event.js';
import { getEnumValues } from '../utils.js';
import { BaseModal } from './base_modal.js';
import { BooleanPicker } from './boolean_picker.js';
import { EnumPicker } from './enum_picker.js';
import { NumberPicker } from './number_picker.js';

const factionRestrictionsToLabels: Record<UIItem_FactionRestriction, string> = {
	[UIItem_FactionRestriction.UNSPECIFIED]: '无限制',
	[UIItem_FactionRestriction.ALLIANCE_ONLY]: '仅联盟',
	[UIItem_FactionRestriction.HORDE_ONLY]: '仅部落',
};

export class FiltersMenu extends BaseModal {
	constructor(rootElem: HTMLElement, player: Player<any>, slot: ItemSlot) {
		super(rootElem, 'filters-menu', { size: 'md', title: '筛选' });

		let section = this.newSection('阵营');

		new EnumPicker(section, player.sim, {
			label: '阵营限制',
			values: [
				UIItem_FactionRestriction.UNSPECIFIED,
				UIItem_FactionRestriction.ALLIANCE_ONLY,
				UIItem_FactionRestriction.HORDE_ONLY
			].map(restriction => {
				return {
					name: factionRestrictionsToLabels[restriction],
					value: restriction,
				};
			}),
			changedEvent: sim => sim.filtersChangeEmitter,
			getValue: (sim: Sim) => sim.getFilters().factionRestriction,
			setValue: (eventID: EventID, sim: Sim, newValue: UIItem_FactionRestriction) => {
				const newFilters = sim.getFilters();
				newFilters.factionRestriction = newValue;
				sim.setFilters(eventID, newFilters);
			},
		});

		section = this.newSection('来源');
		section.classList.add('filters-menu-section-bool-list');
		const sources = Sim.ALL_SOURCES.filter(s => s != SourceFilterOption.SourceUnknown);
		sources.forEach(source => {
			new BooleanPicker<Sim>(section, player.sim, {
				label: sourceNames.get(source),
				inline: true,
				changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
				getValue: (sim: Sim) => sim.getFilters().sources.includes(source),
				setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
					const filters = sim.getFilters();
					if (newValue) {
						filters.sources.push(source);
					} else {
						filters.sources = filters.sources.filter(v => v != source);
					}
					sim.setFilters(eventID, filters);
				},
			});
		});

		section = this.newSection('团队副本');
		section.classList.add('filters-menu-section-bool-list');
		const raids = Sim.ALL_RAIDS.filter(s => s != RaidFilterOption.RaidUnknown);
		raids.forEach(raid => {
			new BooleanPicker<Sim>(section, player.sim, {
				label: raidNames.get(raid),
				inline: true,
				changedEvent: (sim: Sim) => sim.filtersChangeEmitter,
				getValue: (sim: Sim) => sim.getFilters().raids.includes(raid),
				setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
					const filters = sim.getFilters();
					if (newValue) {
						filters.raids.push(raid);
					} else {
						filters.raids = filters.raids.filter(v => v != raid);
					}
					sim.setFilters(eventID, filters);
				},
			});
		});

		if (Player.ARMOR_SLOTS.includes(slot)) {
			const maxArmorType = classToMaxArmorType[player.getClass()];
			if (maxArmorType >= ArmorType.ArmorTypeLeather) {
				const section = this.newSection('护甲类型');
				section.classList.add('filters-menu-section-bool-list');
				const armorTypes = (getEnumValues(ArmorType) as Array<ArmorType>)
					.filter(at => at != ArmorType.ArmorTypeUnknown)
					.filter(at => at <= maxArmorType);

				armorTypes.forEach(armorType => {
					new BooleanPicker<Sim>(section, player.sim, {
						label: armorTypeNames.get(armorType),
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
			const weaponTypeSection = this.newSection('武器类型');
			weaponTypeSection.classList.add('filters-menu-section-bool-list');
			const weaponTypes = classToEligibleWeaponTypes[player.getClass()].map(ewt => ewt.weaponType);

			weaponTypes.forEach(weaponType => {
				new BooleanPicker<Sim>(weaponTypeSection, player.sim, {
					label: weaponTypeNames.get(weaponType),
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

			const weaponSpeedSection = this.newSection('武器速度');
			weaponSpeedSection.classList.add('filters-menu-section-number-list');
			new NumberPicker<Sim>(weaponSpeedSection, player.sim, {
				label: '最低主手攻速',
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
				label: '最高主手攻速',
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
					label: '最低副手攻速',
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
					label: '最高副手攻速',
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
			const rangedWeaponTypeSection = this.newSection('远程武器类型');
			rangedWeaponTypeSection.classList.add('filters-menu-section-bool-list');

			rangedWeaponTypes.forEach(rangedWeaponType => {
				new BooleanPicker<Sim>(rangedWeaponTypeSection, player.sim, {
					label: rangedWeaponTypeNames.get(rangedWeaponType),
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

			const rangedWeaponSpeedSection = this.newSection('远程武器速度');
			rangedWeaponSpeedSection.classList.add('filters-menu-section-number-list');
			new NumberPicker<Sim>(rangedWeaponSpeedSection, player.sim, {
				label: '最低远程武器攻速',
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
				label: '最高远程武器攻速',
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
		this.body.appendChild(section);
		section.innerHTML = `
			<div class="menu-section-header">
				<h6 class="menu-section-title">${name}</h6>
			</div>
			<div class="menu-section-content"></div>
		`;
		return section.getElementsByClassName('menu-section-content')[0] as HTMLElement;
	}
}
