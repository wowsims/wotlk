import { Component } from './component';
import { SettingsMenu } from './settings_menu';
import { SimUI } from '../sim_ui';
import { TypedEvent } from '../typed_event';
import { Tooltip } from 'bootstrap';
import { SimTab } from './sim_tab';

// Config for displaying a warning to the user whenever a condition is met.
interface SimWarning {
	updateOn: TypedEvent<any>,
	getContent: () => string | Array<string>,
}

interface ToolbarLinkArgs {
	parent: HTMLElement,
	href?: string,
	text?: string,
	icon?: string,
	tooltip?: string,
	classes?: string,
	onclick?: Function
}

export class SimHeader extends Component {
  private simUI: SimUI;

  private simTabsContainer: HTMLElement;
	private simToolbar: HTMLElement;
  private warningsLink: HTMLElement;
	private knownIssuesLink: HTMLElement;

  private warnings: Array<SimWarning> = [];

  constructor(parentElem: HTMLElement, simUI: SimUI) {
    super(parentElem, 'sim-header');

    this.simUI = simUI;

    this.simTabsContainer = this.rootElem.querySelector('.sim-tabs') as HTMLElement;
		this.simToolbar = this.rootElem.querySelector('.sim-toolbar') as HTMLElement;	

		this.warningsLink = this.addWarningsLink();
		this.updateWarnings();

		this.knownIssuesLink = this.addKnownIssuesLink();
    this.addBugReportLink();
    this.addDownloadBinaryLink();
    this.addSimOptionsLink();
		this.addSocialLinks();
  }

  activateTab(className: string) {
	(this.simTabsContainer.getElementsByClassName(className)[0] as HTMLElement).click();
  }

