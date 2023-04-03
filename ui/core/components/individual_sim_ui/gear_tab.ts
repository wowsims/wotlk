import { ContentBlock } from "../content_block";
import { Database } from '../../proto_utils/database';
import { Importer } from "../importers";
import { Component } from '../component';

import { IndividualSimUI } from "../../individual_sim_ui";
import { Player } from "../../player";
import { EventID, TypedEvent } from "../../typed_event";

import { EquipmentSpec, BulkEquipmentSpec, Spec, UnitStats, BulkEquipmentSpec_ItemSpecWithSlots, ItemSpec, ItemSlot } from "../../proto/common";
import { SavedGearSet } from "../../proto/ui";
import { Stats } from "../../proto_utils/stats";

import { GearPicker, ItemRenderer } from "../gear_picker";
import { SavedDataManager } from "../saved_data_manager";
import { SimTab } from "../sim_tab";

import { getEligibleItemSlots } from '../../proto_utils/utils.js';
import { EquippedItem } from "ui/core/proto_utils/equipped_item";


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

class BulkSimResultRenderer extends Component {
  private readonly simUI: IndividualSimUI<Spec>;

	constructor(parent: HTMLElement, simUI: IndividualSimUI<Spec>, spec: BulkEquipmentSpec) {
		super(parent, 'bulk-result');
    this.rootElem.style.flexDirection = 'row';
    this.rootElem.style.display = 'flex';
		this.simUI = simUI;

    for (const is of spec.items) {
      const item = this.simUI.sim.db.lookupItemSpec(is.item!)
      const renderer = new ItemRenderer(this.rootElem, this.simUI, this.simUI.player);
      renderer.update(item!);
      const p = document.createElement('a');
      p.classList.add('bulk-result-item-slot');
      p.textContent = JSON.parse(BulkEquipmentSpec_ItemSpecWithSlots.toJsonString(is))['slots'][0];
      renderer.nameElem.appendChild(p); 
    }

    if (spec.items.length == 0) {
      const p = document.createElement('p');
      p.textContent = 'No changes - this is your currently equipped gear!';
      this.rootElem.appendChild(p);
    }
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
      header: {title: 'Items', extraCssClasses: ['bulk-header']}
    });
    itemsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-left', 'tab-panel-col');
    itemsBlock.bodyElement.style.flexDirection = 'row';
    itemsBlock.bodyElement.style.display = 'flex';
    let resultsBlock = new ContentBlock(this.column1, 'bulk-results', {
      header: {title: 'Results', extraCssClasses: ['bulk-header']}
    });
    resultsBlock.bodyElement.classList.add('gear-picker-root', 'gear-picker-left', 'tab-panel-col');
    this.simUI.sim.simResultEmitter.on((idx, simResult) => {
      let i = 1;
      resultsBlock.bodyElement.innerHTML = '';
      for (const r of simResult.result.bulkResults) {
        const resultBlock = new ContentBlock(resultsBlock.bodyElement, 'bulk-result', {
          header: {
            title: 'Rank ' + i + ': ' + (Math.round(r.raidMetrics?.dps?.avg! * 100) / 100).toFixed(2) + 'DPS',
            extraCssClasses: ['bulk-item-header'],
          }
        });
        //const text = 'Rank ' + i + ' with ' + r.raidMetrics?.dps?.avg + 'DPS, changes: ' + BulkEquipmentSpec.toJsonString(r.itemsAdded!);
        new BulkSimResultRenderer(resultBlock.bodyElement, this.simUI, r.itemsAdded!)
        i++;
      }
    });

    const settingsBlock = new ContentBlock(this.rightPanel, 'bulk-settings', {
      header: {title: 'Settings'}
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
      resultsBlock.bodyElement.innerHTML = '';
    });
    settingsBlock.bodyElement.appendChild(clearButton);

    this.simUI.player.bulkGearChangeEmitter.on(() => {
      const bulkEquipmentSpec = this.simUI.player.getBulkEquipmentSpec();
      if (bulkEquipmentSpec) {
        itemsBlock.bodyElement.innerHTML = '';
        // Note: if this were to be fired from the player's fromProto() call,
        // we would have to load the database here with the missing items.
        // For now we assume that bulk items are not persisted. Bulk simming should be
        // a very conscious choice.
        for (const itemWithSlots of bulkEquipmentSpec.items) {
          const item = this.simUI.sim.db.lookupItemSpec(itemWithSlots.item!);
          if (item) {
            const itemRenderer = new ItemRenderer(itemsBlock.bodyElement, this.simUI, this.simUI.player);
            itemRenderer.update(item);
          }
        }
      }
    });
  }
}

export class GearTab extends SimTab {
  protected simUI: IndividualSimUI<Spec>;

  readonly leftPanel: HTMLElement;
  readonly rightPanel: HTMLElement;

  constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
    super(parentElem, simUI, {identifier: 'gear-tab', title: 'Gear'});
    this.simUI = simUI;

    this.leftPanel = document.createElement('div');
    this.leftPanel.classList.add('gear-tab-left', 'tab-panel-left');

    this.rightPanel = document.createElement('div');
    this.rightPanel.classList.add('gear-tab-right', 'tab-panel-right');

    this.contentContainer.appendChild(this.leftPanel);
    this.contentContainer.appendChild(this.rightPanel);

    this.buildTabContent();
  }

  protected buildTabContent() {
    this.buildGearPickers();

    this.buildSavedGearsetPicker();
  }

  private buildGearPickers() {
    new GearPicker(this.leftPanel, this.simUI, this.simUI.player);
  }

  private buildSavedGearsetPicker() {
    const savedGearManager = new SavedDataManager<Player<any>, SavedGearSet>(this.rightPanel, this.simUI, this.simUI.player, {
      header: {title: "Gear Sets"},
      label: 'Gear Set',
      storageKey: this.simUI.getSavedGearStorageKey(),
      getData: (player: Player<any>) => {
        return SavedGearSet.create({
          gear: player.getGear().asSpec(),
          bonusStatsStats: player.getBonusStats().toProto(),
        });
      },
      setData: (eventID: EventID, player: Player<any>, newSavedGear: SavedGearSet) => {
        TypedEvent.freezeAllAndDo(() => {
          player.setGear(eventID, this.simUI.sim.db.lookupEquipmentSpec(newSavedGear.gear || EquipmentSpec.create()));
          if (newSavedGear.bonusStats && newSavedGear.bonusStats.some(s => s != 0)) {
            player.setBonusStats(eventID, new Stats(newSavedGear.bonusStats));
          } else {
            player.setBonusStats(eventID, Stats.fromProto(newSavedGear.bonusStatsStats || UnitStats.create()));
          }
        });
      },
      changeEmitters: [this.simUI.player.changeEmitter],
      equals: (a: SavedGearSet, b: SavedGearSet) => SavedGearSet.equals(a, b),
      toJson: (a: SavedGearSet) => SavedGearSet.toJson(a),
      fromJson: (obj: any) => SavedGearSet.fromJson(obj),
    });

		this.simUI.sim.waitForInit().then(() => {
			savedGearManager.loadUserData();
			this.simUI.individualConfig.presets.gear.forEach(presetGear => {
				savedGearManager.addSavedData({
					name: presetGear.name,
					tooltip: presetGear.tooltip,
					isPreset: true,
					data: SavedGearSet.create({
						// Convert to gear and back so order is always the same.
						gear: this.simUI.sim.db.lookupEquipmentSpec(presetGear.gear).asSpec(),
						bonusStatsStats: new Stats().toProto(),
					}),
					enableWhen: presetGear.enableWhen,
				});
			});
		});
  }
}
