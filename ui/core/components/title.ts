import { Component } from '/wotlk/core/components/component.js';
import { getLaunchedSims, raidSimLaunched } from '/wotlk/core/launched_sims.js';
import { Spec } from '/wotlk/core/proto/common.js';
import {
	classColors,
	getSpecSiteUrl,
	naturalSpecOrder,
	raidSimSiteUrl,
	specNames,
	specToClass,
	titleIcons,
	raidSimIcon,
} from '/wotlk/core/proto_utils/utils.js';

declare var tippy: any;

export interface SimLinkOption {
	iconUrl: string,
	href: string,
	text: string,
	color: string,
};

// Dropdown menu for selecting a player.
export class Title extends Component {
	private readonly buttonElem: HTMLElement;

	constructor(parent: HTMLElement, currentSpec: Spec | null) {
		super(parent, 'sim-title-root');
		this.rootElem.classList.add('dropdown-root');

		this.rootElem.innerHTML = `
			<div class="dropdown-button sim-title-button"></div>
			<div class="dropdown-panel sim-title-dropdown within-raid-sim-hide"></div>
    `;

		this.buttonElem = this.rootElem.getElementsByClassName('sim-title-button')[0] as HTMLElement;
		const dropdownPanel = this.rootElem.getElementsByClassName('dropdown-panel')[0] as HTMLElement;

		this.buttonElem.addEventListener('click', event => {
			event.preventDefault();
		});

		const orderedLaunchedSpecs: Array<Spec | null> = (naturalSpecOrder
			.filter(spec => getLaunchedSims().includes(spec)) as Array<Spec | null>)
			.concat(raidSimLaunched ? [null] : []); // Null represents the raid sim.

		dropdownPanel.style.gridTemplateRows = `repeat(${Math.ceil(orderedLaunchedSpecs.length / 2)}, 1fr)`;

		const currentOption = this.makeOptionData(currentSpec, true);
		const otherOptions = orderedLaunchedSpecs.map(spec => this.makeOptionData(spec, false));

		this.buttonElem.appendChild(Title.makeOptionElem(currentOption));

		const isWithinRaidSim = this.rootElem.closest('.within-raid-sim') != null;
		if (!isWithinRaidSim) {
			otherOptions.forEach((option, i) => dropdownPanel.appendChild(this.makeOption(option)));
		}
	}

	private makeOptionData(spec: Spec | null, isButton: boolean) {
		if (spec == null) {
			return {
				iconUrl: raidSimIcon,
				href: raidSimSiteUrl,
				text: 'RAID',
				color: isButton ? '' : 'black',
			};
		} else {
			return {
				iconUrl: titleIcons[spec],
				href: getSpecSiteUrl(spec),
				text: specNames[spec].toUpperCase(),
				color: isButton ? '' : classColors[specToClass[spec]],
			};
		}
	}

	private makeOption(data: SimLinkOption): HTMLElement {
		const option = Title.makeOptionElem(data);

		option.addEventListener('click', event => {
			event.preventDefault();
			window.location.href = data.href;
		});

		return option;
	}

	static makeOptionElem(data: SimLinkOption): HTMLElement {
		const optionContainer = document.createElement('a');
		optionContainer.href = data.href;
		optionContainer.classList.add('sim-title-dropdown-option-container', 'dropdown-option-container');

		const option = document.createElement('div');
		option.classList.add('sim-title-option', 'dropdown-option');
		optionContainer.appendChild(option);

		if (data.color) {
			option.style.backgroundColor = data.color;
		}

		const icon = document.createElement('img');
		icon.src = data.iconUrl;
		icon.classList.add('sim-title-icon');
		option.appendChild(icon);

		const labelDiv = document.createElement('div');
		labelDiv.classList.add('sim-title-label-container');
		option.appendChild(labelDiv);

		if (!data.color) { // Hacky check for 'isButton'
			const simLabel = document.createElement('span');
			simLabel.textContent = 'WotLK Simulator';
			simLabel.classList.add('sim-title-sim-label', 'sim-title-label');
			labelDiv.appendChild(simLabel);
		}

		const label = document.createElement('span');
		label.textContent = data.text;
		label.classList.add('sim-title-spec-label', 'sim-title-label');
		labelDiv.appendChild(label);

		return optionContainer;
	}
}
