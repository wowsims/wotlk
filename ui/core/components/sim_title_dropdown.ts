import { Component } from './component.js';
import { getLaunchedSimsForClass, raidSimLaunched } from '../launched_sims.js';
import { Class, Spec } from '../proto/common.js';
import {
  getSpecSiteUrl,
  raidSimSiteUrl,
  specNames,
  classNames,
  specToClass,
  raidSimIcon,
  raidSimLabel
} from '../proto_utils/utils.js';
import { classList, getClassIcon, getSpecIcon } from '../proto_utils/class_spec_utils.js';

interface ClassOptions {
  type: 'Class',
  index: Class
}

interface SpecOptions {
  type: 'Spec'
  index: Spec
}

interface RaidOptions {
  type: 'Raid'
}

// Dropdown menu for selecting a player.
export class SimTitleDropdown extends Component {

  private readonly dropdownMenu: HTMLElement;

  private readonly specLabels: Record<Spec, string> = {
    [Spec.SpecBalanceDruid]:       'Balance',
    [Spec.SpecElementalShaman]:    'Elemental',
    [Spec.SpecEnhancementShaman]:  'Enhancement',
    [Spec.SpecFeralDruid]:         'Feral',
    [Spec.SpecFeralTankDruid]:     'Feral Tank',
    [Spec.SpecHunter]:             'Hunter',
    [Spec.SpecMage]:               'Mage',
    [Spec.SpecRogue]:              'Rogue',
    [Spec.SpecRetributionPaladin]: 'Retribution',
    [Spec.SpecProtectionPaladin]:  'Protection',
    [Spec.SpecHealingPriest]:      'Priest',
    [Spec.SpecShadowPriest]:       'Shadow',
    [Spec.SpecSmitePriest]:        'Smite',
    [Spec.SpecWarlock]:            'Warlock',
    [Spec.SpecWarrior]:            'DPS',
    [Spec.SpecProtectionWarrior]:  'Protection',
    [Spec.SpecDeathknight]:        'DPS',
    [Spec.SpecTankDeathknight]:    'Tank',
  }

  constructor(parent: HTMLElement, currentSpecIndex: Spec | null) {
    super(parent, 'sim-title-dropdown-root');

    let rootLinkArgs: SpecOptions|RaidOptions = currentSpecIndex === null ? {type: 'Raid'} : {type: 'Spec', index: currentSpecIndex}
    let rootLink = this.buildRootSimLink(rootLinkArgs);

    this.rootElem.innerHTML = `
      <div class="dropdown sim-link-dropdown">
        ${rootLink.outerHTML}
        <ul class="dropdown-menu"></ul>
      </div>
    `;

    this.dropdownMenu = this.rootElem.getElementsByClassName('dropdown-menu')[0] as HTMLElement;
    this.buildDropdown();

    // Prevent Bootstrap from closing the menu instead of opening class menus
    this.dropdownMenu.addEventListener('click', (event) => {
      let target = event.target as HTMLElement;
      let link = target.closest('a:not([href="javascript:void(0)"]');

      if (!link) {
        event.stopPropagation();
        event.preventDefault();
      }
    })
  }
  
  private buildDropdown() {
    if (raidSimLaunched) {
      // Add the raid sim to the top of the dropdown
      let raidListItem = document.createElement('li');
      raidListItem.appendChild(this.buildRaidLink());
      this.dropdownMenu.appendChild(raidListItem);
    }

    classList.forEach( (classIndex) => {
      let listItem = document.createElement('li');
      let sims = getLaunchedSimsForClass(classIndex);

      if (sims.length == 1) {
        // The class only has one listed sim so make a direct link to the sim
        listItem.appendChild(this.buildClassLink(classIndex));
        this.dropdownMenu.appendChild(listItem);
      } else if (sims.length > 1) {
        // Add the class to the dropdown with an additional spec dropdown
        listItem.appendChild(this.buildClassDropdown(classIndex));
        this.dropdownMenu.appendChild(listItem);
      }
    });
  }

