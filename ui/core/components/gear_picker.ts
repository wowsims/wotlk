import { difficultyNames, professionNames, slotNames } from '../proto_utils/names.js';
import { BaseModal } from './base_modal';
import { Component } from './component';
import { FiltersMenu } from './filters_menu';
import { Input, InputConfig } from './input';
import {
	makePhaseSelector,
	makeShow1hWeaponsSelector,
	makeShow2hWeaponsSelector,
	makeShowMatchingGemsSelector,
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
import { ItemSwapGear } from '../proto_utils/gear'
import { getEmptyGemSocketIconUrl, gemMatchesSocket } from '../proto_utils/gems';
import { Stats } from '../proto_utils/stats';

import {
	Class,
	Spec,
	GemColor,
	ItemQuality,
	ItemSlot,
	ItemSpec,
	ItemSwap,
	ItemType,
} from '../proto/common';
import {
	UIEnchant as Enchant,
	UIGem as Gem,
	UIItem as Item,
} from '../proto/ui.js';

declare var tippy: any;
declare var WowSim: any;

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

class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;
	private readonly iconElem: HTMLAnchorElement;
	private readonly nameElem: HTMLAnchorElement;
	private readonly enchantElem: HTMLAnchorElement;
	private readonly socketsContainerElem: HTMLElement;
	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];

	private _equippedItem: EquippedItem | null = null;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');
		this.slot = slot;
		this.simUI = simUI;
		this.player = player;

		this.rootElem.innerHTML = `
      <a class="item-picker-icon" href="javascript:void(0)" role="button">
        <div class="item-picker-sockets-container"></div>
      </a>
      <div class="item-picker-labels-container">
        <a class="item-picker-name" href="javascript:void(0)" role="button"></a><br>
        <a class="item-picker-enchant" href="javascript:void(0)" role="button"></a>
      </div>
    `;

		this.iconElem = this.rootElem.querySelector('.item-picker-icon') as HTMLAnchorElement;
		this.nameElem = this.rootElem.querySelector('.item-picker-name') as HTMLAnchorElement;
		this.enchantElem = this.rootElem.querySelector('.item-picker-enchant') as HTMLAnchorElement;
		this.socketsContainerElem = this.rootElem.querySelector('.item-picker-sockets-container') as HTMLElement;

		this.item = player.getEquippedItem(slot);
		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(this.slot);
			this._enchants = this.player.getEnchants(this.slot);

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
			const onClickEnd = (event: Event) => {
				event.preventDefault();
			};

			// Make icon open gear selector
			this.iconElem.addEventListener('click', openGearSelector);
			this.iconElem.addEventListener('touchstart', openGearSelector);
			this.iconElem.addEventListener('touchend', onClickEnd);

			// Make item name open gear selector
			this.nameElem.addEventListener('click', openGearSelector);
			this.nameElem.addEventListener('touchstart', openGearSelector);
			this.nameElem.addEventListener('touchend', onClickEnd);

			// Make enchant name open enchant selector
			this.enchantElem.addEventListener('click', openEnchantSelector);
			this.enchantElem.addEventListener('touchstart', openEnchantSelector);
			this.enchantElem.addEventListener('touchend', onClickEnd);
		});

		player.gearChangeEmitter.on(() => {
			this.item = player.getEquippedItem(slot);
		});
		player.professionChangeEmitter.on(() => {
			if (this._equippedItem != null) {
				this.player.setWowheadData(this._equippedItem, this.iconElem);
			}
		});

		// Use hacky wowhead xhr override to 'preprocess' tooltips
		WowSim.WhOnLoadHook = (a: any) => {
			if (a.tooltip) {
				// This fixes wowhead being able to parse 'pcs' aka set bonus highlighting in tooltip
				// Their internal regex looks for 'href="/item=' but for wotlk we get 'href="/wotlk/item="'
				a.tooltip = (<String>a.tooltip).replaceAll("href=\"/wotlk/item", "href=\"/item");
			}
			return a;
		}
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.nameElem.removeAttribute('data-wowhead');
		this.nameElem.removeAttribute('href');
		this.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconElem.removeAttribute('data-wowhead');
		this.iconElem.removeAttribute('href');
		this.enchantElem.removeAttribute('data-wowhead');
		this.enchantElem.removeAttribute('href');
		this.iconElem.removeAttribute('href');

		this.nameElem.textContent = slotNames[this.slot];
		setItemQualityCssClass(this.nameElem, null);

		this.enchantElem.innerHTML = '';
		this.socketsContainerElem.innerHTML = '';

		if (newItem != null) {
			this.nameElem.textContent = newItem.item.name;
			if (newItem.item.heroic) {
				var heroic_span = document.createElement('span');
				heroic_span.style.color = "green";
				heroic_span.style.marginLeft = "3px";
				heroic_span.innerText = "[H]";
				this.nameElem.appendChild(heroic_span);
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
					this.enchantElem.setAttribute('data-wowhead', `domain=wotlk&spell=${newItem.enchant.spellId}`);
				} else {
					this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
					this.enchantElem.setAttribute('data-wowhead', `domain=wotlk&item=${newItem.enchant.itemId}`);
				}
			}

			newItem.allSocketColors().forEach((socketColor, gemIdx) => {
				let gemFragment = document.createElement('fragment');
				gemFragment.innerHTML = `
					<div class="gem-socket-container">
						<img class="gem-icon" />
						<img class="socket-icon" />
					</div>
				`;

				const gemContainer = gemFragment.children[0] as HTMLElement;
				const gemIconElem = gemContainer.querySelector('.gem-icon') as HTMLImageElement;
				const socketIconElem = gemContainer.querySelector('.socket-icon') as HTMLImageElement;
				socketIconElem.src = getEmptyGemSocketIconUrl(socketColor);

				if (newItem.gems[gemIdx] == null) {
					gemIconElem.classList.add('hide');
				} else {
					gemIconElem.classList.remove('hide');
					ActionId.fromItemId(newItem.gems[gemIdx]!.id).fill().then(filledId => {
						gemIconElem.src = filledId.iconUrl;
					});
				}

				this.socketsContainerElem.appendChild(gemContainer);

				if (gemIdx == newItem.numPossibleSockets - 1 && [ItemType.ItemTypeWrist, ItemType.ItemTypeHands].includes(newItem.item.type)) {
					const updateProfession = () => {
						if (this.player.isBlacksmithing()) {
							gemContainer.classList.remove('hide');
						} else {
							gemContainer.classList.add('hide');
						}
					};
					this.player.professionChangeEmitter.on(updateProfession);
					updateProfession();
				}
			});
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
			gearData: gearData
		})
	}
}

