import { difficultyNames, professionNames, slotNames } from '../proto_utils/names.js';
import { BaseModal } from './base_modal';
import { Component } from './component';
import { FiltersMenu } from './filters_menu';
import { Input, InputConfig } from './input';
import {
	makePhaseSelector,
	makeShow1hWeaponsSelector,
	makeShow2hWeaponsSelector,
	makeShowEPValuesSelector,
} from './other_inputs';

import { setItemQualityCssClass } from '../css_utils';
import { Player } from '../player';
import { Sim } from '../sim.js';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { formatDeltaTextElem } from '../utils';

import { ActionId } from '../proto_utils/action_id';
import { getEnchantDescription, getUniqueEnchantString } from '../proto_utils/enchants';
import { EquippedItem } from '../proto_utils/equipped_item';
import { Stats } from '../proto_utils/stats';

import { Tooltip } from 'bootstrap';
import { IndividualSimUI } from '../individual_sim_ui.js';
import {
	Class,
	ItemQuality,
	ItemSlot,
	ItemSpec,
	Spec,
} from '../proto/common';
import {
	DatabaseFilters,
	UIEnchant as Enchant,
	UIItem as Item,
	UIRune as Rune,
} from '../proto/ui.js';
// eslint-disable-next-line unused-imports/no-unused-imports
import { element, fragment, ref } from 'tsx-vanilla';

import { itemTypeToSlotsMap } from '../proto_utils/utils.js';
import { Clusterize } from './virtual_scroll/clusterize.js';

const EP_TOOLTIP = `
	EP (Equivalence Points) is way of comparing items by multiplying the raw stats of an item with your current stat weights.
	More EP does not necessarily mean more DPS, as EP doesn't take into account stat caps and non-linear stat calculations.
`

const createHeroicLabel = () => {
	return (<span className='heroic-label'>[H]</span>);
}

export class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSide = document.createElement('div');
		leftSide.classList.add('gear-picker-left', 'tab-panel-col');
		this.rootElem.appendChild(leftSide);

		const rightSide = document.createElement('div');
		rightSide.classList.add('gear-picker-right', 'tab-panel-col');
		this.rootElem.appendChild(rightSide);

		const leftItemPickers = [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotNeck,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotBack,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotMainHand,
			ItemSlot.ItemSlotOffHand,
			ItemSlot.ItemSlotRanged,
		].map(slot => new ItemPicker(leftSide, simUI, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
		].map(slot => new ItemPicker(rightSide, simUI, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);
	}
}

export class ItemRenderer extends Component {
	private readonly player: Player<any>;

	readonly iconElem: HTMLAnchorElement;
	readonly nameElem: HTMLAnchorElement;
	readonly enchantElem: HTMLAnchorElement;
	readonly runeElem: HTMLAnchorElement;

	constructor(parent: HTMLElement, player: Player<any>) {
		super(parent, 'item-picker-root');
		this.player = player;

		let iconElem = ref<HTMLAnchorElement>();
		let nameElem = ref<HTMLAnchorElement>();
		let enchantElem = ref<HTMLAnchorElement>();
		let runeElem = ref<HTMLAnchorElement>();
		this.rootElem.appendChild(
			<>
				<a ref={iconElem} className="item-picker-icon" href="javascript:void(0)" attributes={{role:"button"}}></a>
				<div className="item-picker-labels-container">
					<a ref={nameElem} className="item-picker-name" href="javascript:void(0)" attributes={{role:"button"}}></a>
					<a ref={enchantElem} className="item-picker-enchant" href="javascript:void(0)" attributes={{role:"button"}}></a>
					<a ref={runeElem} className="item-picker-rune" href="javascript:void(0)" attributes={{role:"button"}}></a>
				</div>
			</>
		);

		this.iconElem = iconElem.value!;
		this.nameElem = nameElem.value!;
		this.enchantElem = enchantElem.value!;
		this.runeElem = runeElem.value!;
	}

	clear() {
		this.nameElem.removeAttribute('data-wowhead');
		this.nameElem.removeAttribute('href');
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');
		this.enchantElem.removeAttribute('data-wowhead');
		this.enchantElem.removeAttribute('href');
		this.iconElem.removeAttribute('href');

		this.iconElem.style.backgroundImage = '';
		this.enchantElem.innerText = '';
		this.runeElem.innerText = '';
		this.nameElem.textContent = '';
	}

