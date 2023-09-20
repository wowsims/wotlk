import { Tooltip } from 'bootstrap';
import { Component } from './component.js';

import { element, fragment } from 'tsx-vanilla';

export interface ContentBlockHeaderConfig {
	title: string,
	extraCssClasses?: Array<string>,
	titleTag?: string,
	tooltip?: string,
}

export interface ContentBlockConfig {
	bodyClasses?: Array<string>,
	extraCssClasses?: Array<string>,
	rootElem?: HTMLElement,
	header?: ContentBlockHeaderConfig,
}

export class ContentBlock extends Component {
	readonly headerElement: HTMLElement | null;
	readonly bodyElement: HTMLElement;

	readonly config: ContentBlockConfig;

	constructor(parent: HTMLElement, cssClass: string, config: ContentBlockConfig) {
		super(parent, 'content-block', config.rootElem);
		this.config = config;
		this.rootElem.classList.add(cssClass);

		if (config.extraCssClasses) {
			this.rootElem.classList.add(...config.extraCssClasses);
		}

		this.headerElement = this.buildHeader();
		this.bodyElement = this.buildBody();
		config.bodyClasses?.forEach((cl) => {
			this.bodyElement.classList.add(cl);
		})
	}

	private buildHeader(): HTMLElement | null {
		if (this.config.header && Object.keys(this.config.header).length) {
			let TitleTag = this.config.header.titleTag || 'h6';
			let header =(
				<div className="content-block-header">
					<TitleTag className="content-block-title">
						{this.config.header.title}
					</TitleTag>
				</div>
			);

			if (this.config.header.extraCssClasses) {
				header.classList.add(...this.config.header.extraCssClasses);
			}

			if (this.config.header.tooltip)
				Tooltip.getOrCreateInstance(header.querySelector('.content-block-title') as HTMLElement,  {
					html: true,
					title: this.config.header.tooltip
				});

			this.rootElem.appendChild(header);

			return header as HTMLElement;
		} else {
			return null;
		}
	}

	private buildBody(): HTMLElement {
		let bodyElem = document.createElement('div');
		bodyElem.classList.add('content-block-body');

		this.rootElem.appendChild(bodyElem);

		return bodyElem;
	}
}
