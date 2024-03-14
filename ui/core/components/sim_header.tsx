import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element } from 'tsx-vanilla';

import { SimUI } from '../sim_ui';
import { Component } from './component';
import { SettingsMenu } from './settings_menu';
import { SimTab } from './sim_tab';
import { SocialLinks } from './social_links';

interface ToolbarLinkArgs {
	parent: HTMLElement;
	href?: string;
	text?: string;
	icon?: string;
	tooltip?: string | HTMLElement;
	classes?: string;
	onclick?: () => void;
}

export class SimHeader extends Component {
	private simUI: SimUI;

	private simTabsContainer: HTMLElement;
	private simToolbar: HTMLElement;
	private knownIssuesLink: HTMLElement;
	private knownIssuesContent: HTMLElement;

	constructor(parentElem: HTMLElement, simUI: SimUI) {
		super(parentElem, 'sim-header');
		this.simUI = simUI;
		this.simTabsContainer = this.rootElem.querySelector('.sim-tabs') as HTMLElement;
		this.simToolbar = this.rootElem.querySelector('.sim-toolbar') as HTMLElement;

		this.knownIssuesContent = (<ul className="text-start ps-3 mb-0"></ul>) as HTMLElement;
		this.knownIssuesLink = this.addKnownIssuesLink();
		this.addBugReportLink();
		this.addDownloadBinaryLink();
		this.addSimOptionsLink();
		this.addSocialLinks();

		// Allow styling the sticky header
		new IntersectionObserver(
			([e]) => e.target.classList.toggle('stuck', e.intersectionRatio < 1),
			{ threshold: [1] },
		).observe(this.rootElem);
	}

	activateTab(className: string) {
		(this.simTabsContainer.getElementsByClassName(className)[0] as HTMLElement).click();
	}

	addTab(title: string, contentId: string) {
		const isFirstTab = this.simTabsContainer.children.length == 0;

		const classes = `${contentId} nav-item`;
		const tab = (
			<li className={classes} attributes={{ role: 'presentation' }}>
				<a
					className={`nav-link ${isFirstTab && 'active'}`}
					dataset={{
						bsToggle: 'tab',
						bsTarget: `#${contentId}`,
					}}
					attributes={{
						role: 'tab',
						'aria-selected': isFirstTab,
					}}
					type="button">
					{title}
				</a>
			</li>
		);
		tab.setAttribute('aria-controls', contentId);

		this.simTabsContainer.appendChild(tab);
	}

	addSimTabLink(tab: SimTab) {
		const isFirstTab = this.simTabsContainer.children.length == 0;

		tab.navLink.setAttribute('aria-selected', isFirstTab.toString());

		if (isFirstTab) tab.navLink.classList.add('active', 'show');

		this.simTabsContainer.appendChild(tab.navItem);
	}

	addImportLink(label: string, onClick: (parent: HTMLElement) => void, hideInRaidSim?: boolean) {
		this.addImportExportLink('import-dropdown', label, onClick, hideInRaidSim);
	}
	addExportLink(label: string, onClick: (parent: HTMLElement) => void, hideInRaidSim?: boolean) {
		this.addImportExportLink('export-dropdown', label, onClick, hideInRaidSim);
	}
	private addImportExportLink(
		cssClass: string,
		label: string,
		onClick: (parent: HTMLElement) => void,
		hideInRaidSim?: boolean,
	) {
		const dropdownElem = this.rootElem.getElementsByClassName(cssClass)[0] as HTMLElement;
		const menuElem = dropdownElem.getElementsByClassName('dropdown-menu')[0] as HTMLElement;

		const itemElem = (
			<li>
				<a
					href="javascript:void(0)"
					className="dropdown-item"
					attributes={{
						role: 'button',
					}}>
					{label}
				</a>
			</li>
		);

		const linkElem = itemElem.children[0];
		linkElem.addEventListener('click', () => onClick(menuElem));
		menuElem.appendChild(itemElem);
	}

	private addToolbarLink(args: ToolbarLinkArgs): HTMLElement {
		const item = (
			<div className="sim-toolbar-item">
				<a
					href={args.href ? args.href : 'javascript:void(0)'}
					className={args.classes}
					target={args.href ? '_blank' : '_self'}>
					{args.icon && <i className={args.icon}></i>}
					{args.text ? ` ${args.text} ` : ''}
				</a>
			</div>
		);

		const link = item.children[0];

		if (args.onclick) {
			link.addEventListener('click', () => {
				// Typescript is requiring this even though the condition is being done already above
				if (args.onclick) args.onclick();
			});
		}

		if (args.tooltip) {
			new Tooltip(link, {
				placement: 'bottom',
				title: args.tooltip,
				html: true,
			});
		}

		return args.parent.appendChild(item) as HTMLElement;
	}

