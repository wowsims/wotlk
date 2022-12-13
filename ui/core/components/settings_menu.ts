import { Class } from '../proto/common.js';
import { EquipmentSpec } from '../proto/common.js';
import { ItemSpec } from '../proto/common.js';
import { Race } from '../proto/common.js';
import { Spec } from '../proto/common.js';
import { Stat } from '../proto/common.js';
import { SimUI } from '../sim_ui.js';
import { IndividualSimUI } from '../individual_sim_ui.js';
import { Sim } from '../sim.js';
import { Player } from '../player.js';
import { classNames, nameToClass, nameToRace, statNames } from '../proto_utils/names.js';
import { talentSpellIdsToTalentString } from '../talents/factory.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { wowheadSupportedLanguages } from '../constants/lang.js';
import { BooleanPicker } from '../components/boolean_picker.js';
import { EnumPicker } from '../components/enum_picker.js';
import { NumberPicker } from '../components/number_picker.js';
import { getEnumValues } from '../utils.js';

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
			<div class="menu-title">
				<span>OPTIONS</span>
			</div>
			<div class="menu-content">
				<div class="settings-menu-content-left">
					<button class="restore-defaults-button btn btn-primary">Restore Defaults</button>
					<div class="fixed-rng-seed"></div>
					<div class="last-used-rng-seed-container form-text">
						<small>Last used RNG seed:</small><small class="last-used-rng-seed">0</small>
					</div>
					<div class="show-threat-metrics-picker"></div>
					<div class="show-experimental-picker"></div>
				</div>
				<div class="settings-menu-content-right">
					<div class="language-picker within-raid-sim-hide"></div>
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
			inline: true,
			changedEvent: (sim: Sim) => sim.showThreatMetricsChangeEmitter,
			getValue: (sim: Sim) => sim.getShowThreatMetrics(),
			setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
				sim.setShowThreatMetrics(eventID, newValue);
			},
		});

		const showExperimental = this.rootElem.getElementsByClassName('show-experimental-picker')[0] as HTMLElement;
		new BooleanPicker(showExperimental, this.simUI.sim, {
			label: 'Show Experimental',
			labelTooltip: 'Shows experimental options, if there are any active experiments.',
			inline: true,
			changedEvent: (sim: Sim) => sim.showExperimentalChangeEmitter,
			getValue: (sim: Sim) => sim.getShowExperimental(),
			setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
				sim.setShowExperimental(eventID, newValue);
			},
		});

		const language = this.rootElem.getElementsByClassName('language-picker')[0] as HTMLElement;
		const langs = Object.keys(wowheadSupportedLanguages);
		const defaultLang = langs.indexOf('en');
		const languagePicker = new EnumPicker(language, this.simUI.sim, {
			label: 'Language',
			labelTooltip: 'Controls the language for Wowhead tooltips.',
			values: langs.map((lang, i) => {
				return {
					name: wowheadSupportedLanguages[lang],
					value: i,
				};
			}),
			changedEvent: (sim: Sim) => sim.languageChangeEmitter,
			getValue: (sim: Sim) => {
				const idx = langs.indexOf(sim.getLanguage());
				return idx == -1 ? defaultLang : idx;
			},
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setLanguage(eventID, langs[newValue] || 'en');
			},
		});
		// Refresh page after language change, to apply the changes.
		languagePicker.changeEmitter.on(() => setTimeout(() => location.reload(), 100));
	}
}