	update(newItem: EquippedItem) {
		this.nameElem.textContent = newItem.item.name;
		if (newItem.item.heroic) {
			this.nameElem.insertAdjacentElement('beforeend', createHeroicLabel());
		} else {
			this.nameElem.querySelector('.heroic-label')?.remove();
		}

		setItemQualityCssClass(this.nameElem, newItem.item.quality);

		this.player.setWowheadData(newItem, this.iconElem);
		this.player.setWowheadData(newItem, this.nameElem);
		newItem.asActionId().fill().then(filledId => {
			filledId.setBackgroundAndHref(this.iconElem);
			filledId.setWowheadHref(this.nameElem);
		});

		if (newItem.enchant) {
			getEnchantDescription(newItem.enchant).then(description => {
				this.enchantElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.enchant.spellId) {
				this.enchantElem.href = ActionId.makeSpellUrl(newItem.enchant.spellId);
				this.enchantElem.dataset.wowhead = `domain=classic&spell=${newItem.enchant.spellId}`;
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
				this.enchantElem.dataset.wowhead = `domain=classic&item=${newItem.enchant.itemId}`;
			}
			this.enchantElem.dataset.whtticon = 'false';
		}

		const isRuneSlot = itemTypeToSlotsMap[newItem._item.type]?.some(slot => this.player.sim.db.hasRuneBySlot(slot, this.player.getClass()));
		if (isRuneSlot){
			this.iconElem.appendChild(this.createRuneContainer(newItem.rune))

			if (newItem.rune) {
				this.runeElem.textContent = newItem.rune.name
				this.runeElem.href = ActionId.makeSpellUrl(newItem.rune.id)
				this.enchantElem.dataset.wowhead = `domain=classic&spell=${newItem.rune.id}`;
			}
		}
	}

	private createRuneContainer = (rune : Rune|null) => {
		const runeIconElem = ref<HTMLImageElement>();
		let runeContainer = (
			<div className="item-picker-rune-container">
				<img ref={runeIconElem} className="item-picker-rune-icon" />
			</div>
		);
	
		if (rune) {
			ActionId.fromSpellId(rune.id).fill().then(filledId => runeIconElem.value!.src = filledId.iconUrl);
		} else {
			runeIconElem.value!.src = "https://wow.zamimg.com/images/wow/icons/medium/inventoryslot_empty.jpg"
		}

		return runeContainer;
	}
}

export class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly itemElem: ItemRenderer;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];
	private _runes: Array<Rune> = [];
	private _equippedItem: EquippedItem | null = null;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');
		this.slot = slot;
		this.simUI = simUI;
		this.player = player;
		this.itemElem = new ItemRenderer(this.rootElem, player);
		this.item = player.getEquippedItem(slot);

		const loadItems = () => this._items = this.player.getItems(this.slot);

		player.sim.waitForInit().then(() => {
			this._enchants = this.player.getEnchants(this.slot);
			this._runes = this.player.getRunes(this.slot);
			loadItems();

			const gearData = {
				equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
					this.player.equipItem(eventID, this.slot, equippedItem);
				},
				getEquippedItem: () => this.player.getEquippedItem(this.slot),
				changeEvent: player.gearChangeEmitter,
			};

			const openGearSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Items, gearData);
			};
			const openEnchantSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Enchants, gearData);
			};

			this.itemElem.iconElem.addEventListener('click', openGearSelector);
			this.itemElem.nameElem.addEventListener('click', openGearSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantSelector);
		});

		player.levelChangeEmitter.on(loadItems)

		player.gearChangeEmitter.on(() => {
			this.item = player.getEquippedItem(slot);
		});
		player.professionChangeEmitter.on(() => {
			if (this._equippedItem != null) {
				this.player.setWowheadData(this._equippedItem, this.itemElem.iconElem);
			}
		});
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.itemElem.clear();
		this.itemElem.nameElem.textContent = slotNames.get(this.slot) ?? '';
		setItemQualityCssClass(this.itemElem.nameElem, null);

		if (newItem != null) {
			this.itemElem.update(newItem);
		} else {
			this.itemElem.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		}

		this._equippedItem = newItem;
	}

	private openSelectorModal(tab: SelectorModalTabs, gearData: GearData) {
		new SelectorModal(this.simUI.rootElem, this.simUI, this.player, {
			selectedTab: tab,
			slot: this.slot,
			equippedItem: this._equippedItem,
			eligibleItems: this._items,
			eligibleEnchants: this._enchants,
			eligibleRunes: this._runes,
			gearData: gearData
		})
	}
}

