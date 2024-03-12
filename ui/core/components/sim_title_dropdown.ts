import {
	getLaunchedSimsForClass,
	LaunchStatus,
	raidSimStatus,
	simLaunchStatuses,
} from '../launched_sims.js';
import { Class, Spec } from '../proto/common.js';
import {
	classNames,
	getSpecSiteUrl,
	naturalClassOrder,
	raidSimIcon,
	raidSimLabel,
	raidSimSiteUrl,
	specNames,
	specToClass,
	textCssClassForClass,
	textCssClassForSpec,
	titleIcons,
} from '../proto_utils/utils.js';
import { Component } from './component.js';

interface ClassOptions {
	type: 'Class';
	index: Class;
}

interface SpecOptions {
	type: 'Spec';
	index: Spec;
}

interface RaidOptions {
	type: 'Raid';
}

type SimTitleDropdownConfig = {
	noDropdown?: boolean;
};

// Dropdown menu for selecting a player.
export class SimTitleDropdown extends Component {
	private readonly dropdownMenu: HTMLElement | undefined;

	private readonly specLabels: Record<Spec, string> = {
		[Spec.SpecBalanceDruid]: 'Balance',
		[Spec.SpecFeralDruid]: 'Feral DPS',
		[Spec.SpecFeralTankDruid]: 'Feral Tank',
		[Spec.SpecRestorationDruid]: 'Restoration',
		[Spec.SpecElementalShaman]: 'Elemental',
		[Spec.SpecEnhancementShaman]: 'Enhancement',
		[Spec.SpecRestorationShaman]: 'Restoration',
		[Spec.SpecHunter]: 'Hunter',
		[Spec.SpecMage]: 'Mage',
		[Spec.SpecRogue]: 'Rogue',
		[Spec.SpecHolyPaladin]: 'Holy',
		[Spec.SpecProtectionPaladin]: 'Protection',
		[Spec.SpecRetributionPaladin]: 'Retribution',
		[Spec.SpecHealingPriest]: 'Healing',
		[Spec.SpecShadowPriest]: 'Shadow',
		[Spec.SpecSmitePriest]: 'Smite',
		[Spec.SpecWarlock]: 'Warlock',
		[Spec.SpecWarrior]: 'DPS',
		[Spec.SpecProtectionWarrior]: 'Protection',
		[Spec.SpecDeathknight]: 'DPS',
		[Spec.SpecTankDeathknight]: 'Tank',
	};

	constructor(
		parent: HTMLElement,
		currentSpecIndex: Spec | null,
		config: SimTitleDropdownConfig = {},
	) {
		super(parent, 'sim-title-dropdown-root');

		const rootLinkArgs: SpecOptions | RaidOptions =
			currentSpecIndex === null
				? { type: 'Raid' }
				: { type: 'Spec', index: currentSpecIndex };
		const rootLink = this.buildRootSimLink(rootLinkArgs);

		if (config.noDropdown) {
			this.rootElem.innerHTML = rootLink.outerHTML;
			return;
		}

		this.rootElem.innerHTML = `
			<div class="dropdown sim-link-dropdown">
				${rootLink.outerHTML}
				<ul class="dropdown-menu"></ul>
			</div>
		`;

		this.dropdownMenu = this.rootElem.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
		this.buildDropdown();

		// Prevent Bootstrap from closing the menu instead of opening class menus
		this.dropdownMenu.addEventListener('click', event => {
			const target = event.target as HTMLElement;
			const link = target.closest('a:not([href="javascript:void(0)"]');

			if (!link) {
				event.stopPropagation();
				event.preventDefault();
			}
		});
	}

	private buildDropdown() {
		if (raidSimStatus >= LaunchStatus.Alpha) {
			// Add the raid sim to the top of the dropdown
			const raidListItem = document.createElement('li');
			raidListItem.appendChild(this.buildRaidLink());
			this.dropdownMenu?.appendChild(raidListItem);
		}

		naturalClassOrder.forEach(classIndex => {
			const listItem = document.createElement('li');
			const sims = getLaunchedSimsForClass(classIndex);

			if (sims.length == 1) {
				// The class only has one listed sim so make a direct link to the sim
				listItem.appendChild(this.buildClassLink(classIndex));
				this.dropdownMenu?.appendChild(listItem);
			} else if (sims.length > 1) {
				// Add the class to the dropdown with an additional spec dropdown
				listItem.appendChild(this.buildClassDropdown(classIndex));
				this.dropdownMenu?.appendChild(listItem);
			}
		});
	}

	private buildClassDropdown(classIndex: Class) {
		const sims = getLaunchedSimsForClass(classIndex);
		const dropdownFragment = document.createElement('fragment');
		const dropdownMenu = document.createElement('ul');
		dropdownMenu.classList.add('dropdown-menu');

		// Generate the class link to act as a dropdown toggle for the spec dropdown
		const classLink = this.buildClassLink(classIndex);

		// Generate links for a class's specs
		sims.forEach(specIndex => {
			const listItem = document.createElement('li');
			const link = this.buildSpecLink(specIndex);

			listItem.appendChild(link);
			dropdownMenu.appendChild(listItem);
		});

		dropdownFragment.innerHTML = `
			<div class="dropend sim-link-dropdown">
				${classLink.outerHTML}
				${dropdownMenu.outerHTML}
			</div>
    	`;

		return dropdownFragment.children[0] as HTMLElement;
	}

