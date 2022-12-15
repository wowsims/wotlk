import { EquippedItem } from '../proto_utils/equipped_item.js';
import { getEmptyGemSocketIconUrl } from '../proto_utils/gems.js';
import { setGemSocketCssClass } from '../proto_utils/gems.js';
import { Class, GemColor, Spec, ItemSwap, ItemSpec } from '../proto/common.js';
import { ItemQuality } from '../proto/common.js';
import { ItemSlot } from '../proto/common.js';
import { getUniqueEnchantString } from '../proto_utils/enchants.js';
import { ActionId } from '../proto_utils/action_id.js';
import { setItemQualityCssClass } from '../css_utils.js';
import { Player } from '../player.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { formatDeltaTextElem } from '../utils.js';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';

import { Input, InputConfig } from './input.js';
import { FiltersMenu } from './filters_menu.js';
import { Popup } from './popup.js';
import { makePhaseSelector } from './other_inputs.js';
import { makeShow1hWeaponsSelector } from './other_inputs.js';
import { makeShow2hWeaponsSelector } from './other_inputs.js';
import { makeShowMatchingGemsSelector } from './other_inputs.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { ContentBlock } from './content_block.js';

declare var tippy: any;
declare var WowSim: any;

export function ItemSwapSection(parentElem: HTMLElement, simUI: IndividualSimUI<Spec.SpecEnhancementShaman>): ContentBlock {
	let contentBlock = new ContentBlock(parentElem, 'item-swap-settings', {
		header: {title: 'Item Swap'}
	});

	let itemSwapContianer = Input.newGroupContainer();
	itemSwapContianer.classList.add('item-swap-inputs-container', 'icon-group');
	contentBlock.bodyElement.appendChild(itemSwapContianer);

	new IconItemSwapPicker(itemSwapContianer, simUI.player, ItemSlot.ItemSlotMainHand, {
		// Returns the event indicating the mapped value has changed.
		changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,

		// Get and set the mapped value.
		getValue: (player: Player<Spec.SpecEnhancementShaman>) => {
			return player.getSpecOptions().weaponSwap?.mhItem
		},
		setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: ItemSpec | undefined) => {
			const options = player.getSpecOptions()
			options.weaponSwap!.mhItem = newValue;
			player.setSpecOptions(eventID, options)
		},
	})

	new IconItemSwapPicker(itemSwapContianer, simUI.player, ItemSlot.ItemSlotMainHand, {
		// Returns the event indicating the mapped value has changed.
		changedEvent: (player: Player<Spec.SpecEnhancementShaman>) => player.specOptionsChangeEmitter,

		// Get and set the mapped value.
		getValue: (player: Player<Spec.SpecEnhancementShaman>) => {
			return player.getSpecOptions().weaponSwap?.mhItem
		},
		setValue: (eventID: EventID, player: Player<Spec.SpecEnhancementShaman>, newValue: ItemSpec | undefined) => {
			const options = player.getSpecOptions()
			options.weaponSwap!.mhItem = newValue;
			player.setSpecOptions(eventID, options)
		},
	})

	return contentBlock
}

class IconItemSwapPicker<SpecType extends Spec, ValueType> extends Input<Player<SpecType>, ValueType> {
	private readonly config: InputConfig<Player<SpecType>, ValueType>;
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly player: Player<SpecType>;
	private readonly slot: number;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];
	private currentValue: ItemSpec; 
	
	constructor(parent: HTMLElement, player: Player<SpecType>, slot: number,  config: InputConfig<Player<SpecType>, ValueType>) {
		super(parent, 'icon-picker-root', player, config)
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.currentValue = config.defaultValue as unknown as ItemSpec;
		this.config = config;
		this.slot = slot;

		this.iconAnchor = document.createElement('a');
		this.iconAnchor.classList.add('icon-picker-button');
		this.iconAnchor.target = '_blank';
		this.rootElem.prepend(this.iconAnchor);

		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(slot);
			this._enchants = this.player.getEnchants(slot);

			this.init();

			const onClickStart = (event: Event) => {
				event.preventDefault();
				let equippedItem = null
				if (this.currentValue) {
					console.log("Hello");
					//console.log(this._items.filter(item => item.id = this.currentValue.id)[0])
					//equippedItem = new EquippedItem(this._items.filter(item => item.id = this.currentValue.id)[0])
				}
				new SelectorModal(this.rootElem.closest('.individual-sim-ui')!, this.player, this.slot, equippedItem, this._items, this._enchants)
			};

			this.iconAnchor.addEventListener('click', onClickStart);
			this.iconAnchor.addEventListener('touchstart', onClickStart);
		});
	}
	
	getInputElem(): HTMLElement {
		return this.iconAnchor;
	}
	getInputValue(): ValueType {
		return this.currentValue as unknown as ValueType;
	}

	setInputValue(newValue: ValueType): void {
		this.iconAnchor.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconAnchor.removeAttribute('data-wowhead');
		this.iconAnchor.href = "#"

		this.currentValue = newValue as unknown as ItemSpec;
		if (this.currentValue) {
			ActionId.fromItemId(this.currentValue.id).fillAndSet(this.iconAnchor, true, true);
			let item = this._items.filter(item => item.id = this.currentValue.id)[0];
			console.log(item)
			//this.player.setWowheadData(equippedItem, this.iconAnchor);
			this.iconAnchor.classList.add("active")
		} else {
			this.iconAnchor.classList.remove("active")
		}
	}

}