export class IconItemSwapPicker<SpecType extends Spec, ValueType> extends Input<Player<SpecType>, ValueType> {
	private readonly config: InputConfig<Player<SpecType>, ValueType>;
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly player: Player<SpecType>;
	private readonly slot: ItemSlot;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];
	private _runes: Array<Rune> = [];

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<SpecType>, slot: ItemSlot, config: InputConfig<Player<SpecType>, ValueType>) {
		super(parent, 'icon-picker-root', player, config)
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.config = config;
		this.slot = slot;

		this.iconAnchor = document.createElement('a');
		this.iconAnchor.classList.add('icon-picker-button');
		this.iconAnchor.target = '_blank';
		this.rootElem.prepend(this.iconAnchor);

		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(slot);
			this._enchants = this.player.getEnchants(slot);
			this._runes = this.player.getRunes(slot);
			const gearData = {
				equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
					let isg = this.player.getItemSwapGear();
					this.player.setItemSwapGear(eventID, isg.withEquippedItem(this.slot, equippedItem, player.canDualWield2H()));
					this.inputChanged(eventID);
				},
				getEquippedItem: () => this.player.getItemSwapGear().getEquippedItem(this.slot),
				changeEvent: config.changedEvent(player),
			}

			this.iconAnchor.addEventListener('click', (event: Event) => {
				event.preventDefault();
				new SelectorModal(simUI.rootElem, simUI, this.player, {
					selectedTab: SelectorModalTabs.Items,
					slot: this.slot,
					equippedItem: this.player.getItemSwapGear().getEquippedItem(slot),
					eligibleItems: this._items,
					eligibleEnchants: this._enchants,
					eligibleRunes: this._runes,
					gearData: gearData,
				});
			});
		}).finally(() => this.init());

	}

	getInputElem(): HTMLElement {
		return this.iconAnchor;
	}
	getInputValue(): ValueType {
		return this.player.getItemSwapGear().toProto() as unknown as ValueType
	}

	setInputValue(_: ValueType): void {
		this.iconAnchor.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconAnchor.removeAttribute('data-wowhead');
		this.iconAnchor.href = "#";

		const equippedItem = this.player.getItemSwapGear().getEquippedItem(this.slot);
		if (equippedItem) {
			this.iconAnchor.classList.add("active")

			equippedItem.asActionId().fillAndSet(this.iconAnchor, true, true);
			this.player.setWowheadData(equippedItem, this.iconAnchor);

		} else {
			this.iconAnchor.classList.remove("active")
		}
	}

}

export interface GearData {
	equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => void,
	getEquippedItem: () => EquippedItem | null,
	changeEvent: TypedEvent<any>,
}

export enum SelectorModalTabs {
	Items = 'Items',
	Enchants = 'Enchants',
	Rune = 'Rune',
}

interface SelectorModalConfig {
	selectedTab: SelectorModalTabs
	slot: ItemSlot,
	equippedItem: EquippedItem | null,
	eligibleItems: Array<Item>,
	eligibleEnchants: Array<Enchant>,
	eligibleRunes: Array<Rune>,
	gearData: GearData
}

export class SelectorModal extends BaseModal {
	private readonly simUI: SimUI;
	private player: Player<any>;
	private config: SelectorModalConfig;
	private ilists: ItemList<any>[];

