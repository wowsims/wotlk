import { ContentBlock } from "../content_block";
import { Database } from '../../proto_utils/database';
import { Importer } from "../importers";
import { Component } from '../component';

import { IndividualSimUI } from "../../individual_sim_ui";
import { TypedEvent } from "../../typed_event";

import { EquipmentSpec, BulkEquipmentSpec, Spec, BulkEquipmentSpec_ItemSpecWithSlots } from "../../proto/common";

import { ItemRenderer } from "../gear_picker";
import { SimTab } from "../sim_tab";

import { getEligibleItemSlots } from '../../proto_utils/utils.js';
import { BulkSimResultWithSubstitutions, RaidSimResult } from "ui/core/proto/api";

export class BulkGearJsonImporter<SpecType extends Spec> extends Importer {
  private readonly simUI: IndividualSimUI<SpecType>;
  constructor(parent: HTMLElement, simUI: IndividualSimUI<SpecType>) {
    super(parent, simUI, 'Bag Item Import', true);
    this.simUI = simUI;

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
        const bulkEquipment = BulkEquipmentSpec.create();
        for (const itemSpec of equipment.items) {
          if (itemSpec.id == 0) {
            continue;
          }

          const item = db.lookupItemSpec(itemSpec)
          if (!item) {
            throw new Error("cannot find item with ID " + itemSpec.id);
          }

          const itemWithSlot = BulkEquipmentSpec_ItemSpecWithSlots.create();
          itemWithSlot.item = itemSpec;
          itemWithSlot.slots = getEligibleItemSlots(item.item);
          bulkEquipment.items.push(itemWithSlot);
        }

        const eventID = TypedEvent.nextEventID();
        this.simUI.player.setBulkEquipmentSpec(eventID, bulkEquipment);
      }
      this.close();
    } catch (e: any) {
      alert(e.toString());
    }
  }
}

class BulkSimResultRenderer {

  constructor(parent: ContentBlock, simUI: IndividualSimUI<Spec>, result: BulkSimResultWithSubstitutions, rank: number, baseResult: RaidSimResult) {
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
    dpsNumber.textContent = this.formatDps(result.raidMetrics?.dps?.avg!);
    dpsNumber.classList.add('topline-result-avg');
    dpsDiv.appendChild(dpsNumber);

    const dpsDelta = result.raidMetrics?.dps?.avg! - baseResult.raidMetrics?.dps?.avg!;
    const dpsDeltaSpan = document.createElement('span'); 
    dpsDeltaSpan.textContent = `${this.formatDpsDelta(dpsDelta)}`;
    dpsDeltaSpan.classList.add(dpsDelta >= 0 ? 'bulk-result-header-positive' : 'bulk-result-header-negative');
    dpsDiv.appendChild(dpsDeltaSpan);

    const itemsContainer = document.createElement('div');
    itemsContainer.style.flexDirection = 'row';
    itemsContainer.style.display = 'flex';
    parent.bodyElement.appendChild(itemsContainer);

    if (result.itemsAdded && result.itemsAdded.items && result.itemsAdded.items.length > 0) {
      for (const is of result.itemsAdded.items) {
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

  private itemSlotName(is: BulkEquipmentSpec_ItemSpecWithSlots): string {
    return JSON.parse(BulkEquipmentSpec_ItemSpecWithSlots.toJsonString(is))['slots'][0].replace('ItemSlot', '')
  }
}

export class BulkTab extends SimTab {
  protected simUI: IndividualSimUI<Spec>;

  readonly leftPanel: HTMLElement;
  readonly rightPanel: HTMLElement;

  readonly column1: HTMLElement = this.buildColumn(1, 'raid-settings-col');

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

    this.buildTabContent();
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
    this.simUI.sim.simRequestStartedEmitter.on(() => {
      resultsBlock.rootElem.hidden = true;
    });

    this.simUI.sim.simResultEmitter.on((_, simResult) => {
      resultsBlock.rootElem.hidden = simResult.result.bulkResults.length == 0;
      resultsBlock.bodyElement.innerHTML = '';

      let rank = 1;
      for (const r of simResult.result.bulkResults) {
        const resultBlock = new ContentBlock(resultsBlock.bodyElement, 'bulk-result', {header: {title: ''}});
        new BulkSimResultRenderer(resultBlock, this.simUI, r, rank, simResult.result);
        rank++;
      }
    });

    const settingsBlock = new ContentBlock(this.rightPanel, 'bulk-settings', {
      header: {title: 'Import'}
    });

    const importButton = document.createElement('button');
    importButton.classList.add('btn', 'btn-primary', 'w-100', 'bulk-settings-button');
    importButton.textContent = 'Import From Bags';
    importButton.addEventListener('click', () => new BulkGearJsonImporter(this.simUI.rootElem, this.simUI));
    settingsBlock.bodyElement.appendChild(importButton);

    const clearButton = document.createElement('button');
    clearButton.classList.add('btn', 'btn-primary', 'w-100', 'bulk-settings-button');
    clearButton.textContent = 'Clear All';
    clearButton.addEventListener('click', () => {
      const eventID = TypedEvent.nextEventID();
      this.simUI.player.setBulkEquipmentSpec(eventID, BulkEquipmentSpec.create());
      resultsBlock.rootElem.hidden = true;
      resultsBlock.bodyElement.innerHTML = '';
    });
    settingsBlock.bodyElement.appendChild(clearButton);

    this.simUI.player.bulkGearChangeEmitter.on(() => {
      const bulkEquipmentSpec = this.simUI.player.getBulkEquipmentSpec();
      if (bulkEquipmentSpec) {
        itemList.innerHTML = '';
        // Note: if this were to be fired from the player's fromProto() call,
        // we would have to load the database here with the missing items.
        // For now we assume that bulk items are not persisted. Bulk simming should be
        // a very conscious choice.
        if (bulkEquipmentSpec.items.length > 0) {
          notice.textContent = 'The following items will be simmed in all possible combinations together with your equipped gear.';
        }

        for (const itemWithSlots of bulkEquipmentSpec.items) {
          const item = this.simUI.sim.db.lookupItemSpec(itemWithSlots.item!);
          if (item) {
            const itemRenderer = new ItemRenderer(itemList, this.simUI, this.simUI.player);
            itemRenderer.update(item);
          }
        }
      }
    });
  }
}
