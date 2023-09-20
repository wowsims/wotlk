import { Tooltip } from 'bootstrap';
import { Component } from '../components/component.js';
import { TypedEvent } from '../typed_event.js';

import { element, fragment } from 'tsx-vanilla';

// Config for displaying a warning to the user whenever a condition is met.
interface SimWarning {
	updateOn: TypedEvent<any>,
	getContent: () => string | Array<string>,
}

interface WarningLinkArgs {
	parent: HTMLElement,
	href?: string,
	text?: string,
	icon?: string,
	tooltip?: string,
	classes?: string,
	onclick?: Function
}

export class ResultsViewer extends Component {
	readonly pendingElem: HTMLElement;
	readonly contentElem: HTMLElement;
	readonly warningElem: HTMLElement;
	private warningsLink: HTMLElement;

	private warnings: Array<SimWarning> = [];

	constructor(parentElem: HTMLElement) {
		super(parentElem, 'results-viewer');
		this.rootElem.appendChild(
			<>
				<div className="results-pending">
				<div className="loader"></div>
				</div>
				<div className="results-content">
				</div>
				<div className="warning-zone" style="text-align: center">
				</div>
			</>
		);
		this.pendingElem = this.rootElem.getElementsByClassName('results-pending')[0] as HTMLElement;
		this.contentElem = this.rootElem.getElementsByClassName('results-content')[0] as HTMLElement;
		this.warningElem = this.rootElem.getElementsByClassName('warning-zone')[0] as HTMLElement;


		this.warningsLink = this.addWarningsLink();
		this.updateWarnings();

		this.hideAll();
	}

	private addWarningLink(args: WarningLinkArgs): HTMLElement {
		let item = (
			<div className="sim-toolbar-item">
				<a
					href={args.href ? args.href : 'javascript:void(0)'}
					target={args.href ? '_blank' : '_self'}
					className={args.classes}
				>
					{args.icon && <i className={args.icon}></i>}
					{args.text ? args.text : ''}
				</a>
			</div>
		);

		let link = item.children[0] as HTMLElement;

		if (args.onclick) {
			link.addEventListener('click', () => {
				if (args.onclick)
					args.onclick();
			});
		}

		let cfg = {};
		if (args.tooltip) {
			cfg = {
				title: args.tooltip,
				html: true,
				placement: 'bottom',
			};
			link.setAttribute('data-bs-title', args.tooltip);
		}
		new Tooltip(link, cfg);
		args.parent.appendChild(item);

		return item as HTMLElement;
	}

	private addWarningsLink(): HTMLElement {
		return this.addWarningLink({
			parent: this.warningElem,
			icon: 'fas fa-exclamation-triangle fa-3x',
			tooltip: "<ul class='text-start ps-3 mb-0'></ul>",
			classes: 'warning link-warning'
		}).children[0] as HTMLElement;
	}

	addWarning(warning: SimWarning) {
		this.warnings.push(warning);
		warning.updateOn.on(() => this.updateWarnings());
		this.updateWarnings();
	}

	private updateWarnings() {
		const activeWarnings = this.warnings.map(warning => warning.getContent()).flat().filter(content => content != '');
		let tooltipFragment = document.createElement('fragment');
		tooltipFragment.innerHTML = this.warningsLink.getAttribute('data-bs-title') as string;
		let list = tooltipFragment.children[0] as HTMLElement;
		list.innerHTML = '';
		if (activeWarnings.length == 0) {
			this.warningsLink.parentElement?.classList?.add('hide');
		} else {
			this.warningsLink.parentElement?.classList?.remove('hide');
			activeWarnings.forEach(warning => {
				let listItem = document.createElement('li');
				listItem.innerHTML = warning;
				list.appendChild(listItem);
			});
		}
		this.warningsLink.setAttribute('data-bs-title', list.outerHTML);
		new Tooltip(this.warningsLink, {
			title: list.outerHTML,
			html: true,
			placement: 'bottom',
		});
	}

	hideAll() {
		this.contentElem.style.display = 'none';
		this.pendingElem.style.display = 'none';
	}

	setPending() {
		this.contentElem.style.display = 'none';
		this.pendingElem.style.display = 'block';
	}

	setContent(innerHTML: string) {
		this.contentElem.innerHTML = innerHTML;
		this.contentElem.style.display = 'block';
		this.pendingElem.style.display = 'none';
	}
}