	private readonly tabsElem: HTMLElement;
	private readonly contentElem: HTMLElement;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, config: SelectorModalConfig) {
		super(parent, 'selector-modal');

		this.simUI = simUI;
		this.player = player;
		this.config = config;
		this.ilists = [];

		window.scrollTo({ top: 0 });

		this.header!.insertAdjacentElement('afterbegin', <ul className="nav nav-tabs selector-modal-tabs"></ul>);

		this.body.appendChild(<div className='tab-content selector-modal-tab-content'></div>);

		this.tabsElem = this.rootElem.querySelector('.selector-modal-tabs') as HTMLElement;
		this.contentElem = this.rootElem.querySelector('.selector-modal-tab-content') as HTMLElement;

		this.setData();
	}

	// Could be 'Items' 'Enchants' or 'Rune'
	openTabName(name: string) {
		Array.from(this.tabsElem.getElementsByClassName("selector-modal-item-tab")).forEach(elem => {
			if (elem.getAttribute("data-content-id") == name + "-tab") {
				(elem as HTMLElement).click();
			}
		});
	}

	openTab(idx: number) {
		const elems = this.tabsElem.getElementsByClassName("selector-modal-item-tab");
		(elems[idx] as HTMLElement).click();
	}

	setData() {
		this.tabsElem.innerText = '';
		this.contentElem.innerText = '';

		const { slot, eligibleItems, eligibleEnchants, eligibleRunes, gearData } = this.config;

		this.addTab<Item>(
			'Items',
			eligibleItems.map(item => {
				return {
					item: item,
					id: item.id,
					actionId: ActionId.fromItem(item),
					name: item.name,
					quality: item.quality,
					heroic: item.heroic,
					phase: item.phase,
					baseEP: this.player.computeItemEP(item, slot),
					ignoreEPFilter: false,
					onEquip: (eventID, item: Item) => {
						const equippedItem = gearData.getEquippedItem();
						if (equippedItem) {
							gearData.equipItem(eventID, equippedItem.withItem(item));
						} else {
							gearData.equipItem(eventID, new EquippedItem({item}));
						}
					},
				};
			}),
			item => this.player.computeItemEP(item, slot),
			equippedItem => equippedItem?.item,
			eventID => {
				gearData.equipItem(eventID, null);
			});

		this.addTab<Enchant>(
			'Enchants',
			eligibleEnchants.map(enchant => {
				return {
					item: enchant,
					id: enchant.effectId,
					actionId: enchant.spellId ? ActionId.fromSpellId(enchant.spellId) : ActionId.fromItemId(enchant.itemId),
					name: enchant.name,
					quality: enchant.quality,
					phase: enchant.phase || 1,
					baseEP: this.player.computeStatsEP(new Stats(enchant.stats)),
					ignoreEPFilter: true,
					heroic: false,
					onEquip: (eventID, enchant: Enchant) => {
						const equippedItem = gearData.getEquippedItem();
						if (equippedItem)
							gearData.equipItem(eventID, equippedItem.withEnchant(enchant));
					},
				};
			}),
			enchant => this.player.computeEnchantEP(enchant),
			equippedItem => equippedItem?.enchant,
			eventID => {
				const equippedItem = gearData.getEquippedItem();
				if (equippedItem)
					gearData.equipItem(eventID, equippedItem.withEnchant(null));
			});

			this.addTab<Rune>(
				'Rune',
				eligibleRunes.map(rune => {
					return {
						item: rune,
						id: rune.id,
						actionId: ActionId.fromSpellId(rune.id),
						name: rune.name,
						quality: ItemQuality.ItemQualityCommon,
						phase: 1,
						baseEP: 1,
						ignoreEPFilter: true,
						heroic: false,
						onEquip: (eventID, rune: Rune) => {
							const equippedItem = gearData.getEquippedItem();
							if (equippedItem)
								gearData.equipItem(eventID, equippedItem.withRune(rune));
						},
					};
				}),
				() => 1,
				equippedItem => equippedItem?.rune,
				eventID => {
					const equippedItem = gearData.getEquippedItem();
					if (equippedItem)
						gearData.equipItem(eventID, equippedItem.withRune(null));
				});
	}

	protected override onShow(e: Event) {
		// Only refresh opened tab
		let t = e.target! as HTMLElement;
		let tab = t.querySelector<HTMLElement>('.active')!.dataset.contentId!;
		if (tab.includes('Item')) {
			this.ilists[0].sizeRefresh();
		}
		else if (tab.includes('Enchant')) {
			this.ilists[1].sizeRefresh();
		}
	}

	/**
	 * Adds one of the tabs for the item selector menu.
	 *
	 * T is expected to be Item, Enchant, or Rune. Tab menus for all 3 looks extremely
	 * similar so this function uses extra functions to do it generically.
	 */
	private addTab<T>(
		label: string,
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => (T | null | undefined),
		onRemove: (eventID: EventID) => void,
		setTabContent?: (tabElem: HTMLAnchorElement) => void) {
		if (itemData.length == 0) {
			return;
		}

		const { gearData } = this.config;
		const tabContentId = (label + '-tab').split(' ').join('');
		const selected = label === this.config.selectedTab;

		const tabAnchor = ref<HTMLAnchorElement>();
		this.tabsElem.appendChild(
			<li className="nav-item">
				<a
					ref={tabAnchor}
					className={`nav-link selector-modal-item-tab ${selected ? 'active' : ''}`}
					dataset={{
						label: label,
						contentId:tabContentId,
						bsToggle:'tab',
						bsTarget:`#${tabContentId}`,
					}}
					attributes={{
						role:'tab',
						'aria-selected':selected,
					}}
					type="button"
				></a>
			</li>
		);

		if (setTabContent) {
			setTabContent(tabAnchor.value!);
		} else {
			tabAnchor.value!.textContent = label;
		}

		let ilist = new ItemList<T>(
			this.contentElem,
			this.simUI,
			this.config,
			this.player,
			label,
			itemData,
			computeEP,
			equippedToItemFn,
			onRemove,
			(itemData: ItemData<T>) => {
				const item = itemData.item;
				itemData.onEquip(TypedEvent.nextEventID(), item);
			},
		)

		let invokeUpdate = () => { ilist.updateSelected() }
		let applyFilter = () => { ilist.applyFilters() }
		let hideOrShowEPValues = () => { ilist.hideOrShowEPValues() }
		// Add event handlers
		gearData.changeEvent.on(invokeUpdate);

		this.player.sim.phaseChangeEmitter.on(applyFilter);
		this.player.sim.filtersChangeEmitter.on(applyFilter);
		this.player.sim.showEPValuesChangeEmitter.on(hideOrShowEPValues);

		this.addOnDisposeCallback(() => {
			gearData.changeEvent.off(invokeUpdate)
			this.player.sim.phaseChangeEmitter.off(applyFilter);
			this.player.sim.filtersChangeEmitter.off(applyFilter);
			this.player.sim.showEPValuesChangeEmitter.off(hideOrShowEPValues);
			ilist.dispose();
		});

		tabAnchor.value!.addEventListener('shown.bs.tab', (_) => {
			ilist.sizeRefresh()
		});

		this.ilists.push(ilist);
	}

	private removeTabs(labelSubstring: string) {
		const tabElems = Array.prototype.slice.call(this.tabsElem.getElementsByClassName('selector-modal-item-tab'))
			.filter(tab => tab.dataset.label.includes(labelSubstring));

		const contentElems = tabElems
			.map(tabElem => document.getElementById(tabElem.dataset.contentId!))
			.filter(tabElem => Boolean(tabElem));

		tabElems.forEach(elem => elem.parentElement.remove());
		contentElems.forEach(elem => elem!.remove());
	}
}