	private buildRootSimLink(data: SpecOptions | RaidOptions): HTMLElement {
		const iconPath = this.getSimIconPath(data);
		const textKlass = this.getContextualKlass(data);
		let label;

		if (data.type == 'Raid') label = raidSimLabel;
		else {
			const classIndex = specToClass[data.index];
			if (getLaunchedSimsForClass(classIndex).length > 1)
				// If the class has multiple sims, use the spec name
				label = specNames[data.index];
			// If the class has only 1 sim, use the class name
			else label = classNames[classIndex];
		}

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="javascript:void(0)" class="sim-link ${textKlass}" role="button" data-bs-toggle="dropdown" data-bs-trigger="click" aria-expanded="false">
				<div class="sim-link-content">
					<img src="${iconPath}" class="sim-link-icon">
					<div class="d-flex flex-column">
						<span class="sim-link-label text-white">WoWSims - WOTLK</span>
						<span class="sim-link-title">${label}</span>
						${this.launchStatusLabel(data)}
					</div>
				</div>
			</a>
		`;

		return fragment.children[0] as HTMLElement;
	}

	private buildRaidLink(): HTMLElement {
		const href = raidSimSiteUrl;
		const textKlass = this.getContextualKlass({ type: 'Raid' });
		const iconPath = this.getSimIconPath({ type: 'Raid' });
		const label = raidSimLabel;

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="${raidSimSiteUrl}" class="sim-link ${textKlass}">
				<div class="sim-link-content">
					<img src="${iconPath}" class="sim-link-icon">
					<div class="d-flex flex-column">
						<span class="sim-link-title">${label}</span>
						${this.launchStatusLabel({ type: 'Raid' })}
					</div>
				</div>
			</a>
    	`;

		return fragment.children[0] as HTMLElement;
	}

	private buildClassLink(classIndex: Class): HTMLElement {
		const specIndexes = getLaunchedSimsForClass(classIndex);
		const href = specIndexes.length > 1 ? 'javascript:void(0)' : getSpecSiteUrl(specIndexes[0]);
		const textKlass = this.getContextualKlass({ type: 'Class', index: classIndex });
		const iconPath = this.getSimIconPath({ type: 'Class', index: classIndex });
		const label = classNames[classIndex];

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="${href}" class="sim-link ${textKlass}" ${
				specIndexes.length > 1
					? 'role="button" data-bs-toggle="dropdown" aria-expanded="false"'
					: ''
			}>
				<div class="sim-link-content">
					<img src="${iconPath}" class="sim-link-icon">
					<div class="d-flex flex-column">
						<span class="sim-link-title">${label}</span>
						${specIndexes.length == 1 ? this.launchStatusLabel({ type: 'Spec', index: specIndexes[0] }) : ''}
					</div>
				</div>
			</a>
    	`;

		return fragment.children[0] as HTMLElement;
	}

	private buildSpecLink(specIndex: Spec): HTMLElement {
		const href = getSpecSiteUrl(specIndex);
		const textKlass = this.getContextualKlass({ type: 'Spec', index: specIndex });
		const iconPath = this.getSimIconPath({ type: 'Spec', index: specIndex });
		const className = classNames[specToClass[specIndex]];
		const specLabel = this.specLabels[specIndex];

		const fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<a href="${href}" class="sim-link ${textKlass}" role="button">
				<div class="sim-link-content">
					<img src="${iconPath}" class="sim-link-icon">
					<div class="d-flex flex-column">
						<span class="sim-link-label">${className}</span>
						<span class="sim-link-title">${specLabel}</span>
						${this.launchStatusLabel({ type: 'Spec', index: specIndex })}
					</div>
				</div>
			</a>
    	`;

		return fragment.children[0] as HTMLElement;
	}

	private launchStatusLabel(data: SpecOptions | RaidOptions): string {
		if (
			(data.type == 'Raid' && raidSimStatus == LaunchStatus.Launched) ||
			(data.type == 'Spec' && simLaunchStatuses[data.index] == LaunchStatus.Launched)
		)
			return '';

		const label =
			data.type == 'Raid'
				? LaunchStatus[raidSimStatus]
				: LaunchStatus[simLaunchStatuses[data.index]];
		const elem = document.createElement('span');
		elem.classList.add('launch-status-label', 'text-brand');
		elem.textContent = label;

		return elem.outerHTML;
	}

	private getSimIconPath(data: ClassOptions | SpecOptions | RaidOptions): string {
		let iconPath: string;

		if (data.type == 'Raid') {
			iconPath = raidSimIcon;
		} else if (data.type == 'Class') {
			const className = classNames[data.index];
			iconPath = `/wotlk/assets/img/${className.toLowerCase().replace(/\s/g, '_')}_icon.png`;
		} else {
			iconPath = titleIcons[data.index];
		}

		return iconPath;
	}

	private getContextualKlass(data: ClassOptions | SpecOptions | RaidOptions): string {
		if (data.type == 'Raid')
			// Raid link
			return 'text-white';
		else if (data.type == 'Class')
			// Class links
			return textCssClassForClass(data.index);
		else return textCssClassForSpec(data.index);
	}
}
