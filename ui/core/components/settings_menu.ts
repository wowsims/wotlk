import { Tooltip } from 'bootstrap';

import { BooleanPicker } from '../components/boolean_picker.js';
import { EnumPicker } from '../components/enum_picker.js';
import { NumberPicker } from '../components/number_picker.js';
import { wowheadSupportedLanguages } from '../constants/lang.js';
import { Sim } from '../sim.js';
import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';
import { BaseModal } from './base_modal.js';

export class SettingsMenu extends BaseModal {
	private readonly simUI: SimUI;

	constructor(parent: HTMLElement, simUI: SimUI) {
		super(parent, 'settings-menu', { title: "选项", footer: true });
		this.simUI = simUI;

		this.body.innerHTML = `
			<div class="picker-group">
				<div class="fixed-rng-seed-container">
					<div class="fixed-rng-seed"></div>
					<div class="form-text">
						<span>上次使用RNG种子:</span>&nbsp;<span class="last-used-rng-seed">0</span>
					</div>
				</div>
				<div class="language-picker within-raid-sim-hide"></div>
			</div>
			<div class="show-threat-metrics-picker w-50 pe-2"></div>
			<div class="show-experimental-picker w-50 pe-2"></div>
		`;
		this.footer!.innerHTML = `
			<button
				class="restore-defaults-button btn btn-primary"
			>重置所有选项</button>
		`

		const restoreDefaultsButton = this.rootElem.getElementsByClassName('restore-defaults-button')[0] as HTMLElement;
		Tooltip.getOrCreateInstance(restoreDefaultsButton, {
			title: "恢复所有默认设置（装备、消耗品、增益、天赋、EP权重等）。保存的设置会被保留。"
		});
		restoreDefaultsButton.addEventListener('click', event => {
			this.simUI.applyDefaults(TypedEvent.nextEventID());
		});

		const fixedRngSeed = this.rootElem.getElementsByClassName('fixed-rng-seed')[0] as HTMLElement;
		new NumberPicker(fixedRngSeed, this.simUI.sim, {
			label: '固定RNG种子',
			labelTooltip: '用于模拟期间的随机数生成器的种子值，或设置为0以在每次运行时使用不同的随机性。使用此值可以共享精确的模拟结果或用于调试。',
			extraCssClasses: ['mb-0'],
			changedEvent: (sim: Sim) => sim.fixedRngSeedChangeEmitter,
			getValue: (sim: Sim) => sim.getFixedRngSeed(),
			setValue: (eventID: EventID, sim: Sim, newValue: number) => {
				sim.setFixedRngSeed(eventID, newValue);
			},
		});

		const lastUsedRngSeed = this.rootElem.getElementsByClassName('last-used-rng-seed')[0] as HTMLElement;
		lastUsedRngSeed.textContent = String(this.simUI.sim.getLastUsedRngSeed());
		this.simUI.sim.lastUsedRngSeedChangeEmitter.on(() => lastUsedRngSeed.textContent = String(this.simUI.sim.getLastUsedRngSeed()));

		// const language = this.rootElem.getElementsByClassName('language-picker')[0] as HTMLElement;
		// const langs = Object.keys(wowheadSupportedLanguages);
		// const defaultLang = langs.indexOf('cn');
		// const languagePicker = new EnumPicker(language, this.simUI.sim, {
		// 	label: 'Language',
		// 	labelTooltip: 'Controls the language for Wowhead tooltips.',
		// 	values: langs.map((lang, i) => {
		// 		return {
		// 			name: wowheadSupportedLanguages[lang],
		// 			value: i,
		// 		};
		// 	}),
		// 	changedEvent: (sim: Sim) => sim.languageChangeEmitter,
		// 	getValue: (sim: Sim) => {
		// 		const idx = langs.indexOf(sim.getLanguage());
		// 		return idx == -1 ? defaultLang : idx;
		// 	},
		// 	setValue: (eventID: EventID, sim: Sim, newValue: number) => {
		// 		sim.setLanguage(eventID, langs[newValue] || 'en');
		// 	},
		// });
		// // Refresh page after language change, to apply the changes.
		// languagePicker.changeEmitter.on(() => setTimeout(() => location.reload(), 100));

		const showThreatMetrics = this.rootElem.getElementsByClassName('show-threat-metrics-picker')[0] as HTMLElement;
		new BooleanPicker(showThreatMetrics, this.simUI.sim, {
			label: '显示威胁/坦克选项指标',
			labelTooltip: '显示与坦克相关的所有选项和指标，例如TPS/DTPS。',
			inline: true,
			changedEvent: (sim: Sim) => sim.showThreatMetricsChangeEmitter,
			getValue: (sim: Sim) => sim.getShowThreatMetrics(),
			setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
				sim.setShowThreatMetrics(eventID, newValue);
			},
		});

		const showExperimental = this.rootElem.getElementsByClassName('show-experimental-picker')[0] as HTMLElement;
		new BooleanPicker(showExperimental, this.simUI.sim, {
			label: '显示实验功能',
			labelTooltip: '如果有任何活跃的实验功能，显示实验功能。',
			inline: true,
			changedEvent: (sim: Sim) => sim.showExperimentalChangeEmitter,
			getValue: (sim: Sim) => sim.getShowExperimental(),
			setValue: (eventID: EventID, sim: Sim, newValue: boolean) => {
				sim.setShowExperimental(eventID, newValue);
			},
		});

	}
}