export interface ItemData<T> {
	item: T,
	name: string,
	id: number,
	actionId: ActionId,
	quality: ItemQuality,
	phase: number,
	baseEP: number,
	ignoreEPFilter: boolean,
	heroic: boolean,
	onEquip: (eventID: EventID, item: T) => void,
}

interface ItemDataWithIdx<T> {
	idx: number,
	data: ItemData<T>,
}

const emptySlotIcons: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: '/sod/assets/item_slots/head.jpg',
	[ItemSlot.ItemSlotNeck]: '/sod/assets/item_slots/neck.jpg',
	[ItemSlot.ItemSlotShoulder]: '/sod/assets/item_slots/shoulders.jpg',
	[ItemSlot.ItemSlotBack]: '/sod/assets/item_slots/shirt.jpg',
	[ItemSlot.ItemSlotChest]: '/sod/assets/item_slots/chest.jpg',
	[ItemSlot.ItemSlotWrist]: '/sod/assets/item_slots/wrists.jpg',
	[ItemSlot.ItemSlotHands]: '/sod/assets/item_slots/hands.jpg',
	[ItemSlot.ItemSlotWaist]: '/sod/assets/item_slots/waist.jpg',
	[ItemSlot.ItemSlotLegs]: '/sod/assets/item_slots/legs.jpg',
	[ItemSlot.ItemSlotFeet]: '/sod/assets/item_slots/feet.jpg',
	[ItemSlot.ItemSlotFinger1]: '/sod/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotFinger2]: '/sod/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: '/sod/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: '/sod/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: '/sod/assets/item_slots/mainhand.jpg',
	[ItemSlot.ItemSlotOffHand]: '/sod/assets/item_slots/offhand.jpg',
	[ItemSlot.ItemSlotRanged]: '/sod/assets/item_slots/ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}

export class ItemList<T> {
	private listElem: HTMLElement;
	private readonly player: Player<any>;
	private label: string;
	private slot: ItemSlot;
	private itemData: Array<ItemData<T>>;
	private itemsToDisplay: Array<number>;
	private currentFilters: DatabaseFilters;
	private searchInput: HTMLInputElement;
	private computeEP: (item: T) => number;
	private equippedToItemFn: (equippedItem: EquippedItem | null) => (T | null | undefined);
	private gearData: GearData;
	private tabContent: Element;
	private onItemClick: (itemData: ItemData<T>) => void;
	private scroller: Clusterize;

