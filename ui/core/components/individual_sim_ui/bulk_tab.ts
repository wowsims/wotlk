import { ContentBlock } from "../content_block";
import { Database } from '../../proto_utils/database';
import { Importer } from "../importers";

import { IndividualSimUI } from "../../individual_sim_ui";
import { TypedEvent } from "../../typed_event";

import { EquipmentSpec, ItemSpec, SimDatabase, SimEnchant, SimGem, SimItem, Spec } from "../../proto/common";
import { BulkComboResult, BulkSettings, ItemSpecWithSlot, ProgressMetrics, RaidSimResult } from "../../proto/api";

import { ItemRenderer } from "../gear_picker";
import { SimTab } from "../sim_tab";

import { UIEnchant, UIGem, UIItem } from "../../proto/ui";

export class BulkGearJsonImporter<SpecType extends Spec> extends Importer {
  private readonly simUI: IndividualSimUI<SpecType>;
  private readonly bulkUI: BulkTab
  constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>, bulkUI: BulkTab) {
    super(parent, simUI, 'Bag Item Import', true);
    this.simUI = simUI;
    this.bulkUI = bulkUI;
    this.descriptionElem.innerHTML = `
      <p>Import bag items from a JSON file, which can be created by the WowSimsExporter in-game AddOn.</p>
      <p>To import, upload the file or paste the text below, then click, 'Import'.</p>
    `;
  }

  async onImport(data: string) {
    try {
      const equipment = EquipmentSpec.fromJsonString(data, { ignoreUnknownFields: true });
      if (equipment?.items?.length > 0) {
        const db = await Database.loadLeftoversIfNecessary(equipment);
        const items = equipment.items.filter((spec) => spec.id > 0);
        if (items.length > 0) {
          for (const itemSpec of items) {
            if (itemSpec.id == 0) {
              continue;
            }
            if (!db.lookupItemSpec(itemSpec)) {
              throw new Error("cannot find item with ID " + itemSpec.id);
            }
          }
          this.bulkUI.importItems(items);
        }
      }
      this.close();
    } catch (e: any) {
      alert(e.toString());
    }
  }
}

class BulkSimResultRenderer {

  constructor(parent: ContentBlock, simUI: IndividualSimUI<Spec>, result: BulkComboResult, rank: number, baseResult: BulkComboResult) {
    if (parent.headerElement) {
      parent.headerElement.innerHTML = `Rank ${rank}`;
    }

    const dpsDivParent = document.createElement('div');
    dpsDivParent.classList.add('results-sim');
    parent.bodyElement.appendChild(dpsDivParent);

    const dpsDiv = document.createElement('div');
    dpsDiv.classList.add('bulk-result-body-dps', 'bulk-items-text-line', 'results-sim-dps', 'damage-metrics');
    dpsDivParent.appendChild(dpsDiv);
    
    const dpsNumber = document.createElement('span');
    dpsNumber.textContent = this.formatDps(result.unitMetrics?.dps?.avg!);
    dpsNumber.classList.add('topline-result-avg');
    dpsDiv.appendChild(dpsNumber);

    const dpsDelta = result.unitMetrics?.dps?.avg! - baseResult.unitMetrics?.dps?.avg!;
    const dpsDeltaSpan = document.createElement('span'); 
    dpsDeltaSpan.textContent = `${this.formatDpsDelta(dpsDelta)}`;
    dpsDeltaSpan.classList.add(dpsDelta >= 0 ? 'bulk-result-header-positive' : 'bulk-result-header-negative');
    dpsDiv.appendChild(dpsDeltaSpan);

    const itemsContainer = document.createElement('div');
    itemsContainer.style.flexDirection = 'row';
    itemsContainer.style.display = 'flex';
    parent.bodyElement.appendChild(itemsContainer);

    if (result.itemsAdded && result.itemsAdded.length > 0) {
      for (const is of result.itemsAdded) {
        const item = simUI.sim.db.lookupItemSpec(is.item!)
        const renderer = new ItemRenderer(itemsContainer, simUI, simUI.player);
        renderer.update(item!);
  
        const p = document.createElement('a');
        p.classList.add('bulk-result-item-slot');
        p.textContent = this.itemSlotName(is);
        renderer.nameElem.appendChild(p); 
      }
    } else {
      const p = document.createElement('p');
      p.textContent = 'No changes - this is your currently equipped gear!';
      parent.bodyElement.appendChild(p);
      dpsDeltaSpan.textContent = '';
    }
  }

  private formatDps(dps: number): string {
    return (Math.round(dps * 100) / 100).toFixed(2);
  }

  private formatDpsDelta(delta: number): string {
    return ((delta >= 0) ? "+" : "") + this.formatDps(delta); 
  }

  private itemSlotName(is: ItemSpecWithSlot): string {
    return JSON.parse(ItemSpecWithSlot.toJsonString(is, {emitDefaultValues: true}))['slot'].replace('ItemSlot', '')
  }
}

export class BulkTab extends SimTab {
  readonly simUI: IndividualSimUI<Spec>;
  
	readonly itemsChangedEmitter = new TypedEvent<void>();

  readonly leftPanel: HTMLElement;
  readonly rightPanel: HTMLElement;

  readonly column1: HTMLElement = this.buildColumn(1, 'raid-settings-col');

  protected items: Array<ItemSpec> = new Array<ItemSpec>();

  // TODO: Make a real options probably
  private doCombos: boolean;
  private fastMode: boolean;

  constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
    super(parentElem, simUI, {identifier: 'bulk-tab', title: 'Bulk'});
    this.simUI = simUI;

    this.leftPanel = document.createElement('div');
    this.leftPanel.classList.add('bulk-tab-left', 'tab-panel-left');
    this.leftPanel.appendChild(this.column1);