class SelectorModal extends Popup {
	private player: Player<Spec.SpecEnhancementShaman>;
	private readonly tabsElem: HTMLElement;
	private readonly contentElem: HTMLElement;

	constructor(parent: HTMLElement, player:  Player<any>, slot: ItemSlot, equippedItem: EquippedItem | null, eligibleItems: Array<Item>, eligibleEnchants: Array<Enchant>) {
		super(parent);
		this.player = player;

		this.rootElem.classList.add('selector-modal');
		this.rootElem.innerHTML = `
			<ul class="nav nav-tabs selector-modal-tabs">
			</ul>
			<div class="tab-content selector-modal-tab-content">
			</div>
		`;

		this.addCloseButton();
		this.tabsElem = this.rootElem.getElementsByClassName('selector-modal-tabs')[0] as HTMLElement;
		this.contentElem = this.rootElem.getElementsByClassName('selector-modal-tab-content')[0] as HTMLElement;

		this.setData(slot, equippedItem, eligibleItems, eligibleEnchants);
	}

	openTab(idx: number) {
		const elems = this.tabsElem.getElementsByClassName("selector-modal-item-tab");
		(elems[idx] as HTMLElement).click();
	}

	setData(slot: ItemSlot, equippedItem: EquippedItem | null, eligibleItems: Array<Item>, eligibleEnchants: Array<Enchant>) {
		this.tabsElem.innerHTML = '';
		this.contentElem.innerHTML = '';

		this.addTab(
			'Items',
			slot,
			equippedItem,
			eligibleItems.map(item => {
				return {
					item: item,
					id: item.id,
					actionId: ActionId.fromItem(item),
					name: item.name,
					quality: item.quality,
					heroic: item.heroic,
					phase: item.phase,
					baseEP: this.player.computeItemEP(item),
					ignoreEPFilter: false,
					onEquip: (eventID, item: Item) => {
						const equippedItem = new EquippedItem(item)
						const options = this.player.getSpecOptions()
						options.weaponSwap!.mhItem = equippedItem.asSpec()
						this.player.setSpecOptions(eventID, options)
					},
				};
			}),
			item => this.player.computeItemEP(item),
			equippedItem => equippedItem?.item,
			GemColor.GemColorUnknown,
			eventID => {
					const options = this.player.getSpecOptions()
					options.weaponSwap!.mhItem = undefined
					this.player.setSpecOptions(eventID, options)
			});

	}