  addTab(title: string, contentId: string) {
		const isFirstTab = this.simTabsContainer.children.length == 0;

		const tabFragment = document.createElement('fragment');
		tabFragment.innerHTML = `
			<li class="${contentId} nav-item" role="presentation">
				<a
					class="nav-link ${isFirstTab ? 'active' : ''}"
					data-bs-toggle="tab"
					data-bs-target="#${contentId}"
					type="button"
					role="tab"
					aria-controls="${contentId}"
					aria-selected="${isFirstTab}"
				>${title}</a>
			</li>
		`;

		this.simTabsContainer.appendChild(tabFragment.children[0] as HTMLElement);
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
	private addImportExportLink(cssClass: string, label: string, onClick: (parent: HTMLElement) => void, hideInRaidSim?: boolean) {
		const dropdownElem = this.rootElem.getElementsByClassName(cssClass)[0] as HTMLElement;
		const menuElem = dropdownElem.getElementsByClassName('dropdown-menu')[0] as HTMLElement;

		const itemFragment = document.createElement('fragment');
		itemFragment.innerHTML = `
			<li>
				<a
					href="javascript:void(0)"
					class="dropdown-item"
					role="button"
				>${label}</a>
			</li>
		`;
		const itemElem = itemFragment.children[0] as HTMLElement;
		const linkElem = itemElem.children[0] as HTMLElement;
		linkElem.addEventListener('click', () => onClick(menuElem));
		menuElem.appendChild(itemElem);
	}

	private addToolbarLink(args: ToolbarLinkArgs): HTMLElement {
		let fragment = document.createElement('fragment');
		fragment.innerHTML = `
			<div class="sim-toolbar-item">
				<a
					href="${args.href ? args.href : 'javascript:void(0)'}"
					${args.href ? 'target="_blank"' : ''}
					class="${args.classes}"
					${args.tooltip ? 'data-bs-toggle="tooltip"' : ''}
					${args.tooltip ? 'data-bs-placement="bottom"' : ''}
					${args.tooltip ? `data-bs-title="${args.tooltip}"` : ''}
					${args.tooltip ? 'data-bs-html="true"' : ''}
				>
					${args.icon ? `<i class="${args.icon}"></i>` : ''}
					${args.text ? args.text : ''}
				</a>
			</div>
		`;

		let item = fragment.children[0] as HTMLElement;
		let link = item.children[0] as HTMLElement;

		if (args.onclick) {
			link.addEventListener('click', () => {
				// Typescript is requiring this even though the condition is being done already above
				if (args.onclick)
					args.onclick();
			});
		}

		new Tooltip(link);
		args.parent.appendChild(item);

		return item;
	}

	private addWarningsLink(): HTMLElement {
		return this.addToolbarLink({
			parent: this.simToolbar,
			icon: 'fas fa-exclamation-triangle fa-3x',
			tooltip: "<ul class='text-start ps-3 mb-0'></ul>",
			classes: 'warnings link-warning'
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
    new Tooltip(this.warningsLink);
	}

	private addKnownIssuesLink(): HTMLElement {
		return this.addToolbarLink({
			parent: this.simToolbar,
			text: "Known Issues",
			tooltip: "<ul class='text-start ps-3 mb-0'></ul>",
			classes: "known-issues link-danger"
		}).children[0] as HTMLElement;
	}

  addKnownIssue(issue: string) {
		let tooltipFragment = document.createElement('fragment');
		tooltipFragment.innerHTML = this.knownIssuesLink.getAttribute('data-bs-title') as string;
		let list = tooltipFragment.children[0] as HTMLElement;
		let listItem = document.createElement('li');
		listItem.innerHTML = issue;
		list.appendChild(listItem);
		this.knownIssuesLink.setAttribute('data-bs-title', list.outerHTML);
    new Tooltip(this.knownIssuesLink);
	}

  private addBugReportLink() {
		this.addToolbarLink({
			href: "https://github.com/wowsims/wotlk/issues/new/choose",
			parent: this.simToolbar,
			icon: "fas fa-bug fa-lg",
			tooltip: "Report a bug or<br>Request a feature"
		})
	}

  private addDownloadBinaryLink() {
		let href = "https://github.com/wowsims/wotlk/releases";
		let icon = "fas fa-gauge-high fa-lg"
		let parent = this.simToolbar;

		if (document.location.href.includes("localhost")) {
			fetch("/version").then(resp => {
				resp.json()
					.then((versionInfo) => {
						if (versionInfo.outdated == 2) {
							this.addToolbarLink({
								href: href,
								parent: parent,
								icon: icon,
								tooltip: "Newer version of simulator available for download",
								classes: "downbin link-danger",
							})
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
				tooltip: "Download simulator for faster simulating",
				classes: "downbin",
			})
		}
  }

  private addSimOptionsLink() {
		this.addToolbarLink({
			parent: this.simToolbar,
			icon: "fas fa-cog fa-lg",
			tooltip: "Show Sim Options",
			classes: 'sim-options',
			onclick: () => new SettingsMenu(this.simUI.rootElem, this.simUI)
		})
	}

	private addSocialLinks() {
		let container = document.createElement('div');
		container.classList.add('sim-toolbar-socials')
		this.simToolbar.appendChild(container);

		this.addDiscordLink(container);
		this.addGitHubLink(container);
		this.addPatreonLink(container);
	}

	private addDiscordLink(container: HTMLElement) {
		this.addToolbarLink({
			href: "https://discord.gg/p3DgvmnDCS",
			parent: container,
			icon: "fab fa-discord fa-lg",
			tooltip: "Join us on Discord",
			classes: "discord-link link-alt"
		})
	}

	private addGitHubLink(container: HTMLElement) {
		this.addToolbarLink({
			href: "https://github.com/wowsims/wotlk",
			parent: container,
			icon: "fab fa-github fa-lg",
			tooltip: "Contribute on GitHub",
			classes: "github-link link-alt"
		})
	}

	private addPatreonLink(container: HTMLElement) {
		this.addToolbarLink({
			href: "https://patreon.com/wowsims",
			parent: container,
			text: "Support our devs",
			icon: "fab fa-patreon fa-lg",
			classes: "patreon-link link-alt"
		})
	}

	protected customRootElement(): HTMLElement {
		let headerFragment = document.createElement('fragment');
		headerFragment.innerHTML = `
			<header class="sim-header">
				<ul class="sim-tabs nav nav-tabs" role="tablist"></ul>
				<div class="import-export within-raid-sim-hide">
					<div class="dropdown sim-dropdown-menu import-dropdown">
						<a href="javascript:void(0)" class="import-link" role="button" data-bs-toggle="dropdown" data-bs-display="dynamic" aria-expanded="false">
							<i class="fa fa-download"></i>
							Import
						</a>
						<ul class="dropdown-menu"></ul>
					</div>
					<div class="dropdown sim-dropdown-menu export-dropdown">
						<a href="javascript:void(0)" class="export-link" role="button" data-bs-toggle="dropdown" data-bs-display="dynamic" aria-expanded="false">
							<i class="fa fa-right-from-bracket"></i>
							Export
						</a>
						<ul class="dropdown-menu"></ul>
					</div>
				</div>
				<div class="sim-toolbar"></div>
			</header>
		`;

		return headerFragment.children[0] as HTMLElement;
	}
}
