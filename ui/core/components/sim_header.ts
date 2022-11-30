import { Component } from './component';
import { SettingsMenu } from './settings_menu';
import { SimUI } from '../sim_ui';
import { EventID, TypedEvent } from '../typed_event';
import { Tooltip } from 'bootstrap';

// Config for displaying a warning to the user whenever a condition is met.
interface SimWarning {
	updateOn: TypedEvent<any>,
	getContent: () => string | Array<string>,
}

export class SimHeader extends Component {
  private simUI: SimUI;

  private simTabsContainer: HTMLElement;
	private importExportContainer: HTMLElement;
	private simToolbar: HTMLElement;
  private warningsLink: HTMLElement;
	private knownIssuesLink: HTMLElement;

  private warnings: Array<SimWarning> = [];

  constructor(parentElem: HTMLElement, simUI: SimUI) {
    super(parentElem, 'sim-header', headerRoot);

    this.simUI = simUI;

    this.simTabsContainer = this.rootElem.querySelector('.sim-tabs') as HTMLElement;
		this.importExportContainer = this.rootElem.querySelector('.import-export') as HTMLElement;
		this.simToolbar = this.rootElem.querySelector('.sim-toolbar') as HTMLElement;	
    this.warningsLink = this.rootElem.querySelector('.sim-toolbar .warnings') as HTMLElement;
    this.knownIssuesLink = this.rootElem.querySelector('.sim-toolbar .known-issues') as HTMLElement;

    this.addBugReportLink();
    this.addDownloadBinaryLink();
    this.addSimOptionsLink();
    this.addPatreonLink();
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

  addImportLink(importElem: HTMLElement) {
    this.importExportContainer.appendChild(importElem);
  }

  addExportLink(exportElem: HTMLElement) {
    this.importExportContainer.appendChild(exportElem);
  }

  addToolbarItem(elem: HTMLElement) {
		const toolbarItem = document.createElement('div');
		toolbarItem.appendChild(elem);
		toolbarItem.classList.add('sim-toolbar-item');
		this.simToolbar.appendChild(toolbarItem);
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

  addBugReportLink() {
		let bugReportFragment = document.createElement('fragment');
		bugReportFragment.innerHTML = `
			<a
				href="https://github.com/wowsims/wotlk/issues/new/choose"
				target="_blank"
				onclick="${onclick}"
				data-bs-toggle="tooltip"
				data-bs-placement="bottom"
				data-bs-title="Report a bug or<br>Request a feature"
				data-bs-html="true"
			>
				<i class="fas fa-bug fa-lg"></i>
			</a>
		`;

		let bugReportLink = bugReportFragment.children[0] as HTMLElement;
		new Tooltip(bugReportLink);
		this.addToolbarItem(bugReportLink);
	}

  addDownloadBinaryLink() {
    let downloadFragment = document.createElement('fragment');
    downloadFragment.innerHTML = `
      <a
        href="https://github.com/wowsims/wotlk/releases"
        target="_blank"
        class="downbin"
        data-bs-toggle="tooltip"
      >
        <i class="fas fa-gauge-high fa-lg"></i>
      </a>
    `;

    let downloadBinary = downloadFragment.children[0] as HTMLElement;

		if (document.location.href.includes("localhost")) {
			fetch(document.location.protocol + "//" + document.location.host + "/version").then(resp => {
				resp.json()
					.then((versionInfo) => {
						if (versionInfo.outdated == 2) {
              downloadBinary.setAttribute('data-bs-title', 'Newer version of simulator available for download')
							downloadBinary.classList.add('link-danger')
              new Tooltip(downloadBinary);
              this.addToolbarItem(downloadBinary);
						}
					})
					.catch(error => {
						console.warn('No version info found!');
					});
			});
		} else {
      downloadBinary.setAttribute('data-bs-title', 'Download simulator for faster simulating')
      new Tooltip(downloadBinary);
      this.addToolbarItem(downloadBinary);
		}
  }

  addSimOptionsLink() {
		let optionsFragment = document.createElement('fragment');
		optionsFragment.innerHTML = `
			<a
				href="javascript:void(0)"
				class="sim-options"
				role="button"
				onmousedown="${onclick}"
				data-bs-toggle="tooltip"
				data-bs-placement="bottom"
				data-bs-title="Show Sim Options"
			>
				<i class="fas fa-cog fa-lg"></i>
			</a>
		`;
		let optionsLink = optionsFragment.children[0] as HTMLElement;
		optionsLink.addEventListener('click', () => new SettingsMenu(this.rootElem, this.simUI));
		new Tooltip(optionsLink);
		this.addToolbarItem(optionsLink);
	}

	addPatreonLink() {
		let patreonFragment = document.createElement('fragment');
		patreonFragment.innerHTML = `
			<a href="https://patreon.com/wowsims" target="_blank" class="link-alt patreon-link">
				<i class="fab fa-patreon fa-lg"></i>
				<span>Support our devs</span>
			</a>
		`;
		let patreonLink = patreonFragment.children[0] as HTMLElement;
		this.addToolbarItem(patreonLink);
    patreonLink.parentElement?.classList.add('py-2', 'border-start');
	}
}

let headerFragment = document.createElement('fragment');
headerFragment.innerHTML = `
  <header id="simHeader">
    <ul class="sim-tabs nav nav-tabs" role="tablist"></ul>
    <div class="import-export"></div>
    <div class="sim-toolbar">
    <div class="sim-toolbar-item hide">
      <a
        href="javascript:void(0)"
        class="warnings link-warning"
        role="button"
        data-bs-toggle="tooltip"
        data-bs-placement="bottom"
        data-bs-html="true"
        data-bs-title="<ul class='text-start ps-3 mb-0'></ul>"
      >
        <i class="fas fa-exclamation-triangle fa-3x"></i>
      </a>
    </div>
      <div class="sim-toolbar-item">
        <a
          href="javascript:void(0)"
          class="known-issues link-danger"
          role="button"
          data-bs-toggle="tooltip"
          data-bs-placement="bottom"
          data-bs-html="true"
          data-bs-title="<ul class='text-start ps-3 mb-0'></ul>"
        >Known Issues</a>
      </div>
    </div>
  </header>
`;
let headerRoot = headerFragment.children[0] as HTMLElement;
