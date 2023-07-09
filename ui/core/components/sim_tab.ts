import { SimUI } from "../sim_ui";
import { Component } from "./component";

export interface SimTabConfig {
	identifier: string,
	title: string,
}

export abstract class SimTab extends Component {
	protected simUI: SimUI;
	protected config: SimTabConfig;

	readonly navItem: HTMLElement;
	readonly navLink: HTMLElement;
	readonly contentContainer: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: SimUI, config: SimTabConfig) {
		super(parentElem, 'sim-tab');

		this.rootElem.classList.add(config.identifier);

		this.simUI = simUI;
		this.config = config;

		this.rootElem.id = this.config.identifier;
		this.rootElem.classList.add('tab-pane', 'fade');

		if (parentElem.childNodes.length == 0)
			this.rootElem.classList.add('active', 'show');

		this.navItem = this.buildNavItem();
		this.navLink = this.navItem.children[0] as HTMLElement;
		this.contentContainer = document.createElement('div');
		this.contentContainer.classList.add('tab-pane-content-container');
		this.rootElem.appendChild(this.contentContainer);

		this.simUI.simHeader.addSimTabLink(this);
	}

	private buildNavItem(): HTMLElement {
		const tabFragment = document.createElement('fragment');
		tabFragment.innerHTML = `
			<li class="${this.config.identifier} nav-item" role="presentation">
				<a
					class="nav-link"
					data-bs-toggle="tab"
					data-bs-target="#${this.config.identifier}"
					type="button"
					role="tab"
					aria-controls="${this.config.identifier}"
				>${this.config.title}</a>
			</li>
		`;

		return tabFragment.children[0] as HTMLElement;
	}

	protected abstract buildTabContent(): void;

	protected buildColumn(index: number, customCssClass: string): HTMLElement {
		let column = document.createElement('div');
		column.classList.add('tab-panel-col', `${customCssClass}-${index}`)
		return column;
	}
}