export class IconItemSwapPicker<SpecType extends Spec, ValueType> extends Input<Player<SpecType>, ValueType> {
	private readonly config: InputConfig<Player<SpecType>, ValueType>;
	private readonly iconAnchor: HTMLAnchorElement;
	private readonly socketsContainerElem: HTMLElement;
	private readonly player: Player<SpecType>;
	private readonly slot: ItemSlot;
	private readonly gear: ItemSwapGear;

	// All items and enchants that are eligible for this slot
	private _items: Array<Item> = [];
	private _enchants: Array<Enchant> = [];

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<SpecType>, slot: ItemSlot, config: InputConfig<Player<SpecType>, ValueType>) {
		super(parent, 'icon-picker-root', player, config)
		this.rootElem.classList.add('icon-picker');
		this.player = player;
		this.config = config;
		this.slot = slot;
		this.gear = this.player.getItemSwapGear();

		this.iconAnchor = document.createElement('a');
		this.iconAnchor.classList.add('icon-picker-button');
		this.iconAnchor.target = '_blank';
		this.rootElem.prepend(this.iconAnchor);

		this.socketsContainerElem = document.createElement('div')
		this.socketsContainerElem.classList.add('item-picker-sockets-container')
		this.iconAnchor.appendChild(this.socketsContainerElem);

		player.sim.waitForInit().then(() => {
			this._items = this.player.getItems(slot);
			this._enchants = this.player.getEnchants(slot);
			this.addItemSpecToGear();
			const gearData = {
				equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
					this.gear.equipItem(this.slot, equippedItem, player.canDualWield2H());
					this.inputChanged(eventID);
				},
				getEquippedItem: () => this.gear.getEquippedItem(this.slot),
				changeEvent: config.changedEvent(player),
			}

			const onClickStart = (event: Event) => {
				event.preventDefault();
				new SelectorModal(simUI.rootElem, simUI, this.player, {
					selectedTab: SelectorModalTabs.Items,
					slot: this.slot,
					equippedItem: this.gear.getEquippedItem(slot),
					eligibleItems: this._items,
					eligibleEnchants: this._enchants,
					gearData: gearData,
				})
			};

			this.iconAnchor.addEventListener('click', onClickStart);
			this.iconAnchor.addEventListener('touchstart', onClickStart);
		}).finally(() => this.init());

		// Use hacky wowhead xhr override to 'preprocess' tooltips
		WowSim.WhOnLoadHook = (a: any) => {
			if (a.tooltip) {
				// This fixes wowhead being able to parse 'pcs' aka set bonus highlighting in tooltip
				// Their internal regex looks for 'href="/item=' but for wotlk we get 'href="/wotlk/item="'
				a.tooltip = (<String>a.tooltip).replaceAll("href=\"/wotlk/item", "href=\"/item");
			}
			return a;
		}

	}

	private addItemSpecToGear() {
		const itemSwap = this.config.getValue(this.player) as unknown as ItemSwap
		const fieldName = this.getFieldNameFromItemSlot(this.slot) 

		if (!fieldName)
			return;

		const itemSpec = itemSwap[fieldName] as unknown as ItemSpec
		
		if (!itemSpec)
			return;

		const equippedItem = this.player.sim.db.lookupItemSpec(itemSpec);

		if (equippedItem) {
			this.gear.equipItem(this.slot, equippedItem, this.player.canDualWield2H());
		}
	}

	private getFieldNameFromItemSlot(slot: ItemSlot): keyof ItemSwap | undefined {
		switch (slot) {
			case ItemSlot.ItemSlotMainHand:
				return 'mhItem';
			case ItemSlot.ItemSlotOffHand:
				return 'ohItem';
			case ItemSlot.ItemSlotRanged:
				return 'rangedItem';
		}

		return undefined;
	}

	getInputElem(): HTMLElement {
		return this.iconAnchor;
	}
	getInputValue(): ValueType {
		return this.gear.toProto() as unknown as ValueType
	}

	setInputValue(newValue: ValueType): void {
		this.iconAnchor.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		this.iconAnchor.removeAttribute('data-wowhead');
		this.iconAnchor.href = "#";
		this.socketsContainerElem.innerHTML = '';

		const equippedItem = this.gear.getEquippedItem(this.slot);
		if (equippedItem) {
			this.iconAnchor.classList.add("active")

			equippedItem.asActionId().fillAndSet(this.iconAnchor, true, true);
			this.player.setWowheadData(equippedItem, this.iconAnchor);

			equippedItem.allSocketColors().forEach((socketColor, gemIdx) => {
				const gemIconElem = document.createElement('img');
				gemIconElem.classList.add('item-picker-gem-icon');

				if (equippedItem!.gems[gemIdx] == null) {
					gemIconElem.src = getEmptyGemSocketIconUrl(socketColor);
				} else {
					ActionId.fromItemId(equippedItem!.gems[gemIdx]!.id).fill().then(filledId => {
						gemIconElem.src = filledId.iconUrl;
					});
				}
				this.socketsContainerElem.appendChild(gemIconElem);
			});

		} else {
			this.iconAnchor.classList.remove("active")
		}
	}

}

