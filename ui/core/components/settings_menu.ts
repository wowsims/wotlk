import { Class } from '/tbc/core/proto/common.js';
import { EquipmentSpec } from '/tbc/core/proto/common.js';
import { ItemSpec } from '/tbc/core/proto/common.js';
import { Race } from '/tbc/core/proto/common.js';
import { Spec } from '/tbc/core/proto/common.js';
import { Stat } from '/tbc/core/proto/common.js';
import { SimUI } from '/tbc/core/sim_ui.js';
import { IndividualSimUI } from '/tbc/core/individual_sim_ui.js';
import { Sim } from '/tbc/core/sim.js';
import { Player } from '/tbc/core/player.js';
import { classNames, nameToClass, nameToRace, statNames } from '/tbc/core/proto_utils/names.js';
import { talentSpellIdsToTalentString } from '/tbc/core/talents/factory.js';
import { EventID, TypedEvent } from '/tbc/core/typed_event.js';
import { BooleanPicker } from '/tbc/core/components/boolean_picker.js';
import { NumberPicker } from '/tbc/core/components/number_picker.js';
import { getEnumValues } from '/tbc/core/utils.js';

import { Popup } from './popup.js';

declare var $: any;
declare var tippy: any;

export class SettingsMenu extends Popup {
	private readonly simUI: SimUI;

	constructor(parent: HTMLElement, simUI: SimUI) {
		super(parent);
		this.rootElem.classList.add('settings-menu');
		this.simUI = simUI;

		this.rootElem.innerHTML = `
			<div class="settings-menu-title">
				<span>OPTIONS</span>
			</div>
			<div class="settings-menu-content">
				<div class="settings-menu-content-left">
					<button class="restore-defaults-button sim-button">RESTORE DEFAULTS</button>
					<div class="settings-menu-section">
						<div class="fixed-rng-seed">
						</div>
						<div class="last-used-rng-seed-container">
							<span>Last used RNG seed:</span><span class="last-used-rng-seed">0</span>
						</div>
						<div class="show-threat-metrics-picker">
						</div>
						<div class="show-experimental-picker">
						</div>
					</div>
				</div>
				<div class="settings-menu-content-right">
					<div class="settings-menu-section settings-menu-ep-weights within-raid-sim-hide">
					</div>
				</div>
			</div>
		`;
		this.addCloseButton();

		const restoreDefaultsButton = this.rootElem.getElementsByClassName('restore-defaults-button')[0] as HTMLElement;
		restoreDefaultsButton.addEventListener('click', event => {
			this.simUI.applyDefaults(TypedEvent.nextEventID());
		});
		tippy(restoreDefaultsButton, {
			'content': 'Restores all default settings (gear, consumes, buffs, talents, EP weights, etc).',
			'allowHTML': true,
		});

		const fixedRngSeed = this.rootElem.getElementsByClassName('fixed-rng-seed')[0] as HTMLElement;
		new NumberPicker(fixedRngSeed, this.simUI.sim, {
			label: 'Fixed RNG Seed',
			labelTooltip: 'Seed value for the random number generator used during sims, or 0 to use different randomness each run. Use this to share exact sim results or for debugging.',
			changedEvent: (sim: Sim) => sim.fixedRngSeedChangeEmitter,
			getValue: (sim: Sim) => sim.getFixedRngSeed(),
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setFixedRngSeed(eventID, newValue);
			},
		});

		const lastUsedRngSeed = this.rootElem.getElementsByClassName('last-used-rng-seed')[0] as HTMLElement;
		lastUsedRngSeed.textContent = String(this.simUI.sim.getLastUsedRngSeed());
		this.simUI.sim.lastUsedRngSeedChangeEmitter.on(() => lastUsedRngSeed.textContent = String(this.simUI.sim.getLastUsedRngSeed()));

		const showThreatMetrics = this.rootElem.getElementsByClassName('show-threat-metrics-picker')[0] as HTMLElement;
		new BooleanPicker(showThreatMetrics, this.simUI.sim, {
			label: 'Show Threat/Tank Options',
			labelTooltip: 'Shows all options and metrics relevant to tanks, like TPS/DTPS.',
			changedEvent: (sim: Sim) => sim.showThreatMetricsChangeEmitter,
			getValue: (sim: Sim) => sim.getShowThreatMetrics(),
			setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
				sim.setShowThreatMetrics(eventID, newValue);
			},
		});

		// Comment this out when there are no experiments to show.
		const showExperimental = this.rootElem.getElementsByClassName('show-experimental-picker')[0] as HTMLElement;
		new BooleanPicker(showExperimental, this.simUI.sim, {
			label: 'Show Experimental',
			changedEvent: (sim: Sim) => sim.showExperimentalChangeEmitter,
			getValue: (sim: Sim) => sim.getShowExperimental(),
			setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
				sim.setShowExperimental(eventID, newValue);
			},
		});
	}
}
