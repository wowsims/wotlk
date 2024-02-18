import { IndividualSimUI } from "../../individual_sim_ui";
import { Player } from "../../player";
import { EventID, TypedEvent } from "../../typed_event";

import { Class, Glyphs, Spec } from "../../proto/common";
import { SavedTalents } from "../../proto/ui";

import { classGlyphsConfig, classTalentsConfig } from "../../talents/factory";
import { GlyphsPicker } from "../../talents/glyphs_picker";
import { HunterPetTalentsPicker, makePetTypeInputConfig } from "../../talents/hunter_pet";
import { TalentsPicker } from "../../talents/talents_picker";

import { IconEnumPicker } from "../icon_enum_picker";
import { SavedDataManager } from "../saved_data_manager";
import { SimTab } from "../sim_tab";

import * as Mechanics from '../../constants/mechanics';

export class TalentsTab extends SimTab {
	protected simUI: IndividualSimUI<Spec>;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: IndividualSimUI<Spec>) {
		super(parentElem, simUI, { identifier: 'talents-tab', title: 'Talents' });
		this.simUI = simUI;

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('talents-tab-left', 'tab-panel-left');
		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('talents-tab-right', 'tab-panel-right', 'within-raid-sim-hide');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
    if (this.simUI.player.getClass() == Class.ClassHunter) {
      this.buildHunterPickers();
    } else {
      this.buildTalentsPicker(this.leftPanel);
      this.buildGlyphsPicker(this.leftPanel);
    }

    this.buildSavedTalentsPicker();
	}

  private buildTalentsPicker(parentElem: HTMLElement) {
    new TalentsPicker(parentElem, this.simUI.player, {
      klass: this.simUI.player.getClass(),
      trees: classTalentsConfig[this.simUI.player.getClass()],
      changedEvent: (player: Player<any>) => player.talentsChangeEmitter,
      getValue: (player: Player<any>) => player.getTalentsString(),
      setValue: (eventID: EventID, player: Player<any>, newValue: string) => {
        player.setTalentsString(eventID, newValue);
      },
      pointsPerRow: 5,
      maxPoints: Mechanics.MAX_TALENT_POINTS,
    });
  }

  private buildGlyphsPicker(parentElem: HTMLElement) {
    new GlyphsPicker(parentElem, this.simUI.player, classGlyphsConfig[this.simUI.player.getClass()]);
  }

  private buildHunterPickers() {
    this.leftPanel.innerHTML = `
      <div class="hunter-talents-pickers-container tab-content">
        <ul class="nav nav-tabs" role="tablist">
          <li class="nav-item" role="presentation">
            <a
              class="nav-link active"
              type="button"
              role="tab"
              data-bs-toggle="tab"
              data-bs-target="#player-talents-tab"
              aria-controls="#player-talents-tab"
              aria-selected="true"
            >
              Player
            </a>
          </li>
          <li class="nav-item" role="presentation">
            <a
              class="nav-link"
              type="button"
              role="tab"
              data-bs-toggle="tab"
              data-bs-target="#pet-talents-tab"
              aria-controls="#pet-talents-tab"
              aria-selected="false"
            >
              Pet</a
            >
          </li>
        </ul>
        <div id="player-talents-tab" class="tab-pane fade active show" role="tabpanel">
        </div>
        <div id="pet-talents-tab" class="tab-pane fade" role="tabpanel">
        </div>
      </div>
    `

    const playerTab = this.leftPanel.querySelector('#player-talents-tab') as HTMLElement;
    const petTab = this.leftPanel.querySelector('#pet-talents-tab') as HTMLElement;

    this.buildTalentsPicker(playerTab);
    this.buildGlyphsPicker(playerTab);
    this.buildHunterPetPicker(petTab);
  }

  private buildHunterPetPicker(parentElem: HTMLElement) {
    new HunterPetTalentsPicker(parentElem, this.simUI, this.simUI.player as Player<Spec.SpecHunter>);
    new IconEnumPicker(parentElem, this.simUI.player as Player<Spec.SpecHunter>, makePetTypeInputConfig());
  }

  private buildSavedTalentsPicker() {
    const savedTalentsManager = new SavedDataManager<Player<any>, SavedTalents>(this.rightPanel, this.simUI.player, {
			label: 'Talents',
			header: { title: 'Saved Talents' },
			storageKey: this.simUI.getSavedTalentsStorageKey(),
			getData: (player: Player<any>) => SavedTalents.create({
				talentsString: player.getTalentsString(),
				glyphs: player.getGlyphs(),
			}),
			setData: (eventID: EventID, player: Player<any>, newTalents: SavedTalents) => {
				TypedEvent.freezeAllAndDo(() => {
					player.setTalentsString(eventID, newTalents.talentsString);
					player.setGlyphs(eventID, newTalents.glyphs || Glyphs.create());
				});
			},
			changeEmitters: [this.simUI.player.talentsChangeEmitter, this.simUI.player.glyphsChangeEmitter],
			equals: (a: SavedTalents, b: SavedTalents) => SavedTalents.equals(a, b),
			toJson: (a: SavedTalents) => SavedTalents.toJson(a),
			fromJson: (obj: any) => SavedTalents.fromJson(obj),
		});

    this.simUI.sim.waitForInit().then(() => {
			savedTalentsManager.loadUserData();
			this.simUI.individualConfig.presets.talents.forEach(config => {
				config.isPreset = true;
				savedTalentsManager.addSavedData({
					name: config.name,
					isPreset: true,
					data: config.data,
				});
			});
		});
  } 
}