	/**
	 * Adds one of the tabs for the item selector menu.
	 *
	 * T is expected to be Item, Enchant, or Gem. Tab menus for all 3 looks extremely
	 * similar so this function uses extra functions to do it generically.
	 */
	private addTab<T>(
		label: string,
		slot: ItemSlot,
		equippedItem: EquippedItem | null,
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => (T | null | undefined),
		socketColor: GemColor,
		onRemove: (eventID: EventID) => void,
		setTabContent?: (tabElem: HTMLAnchorElement) => void) {
		if (itemData.length == 0) {
			return;
		}

		if (slot == ItemSlot.ItemSlotTrinket1 || slot == ItemSlot.ItemSlotTrinket2) {
			// Trinket EP is weird so just sort by ilvl instead.
			itemData.sort((dataA, dataB) => (dataB.item as unknown as Item).ilvl - (dataA.item as unknown as Item).ilvl);
		} else {
			itemData.sort((dataA, dataB) => {
				const diff = computeEP(dataB.item) - computeEP(dataA.item);
				// if EP is same, sort by ilvl
				if (Math.abs(diff) < 0.01) {
					return (dataB.item as unknown as Item).ilvl - (dataA.item as unknown as Item).ilvl;
				}
				return diff;
			});
		}

		const tabContentId = (label + '-tab').split(' ').join('');
		const selected = label === 'Items';

		const tabFragment = document.createElement('fragment');
		tabFragment.innerHTML = `
			<li class="nav-item">
				<a
					class="nav-link selector-modal-item-tab ${selected ? 'active' : ''}"
					data-content-id="${tabContentId}"
					data-bs-toggle="tab"
					data-bs-target="#${tabContentId}"
					type="button"
					role="tab"
					aria-controls="${tabContentId}"
					aria-selected="${selected}"
				></a>
			</li>
		`;

		const tabElem = tabFragment.children[0] as HTMLElement;
		const tabAnchor = tabElem.getElementsByClassName('selector-modal-item-tab')[0] as HTMLAnchorElement;
		tabAnchor.dataset.label = label;
		if (setTabContent) {
			setTabContent(tabAnchor);
		} else {
			tabAnchor.textContent = label;
		}

		this.tabsElem.appendChild(tabElem);

		const tabContentFragment = document.createElement('fragment');
		tabContentFragment.innerHTML = `
			<div
				id="${tabContentId}"
				class="selector-modal-tab-pane tab-pane fade ${selected ? 'active show' : ''}"
			>
				<div class="selector-modal-tab-content-header">
					<button class="selector-modal-remove-button sim-button">Remove</button>
					<input class="selector-modal-search" type="text" placeholder="Search...">
					<div class="selector-modal-filter-bar-filler"></div>
					<div class="sim-input selector-modal-boolean-option selector-modal-show-1h-weapons"></div>
					<div class="sim-input selector-modal-boolean-option selector-modal-show-2h-weapons"></div>
					<div class="sim-input selector-modal-boolean-option selector-modal-show-matching-gems"></div>
					<div class="selector-modal-phase-selector"></div>
					<button class="selector-modal-filters-button sim-button">Filters</button>
				</div>
				<div style="width: 100%;height: 30px;font-size: 18px;">
					<span style="float:left">Item</span>
					<span style="float:right">EP(+/-)<span class="ep-help fas fa-search" style="font-size:10px"></span></span>
				</div>
				<ul class="selector-modal-list"></ul>
			</div>
		`;
		
		const tabContent = tabContentFragment.children[0] as HTMLElement;

		this.contentElem.appendChild(tabContent);

		const helpIcon = tabContent.getElementsByClassName("ep-help").item(0);
		tippy(helpIcon, {'content': 'These values are computed using stat weights which can be edited using the "Stat Weights" button.'});
		const show1hWeaponsSelector = makeShow1hWeaponsSelector(tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement, this.player.sim);
		const show2hWeaponsSelector = makeShow2hWeaponsSelector(tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement, this.player.sim);
		if (!(label == 'Items' && (slot == ItemSlot.ItemSlotMainHand || (slot == ItemSlot.ItemSlotOffHand && this.player.getClass() == Class.ClassWarrior)))) {
			(tabContent.getElementsByClassName('selector-modal-show-1h-weapons')[0] as HTMLElement).style.display = 'none';
			(tabContent.getElementsByClassName('selector-modal-show-2h-weapons')[0] as HTMLElement).style.display = 'none';
		}

		const showMatchingGemsSelector = makeShowMatchingGemsSelector(tabContent.getElementsByClassName('selector-modal-show-matching-gems')[0] as HTMLElement, this.player.sim);
		if (!label.startsWith('Gem')) {
			(tabContent.getElementsByClassName('selector-modal-show-matching-gems')[0] as HTMLElement).style.display = 'none';
		}

		const phaseSelector = makePhaseSelector(tabContent.getElementsByClassName('selector-modal-phase-selector')[0] as HTMLElement, this.player.sim);

		const filtersButton = tabContent.getElementsByClassName('selector-modal-filters-button')[0] as HTMLElement;
		if (FiltersMenu.anyFiltersForSlot(slot)) {
			filtersButton.addEventListener('click', () => new FiltersMenu(this.rootElem, this.player, slot));
		} else {
			filtersButton.style.display = 'none';
		}

		if (label == 'Items') {
			tabElem.classList.add('active', 'in');
			tabContent.classList.add('active', 'in');
		}

		const listElem = tabContent.getElementsByClassName('selector-modal-list')[0] as HTMLElement;
		const initialFilters = this.player.sim.getFilters();
		let lastFavElem: HTMLElement|null = null;

		const listItemElems = itemData.map((itemData, itemIdx) => {
			const item = itemData.item;
			const itemEP = computeEP(item);

			const listItemElem = document.createElement('li');
			listItemElem.classList.add('selector-modal-list-item');
			listElem.appendChild(listItemElem);

			listItemElem.dataset.idx = String(itemIdx);

			listItemElem.innerHTML = `
				<div class="selector-modal-list-label-cell">
					<a class="selector-modal-list-item-icon"></a>
					<a class="selector-modal-list-item-name">${itemData.heroic ? itemData.name + "<span style=\"color:green\">[H]</span>" : itemData.name}</a>
				</div>
				<div>
					<span class="selector-modal-list-item-favorite fa-star"></span>
				</div>
				<div class="selector-modal-list-item-ep">
					<span class="selector-modal-list-item-ep-value">${itemEP < 9.95 ? itemEP.toFixed(1) : Math.round(itemEP)}</span>
				</div>
				<div class="selector-modal-list-item-ep">
					<span class="selector-modal-list-item-ep-delta"></span>
				</div>
      `;

			if (slot == ItemSlot.ItemSlotTrinket1 || slot == ItemSlot.ItemSlotTrinket2) {
				const epElem = listItemElem.getElementsByClassName('selector-modal-list-item-ep')[0] as HTMLElement;
				epElem.style.display = 'none';
			}

			const iconElem = listItemElem.getElementsByClassName('selector-modal-list-item-icon')[0] as HTMLAnchorElement;
			const nameElem = listItemElem.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLAnchorElement;
			itemData.actionId.fill().then(filledId => {
				filledId.setWowheadHref(iconElem);
				filledId.setWowheadHref(nameElem);
				iconElem.style.backgroundImage = `url('${filledId.iconUrl}')`;
			});

			setItemQualityCssClass(nameElem, itemData.quality);

			const onclick = (event: Event) => {
				event.preventDefault();
				itemData.onEquip(TypedEvent.nextEventID(), item);

				// If the item changes, the gem slots might change, so remove and recreate the gem tabs
				// if (Item.is(item)) {
				// 	this.removeTabs('Gem');
				// 	this.addGemTabs(slot, this.player.getEquippedItem(slot));
				// }
			};
			nameElem.addEventListener('click', onclick);
			iconElem.addEventListener('click', onclick);

			const favoriteElem = listItemElem.getElementsByClassName('selector-modal-list-item-favorite')[0] as HTMLElement;
			tippy(favoriteElem, {'content': 'Add to Favorites'});
			const setFavorite = (isFavorite: boolean) => {
				const filters = this.player.sim.getFilters();
				if (label == 'Items') {
					const favId = itemData.id;
					if (isFavorite) {
						filters.favoriteItems.push(favId);
					} else {
						const favIdx = filters.favoriteItems.indexOf(favId);
						if (favIdx != -1) {
							filters.favoriteItems.splice(favIdx, 1);
						}
					}
				} else if (label == 'Enchants') {
					const favId = getUniqueEnchantString(item as unknown as Enchant);
					if (isFavorite) {
						filters.favoriteEnchants.push(favId);
					} else {
						const favIdx = filters.favoriteEnchants.indexOf(favId);
						if (favIdx != -1) {
							filters.favoriteEnchants.splice(favIdx, 1);
						}
					}
				} else if (label.startsWith('Gem')) {
					const favId = itemData.id;
					if (isFavorite) {
						filters.favoriteGems.push(favId);
					} else {
						const favIdx = filters.favoriteGems.indexOf(favId);
						if (favIdx != -1) {
							filters.favoriteGems.splice(favIdx, 1);
						}
					}
				}
				this.player.sim.setFilters(TypedEvent.nextEventID(), filters);

				// Reorder and update this element.
				const curItemElems = Array.from(listElem.children) as Array<HTMLElement>;
				if (isFavorite) {
					// Use same sorting order (based on idx) among the favorited elems.
					const nextElem = curItemElems.find(elem => elem.dataset.fav == 'false' || parseInt(elem.dataset.idx!) > itemIdx);
					if (nextElem) {
						listElem.insertBefore(listItemElem, nextElem);
					} else {
						listElem.appendChild(listItemElem);
					}

					favoriteElem.classList.add('fa-solid');
					favoriteElem.classList.remove('fa-regular');
					listItemElem.dataset.fav = 'true';
				} else {
					// Put back in original spot. itemIdx will usually be a very good starting point for the search.
					// Need to search in both directions to handle all cases of favorited elems / itemIdx location.
					let curIdx = itemIdx;
					while (curIdx > 0 && curItemElems[curIdx].dataset.fav == 'false' && parseInt(curItemElems[curIdx].dataset.idx!) > itemIdx) {
						curIdx--;
					}
					while (curIdx < curItemElems.length && (curItemElems[curIdx].dataset.fav == 'true' || parseInt(curItemElems[curIdx].dataset.idx!) < itemIdx)) {
						curIdx++;
					}
					if (curIdx == curItemElems.length) {
						listElem.appendChild(listItemElem);
					} else {
						listElem.insertBefore(listItemElem, curItemElems[curIdx]);
					}

					favoriteElem.classList.remove('fa-solid');
					favoriteElem.classList.add('fa-regular');
					listItemElem.dataset.fav = 'false';
				}
			};
			favoriteElem.addEventListener('click', () => setFavorite(listItemElem.dataset.fav == 'false'));

			let isFavorite = false;
			if (label == 'Items') {
				isFavorite = initialFilters.favoriteItems.includes(itemData.id);
			} else if (label == 'Enchants') {
				isFavorite = initialFilters.favoriteEnchants.includes(getUniqueEnchantString(item as unknown as Enchant));
			} else if (label.startsWith('Gem')) {
				isFavorite = initialFilters.favoriteGems.includes(itemData.id);
			}
			if (isFavorite) {
				favoriteElem.classList.add('fa-solid');
				listItemElem.dataset.fav = 'true';
				if (lastFavElem == null) {
					listElem.prepend(listItemElem);
				} else {
					lastFavElem.after(listItemElem)
				}
				lastFavElem = listItemElem;
			} else {
				favoriteElem.classList.add('fa-regular');
				listItemElem.dataset.fav = 'false';
			}

			return listItemElem;
		});

		const removeButton = tabContent.getElementsByClassName('selector-modal-remove-button')[0] as HTMLButtonElement;
		removeButton.addEventListener('click', event => {
			listItemElems.forEach(elem => elem.classList.remove('active'));
			onRemove(TypedEvent.nextEventID());
		});

		const updateSelected = () => {
			const wepSwap = this.player.getSpecOptions().weaponSwap
			let newEquippedItem = null
			if (wepSwap){
				newEquippedItem = new EquippedItem(itemData.filter(i => i.id == wepSwap.mhItem?.id)[0].item as unknown as Item)
			}
			const newItem = equippedToItemFn(newEquippedItem);

			const newItemId = newItem ? (label == 'Enchants' ? (newItem as unknown as Enchant).effectId : (newItem as unknown as Item|Gem).id) : 0;
			const newEP = newItem ? computeEP(newItem) : 0;

			listItemElems.forEach(elem => {
				const listItemIdx = parseInt(elem.dataset.idx!);
				const listItemData = itemData[listItemIdx];
				const listItem = listItemData.item;

				elem.classList.remove('active');
				if (listItemData.id == newItemId) {
					elem.classList.add('active');
				}

				const epDeltaElem = elem.getElementsByClassName('selector-modal-list-item-ep-delta')[0] as HTMLSpanElement;
				epDeltaElem.textContent = '';
				if (listItem) {
					const listItemEP = computeEP(listItem);
					formatDeltaTextElem(epDeltaElem, newEP, listItemEP, 0);
				}
			});
		};
		this.player.specOptionsChangeEmitter.on(updateSelected);
		this.addOnDisposeCallback(() => this.player.specOptionsChangeEmitter.off(updateSelected));
		updateSelected();

		const applyFilters = () => {
			let validItemElems = listItemElems;
			const currentEquippedItem = this.player.getEquippedItem(slot);

			if (label == 'Items') {
				validItemElems = this.player.filterItemData(
						validItemElems,
						elem => itemData[parseInt(elem.dataset.idx!)].item as unknown as Item,
						slot);
			} else if (label == 'Enchants') {
				validItemElems = this.player.filterEnchantData(
						validItemElems,
						elem => itemData[parseInt(elem.dataset.idx!)].item as unknown as Enchant,
						slot,
						currentEquippedItem);
			} else if (label.startsWith('Gem')) {
				validItemElems = this.player.filterGemData(
						validItemElems,
						elem => itemData[parseInt(elem.dataset.idx!)].item as unknown as Gem,
						slot,
						socketColor);
			}

			validItemElems = validItemElems.filter(elem => {
				const listItemData = itemData[parseInt(elem.dataset.idx!)];

				if (listItemData.phase > this.player.sim.getPhase()) {
					return false;
				}

				if (searchInput.value.length > 0) {
					const searchQuery = searchInput.value.toLowerCase().split(" ");
					const name = listItemData.name.toLowerCase();

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

			let numShown = 0;
			listItemElems.forEach(elem => {
				if (validItemElems.includes(elem)) {
					elem.classList.remove('hidden');
					numShown++;
					if (numShown % 2 == 0) {
						elem.classList.remove('odd');
					} else {
						elem.classList.add('odd');
					}
				} else {
					elem.classList.add('hidden');
				}
			});
		};

		const searchInput = tabContent.getElementsByClassName('selector-modal-search')[0] as HTMLInputElement;
		searchInput.addEventListener('input', applyFilters);
		searchInput.addEventListener("keyup", ev => {
			if (ev.key == "Enter") {
				listItemElems.find(ele => {
					if (ele.classList.contains("hidden")) {
						return false;
					}
					const nameElem = ele.getElementsByClassName('selector-modal-list-item-name')[0] as HTMLElement;
					nameElem.click();
					return true;
				});
			}
		});

		this.player.sim.phaseChangeEmitter.on(applyFilters);
		this.player.sim.filtersChangeEmitter.on(applyFilters);
		this.player.specOptionsChangeEmitter.on(applyFilters);
		this.addOnDisposeCallback(() => {
			this.player.sim.phaseChangeEmitter.off(applyFilters);
			this.player.sim.filtersChangeEmitter.off(applyFilters);
			this.player.specOptionsChangeEmitter.off(applyFilters);
		});

		applyFilters();
	}

	private removeTabs(labelSubstring: string) {
		const tabElems = Array.prototype.slice.call(this.tabsElem.getElementsByClassName('selector-modal-item-tab'))
			.filter(tab => tab.dataset.label.includes(labelSubstring));

		const contentElems = tabElems
			.map(tabElem => document.getElementById(tabElem.dataset.contentId!.substring(1)))
			.filter(tabElem => Boolean(tabElem));

		tabElems.forEach(elem => elem.parentElement.remove());
		contentElems.forEach(elem => elem!.remove());
	}
}

interface ItemData<T> {
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

const emptySlotIcons: Record<ItemSlot, string> = {
	[ItemSlot.ItemSlotHead]: 'https://cdn.seventyupgrades.com/item-slots/Head.jpg',
	[ItemSlot.ItemSlotNeck]: 'https://cdn.seventyupgrades.com/item-slots/Neck.jpg',
	[ItemSlot.ItemSlotShoulder]: 'https://cdn.seventyupgrades.com/item-slots/Shoulders.jpg',
	[ItemSlot.ItemSlotBack]: 'https://cdn.seventyupgrades.com/item-slots/Back.jpg',
	[ItemSlot.ItemSlotChest]: 'https://cdn.seventyupgrades.com/item-slots/Chest.jpg',
	[ItemSlot.ItemSlotWrist]: 'https://cdn.seventyupgrades.com/item-slots/Wrists.jpg',
	[ItemSlot.ItemSlotHands]: 'https://cdn.seventyupgrades.com/item-slots/Hands.jpg',
	[ItemSlot.ItemSlotWaist]: 'https://cdn.seventyupgrades.com/item-slots/Waist.jpg',
	[ItemSlot.ItemSlotLegs]: 'https://cdn.seventyupgrades.com/item-slots/Legs.jpg',
	[ItemSlot.ItemSlotFeet]: 'https://cdn.seventyupgrades.com/item-slots/Feet.jpg',
	[ItemSlot.ItemSlotFinger1]: 'https://cdn.seventyupgrades.com/item-slots/Finger.jpg',
	[ItemSlot.ItemSlotFinger2]: 'https://cdn.seventyupgrades.com/item-slots/Finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: 'https://cdn.seventyupgrades.com/item-slots/Trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: 'https://cdn.seventyupgrades.com/item-slots/Trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: 'https://cdn.seventyupgrades.com/item-slots/MainHand.jpg',
	[ItemSlot.ItemSlotOffHand]: 'https://cdn.seventyupgrades.com/item-slots/OffHand.jpg',
	[ItemSlot.ItemSlotRanged]: 'https://cdn.seventyupgrades.com/item-slots/Ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}