  private buildClassDropdown(classIndex: Class) {
    let sims = getLaunchedSimsForClass(classIndex);
    let dropdownFragment = document.createElement('fragment');
    let dropdownMenu = document.createElement('ul');
    dropdownMenu.classList.add('dropdown-menu');

    // Generate the class link to act as a dropdown toggle for the spec dropdown
    let classLink = this.buildClassLink(classIndex);

    // Generate links for a class's specs
    sims.forEach( (specIndex) => {
      let listItem = document.createElement('li');
      let link = this.buildSpecLink(specIndex);

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
    let iconPath  = this.getSimIconPath(data);;
    let textKlass = this.getContextualKlass(data);
    let label;

    if (data.type == 'Raid')
      label = raidSimLabel;
    else {
      let classIndex = specToClass[data.index];
      if (getLaunchedSimsForClass(classIndex).length > 1)
        // If the class has multiple sims, use the spec name
        label = specNames[data.index];
      else
        // If the class has only 1 sim, use the class name
        label = classNames[classIndex];
    }

    let fragment = document.createElement('fragment');
    fragment.innerHTML = `
      <a href="javascript:void(0)" class="sim-link ${textKlass}" role="button" data-bs-toggle="dropdown" aria-expanded="false">
        <div class="sim-link-content">
          <img src="${iconPath}" class="sim-link-icon">
          <div class="d-flex flex-column">
            <span class="sim-link-label text-white">WoWSims - WOTLK</span>
            <span class="sim-link-label">${label}</span>
          </div>
        </div>
      </a>
    `;

    return fragment.children[0] as HTMLElement;
  }

  private buildRaidLink(): HTMLElement {
    let href      = raidSimSiteUrl;
    let textKlass = this.getContextualKlass({type: 'Raid'});
    let iconPath  = this.getSimIconPath({type: 'Raid'});
    let label     = raidSimLabel;

    let fragment = document.createElement('fragment');
    fragment.innerHTML = `
      <a href="${raidSimSiteUrl}" class="sim-link ${textKlass}">
        <div class="sim-link-content">
          <img src="${iconPath}" class="sim-link-icon">
          <div class="d-flex flex-column">
            <span class="sim-link-label">${label}</span>
          </div>
        </div>
      </a>
    `;

    return fragment.children[0] as HTMLElement;
  }

  private buildClassLink(classIndex: Class): HTMLElement {
    let specIndexes = getLaunchedSimsForClass(classIndex);
    let href        = specIndexes.length > 1 ? 'javascript:void(0)' : getSpecSiteUrl(specIndexes[0]);
    let textKlass   = this.getContextualKlass({type: 'Class', index: classIndex});
    let iconPath    = this.getSimIconPath({type: 'Class', index: classIndex});
    let label       = classNames[classIndex];

    let fragment = document.createElement('fragment');
    fragment.innerHTML = `
      <a href="${href}" class="sim-link ${textKlass}" ${specIndexes.length > 1 ? 'role="button" data-bs-toggle="dropdown" aria-expanded="false"' : ''}>
        <div class="sim-link-content">
          <img src="${iconPath}" class="sim-link-icon">
          <div class="d-flex flex-column">
            <span class="sim-link-label">${label}</span>
          </div>
        </div>
      </a>
    `;

    return fragment.children[0] as HTMLElement;
  }

  private buildSpecLink(specIndex: Spec): HTMLElement {
    let href      = getSpecSiteUrl(specIndex);
    let textKlass = this.getContextualKlass({type: 'Spec', index: specIndex});
    let iconPath  = this.getSimIconPath({type: 'Spec', index: specIndex});
    let className = classNames[specToClass[specIndex]];
    let specLabel = this.specLabels[specIndex];

    let fragment = document.createElement('fragment');
    fragment.innerHTML = `
      <a href="${href}" class="sim-link ${textKlass}" role="button">
        <div class="sim-link-content">
          <img src="${iconPath}" class="sim-link-icon">
          <div class="d-flex flex-column">
            <span class="sim-link-label">${className}</span>
            <span class="sim-link-label">${specLabel}</span>
          </div>
        </div>
      </a>
    `;

    return fragment.children[0] as HTMLElement;
  }

  private getSimIconPath(data: ClassOptions | SpecOptions | RaidOptions): string {
    let iconPath: string;

    if (data.type == 'Raid')
      iconPath = raidSimIcon;
    else if (data.type == 'Class')
      iconPath = getClassIcon(data.index);
    else
      iconPath = getSpecIcon(data.index);

    return iconPath;
  }

  private getContextualKlass(data: ClassOptions | SpecOptions | RaidOptions): string {
    let klass: string;

    if (data.type == 'Raid')
      // Raid link
      klass = 'text-white';
    else if (data.type == 'Class')
      // Class links
      klass = `text-${classNames[data.index].toLowerCase().replace(/\s/g, '-')}`;
    else
      // Spec links
      klass = `text-${classNames[specToClass[data.index]].toLowerCase().replace(/\s/g, '-')}`;
    
    return klass;
  }
}
