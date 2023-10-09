import { RaidPicker } from "./raid_picker";
import { RaidSimUI } from "./raid_sim_ui";
import { RaidStats } from "./raid_stats";
import { SavedDataManager } from "../core/components/saved_data_manager";
import { SimTab } from "../core/components/sim_tab";
import { BlessingsAssignments, SavedRaid } from "../core/proto/ui";
import { EventID, TypedEvent } from "../core/typed_event";
import { Raid as RaidProto } from "../core/proto/api";

export class RaidTab extends SimTab {
	protected simUI: RaidSimUI;

	readonly leftPanel: HTMLElement;
	readonly rightPanel: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: RaidSimUI) {
		super(parentElem, simUI, { identifier: 'raid-tab', title: 'Raid' });
		this.simUI = simUI;

		this.rootElem.classList.add('active', 'show')

		this.leftPanel = document.createElement('div');
		this.leftPanel.classList.add('raid-tab-left', 'tab-panel-left');

		this.rightPanel = document.createElement('div');
		this.rightPanel.classList.add('raid-tab-right', 'tab-panel-right');

		this.contentContainer.appendChild(this.leftPanel);
		this.contentContainer.appendChild(this.rightPanel);

		this.buildTabContent();
	}

	protected buildTabContent() {
		this.simUI.raidPicker = new RaidPicker(this.leftPanel, this.simUI);
		new RaidStats(this.leftPanel, this.simUI);

		const savedRaidManager = new SavedDataManager<RaidSimUI, SavedRaid>(this.rightPanel, this.simUI, {
			label: 'Raid',
			header: { title: 'Saved Raid Groups' },
			storageKey: this.simUI.getSavedRaidStorageKey(),
			getData: (raidSimUI: RaidSimUI) => SavedRaid.create({
				raid: this.simUI.sim.raid.toProto(),
				blessings: this.simUI.blessingsPicker!.getAssignments(),
				faction: this.simUI.sim.getFaction(),
				phase: this.simUI.sim.getPhase(),
			}),
			setData: (eventID: EventID, raidSimUI: RaidSimUI, newRaid: SavedRaid) => {
				TypedEvent.freezeAllAndDo(() => {
					this.simUI.sim.raid.fromProto(eventID, newRaid.raid || RaidProto.create());
					this.simUI.blessingsPicker!.setAssignments(eventID, newRaid.blessings || BlessingsAssignments.create());
					if (newRaid.faction) this.simUI.sim.setFaction(eventID, newRaid.faction);
					if (newRaid.phase) this.simUI.sim.setPhase(eventID, newRaid.phase);
				});
			},
			changeEmitters: [this.simUI.changeEmitter, this.simUI.sim.changeEmitter],
			equals: (a: SavedRaid, b: SavedRaid) => {
				return SavedRaid.equals(a, b);
			},
			toJson: (a: SavedRaid) => SavedRaid.toJson(a),
			fromJson: (obj: any) => SavedRaid.fromJson(obj),
		}
		);
		this.simUI.sim.waitForInit().then(() => {
			savedRaidManager.loadUserData();
		});
	}
}