	constructor(
		parent: HTMLElement,
		simUI: SimUI,
		config: SelectorModalConfig,
		player: Player<any>,
		label: string,
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => (T | null | undefined),
		onRemove: (eventID: EventID) => void,
		onItemClick: (itemData: ItemData<T>) => void) {
		this.label = label;
		this.player = player;
		this.itemData = itemData;
		this.computeEP = computeEP;
		this.equippedToItemFn = equippedToItemFn;
		this.onItemClick = onItemClick;

		const { slot, gearData } = config;
		this.slot = slot;
		this.gearData = gearData;
		this.currentFilters = this.player.sim.getFilters();

		const tabContentId = (label + '-tab').split(' ').join('');
		const selected = label === config.selectedTab;

		const epButton = ref<HTMLButtonElement>();
		this.tabContent = (
			<div
				id={tabContentId}
				className={`selector-modal-tab-pane tab-pane fade ${selected ? 'active show' : ''}`}
			>
				<div className="selector-modal-filters">
					<input className="selector-modal-search form-control" type="text" placeholder="Search..."/>
					{label == 'Items' && <button className="selector-modal-filters-button btn btn-primary">Filters</button>}
					<div className="selector-modal-phase-selector"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-1h-weapons"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-2h-weapons"></div>
					<div className="sim-input selector-modal-boolean-option selector-modal-show-ep-values"></div>
					<button className="selector-modal-simall-button btn btn-warning">Add to Batch Sim</button>
					<button className="selector-modal-remove-button btn btn-danger">Unequip Item</button>
				</div>
				<div className="selector-modal-list-labels">
					<label>Item</label>
					<label className="ep-delta-label">
						EP
						<i className="fa-solid fa-plus-minus fa-2xs"></i>
						<button
							ref={epButton}
							className="btn btn-link p-0 ms-1">
							<i className="far fa-question-circle fa-lg"></i>
						</button>
					</label>
				</div>
				<ul className="selector-modal-list"></ul>
			</div>
		);

		parent.appendChild(this.tabContent);

		new Tooltip(epButton.value!, {
			title: EP_TOOLTIP
		});

		makeShow1hWeaponsSelector(this.tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement, player.sim);
		makeShow2hWeaponsSelector(this.tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement, player.sim);
		
		if (!(label == 'Items' && (slot == ItemSlot.ItemSlotMainHand || (slot == ItemSlot.ItemSlotOffHand && player.getClass() == Class.ClassWarrior)))) {
			(this.tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement).style.display = 'none';
			(this.tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement).style.display = 'none';
		}

		makeShowEPValuesSelector(this.tabContent.getElementsByClassName('selector-modal-show-ep-values')[0] as HTMLElement, player.sim);
		makePhaseSelector(this.tabContent.getElementsByClassName('selector-modal-phase-selector')[0] as HTMLElement, player.sim);

		if (label == 'Items') {
			const filtersButton = this.tabContent.getElementsByClassName('selector-modal-filters-button')[0] as HTMLElement;
			filtersButton.addEventListener('click', () => new FiltersMenu(parent, player, slot));
		}

		this.listElem = this.tabContent.getElementsByClassName('selector-modal-list')[0] as HTMLElement;

		this.itemsToDisplay = [];

		this.scroller = new Clusterize({
			getNumberOfRows: () => { return this.itemsToDisplay.length },
			generateRows: (startIdx, endIdx) => {
				let items = [];
				for (let i = startIdx; i < endIdx; ++i) {
					if (i >= this.itemsToDisplay.length)
						break;
					items.push(this.createItemElem({idx:this.itemsToDisplay[i], data:this.itemData[this.itemsToDisplay[i]]}));
				}
				return items;
			}
		}, {
			rows: [],
			scroll_elem: this.listElem,
			content_elem: this.listElem,
			item_height: 56,
			show_no_data_row: false,
			no_data_text: '',
			tag: 'li',
			rows_in_block: 16,
			blocks_in_cluster: 2,
		});

		const removeButton = this.tabContent.getElementsByClassName('selector-modal-remove-button')[0] as HTMLButtonElement;
		removeButton.addEventListener('click', _ => {
			onRemove(TypedEvent.nextEventID());
		});

		if (label.startsWith("Enchants")) {
			removeButton.textContent = 'Remove Enchant';
		} else if (label.startsWith("Rune")) {
			removeButton.textContent = 'Remove Rune';
		}

		this.updateSelected();

		this.searchInput = this.tabContent.getElementsByClassName('selector-modal-search')[0] as HTMLInputElement;
		this.searchInput.addEventListener('input', () => this.applyFilters());

		const simAllButton = this.tabContent.getElementsByClassName('selector-modal-simall-button')[0] as HTMLButtonElement;
		if (label == "Items") {
			simAllButton.hidden = !player.sim.getShowExperimental()
			player.sim.showExperimentalChangeEmitter.on(() => {
				simAllButton.hidden = !player.sim.getShowExperimental();
			});
			simAllButton.addEventListener('click', (_) => {
				if (simUI instanceof IndividualSimUI) {
					let itemSpecs = Array<ItemSpec>();
					const isRangedOrTrinket = this.slot == ItemSlot.ItemSlotRanged ||
						this.slot == ItemSlot.ItemSlotTrinket1 ||
						this.slot == ItemSlot.ItemSlotTrinket2

					const curItem = this.equippedToItemFn(this.player.getEquippedItem(this.slot));
					let curEP = 0;
					if (curItem != null) {
						curEP = this.computeEP(curItem);
					}

					for(let i of this.itemsToDisplay) {
						const idata = this.itemData[i];
						if (!isRangedOrTrinket && curEP > 0 && idata.baseEP < (curEP / 2)) {
							continue; // If we have EPs on current item, dont sim items with less than half the EP.
						}

						// Add any item that is either >0 EP or a trinket/ranged item.
						if (idata.baseEP > 0 || isRangedOrTrinket) {
							itemSpecs.push(ItemSpec.create({ id: idata.id }));
						}

					}
					simUI.bt.addItems(itemSpecs);
					// TODO: should we open the bulk sim UI or should we run in the background showing progress, and then sort the items in the picker?
				}
			});
		} else {
			// always hide non-items from being added to batch.
			simAllButton.hidden = true;
		}
	}

	public sizeRefresh() {
		this.scroller.refresh(true);
		this.applyFilters();
	}

	public dispose() {
		this.scroller.dispose();
	}

	public updateSelected() {
		const newEquippedItem = this.gearData.getEquippedItem();
		const newItem = this.equippedToItemFn(newEquippedItem);

		const newItemId = newItem ? (this.label == 'Enchants' ? (newItem as unknown as Enchant).effectId : (newItem as unknown as Item | Rune).id) : 0;
		const newEP = newItem ? this.computeEP(newItem) : 0;

		this.scroller.elementUpdate((item) => {
			let idx = (item as HTMLElement).dataset.idx!;
			const itemData = this.itemData[parseFloat(idx)];
			if (itemData.id == newItemId)
				item.classList.add('active');
			else
				item.classList.remove('active');

			const epDeltaElem = item.getElementsByClassName('selector-modal-list-item-ep-delta')[0] as HTMLSpanElement;
			if (epDeltaElem) {
				epDeltaElem.textContent = '';
				if (itemData.item) {
					const listItemEP = this.computeEP(itemData.item);
					formatDeltaTextElem(epDeltaElem, newEP, listItemEP, 0);
				}
			}
		});
	}