    this.rightPanel = document.createElement('div');
    this.rightPanel.classList.add('bulk-tab-right', 'tab-panel-right');

    this.contentContainer.appendChild(this.leftPanel);
    this.contentContainer.appendChild(this.rightPanel);

    this.doCombos = true;
    this.fastMode = false;
    this.buildTabContent();
  }

  protected createBulkSettings(): BulkSettings {
    return BulkSettings.create({
      items: this.items,

      // TODO(Riotdog-GehennasEU): Make all of these configurable.
      // For now, it's always constant iteration combinations mode for "sim my bags".
      combinations: this.doCombos,
      fastMode: this.fastMode,
      autoEnchant: false,
      autoGem: false,
      iterationsPerCombo: this.simUI.sim.getIterations(), // TODO(Riotdog-GehennasEU): Define a new UI element for the iteration setting.
    });
  }

  protected createBulkItemsDatabase(): SimDatabase {
    const itemsDb = SimDatabase.create();
    for (const is of this.items) {
      const item = this.simUI.sim.db.lookupItemSpec(is)
      if (!item) {
        throw new Error(`item with ID ${is.id} not found in database`);
      }
      itemsDb.items.push(SimItem.fromJson(UIItem.toJson(item.item), { ignoreUnknownFields: true }))
      if (item.enchant) {
        itemsDb.enchants.push(SimEnchant.fromJson(UIEnchant.toJson(item.enchant), { ignoreUnknownFields: true }));
      }
      for (const gem of item.gems) {
        if (gem) {
          itemsDb.gems.push(SimGem.fromJson(UIGem.toJson(gem), { ignoreUnknownFields: true }));
        }
      }
    }
    return itemsDb;
  }

  importItems(items: Array<ItemSpec>) {
    this.items = items;
    this.itemsChangedEmitter.emit(TypedEvent.nextEventID());
  }

  setCombinations(doCombos: boolean) {
    this.doCombos = doCombos;
  }

  setFastMode(fastMode: boolean) {
    this.fastMode = fastMode;
  }

	protected async runBulkSim(onProgress: Function) {
		try {
			await this.simUI.sim.runBulkSim(this.createBulkSettings(), this.createBulkItemsDatabase(), onProgress);
		} catch (e) {
			this.simUI.handleCrash(e);
		}
	}

  protected buildTabContent() {
    const itemsBlock = new ContentBlock(this.column1, 'bulk-items', {
      header: {title: 'Items'}
    });

    const notice = document.createElement('div');
    notice.classList.add('bulk-items-text-line');
    itemsBlock.bodyElement.appendChild(notice);

    const itemList = document.createElement('div');

    itemList.classList.add('gear-picker-root', 'gear-picker-left', 'tab-panel-col');
    itemList.style.flexDirection = 'row';
    itemList.style.display = 'flex';
    itemsBlock.bodyElement.appendChild(itemList);
    
    let resultsBlock = new ContentBlock(this.column1, 'bulk-results', {header: {
      title: 'Results',
      extraCssClasses: ['bulk-results-header'],
    }});

    resultsBlock.rootElem.hidden = true;
    resultsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-left', 'tab-panel-col');
    
    this.simUI.sim.bulkSimStartEmitter.on(() => {
      resultsBlock.rootElem.hidden = true;
    });

    this.simUI.sim.bulkSimResultEmitter.on((_, bulkSimResult) => {
      resultsBlock.rootElem.hidden = bulkSimResult.results.length == 0;
      resultsBlock.bodyElement.innerHTML = '';

      let rank = 1;
      for (const r of bulkSimResult.results) {
        const resultBlock = new ContentBlock(resultsBlock.bodyElement, 'bulk-result', {header: {title: ''}});
        new BulkSimResultRenderer(resultBlock, this.simUI, r, rank, bulkSimResult.equippedGearResult!);
        rank++;
      }
    });

    const settingsBlock = new ContentBlock(this.rightPanel, 'bulk-settings', {
      header: {title: 'Import'}
    });

    const importButton = document.createElement('button');
    importButton.classList.add('btn', 'btn-primary', 'w-100', 'bulk-settings-button');
    importButton.textContent = 'Import From Bags';
    importButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI, this));
    settingsBlock.bodyElement.appendChild(importButton);

    const bulkSimButton = document.createElement('button');
    bulkSimButton.classList.add('btn', 'btn-primary', 'w-100', 'bulk-settings-button');
    bulkSimButton.textContent = 'Bulk Simulate';
    bulkSimButton.addEventListener('click', () => {
      this.runBulkSim((progressMetrics: ProgressMetrics) => {
        console.log(progressMetrics);
      });
    });
    settingsBlock.bodyElement.appendChild(bulkSimButton);

    const clearButton = document.createElement('button');
    clearButton.classList.add('btn', 'btn-primary', 'w-100', 'bulk-settings-button');
    clearButton.textContent = 'Clear All';
    clearButton.addEventListener('click', () => {
      this.importItems(new Array<ItemSpec>());
      resultsBlock.rootElem.hidden = true;
      resultsBlock.bodyElement.innerHTML = '';
    });
    settingsBlock.bodyElement.appendChild(clearButton);

    this.itemsChangedEmitter.on(() => {
      itemList.innerHTML = '';
      if (this.items.length > 0) {
        notice.textContent = 'The following items will be simmed in all possible combinations together with your equipped gear.';
        for (const spec of this.items) {
          const item = this.simUI.sim.db.lookupItemSpec(spec);
          const itemRenderer = new ItemRenderer(itemList, this.simUI, this.simUI.player);
          itemRenderer.update(item!);
        }
      }
    });
  }
}