interface GearData {
	equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => void,
	getEquippedItem: () => EquippedItem | null,
	changeEvent: TypedEvent<any>,
}

enum SelectorModalTabs {
	Items = 'Items',
	Enchants = 'Enchants',
	Gem1 = 'Gem1',
	Gem2 = 'Gem2',
	Gem3 = 'Gem3',
}

interface SelectorModalConfig {
	selectedTab: SelectorModalTabs
	slot: ItemSlot,
	equippedItem: EquippedItem | null,
	eligibleItems: Array<Item>,
	eligibleEnchants: Array<Enchant>,
	gearData: GearData
}

class SelectorModal extends BaseModal {
	private readonly simUI: SimUI;
	private player: Player<any>;
	private config: SelectorModalConfig;

	private readonly tabsElem: HTMLElement;
	private readonly contentElem: HTMLElement;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>, config: SelectorModalConfig) {
		super(parent, 'selector-modal');

		this.simUI = simUI;
		this.player = player;
		this.config = config;

		window.scrollTo({top: 0});

		this.header!.insertAdjacentHTML('afterbegin', `<ul class="nav nav-tabs selector-modal-tabs"></ul>`);

		this.body.innerHTML = `<div class="tab-content selector-modal-tab-content"></div>`

		this.tabsElem = this.rootElem.querySelector('.selector-modal-tabs') as HTMLElement;
		this.contentElem = this.rootElem.querySelector('.selector-modal-tab-content') as HTMLElement;

		this.setData();
	}

	openTab(idx: number) {
		const elems = this.tabsElem.getElementsByClassName("selector-modal-item-tab");
		(elems[idx] as HTMLElement).click();
	}

	setData() {
		this.tabsElem.innerHTML = '';
		this.contentElem.innerHTML = '';

		const {slot, equippedItem, eligibleItems, eligibleEnchants, gearData } = this.config;

		this.addTab(
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
							gearData.equipItem(eventID, new EquippedItem(item));
						}
					},
				};
			}),
			item => this.player.computeItemEP(item, slot),
			equippedItem => equippedItem?.item,
			GemColor.GemColorUnknown,
			eventID => {
				gearData.equipItem(eventID, null);
			});

		this.addTab(
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
			GemColor.GemColorUnknown,
			eventID => {
				const equippedItem = gearData.getEquippedItem();
				if (equippedItem)
					gearData.equipItem(eventID, equippedItem.withEnchant(null));
			});

		this.addGemTabs(slot, equippedItem, gearData);
	}

	private addGemTabs(slot: ItemSlot, equippedItem: EquippedItem | null, gearData: GearData) {
		if (equippedItem == undefined) {
			return;
		}

		const socketBonusEP = this.player.computeStatsEP(new Stats(equippedItem.item.socketBonus)) / (equippedItem.item.gemSockets.length || 1);
		equippedItem.curSocketColors(this.player.isBlacksmithing()).forEach((socketColor, socketIdx) => {
			this.addTab(
				'Gem ' + (socketIdx + 1),
				this.player.getGems(socketColor).map((gem: Gem) => {
					return {
						item: gem,
						id: gem.id,
						actionId: ActionId.fromItemId(gem.id),
						name: gem.name,
						quality: gem.quality,
						phase: gem.phase,
						heroic: false,
						baseEP: this.player.computeStatsEP(new Stats(gem.stats)),
						ignoreEPFilter: true,
						onEquip: (eventID, gem: Gem) => {
							const equippedItem = gearData.getEquippedItem();
							if (equippedItem)
								gearData.equipItem(eventID, equippedItem.withGem(gem, socketIdx));
						},
					};
				}),
				gem => {
					let gemEP = this.player.computeGemEP(gem);
					if (gemMatchesSocket(gem, socketColor)) {
						gemEP += socketBonusEP;
					}
					return gemEP;
				},
				equippedItem => equippedItem?.gems[socketIdx],
				socketColor,
				eventID => {
					const equippedItem = gearData.getEquippedItem();
					if (equippedItem)
					gearData.equipItem(eventID, equippedItem.withGem(null, socketIdx));
				},
				tabAnchor => {
					tabAnchor.classList.add('selector-modal-tab-gem');
					tabAnchor.innerHTML = `
						<div class="gem-socket-container">
							<img class="gem-icon" />
							<img class="socket-icon" />
						</div>
					`;

					const gemElem = tabAnchor.querySelector('.gem-icon') as HTMLElement;
					const socketElem = tabAnchor.querySelector('.socket-icon') as HTMLElement;
					socketElem.setAttribute('src', getEmptyGemSocketIconUrl(socketColor));

					const updateGemIcon = () => {
						const equippedItem = gearData.getEquippedItem();
						const gem = equippedItem?.gems[socketIdx];

						if (gem) {
							ActionId.fromItemId(gem.id).fill().then(filledId => {
								gemElem.setAttribute('src', filledId.iconUrl);
							});
						}
					};

					gearData.changeEvent.on(updateGemIcon);
					this.addOnDisposeCallback(() => gearData.changeEvent.off(updateGemIcon));
					updateGemIcon();
				});
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
		itemData: Array<ItemData<T>>,
		computeEP: (item: T) => number,
		equippedToItemFn: (equippedItem: EquippedItem | null) => (T | null | undefined),
		socketColor: GemColor,
		onRemove: (eventID: EventID) => void,
		setTabContent?: (tabElem: HTMLAnchorElement) => void) {
		if (itemData.length == 0) {
			return;
		}

		const { slot, gearData } = this.config;

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
		const selected = label === this.config.selectedTab;

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
				<div class="selector-modal-filters">
					<input class="selector-modal-search form-control" type="text" placeholder="Search...">
					${label == 'Items' ? '<button class="selector-modal-filters-button btn btn-primary">Filters</button>' : ''}
					<div class="selector-modal-phase-selector"></div>
					<div class="sim-input selector-modal-boolean-option selector-modal-show-1h-weapons"></div>
					<div class="sim-input selector-modal-boolean-option selector-modal-show-2h-weapons"></div>
					<div class="sim-input selector-modal-boolean-option selector-modal-show-matching-gems"></div>
					<button class="selector-modal-remove-button btn btn-danger">Unequip Item</button>
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

		if (label == 'Items') {
			const filtersButton = tabContent.getElementsByClassName('selector-modal-filters-button')[0] as HTMLElement;
			filtersButton.addEventListener('click', () => new FiltersMenu(this.body, this.player, slot));
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
				<div class="selector-modal-list-item-source-container">
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

			const sourceElem = listItemElem.getElementsByClassName('selector-modal-list-item-source-container')[0] as HTMLDivElement;
			if (label == 'Items') {
				this.fillSourceInfo(item as unknown as Item, sourceElem, this.player.sim);
			} else {
				sourceElem.remove();
			}

			const onclick = (event: Event) => {
				event.preventDefault();
				itemData.onEquip(TypedEvent.nextEventID(), item);

				// If the item changes, the gem slots might change, so remove and recreate the gem tabs
				if (Item.is(item)) {
					this.removeTabs('Gem');
					this.addGemTabs(slot, gearData.getEquippedItem(), gearData);
				}
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
			const newEquippedItem = gearData.getEquippedItem();
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
		gearData.changeEvent.on(updateSelected);
		this.addOnDisposeCallback(() => gearData.changeEvent.off(updateSelected));
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
					//const searchQuery = searchInput.value.toLowerCase().replaceAll(/[^a-zA-Z0-9\s]/g, '').split(" ");
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
		gearData.changeEvent.on(applyFilters);
		this.addOnDisposeCallback(() => {
			this.player.sim.phaseChangeEmitter.off(applyFilters);
			this.player.sim.filtersChangeEmitter.off(applyFilters);
			gearData.changeEvent.off(applyFilters);
		});

		applyFilters();
	}

	private fillSourceInfo(item: Item, container: HTMLDivElement, sim: Sim) {
		if (!item.sources || item.sources.length == 0) {
			return;
		}

		const source = item.sources[0];
		if (source.source.oneofKind == 'crafted') {
			const src = source.source.crafted;
			container.innerHTML = `
				<a href="${ActionId.makeSpellUrl(src.spellId)}">${professionNames[src.profession]}</a>
			`;
		} else if (source.source.oneofKind == 'drop') {
			const src = source.source.drop;
			const zone = sim.db.getZone(src.zoneId);
			const npc = sim.db.getNpc(src.npcId);
			if (!zone) {
				throw new Error('No zone found for item: ' + item);
			}

			let innerHTML = `
				<a href="${ActionId.makeZoneUrl(zone.id)}">${zone.name} (${difficultyNames[src.difficulty]})</a>
			`;

			const category = src.category ? ` - ${src.category}` : '';
			if (npc) {
				innerHTML += `
					<br>
					<a href="${ActionId.makeNpcUrl(npc.id)}">${npc.name + category}</a>
				`;
			} else if (src.otherName) {
				innerHTML += `
					<br>
					<a href="${ActionId.makeZoneUrl(zone.id)}>${src.otherName + category}</a>
				`;
			} else if (category) {
				innerHTML += `
					<br>
					<a href="${ActionId.makeZoneUrl(zone.id)}>${category}</a>
				`;
			}
			container.innerHTML = innerHTML;
		} else if (source.source.oneofKind == 'quest') {
			const src = source.source.quest;
			container.innerHTML = `
				<a href="${ActionId.makeQuestUrl(src.id)}">${src.name}</a>
			`;
		} else if (source.source.oneofKind == 'soldBy') {
			const src = source.source.soldBy;
			container.innerHTML = `
				<a href="${ActionId.makeNpcUrl(src.npcId)}">${src.npcName}</a>
			`;
		}
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
	[ItemSlot.ItemSlotHead]: '/wotlk/assets/item_slots/head.jpg',
	[ItemSlot.ItemSlotNeck]: '/wotlk/assets/item_slots/neck.jpg',
	[ItemSlot.ItemSlotShoulder]: '/wotlk/assets/item_slots/shoulders.jpg',
	[ItemSlot.ItemSlotBack]: '/wotlk/assets/item_slots/shirt.jpg',
	[ItemSlot.ItemSlotChest]: '/wotlk/assets/item_slots/chest.jpg',
	[ItemSlot.ItemSlotWrist]: '/wotlk/assets/item_slots/wrists.jpg',
	[ItemSlot.ItemSlotHands]: '/wotlk/assets/item_slots/hands.jpg',
	[ItemSlot.ItemSlotWaist]: '/wotlk/assets/item_slots/waist.jpg',
	[ItemSlot.ItemSlotLegs]: '/wotlk/assets/item_slots/legs.jpg',
	[ItemSlot.ItemSlotFeet]: '/wotlk/assets/item_slots/feet.jpg',
	[ItemSlot.ItemSlotFinger1]: '/wotlk/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotFinger2]: '/wotlk/assets/item_slots/finger.jpg',
	[ItemSlot.ItemSlotTrinket1]: '/wotlk/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotTrinket2]: '/wotlk/assets/item_slots/trinket.jpg',
	[ItemSlot.ItemSlotMainHand]: '/wotlk/assets/item_slots/mainhand.jpg',
	[ItemSlot.ItemSlotOffHand]: '/wotlk/assets/item_slots/offhand.jpg',
	[ItemSlot.ItemSlotRanged]: '/wotlk/assets/item_slots/ranged.jpg',
};
export function getEmptySlotIconUrl(slot: ItemSlot): string {
	return emptySlotIcons[slot];
}