	public applyFilters() {
		this.currentFilters = this.player.sim.getFilters();
		let itemIdxs = new Array<number>(this.itemData.length);
		for (let i = 0; i < this.itemData.length; ++i)  {
			itemIdxs[i] = i;
		}

		const currentEquippedItem = this.player.getEquippedItem(this.slot);

		if (this.label == 'Items') {
			itemIdxs = this.player.filterItemData(
				itemIdxs,
				i => this.itemData[i].item as unknown as Item,
				this.slot);
		} else if (this.label == 'Enchants') {
			itemIdxs = this.player.filterEnchantData(
				itemIdxs,
				i => this.itemData[i].item as unknown as Enchant,
				this.slot,
				currentEquippedItem);
		}

		itemIdxs = itemIdxs.filter(i => {
			const listItemData = this.itemData[i];

			if (listItemData.phase > this.player.sim.getPhase()) {
				return false;
			}

			if (this.searchInput.value.length > 0) {
				const searchQuery = this.searchInput.value.toLowerCase().replaceAll(/[^a-zA-Z0-9\s]/g, '').split(" ");
				const name = listItemData.name.toLowerCase().replaceAll(/[^a-zA-Z0-9\s]/g, '');

				var include = true;
				searchQuery.forEach(v => {
					if (!name.includes(v))
						include = false;
				});
				if (!include) {
					return false;
				}
			}

			return true;
		});

		let sortFn: (itemA: T, itemB: T) => number;
		if (this.slot == ItemSlot.ItemSlotTrinket1 || this.slot == ItemSlot.ItemSlotTrinket2) {
			// Trinket EP is weird so just sort by ilvl instead.
			sortFn = (itemA, itemB) => (itemB as unknown as Item).ilvl - (itemA as unknown as Item).ilvl;
		} else {
			sortFn = (itemA, itemB) => {
				const diff = this.computeEP(itemB) - this.computeEP(itemA);
				// if EP is same, sort by ilvl
				if (Math.abs(diff) < 0.01) return (itemB as unknown as Item).ilvl - (itemA as unknown as Item).ilvl;
				return diff;
			}
		}

		itemIdxs = itemIdxs.sort((dataA, dataB) => {
			const itemA = this.itemData[dataA];
			const itemB = this.itemData[dataB];
			if (this.isItemFavorited(itemA) && !this.isItemFavorited(itemB)) return -1;
			if (this.isItemFavorited(itemB) && !this.isItemFavorited(itemA)) return 1;

			return sortFn(itemA.item, itemB.item);
		});

		this.itemsToDisplay = itemIdxs;
		this.scroller.update();

		this.hideOrShowEPValues();
	}

	public hideOrShowEPValues() {
		const labels = this.tabContent.getElementsByClassName("ep-delta-label")
		const container = this.tabContent.getElementsByClassName("selector-modal-list")
		const show = this.player.sim.getShowEPValues();
		const display = show ? "" : "none"

		for (let label of labels) {
			(label as HTMLElement).style.display = display;
		}

		for (let c of container) {
			if (show)
				c.classList.remove("hide-ep");
			else
				c.classList.add("hide-ep");
		}
	}