	private addKnownIssuesLink(): HTMLElement {
		return this.addToolbarLink({
			parent: this.simToolbar,
			text: 'Known Issues',
			tooltip: this.knownIssuesContent,
			classes: 'known-issues link-danger hide',
		}).children[0] as HTMLElement;
	}

	addKnownIssue(issue: string) {
		this.knownIssuesContent.appendChild(<li>{issue}</li>);
		this.knownIssuesLink.classList.remove('hide');
		Tooltip.getInstance(this.knownIssuesLink)?.setContent({
			'.tooltip-inner': this.knownIssuesContent,
		});
	}

	private addBugReportLink() {
		this.addToolbarLink({
			href: 'https://github.com/wowsims/wotlk/issues/new/choose',
			parent: this.simToolbar,
			icon: 'fas fa-bug fa-lg',
			tooltip: 'Report a bug or<br>Request a feature',
		});
	}

	private addDownloadBinaryLink() {
		const href = 'https://github.com/wowsims/wotlk/releases';
		const icon = 'fas fa-gauge-high fa-lg';
		const parent = this.simToolbar;

		if (document.location.href.includes('localhost')) {
			fetch('/version').then(resp => {
				resp.json()
					.then(versionInfo => {
						if (versionInfo.outdated == 2) {
							this.addToolbarLink({
								href: href,
								parent: parent,
								icon: icon,
								tooltip: 'Newer version of simulator available for download',
								classes: 'downbin link-danger',
							});
						}
					})
					.catch(error => {
						console.warn('No version info found!');
					});
			});
		} else {
			this.addToolbarLink({
				href: href,
				parent: parent,
				icon: icon,
				tooltip: 'Download simulator for faster simulating',
				classes: 'downbin',
			});
		}
	}

	private addSimOptionsLink() {
		this.addToolbarLink({
			parent: this.simToolbar,
			icon: 'fas fa-cog fa-lg',
			tooltip: 'Show Sim Options',
			classes: 'sim-options',
			onclick: () => new SettingsMenu(this.simUI.rootElem, this.simUI),
		});
	}

	private addSocialLinks() {
		const container = document.createElement('div');
		container.classList.add('sim-toolbar-socials');
		this.simToolbar.appendChild(container);

		this.addDiscordLink(container);
		this.addGitHubLink(container);
		this.addPatreonLink(container);
	}

	private addDiscordLink(container: HTMLElement) {
		container.appendChild(
			<div className="sim-toolbar-item">{SocialLinks.buildDiscordLink()}</div>,
		);
	}

	private addGitHubLink(container: HTMLElement) {
		container.appendChild(
			<div className="sim-toolbar-item">{SocialLinks.buildGitHubLink()}</div>,
		);
	}

	private addPatreonLink(container: HTMLElement) {
		container.appendChild(
			<div className="sim-toolbar-item">{SocialLinks.buildPatreonLink()}</div>,
		);
	}

	protected customRootElement(): HTMLElement {
		return (
			<header className="sim-header">
				<div className="sim-header-container">
					<ul className="sim-tabs nav nav-tabs" attributes={{ role: 'tablist' }}></ul>
					<div className="import-export within-raid-sim-hide">
						<div className="dropdown sim-dropdown-menu import-dropdown">
							<a
								href="javascript:void(0)"
								className="import-link"
								attributes={{ role: 'button', 'aria-expanded': 'false' }}
								dataset={{ bsToggle: 'dropdown', bsDisplay: 'dynamic' }}>
								<i className="fa fa-download"></i>
								{' Import '}
							</a>
							<ul className="dropdown-menu"></ul>
						</div>
						<div className="dropdown sim-dropdown-menu export-dropdown">
							<a
								href="javascript:void(0)"
								className="export-link"
								attributes={{ role: 'button', 'aria-expanded': 'false' }}
								dataset={{ bsToggle: 'dropdown', bsDisplay: 'dynamic' }}>
								<i className="fa fa-right-from-bracket"></i>
								{' Export '}
							</a>
							<ul className="dropdown-menu"></ul>
						</div>
					</div>
					<div className="sim-toolbar"></div>
				</div>
			</header>
		) as HTMLElement;
	}
}