	private createItemElem(item: ItemDataWithIdx<T>): JSX.Element {
		const itemData = item.data;
		const itemEP = this.computeEP(itemData.item);

		const equipedItem = this.equippedToItemFn(this.gearData.getEquippedItem());
		const equipdItemId = equipedItem ? (this.label == 'Enchants' ? (equipedItem as unknown as Enchant).effectId : (equipedItem as unknown as Item).id) : 0;

		const nameElem = ref<HTMLLabelElement>();
		const anchorElem = ref<HTMLAnchorElement>();
		const iconElem = ref<HTMLImageElement>();
		const listItemElem = (
			<li className={`selector-modal-list-item ${equipdItemId == itemData.id ? 'active' : ''}`} dataset={{idx: item.idx.toString()}}>
				<div className='selector-modal-list-label-cell'>
					<a className='selector-modal-list-item-link' ref={anchorElem} dataset={{whtticon:'false'}}>
						<img className='selector-modal-list-item-icon' ref={iconElem}></img>
						<label className='selector-modal-list-item-name' ref={nameElem}>
							{itemData.name}
							{itemData.heroic && createHeroicLabel()}
						</label>
					</a>
				</div>
			</li>
		);

		if (this.label == 'Items') {
			listItemElem.appendChild(
				<div className='selector-modal-list-item-source-container'>
					{this.getSourceInfo(itemData.item as unknown as Item, this.player.sim)}
				</div>
			)
		}

		let favoriteElem = ref<HTMLButtonElement>();

		listItemElem.appendChild(
			<div>
				<button className="selector-modal-list-item-favorite btn btn-link p-0"
						ref={favoriteElem}
						onclick={() => setFavorite(listItemElem.dataset.fav == 'false')}>
					<i className='fa-star fa-xl'></i>
				</button>
			</div>
		)


		if (this.slot != ItemSlot.ItemSlotTrinket1 && this.slot != ItemSlot.ItemSlotTrinket2) {
			listItemElem.appendChild(
				<div className='selector-modal-list-item-ep'>
					<span className='selector-modal-list-item-ep-value'>
						{itemEP < 9.95 ? itemEP.toFixed(1).toString() : Math.round(itemEP).toString()}
					</span>
				</div>
			);
		}

		listItemElem.appendChild(
			<div className='selector-modal-list-item-ep'>
				<span className='selector-modal-list-item-ep-delta'
					ref={(e) => itemData.item && formatDeltaTextElem(e, equipedItem ? this.computeEP(equipedItem) : 0, itemEP, 0)}></span>
			</div>
		);

		anchorElem.value!.addEventListener('click', (event: Event) => {
			event.preventDefault();
			if (event.target === favoriteElem.value) return false;
			this.onItemClick(itemData);
		});

		itemData.actionId.fill().then(filledId => {
			filledId.setWowheadHref(anchorElem.value!);
			iconElem.value!.src = filledId.iconUrl;
		});

		setItemQualityCssClass(nameElem.value!, itemData.quality);

		new Tooltip(favoriteElem.value!, {
			title: 'Add to favorites'
		});
		const setFavorite = (isFavorite: boolean) => {
			const filters = this.player.sim.getFilters();
			if (this.label == 'Items') {
				const favId = itemData.id;
				if (isFavorite) {
					filters.favoriteItems.push(favId);
				} else {
					const favIdx = filters.favoriteItems.indexOf(favId);
					if (favIdx != -1) {
						filters.favoriteItems.splice(favIdx, 1);
					}
				}
			} else if (this.label == 'Enchants') {
				const favId = getUniqueEnchantString(itemData.item as unknown as Enchant);
				if (isFavorite) {
					filters.favoriteEnchants.push(favId);
				} else {
					const favIdx = filters.favoriteEnchants.indexOf(favId);
					if (favIdx != -1) {
						filters.favoriteEnchants.splice(favIdx, 1);
					}
				}
			}
			favoriteElem.value!.children[0].classList.toggle('fas');
			favoriteElem.value!.children[0].classList.toggle('far');
			listItemElem.dataset.fav = isFavorite.toString();

			this.player.sim.setFilters(TypedEvent.nextEventID(), filters);
		};

		let isFavorite = this.isItemFavorited(itemData);
		
		if (isFavorite) {
			favoriteElem.value!.children[0].classList.add('fas');
			listItemElem.dataset.fav = 'true';
		} else {
			favoriteElem.value!.children[0].classList.add('far');
			listItemElem.dataset.fav = 'false';
		}

		return listItemElem;
	}

	private isItemFavorited(itemData: ItemData<T>) : boolean {
		if (this.label == 'Items') {
			return this.currentFilters.favoriteItems.includes(itemData.id);
		} else if (this.label == 'Enchants') {
			return this.currentFilters.favoriteEnchants.includes(getUniqueEnchantString(itemData.item as unknown as Enchant));
		}
		return false;
	}

	private getSourceInfo(item: Item, sim: Sim): JSX.Element {
		if (!item.sources || item.sources.length == 0) {
			return <></>;
		}

		const makeAnchor = (href:string, inner:string) => {
			return <a href={href}><small>{inner}</small></a>;
		}

		const source = item.sources[0];
		if (source.source.oneofKind == 'crafted') {
			const src = source.source.crafted;
			return makeAnchor( ActionId.makeSpellUrl(src.spellId), professionNames.get(src.profession) ?? 'Unknown');
		} else if (source.source.oneofKind == 'drop') {
			const src = source.source.drop;
			const zone = sim.db.getZone(src.zoneId);
			const npc = sim.db.getNpc(src.npcId);
			if (!zone) {
				throw new Error('No zone found for item: ' + item);
			}

			let rtnEl = makeAnchor( ActionId.makeZoneUrl(zone.id), `${zone.name} (${difficultyNames.get(src.difficulty) ?? 'Unknown'})`);

			const category = src.category ? ` - ${src.category}` : '';
			if (npc) {
				rtnEl.appendChild(document.createElement('br'));
				rtnEl.appendChild(makeAnchor(ActionId.makeNpcUrl(npc.id), `${npc.name + category}`));
			} else if (src.otherName) {
				/*innerHTML += `
					<br>
					<a href="${ActionId.makeZoneUrl(zone.id)}"><small>${src.otherName + category}</small></a>
				`;*/
			} else if (category) {
				/*innerHTML += `
					<br>
					<a href="${ActionId.makeZoneUrl(zone.id)}"><small>${category}</small></a>
				`;*/
			}
			return rtnEl;
		} else if (source.source.oneofKind == 'quest') {
			const src = source.source.quest;
			return makeAnchor(ActionId.makeQuestUrl(src.id), src.name);
		} else if (source.source.oneofKind == 'soldBy') {
			const src = source.source.soldBy;
			return makeAnchor(ActionId.makeNpcUrl(src.npcId), src.npcName);
		}
		return <></>;
	}
}
